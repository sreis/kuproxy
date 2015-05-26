
.PHONY: setup gb

all: setup gb
	gb build all

setup: 
	git submodule init
	git submodule update

gb:
	go get github.com/constabulary/gb/...
