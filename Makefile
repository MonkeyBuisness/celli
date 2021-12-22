APP_VERSION=0.2.0
APP_NAME=celli
OUT_DIR=./out
BUILD_FLAGS=-ldflags="-X main.appVersion=$(APP_VERSION) -X main.appName=$(APP_NAME)"

## $(1) - output file name
## $(2) - GOOS
## $(3) - GOARCH
define build
	env GOOS=$(2) GOARCH=$(3) go build $(BUILD_FLAGS) \
	-o $(OUT_DIR)/$(APP_NAME)_$(1)
endef

all: run

run:
	go run $(BUILD_FLAGS) main.go convert t2b --pretty ./example/story.md > ./example/story.javabook

build:
	rm -rf $(OUT_DIR)
	$(call build,darwin_amd64,darwin,amd64)
	$(call build,darwin_arm64,darwin,arm64)
	$(call build,linux_amd64,linux,amd64)
	$(call build,linux_arm64,linux,arm64)
	$(call build,linux_arm,linux,arm)
	$(call build,windows_amd64.exe,windows,amd64)

test:
	go test -count=1 ./notebook/*

lint:
	golangci-lint cache clean
	golangci-lint run --config .golangci.yml -v --timeout=3m

install:
	go install $(BUILD_FLAGS)
