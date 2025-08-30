# Automatically set parallel jobs based on available CPU cores
MAKEFLAGS += -j$(shell nproc)

OUTPUT_DIR := dist

.DEFAULT_GOAL := build

ARCH_MAP_x86_64 := amd64
ARCH_MAP_arm64 := arm64

build: doot-linux-x86_64 doot-darwin-x86_64 doot-windows-x86_64.exe \
       doot-linux-arm64  doot-darwin-arm64  doot-windows-arm64.exe

codegen: lib/common/cache/Colfer.go

test:
	go test ./test -count=1

check:
	staticcheck ./...

create-prerelease:
	.github/create-prerelease.sh

doot-%: codegen
	@GOOS=$(word 1,$(subst -, ,$*)) \
	GOARCH=$(ARCH_MAP_$(word 2,$(subst -, ,$*))) \
	go build -o $(OUTPUT_DIR)/doot-$*

lib/common/cache/Colfer.go: lib/common/cache/cache.colf
	bin/colf -b lib/common Go lib/common/cache/cache.colf

.PHONY: build test
