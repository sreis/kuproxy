
DEPENDECIES=vendor/src/github.com/coreos/go-etcd/

all: setup gb
	gb build all

setup: $(DEPENDECIES)

gb:
	go get github.com/constabulary/gb/...

vendor/src/github.com/coreos/go-etcd/:
	git clone https://github.com/coreos/go-etcd/ vendor/src/github.com/coreos/go-etcd/
