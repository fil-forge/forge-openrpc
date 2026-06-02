.PHONY: all clean

all: docs/upload.openrpc.json

docs/upload.openrpc.json:
	go run ./main.go upload > ./docs/upload.openrpc.json

clean:
	rm -f ./docs/*.openrpc.json
