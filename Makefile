
BIN_DIR=bin

PROG=init-db

all: ${PROG}

init-db: ./cmd/go-git/init-db/main.go
	go build -o ${BIN_DIR}/init-db ./cmd/go-git/init-db/main.go

.PHONY: clean
clean:
	rm -rf ${BIN_DIR}