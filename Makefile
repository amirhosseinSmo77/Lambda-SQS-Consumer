build:
	mkdir -p bin/transactionState
	env GOOS=linux go build -ldflags="-s -w" -o main main.go
	zip bin/transactionState/main.zip main
	rm main
