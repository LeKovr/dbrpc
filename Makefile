##
## Golang application makefile
##
SHELL      = /bin/bash

# application name
PRG       ?= $(shell basename $$PWD)
CFG       ?= .config
SOURCES   ?= *.go */*.go
SOURCEDIR ?= ". workman jwtutil"
SQLSOURCE ?= sql
SQLMASK   ?= [1-8]?_*.sql

# Runtime data
DB_NAME   ?= dbrpc
APP_SITE  ?= localhost:8081
APP_ADDR  ?= $(APP_SITE)

# Default config
OS        ?= linux
ARCH      ?= amd64
DIRDIST   ?= dist
PRGBIN    ?= $(PRG)_$(OS)_$(ARCH)
PRGPATH   ?= $(PRGBIN)
PIDFILE   ?= $(PRGBIN).pid
LOGFILE   ?= $(PRGBIN).log
STAMP     ?= $$(date +%Y-%m-%d_%H:%M.%S)
ALLARCH   ?= "linux/amd64 linux/386 darwin/386 windows/amd64"
#LIBS       = $(shell go list ./... | grep -vE '/(vendor|iface|proto|cmd)/')


JWT_KEY                 ?= $(shell < /dev/urandom tr -dc A-Za-z0-9 | head -c14; echo)
PG_DB_PASS              ?= $(shell < /dev/urandom tr -dc A-Za-z0-9 | head -c14; echo)
AUTH_SAMPLE_SERVICE_KEY ?= $(shell < /dev/urandom tr -dc A-Za-z0-9 | head -c14; echo)
AUTH_SAMPLE_USER_PASS   ?= $(shell < /dev/urandom tr -dc A-Za-z0-9 | head -c14; echo)

# Search .git for commit id fetch
GIT_ROOT  ?= $$([ -d ./.git ] && echo "." || { [ -d ../.git ] && echo ".." ; } || { [ -d ../../.git ] && echo "../.." ; })


##
## Available targets are:
##

# default: show target list
all:
	@grep -A 1 "^##" Makefile

# ------------------------------------------------------------------------------

## build and run daemon
up: build $(PIDFILE)

$(PIDFILE): $(CFG)
	@source $(CFG) && \
  DBC="$$DB_USER:$$DB_PASS@$$DB_ADDR/$$DB_NAME?sslmode=disable" ; \
  nohup ./$(PRGPATH) --log_level debug --db_connect "$$DBC" --http_addr "$$APP_ADDR" --db_schema "$$DBRPC_SCHEMAS" >$(LOGFILE) 2>&1 &

## run in foreground
run: build $(CFG)
	@source $(CFG) && \
  DBC="$$DB_USER:$$DB_PASS@$$DB_ADDR/$$DB_NAME?sslmode=disable" ; \
  ./$(PRGPATH) --log_level debug --db_connect "$$DBC" --http_addr "$$APP_ADDR" --db_schema "$$DBRPC_SCHEMAS" \
  --jwt_key $$JWT_KEY --jwt_issuer dbrpc:login --jwt_issuer dbrpc:api_open

## gracefull code reload
reload: build $(PIDFILE)
	@kill -1 $$(cat $(PIDFILE))

## stop daemon
down:
	@[ -f $(PIDFILE) ] && kill -SIGTERM $$(cat $(PIDFILE)) && rm $(PIDFILE)

## build and show version
ver: build
	./@$(PRGPATH) --version && echo ""

## build app
build: lint vet $(PRGPATH)

## build app for default arch
$(PRGPATH): $(SOURCES)
	@echo "*** $@ ***"
	@[ -d $(GIT_ROOT)/.git ] && GH=`git rev-parse HEAD` || GH=nogit ; \
GOOS=$(OS) GOARCH=$(ARCH) go build -v -o $(PRGBIN) -ldflags \
"-X main.Build=$(STAMP) -X main.Commit=$$GH"

## run go lint
lint:
	@echo "*** $@ ***"
	@for d in "$(SOURCEDIR)" ; do echo $$d && golint $$d/*.go ; done

## run go vet
vet:
	@echo "*** $@ ***"
	@for d in "$(SOURCEDIR)" ; do echo $$d && go vet $$d/*.go ; done
# does not build with go 1.7

# ------------------------------------------------------------------------------
## clean generated files
clean:
	@echo "*** $@ ***"
	@for a in $(ALLARCH) ; do \
  P=$(PRG)_$${a%/*}_$${a#*/} ; \
  [ "$${a%/*}" == "windows" ] && P=$$P.exe ; \
  [ -f $$P ] && rm $$P || true ; \
done ; \
[ -d $(DIRDIST) ] && rm -rf $(DIRDIST) || true

# ------------------------------------------------------------------------------
# Distro making

