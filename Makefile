PROJECT=retry-strategies-go

.SILENT: clean

all: fmt lint build clean

build:
	CGO_ENABLED=0 GO15VENDOREXPERIMENT=1 go build -o $(PROJECT) .

fmt:
	go fmt .

lint:
	golint .

clean:
	-rm -rf bin