start:
	@echo Find the logs at "docker logs --follow pubber-subber"
	docker run -d -p 8080:80 --name=pubber-subber kennethreitz/httpbin
	GOOGLE_APPLICATION_CREDENTIALS=secret/creds.json \
	CONFIG_PATH=config.json \
	RECEIVE_ENDPOINT=http://host.docker.internal:8080/post \
		go run main.go

run:
	@echo
	@echo "** Listening on localhost:8080 **"
	@echo
	@echo "** Sending subscription events to localhost:8081 **"
	@echo
	@docker run --rm -it -p 8080:30980 --name=pubber-subber \
		-v $(PWD)/secret/creds.json:/app/creds.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/app/creds.json \
		-v $(PWD)/config.json:/app/config.json \
		-e CONFIG_PATH=/app/config.json \
		-e RECEIVE_ENDPOINT=http://host.docker.internal:8081 \
 \
		cjea/pubber-subber:latest

build:
	docker build . -t local-registry/pubber-subber

run-local: build
	@echo
	@echo "** Listening on localhost:8080 **"
	@echo
	@echo "** Sending subscription events to localhost:8081 **"
	@echo
	docker run --rm -it -p 8080:30980 --name=pubber-subber \
		-v $(PWD)/secret/creds.json:/app/creds.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/app/creds.json \
		-v $(PWD)/config.json:/app/config.json \
		-e CONFIG_PATH=/app/config.json \
		-e RECEIVE_ENDPOINT=http://host.docker.internal:8081 \
 \
		local-registry/pubber-subber
