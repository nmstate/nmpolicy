.PHONY: all build lint unit-test integration-test docs

CMD=./scripts/make.sh

all: docs
	$(CMD)

build: 
	$(CMD)  --build

lint:
	$(CMD) --lint

unit-test: 
	$(CMD) --unit-test

integration-test:
	$(CMD) --integration-test

docs:
	make -C docs build