## build app for all platforms
buildall: lint vet
	@echo "*** $@ ***"
	@[ -d $(GIT_ROOT)/.git ] && GH=`git rev-parse HEAD` || GH=nogit ; \
for a in $(ALLARCH) ; do \
  echo "** $${a%/*} $${a#*/}" ; \
  P=$(PRG)_$${a%/*}_$${a#*/} ; \
  [ "$${a%/*}" == "windows" ] && P=$$P.exe ; \
  GOOS=$${a%/*} GOARCH=$${a#*/} go build -o $$P -ldflags \
  "-X main.Build=$(STAMP) -X main.Commit=$$GH" ; \
done


## create disro files
dist: clean buildall
	@echo "*** $@ ***"
	@[ -d $(DIRDIST) ] || mkdir $(DIRDIST) ; \
sha256sum $(PRG)* > $(DIRDIST)/SHA256SUMS ; \
for a in $(ALLARCH) ; do \
  echo "** $${a%/*} $${a#*/}" ; \
  P=$(PRG)_$${a%/*}_$${a#*/} ; \
  [ "$${a%/*}" == "windows" ] && P1=$$P.exe || P1=$$P ; \
  zip "$(DIRDIST)/$$P.zip" "$$P1" README.md ; \
done


# ------------------------------------------------------------------------------
# Database sample

## create database schema rpc with used objects
db: $(CFG)
	@source $(CFG) && \
  DBC="$$DB_USER:$$DB_PASS@$$DB_ADDR/$$DB_NAME?sslmode=disable" ; \
  cmd() { echo -e 'BEGIN;\n\set ON_ERROR_STOP 1\n\set SCH rpc\nSET SEARCH_PATH = :SCH, public;' && \
  for f in $(SQLSOURCE)/$(SQLMASK); do echo '\i '$$f; done && \
  echo 'COMMIT;' ; } && cmd | psql -d postgres://$$DBC \
    -v SERVICE_KEY=$$AUTH_SAMPLE_SERVICE_KEY \
    -v USER_PASS=$$AUTH_SAMPLE_USER_PASS


## compile stored functions
db-make: SQLMASK = [56]?_*.sql
db-make: db

## compile auth example functions
db-auth: SQLMASK = samples/60_sample_auth.sql
db-auth: db

## drop database schema rpc
db-clean: $(CFG)
	@source $(CFG) && \
  DBC="$$DB_USER:$$DB_PASS@$$DB_ADDR/$$DB_NAME?sslmode=disable" ; \
  echo "DROP SCHEMA IF EXISTS rpc CASCADE;" | psql -d postgres://$$DBC

## drop database schema rpc
db-clean-auth: $(CFG)
	@source $(CFG) && \
  DBC="$$DB_USER:$$DB_PASS@$$DB_ADDR/$$DB_NAME?sslmode=disable" ; \
  echo "DROP SCHEMA IF EXISTS test_auth CASCADE;" | psql -d postgres://$$DBC

## run psql
psql: $(CFG)
	@source $(CFG) && \
  DBC="$$DB_USER:$$DB_PASS@$$DB_ADDR/$$DB_NAME?sslmode=disable" ; \
  psql -d postgres://$$DBC


# ------------------------------------------------------------------------------

define CONFIG_DEF
# dbrpc config file, generated by make $(CFG)
# API server

# External name
APP_SITE=$(APP_SITE)
# Bind addr
APP_ADDR=$(APP_ADDR)
# JWL key
JWT_KEY=$(JWT_KEY)

# Database

# Search path
DBRPC_SCHEMAS=public

# Host
DB_ADDR=localhost
# Name
DB_NAME=$(DB_NAME)
# User
DB_USER=$(DB_NAME)
# Password
DB_PASS=$(PG_DB_PASS)

# Key form sql/samples auth
AUTH_SAMPLE_SERVICE_KEY=$(AUTH_SAMPLE_SERVICE_KEY)
# User password for sql/samples auth
AUTH_SAMPLE_USER_PASS=$(AUTH_SAMPLE_USER_PASS)

endef
export CONFIG_DEF

## init config
config: $(CFG)


## create config file if none
$(CFG):
	@echo "*** $@ ***"
	@[ -f $@ ] || { echo "$$CONFIG_DEF" > $@ ; echo "Warning: Created default $@" ; }

## build and show program help
help: build
	./$(PRGPATH) --help


.PHONY: all run ver buildall clean dist link vet

# ------------------------------------------------------------------------------

# $$PWD используется для того, чтобы текущий каталог был доступен в контейнере по тому же пути
# и относительные тома новых контейнеров могли его использовать
## run docker-compose
dc: docker-compose.yml
	@docker run --rm -t -i \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $$PWD:$$PWD \
  -w $$PWD \
  docker/compose:1.14.0 \
  $(CMD)

