.DEFAULT_GOAL := build

build:
	@export GO111MODULE=on
	@go mod vendor
	@go build -mod=vendor cmd/background-grabber/background-grabber.go

install:
	@mkdir -p $GOPATH/bin
	@mv background-grabber $GOPATH/bin
	@echo `start on startup` > $HOME/.config/upstart/background-grabber.conf
	@echo `task` >> $HOME/.config/upstart/background-grabber.conf
	@echo "exec $GOPATH/bin/background-grabber.conf" >> $HOME/.config/upstart/background-grabber.con