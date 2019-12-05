
.PHONY: all godoc format fmt lint cov gocov coverage codecov cyclo bench

all: help

## fmt: =`format`, run gofmt tool
fmt: format

## format: run gofmt tool
format: | $(GOBASE)
	@echo "  >  gofmt ..."
	gofmt -l -w -s .

## lint: run golint tool
lint: | $(GOBASE) $(GOLINT)
	@echo "  >  golint ..."
	$(GOLINT) ./...

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

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

