clean:
	rm -fv uwncli-linux uwncli-darwin uwncli-windows.exe

linux:
	GOOS=linux go build -ldflags="-s -w" -o uwncli-linux *.go

darwin:
	GOOS=darwin go build -ldflags="-s -w" -o uwncli-darwin *.go

windows:
	GOOS=windows go build -ldflags="-s -w" -o uwncli-windows.exe *.go

all: clean linux darwin windows
