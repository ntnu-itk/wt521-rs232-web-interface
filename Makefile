ARGS=
F_PROG=./wt521-rs232-web-interface

all: deps build run

deps:
	go get

build:
	go build -o $(F_PROG)

run:
	$(F_PROG) $(ARGS)

verbose: build
	$(F_PROG) $(ARGS) -verbose

forever: build
	while :; do echo "[ERROR] Restarting at $$(date '+%F %T'). Reason: $$($(F_PROG) $(ARGS) 2>&1 | tail -n1)" >&2; sleep 5; done

pandaboard: deps
	GOARCH=arm GOARM=7 GOOS=linux go build -o $(F_PROG)
