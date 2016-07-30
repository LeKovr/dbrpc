// Package workman is a Worker Manager
// Code based on http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/
package workman

import (
	"encoding/json"

	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// Flags is a package flags sample
// in form ready for use with github.com/jessevdk/go-flags
type Flags struct {
	MaxWorkers int `long:"max_worker" default:"2" description:"Number of workers"`
	MaxQueue   int `long:"max_queue"  default:"100" description:"Max queue len"`
}

// -----------------------------------------------------------------------------

// Result stores RPC result
type Result struct {
	Success bool             `json:"success"`
	Error   interface{}      `json:"error,omitempty"`
	Result  *json.RawMessage `json:"result,omitempty"`
}

// Job represents the job to be run
type Job struct {
	Payload string
	Result  chan Result
}

// WorkerFunc is a function to be run by workers
type WorkerFunc func(payload string) Result

// -----------------------------------------------------------------------------

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool  chan chan Job
	WorkerFunc  WorkerFunc
	JobChannel  chan Job
	QuitChannel chan bool
	Log         *logger.Log
}

// NewWorker creates a worker instance
func NewWorker(wf WorkerFunc, workerPool chan chan Job, log *logger.Log) Worker {
	return Worker{
		WorkerFunc:  wf,
		WorkerPool:  workerPool,
		JobChannel:  make(chan Job),
		QuitChannel: make(chan bool),
		Log:         log,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start(id int) {
	go func(id int) {
		w.Log.Debugf("Worker %d started", id)
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				res := w.WorkerFunc(job.Payload)
				job.Result <- res
			case <-w.QuitChannel:
				// we have received a signal to stop
				w.Log.Debugf("Worker %d exiting", id)
				return
			}
		}
	}(id)
}

// -----------------------------------------------------------------------------

// WorkMan struct holds worker manager attributes
type WorkMan struct {
	JobQueue   chan Job       // A buffered channel that we can send work requests on.
	WorkerPool chan chan Job  // A pool of workers channels that are registered with the WorkMan
	QuitPool   chan chan bool // A pool of channels listens for quit
	WorkerFunc WorkerFunc
	Config     *Flags
	Log        *logger.Log
}

// -----------------------------------------------------------------------------
// Functional options

// Config sets store config from flag var
func Config(c *Flags) func(wm *WorkMan) error {
	return func(wm *WorkMan) error {
		return wm.setConfig(c)
	}
}

// Logger sets logger
func Logger(l *logger.Log) func(wm *WorkMan) error {
	return func(wm *WorkMan) error {
		return wm.setLogger(l)
	}
}

// -----------------------------------------------------------------------------
// Internal setters

func (wm *WorkMan) setConfig(c *Flags) error {
	wm.Config = c
	return nil
}

func (wm *WorkMan) setLogger(l *logger.Log) error {
	wm.Log = l.WithField("in", "workman")
	return nil
}

// -----------------------------------------------------------------------------

// Dump object fields
func (wm *WorkMan) Dump() error {
	wm.Log.Debugf("wm: %+v", wm)
	return nil
}

// -----------------------------------------------------------------------------

// New - constructor
func New(workerFunc WorkerFunc, options ...func(*WorkMan) error) (*WorkMan, error) {

	wm := WorkMan{WorkerFunc: workerFunc}
	for _, option := range options {
		err := option(&wm)
		if err != nil {
			return nil, err
		}
	}
	wm.JobQueue = make(chan Job, wm.Config.MaxQueue)
	wm.WorkerPool = make(chan chan Job, wm.Config.MaxWorkers)
	wm.QuitPool = make(chan chan bool, wm.Config.MaxWorkers)

	// If log not set - create default
	if wm.Log == nil {
		log, err := logger.New()
		if err != nil {
			return nil, err
		}
		err = wm.setLogger(log)
		if err != nil {
			return nil, err
		}
	}
	return &wm, nil
}

// -----------------------------------------------------------------------------

// Run starts a set of workers
func (wm *WorkMan) Run() {
	// starting n number of workers
	for i := 0; i < wm.Config.MaxWorkers; i++ {
		worker := NewWorker(wm.WorkerFunc, wm.WorkerPool, wm.Log)
		worker.Start(i)
		wm.QuitPool <- worker.QuitChannel
	}

	go wm.dispatch()
}

// -----------------------------------------------------------------------------

// Stop stops all started workers
func (wm *WorkMan) Stop() {
	// send n number of quits
	wm.Log.Debugf("Stopping workers (%d)", len(wm.QuitPool))
	for q := range wm.QuitPool {
		q <- true
	}
}

// -----------------------------------------------------------------------------

func (wm *WorkMan) dispatch() {
	for {
		select {
		case job := <-wm.JobQueue:
			// a job request has been received

			// try to obtain a worker job channel that is available.
			// this will block until a worker is idle
			jobChannel := <-wm.WorkerPool

			// dispatch the job to the worker job channel
			jobChannel <- job
		}
	}
}
