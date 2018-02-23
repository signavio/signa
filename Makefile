BIN_NAME = signa
CURRENT_TAG := $(shell git describe --tags | sed -e 's/^v//' | sed -e 's/\-.*//')
ARCH = amd64
BUILD_ENTRYPOINT = cmd/signa/main.go

DOCKER_FILE = tools/docker/Dockerfile
DOCKER_REPO = signavio/signa

release:
	GOOS=linux GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-linux-$(ARCH) $(BUILD_ENTRYPOINT)
	GOOS=freebsd GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-freebsd-$(ARCH) $(BUILD_ENTRYPOINT)
	GOOS=darwin GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-darwin-$(ARCH) $(BUILD_ENTRYPOINT)
	GOOS=windows GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-windows-$(ARCH).exe $(BUILD_ENTRYPOINT)

static:
	CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-static-linux-$(ARCH) -a -tags netgo -ldflags '-w' $(BUILD_ENTRYPOINT)
	CGO_ENABLED=0 GOOS=freebsd GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-static-freebsd-$(ARCH) -a -tags netgo -ldflags '-w' $(BUILD_ENTRYPOINT)
	CGO_ENABLED=0 GOOS=darwin GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-static-darwin-$(ARCH) -a -tags netgo -ldflags '-w' $(BUILD_ENTRYPOINT)
	CGO_ENABLED=0 GOOS=windows GOARCH=$(ARCH) go build -o releases/$(BIN_NAME)-$(CURRENT_TAG)-static-windows-$(ARCH).exe -a -tags netgo -ldflags '-w' $(BUILD_ENTRYPOINT)

image:
	docker build . -f $(DOCKER_FILE) -t $(DOCKER_REPO):$(CURRENT_TAG)
	docker tag $(DOCKER_REPO):$(CURRENT_TAG) $(DOCKER_REPO):latest

version:
	@echo $(CURRENT_TAG)

default: release
