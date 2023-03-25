run:
	go run -ldflags "-X main.version=$(shell git describe --tag --abbrev=0)" .
