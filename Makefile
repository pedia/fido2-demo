all: demo.key demo.pem demo

demo.key demo.pem: req.cnf
	openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 -nodes \
		-keyout demo.key -out demo.pem -extensions san -config req.cnf \
		-subj '/CN=wdemo.com'

demo: main.go
	go build -o demo demo

run: demo.key demo.pem
	go run demo
