.PHONY: all
all:
	GOOS=linux CGO_ENABLED=0 go build -o ginta *.go
