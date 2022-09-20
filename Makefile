GO ?= go

.PHONY: binary
binary: dist
	$(GO) version
	$(GO) build -trimpath -ldflags "-s -w" -o dist/sugar ./cmd/sugar

dist:
	mkdir $@


linux:
	$(GO) env -w GOOS=linux
	make binary
	$(GO) env -u GOOS

cli-win: dist
	$(GO) version
	$(GO) build -trimpath -ldflags "-s -w" -o ./dist/cli.exe ./cmd/cli
