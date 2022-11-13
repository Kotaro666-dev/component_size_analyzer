.PHONY: init
init:
	cp config.sample.json config.json

.PHONY: run
run:
	go run main.go
