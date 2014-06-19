GODEP=godep
GOCMD=$(GODEP) go
GOBUILD=$(GOCMD) build

build: deps
	$(GOCMD) build -o pkg/dnsizer

deps:
	$(GODEP) get

clean:
	rm -r pkg

fmt:
	$(GOCMD) fmt

lint:
	golint

vet:
	$(GOCMD) vet
