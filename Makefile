##
## Golang application makefile
##

# application name
PRG       ?= $(shell basename $$PWD)
SOURCES   ?= *.go */*.go
SOURCEDIR ?= ". workman"

# Default config
SHELL      = /bin/bash
OS        ?= linux
ARCH      ?= amd64
DIRDIST   ?= dist
PRGBIN    ?= $(PRG)_$(OS)_$(ARCH)
PRGPATH   ?= $(PRGBIN)
PIDFILE   ?= $(PRGBIN).pid
LOGFILE   ?= $(PRGBIN).log
STAMP     ?= $$(date +%Y-%m-%d_%H:%M.%S)
ALLARCH   ?= "linux/amd64 linux/386 windows/amd64 darwin/386"

# Search .git for commit id fetch
GIT_ROOT  ?= $$([ -d ./.git ] && echo "." || { [ -d ../.git ] && echo ".." ; } || { [ -d ../../.git ] && echo "../.." ; })

DBCONN ?= op:op@localhost/op?sslmode=disable

##
## Available targets are:
##

# default: show target list
all:
	@grep -A 1 "^##" Makefile

## build and run
up: build $(PIDFILE)

$(PIDFILE):
	@nohup ./$(PRGPATH) --log_level debug --db_connect $(DBCONN) >$(LOGFILE) 2>&1 &

run: build
	./$(PRGPATH) --log_level debug --db_connect $(DBCONN)

reload: build $(PIDFILE)
	@kill -1 $$(cat $(PIDFILE))

down:
	@[ -f $(PIDFILE) ] && kill -SIGTERM $$(cat $(PIDFILE)) && rm $(PIDFILE)

## build and show help
help: build
	./$(PRGPATH) --help

## build and show version
ver: build
	./@$(PRGPATH) --version && echo ""

## build app
build: lint vet $(PRGPATH)

## build app for default arch
$(PRGPATH): $(SOURCES)
	@echo "*** $@ ***"
	@[ -d $(GIT_ROOT)/.git ] && GH=`git rev-parse HEAD` || GH=nogit ; \
GOOS=$(OS) GOARCH=$(ARCH) go build -o $(PRGBIN) -ldflags \
"-X main.Build=$(STAMP) -X main.Commit=$$GH"

## build app for all platforms
buildall: lint vet
	@echo "*** $@ ***"
	@[ -d $(GIT_ROOT)/.git ] && GH=`git rev-parse HEAD` || GH=nogit ; \
for a in "$(ALLARCH)" ; do \
  echo "** $${a%/*} $${a#*/}" ; \
  P=$(PRG)_$${a%/*}_$${a#*/} ; \
  [ "$${a%/*}" == "windows" ] && P=$$P.exe ; \
  GOOS=$${a%/*} GOARCH=$${a#*/} go build -o $$P -ldflags \
  "-X main.Build=$(STAMP) -X main.Commit=$$GH" ; \
done

# create database
db:
	PGPASSWORD=op psql -h localhost -U op -f crebas.sql op

## create disro files
dist: clean buildall
	@echo "*** $@ ***"
	@[ -d $(DIRDIST) ] || mkdir $(DIRDIST) ; \
sha256sum $(PRG)* > $(DIRDIST)/SHA256SUMS ; \
for a in "$(ALLARCH)" ; do \
  echo "** $${a%/*} $${a#*/}" ; \
  P=$(PRG)_$${a%/*}_$${a#*/} ; \
  [ "$${a%/*}" == "windows" ] && P1=$$P.exe || P1=$$P ; \
  zip "$(DIRDIST)/$$P.zip" "$$P1" README.md ; \
done

## clean generated files
clean:
	@echo "*** $@ ***"
	@for a in "$(ALLARCH)" ; do \
  P=$(PRG)_$${a%/*}_$${a#*/} ; \
  [ "$${a%/*}" == "windows" ] && P=$$P.exe ; \
  [ -f $$P ] && rm $$P || true ; \
done ; \
[ -d $(DIRDIST) ] && rm -rf $(DIRDIST) || true

## run go lint
lint:
	@echo "*** $@ ***"
	@for d in "$(SOURCEDIR)" ; do echo $$d && golint $$d/*.go ; done

## run go vet
vet:
	@echo "*** $@ ***"
	@for d in "$(SOURCEDIR)" ; do echo $$d && go vet $$d/*.go ; done

.PHONY: all run ver buildall clean dist link vet

