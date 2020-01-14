PROJECTNAME1=$(shell basename "$(PWD)")
PROJECTNAME=$(PROJECTNAME1:go-%=%)
APPNAME=$(PROJECTNAME)
include .env
-include .env.local
GOBASE = $(shell pwd)
# #GOPATH="$(GOBASE)/vendor:$(GOBASE)"
# #GOPATH=$(GOBASE)/vendor:$(GOBASE):$(shell dirname $(GOBASE))
GOPATH2= $(shell dirname $(GOBASE))
GOPATH1= $(shell dirname $(GOPATH2))
GOPATH0= $(shell dirname $(GOPATH1))
GOPATH_= $(shell dirname $(GOPATH0))
GOBIN  = $(GOBASE)/bin
GOFILES= $(wildcard *.go)
BIN    = $(GOPATH)/bin
GOLINT = $(BIN)/golint
GOCYCLO= $(BIN)/gocyclo
GOYOLO = $(BIN)/yolo
GO111MODULE = on
GOPROXY     = $(or $(GOPROXY_CUSTOM),https://athens.azurefd.net)
STDERR      = $(or $(STDERR_CUSTOM),/tmp/.$(PROJECTNAME)-stderr.txt)
PID         = $(or $(PID_CUSTOM),/tmp/.$(PROJECTNAME).pid)
MAKEFLAGS += --silent


.PHONY: godoc format fmt lint cov gocov coverage codecov cyclo bench
all: help

.PHONY: all printvars info help
printvars:
	$(foreach V, $(sort $(filter-out .VARIABLES,$(.VARIABLES))), $(info $(v) = $($(v))) )
	# Simple:
	#   (foreach v, $(filter-out .VARIABLES,$(.VARIABLES)), $(info $(v) = $($(v))) )

print-%:
	@echo $* = $($*)

info:
	@echo "     GOBASE: $(GOBASE)"
	@echo "      GOBIN: $(GOBIN)"
	@echo "     GOROOT: $(GOROOT)"
	@echo "     GOPATH: $(GOPATH)"
	@echo "    GOPATH_: $(GOPATH_)"
	@echo "GO111MODULE: $(GO111MODULE)"
	@echo "    GOPROXY: $(GOPROXY)"
	@echo "PROJECTNAME: $(PROJECTNAME)"
	@echo "    APPNAME: $(APPNAME)"
	@echo "    VERSION: $(VERSION)"
	@echo
	@echo "export GO111MODULE=on"
	@echo "export GOPROXY=$(GOPROXY)"
	@echo "export GOPATH=$(GOPATH)"

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo





## fmt: =`format`, run gofmt tool
fmt: format

## format: run gofmt tool
format: | $(GOBASE)
	@echo "  >  gofmt ..."
	gofmt -l -w -s .

## lint: run golint tool
lint: | $(GOBASE) $(GOLINT)
	@echo "  >  golint ..."
	golint ./...

## cov: =`coverage`, run go coverage test
cov: coverage

## gocov: =`coverage`, run go coverage test
gocov: coverage

## coverage: run go coverage test
coverage:
	@echo "  >  gocov ..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o cover.html
	@open cover.html

## codecov: run go test for codecov; (codecov.io)
codecov:
	@echo "  >  codecov ..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic
	#@bash <(curl -s https://codecov.io/bash) -t $(CODECOV_TOKEN)
	curl -s https://codecov.io/bash | bash -s

## test: run go coverage test
test:
	@echo "  >  test ..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic

## cyclo: run gocyclo tool
cyclo:
	@echo "  >  gocyclo ..."
	gocyclo -top 20 .

## bench: benchmark test
bench:
	@echo "  >  benchmark testing ..."
	go test -bench="." -run=^$ -benchtime=10s ./...

#$(BIN)/godoc: | $(GOBASE)     # # # ❶
#	@echo "  >  installing godoc ..."
#	@GOPATH=$(GOPATH) GO111MODULE=$(GO111MODULE) GOPROXY=$(GOPROXY) \
#	go install golang.org/x/tools/cmd/godoc

## godoc: run godoc server at "localhost;6060"
godoc: # | $(GOBASE) $(BIN)/godoc
	@echo "  >  PWD = $(shell pwd)"
	@echo "  >  started godoc server at :6060: http://localhost:6060/pkg/github.com/hedzr/$(PROJECTNAME1) ..."
	#@echo "  $  GOPATH=$(GOPATH) godoc -http=:6060 -index -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps"
	( cd $(GOPATH_) && pwd && godoc -v -index -http=:6060 -notes '(BUG|TODO|DONE|Deprecated)' -play -timestamps -goroot .; ) 
	# gopkg.in/hedzr/errors.v2.New
	# -goroot $(GOPATH) -index
	# https://medium.com/@elliotchance/godoc-tips-tricks-cda6571549b

