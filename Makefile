BINARY_NAME=websocket-chat

build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}
