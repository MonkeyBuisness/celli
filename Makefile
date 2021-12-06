APP_VERSION=0.1.0
APP_NAME=celli
OUT_DIR=./out

## $(1) - output file name
## $(2) - GOOS
## $(3) - GOARCH
define build
	env GOOS=$(2) GOARCH=$(3) go build \
	-ldflags="-X main.appVersion=$(APP_VERSION) -X main.appName=$(APP_NAME)" \
	-o $(OUT_DIR)/$(APP_NAME)_$(1)
endef

all: run

run:
# go run -ldflags="-X main.appVersion=$(APP_VERSION)" main.go version 
# go run -ldflags="-X main.appVersion=$(APP_VERSION)" main.go convert book2tpl ./test.javabook > test.md  
# go run -ldflags="-X main.appVersion=$(APP_VERSION)" main.go new javabook --dst ./test/data
	go run -ldflags="-X main.appVersion=$(APP_VERSION)" main.go convert tpl2book ./test.md > my.javabook

build:
	rm -rf $(OUT_DIR)
	$(call build,darwin_amd64,darwin,amd64)
	$(call build,darwin_arm64,darwin,arm64)
	$(call build,linux_amd64,linux,amd64)
	$(call build,linux_arm64,linux,arm64)
	$(call build,linux_arm,linux,arm)
	$(call build,windows_amd64.exe,windows,amd64)