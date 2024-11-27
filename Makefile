GO_FILES = $(shell find . -name "*.go") go.mod

bin/tacxble: main.go $(GO_FILES)
	go build -o $@ ./$<

gomod2nix-generate:
	nix run github:nix-community/gomod2nix generate