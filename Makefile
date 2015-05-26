
TARGET=bin/kubectl

.PHONY: setup gb all

all: setup gb $(TARGET)

setup: 
	git submodule init
	git submodule update

gb:
	go get github.com/constabulary/gb/...

bin/kubectl:
	gb build all
