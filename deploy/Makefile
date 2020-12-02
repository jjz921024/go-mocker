NAME=mocker-go
VERSION=v1.1

.PHONY: all build clean

all: build

build:
	go build -tags netgo -o ${NAME} ../main.go	
	docker build -t ${NAME}:${VERSION} --no-cache --rm .

clean:
	rm ${NAME}
	docker rmi ${NAME}:${VERSION}
