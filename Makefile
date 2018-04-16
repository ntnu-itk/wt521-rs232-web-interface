RUN_ARGS=

all: build run

build:
	go build

run:
	./wt521-rs232-web-interface $(RUN_ARGS)

verbose: build
	./wt521-rs232-web-interface $(RUN_ARGS) -verbose