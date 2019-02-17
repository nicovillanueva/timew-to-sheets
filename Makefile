THIS_FILE := $(lastword $(MAKEFILE_LIST))
GEN=go run gsheets/creds/credgen.go
BIN_NAME=to-sheets

all:
	@echo "no target selected, cowardly aborting"

run:
	timew cat | go run main.go

gen:
	$(GEN)

del:
	$(GEN) -erase

build: gen
	go build -o $(BIN_NAME)

release: build
	tar cfz timew-to-sheets-`cat version.txt`.tar.gz $(BIN_NAME)
	@$(MAKE) -f $(THIS_FILE) clean

clean:
	rm $(BIN_NAME)
	@$(MAKE) -f $(THIS_FILE) del

init:
	ln -s /bin/cat ~/.timewarrior/extensions/cat ; \
	git config core.hooksPath .githooks

install: build
	cp $(BIN_NAME) ~/.timewarrior/extensions/
	rm $(BIN_NAME)
