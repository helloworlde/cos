build:
	rm -rf output
	mkdir -p output
	GOOS=linux GOARCH=amd64 go build -o output/main main.go
	zip output/main.zip output/main
image:
	GOOS=linux GOARCH=amd64 go build -o main main.go
	zip main.zip main