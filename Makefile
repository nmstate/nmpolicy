.PHONY: all build lint unit-test integration-test

CMD=./scripts/make.sh

all:
	$(CMD)

build: 
	$(CMD)  --build

lint:
	$(CMD) --lint

unit-test: 
	$(CMD) --unit-test

integration-test:
	$(CMD) --integration-test
