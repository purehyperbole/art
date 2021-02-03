test:
	go test -v --race ./...

bench:
	go test -v -bench=. -benchtime=1000000x

deps:
	go get github.com/stretchr/testify
	go get github.com/google/uuid
