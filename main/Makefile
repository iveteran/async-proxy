GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=fap

VERSION=1.0.1
BUILD=`date +%FT%T%z`

# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildNo=${BUILD}"

all: build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

install: build
	install -v $(BINARY_NAME) ../bin

#deps:
#	$(GOGET) github.com/markbates/goth
#	$(GOGET) github.com/markbates/pop
