all: minidns

minidns: main.go Makefile
	GOOS=linux  GOARCH=amd64 go build -o minidns main.go
