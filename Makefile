GOOGLE_APPLICATION_CREDENTIALS ?= $(PWD)/secret/creds.json
CONFIG_PATH ?= $(PWD)/config.json
RECEIVE_ENDPOINT ?= http://host.docker.internal:8081

start:
	GOOGLE_APPLICATION_CREDENTIALS=$(GOOGLE_APPLICATION_CREDENTIALS)\
	CONFIG_PATH=$(CONFIG_PATH) \
	RECEIVE_ENDPOINT=$(RECEIVE_ENDPOINT) \
		go run main.go

docker:
	@echo
	@echo "** Listening on localhost:8080 **"
	@echo
	@echo "** Sending subscription events to $(RECEIVE_ENDPOINT) **"
	@echo
	@docker run --rm -it -p 8080:30980 --name=pubber-subber \
		-v $(GOOGLE_APPLICATION_CREDENTIALS):/app/creds.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/app/creds.json \
		-v $(CONFIG_PATH):/app/config.json \
		-e CONFIG_PATH=/app/config.json \
		-e RECEIVE_ENDPOINT=$(RECEIVE_ENDPOINT) \
 \
		cjea/pubber-subber:latest

build:
	docker build . -t local-registry/pubber-subber

docker-local: build
	@echo
	@echo "** Listening on localhost:8080 **"
	@echo
	@echo "** Sending subscription events to $(RECEIVE_ENDPOINT) **"
	@echo
	docker run --rm -it -p 8080:30980 --name=pubber-subber \
		-v $(GOOGLE_APPLICATION_CREDENTIALS):/app/creds.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/app/creds.json \
		-v $(CONFIG_PATH):/app/config.json \
		-e CONFIG_PATH=/app/config.json \
		-e RECEIVE_ENDPOINT=$(RECEIVE_ENDPOINT) \
 \
		local-registry/pubber-subber
