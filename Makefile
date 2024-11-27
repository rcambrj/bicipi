GO_FILES = $(shell find . -name "*.go") go.mod

bin/tacxble: main.go $(GO_FILES)
	go build -o $@ ./$<
