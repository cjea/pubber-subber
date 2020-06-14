start:
	GOOGLE_APPLICATION_CREDENTIALS=secret/liveramp-eng-nexus-649fc5bb5f70.json \
	CONFIG_PATH=config.json \
	RECEIVE_ENDPOINT=http://localhost:8080 \
		go run main.go
