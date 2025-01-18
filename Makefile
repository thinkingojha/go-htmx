build: 
	@go build -o ./bin/gohtmx main.go

run: build
	@./bin/gohtmx

test:
	@ go test ./...

tidy:
	@ go fmt ./...
	@ go get -u ./...
	@ go mod tidy