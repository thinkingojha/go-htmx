build: 
	@go build -o ./bin/gohtmx main.go

run: build
	@./bin/gohtmx