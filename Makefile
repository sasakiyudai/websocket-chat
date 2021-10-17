NAME=websocket-chat

build:
	go build -o ${NAME} main.go

run: build
	./${NAME}

test:
	cd trace && go test -cover

clean:
	go clean
	rm ${NAME}

.PHONY: build run test clean
