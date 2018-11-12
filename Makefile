# Default unit executable extension
EXE=
ifeq ($(OS), Windows_NT)
	EXE=.exe
endif
# Go parameters
GOCMD:=go
GOBUILD:=$(GOCMD) build
GOINSTALL:=$(GOCMD) install
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOVET:=$(GOCMD) vet
GOLINT:=golint
BINARY_UNIX:=resiproxy
BINARY_NAME:=$(BINARY_UNIX)$(EXE)
LDFLAGS:=
ifeq (, $(shell which rice))
 $(error "No rice in your PATH, run: go get github.com/GeertJohan/go.rice/rice ")
endif



.PHONY: all build install coverage lint test clean run deps build-linux preso docs docker-build docker-push build test

## Building
all: test build
deps:
		go mod download

build: $(BINARY_NAME)
$(BINARY_NAME): deps
		CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -ldflags  "-s -w"


install:  deps
		CGO_ENABLED=0 $(GOINSTALL) $(TAGS) -ldflags  "-s -w"


## Dev testing
docker-build:
		docker build -t marqub/resiproxy:dev .
docker-run: docker-build
		docker run -ti --rm marqub/resiproxy:dev
docker-push: docker-build
		docker push marqub/resiproxy:dev
helm: docker-push
		helm upgrade --install  resiproxy --namespace resiproxy --recreate-pods --wait --tiller-namespace=resiproxy charts/resproxy
delete:
		helm delete --purge --tiller-namespace=resiproxy resiproxy | true
## Documentation
docs:
		./gendocs.sh

## Testing
coverage: deps
		$(GOTEST) $(TAGS) -coverprofile .testCoverage.txt ./...  && go tool cover -func=.testCoverage.txt
lint:
		$(GOLINT) -set_exit_status .r/...
vet:
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOVET) $(TAGS) ./...
test:  deps
		$(GOTEST) -short -race $(TAGS) ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
		rm -rf rice-box.go


