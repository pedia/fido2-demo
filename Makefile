all: demo.key demo.pem demo

demo.key demo.pem: req.cnf
	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
		-keyout demo.key -out demo.pem -extensions san -config req.cnf \
		-subj '/CN=wdemo.com'

demo: main.go
	go build -o demo demo

run: demo.key demo.pem
	go run demo

demo.linux: main.go
	# brew install FiloSottile/musl-cross/musl-cross
	CGO_ENABLED=1 GOOS=linux  GOARCH=amd64 \
		CC=x86_64-linux-musl-gcc  CXX=x86_64-linux-musl-g++ go \
		build -ldflags="-extldflags=-static" -o $@ demo

pub: demo.linux
	scp -r demo.linux s templates root@ss.ali:demo
