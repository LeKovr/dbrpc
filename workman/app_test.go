package workman

import (
	"encoding/json"
	"github.com/jessevdk/go-flags"
	"reflect"
	"testing"

	"github.com/LeKovr/go-base/logger"
)

func getRaw(j string) *json.RawMessage {
	raw := json.RawMessage([]byte(j))
	return &raw
}

// -----------------------------------------------------------------------------

func worker(log *logger.Log) WorkerFunc {
	return func(payload string) Result {

		var args []string
		json.Unmarshal([]byte(payload), &args)

		var res Result
		if args[0] != "success" {
			res = Result{Success: true, Error: getRaw(payload)}
		} else {
			res = Result{Success: false, Result: getRaw(payload)}
		}
		return res
	}
}

// -----------------------------------------------------------------------------

func TestOne(t *testing.T) {

	var cfg Flags
	p := flags.NewParser(&cfg, flags.Default)
	p.Parse()

	log, _ := logger.New(logger.Disable)

	wm, err := New(
		WorkerFunc(worker(log)),
		Config(&cfg),
		Logger(log),
	)
	if err != nil {
		t.Errorf("WM error: %s", err)
	}
	defer wm.Stop()

	wm.Run()

	respChannel := make(chan Result)

	key := []string{"success", "one"}
	payload, _ := json.Marshal(key)
	work := Job{Payload: string(payload), Result: respChannel}

	// Push the work onto the queue.
	wm.JobQueue <- work

	resp := <-respChannel

	res := []byte(*resp.Result)
	if !reflect.DeepEqual(payload, res) {
		t.Errorf("WM run error, got %s wait %s", res, payload)
	}

}
