BINARY_NAME=websocket-chat

build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

test:
	cd trace && go test -cover

clean:
	go clean
	rm ${BINARY_NAME}
