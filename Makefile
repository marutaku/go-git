
BIN_DIR=bin

PROG=init-db update-cache write-tree commit-tree read-tree cat-file show-diff

all: ${PROG}

init-db: ./cmd/go-git/init-db/main.go
	go build -o ${BIN_DIR}/init-db ./cmd/go-git/init-db/main.go

update-cache: ./cmd/go-git/update-cache/main.go
	go build -o ${BIN_DIR}/update-cache ./cmd/go-git/update-cache/main.go

write-tree: ./cmd/go-git/write-tree/main.go
	go build -o ${BIN_DIR}/write-tree ./cmd/go-git/write-tree/main.go

commit-tree: ./cmd/go-git/commit-tree/main.go
	go build -o ${BIN_DIR}/commit-tree ./cmd/go-git/commit-tree/main.go

read-tree: ./cmd/go-git/read-tree/main.go
	go build -o ${BIN_DIR}/read-tree ./cmd/go-git/read-tree/main.go

cat-file: ./cmd/go-git/cat-file/main.go
	go build -o ${BIN_DIR}/cat-file ./cmd/go-git/cat-file/main.go

show-diff: ./cmd/go-git/show-diff/main.go
	go build -o ${BIN_DIR}/show-diff ./cmd/go-git/show-diff/main.go

.PHONY: clean
clean:
	rm -rf ${BIN_DIR}
