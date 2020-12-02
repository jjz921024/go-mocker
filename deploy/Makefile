NAME=mocker-go
VERSION=v1.0

.PHONY: all build clean

all: build

build:
	go build -o ${NAME} ../main.go	
	docker build -t ${NAME}:${VERSION} --no-cache --rm .

clean:
	rm ${NAME}
	docker rmi ${NAME}:${VERSION}
