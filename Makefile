# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
#latest

download:
	export GOPROXY=https://goproxy.cn
	$(GOCMD) mod download

build:
	$(GOBUILD) cmd/main.go

run:
	export GOPROXY=http://goproxy.cn
	$(GOCMD) run cmd/main.go