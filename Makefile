build:
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go && \
	zip lambda-handler.zip bootstrap && \
	rm bootstrap

clean:
	rm lambda-handler.zip