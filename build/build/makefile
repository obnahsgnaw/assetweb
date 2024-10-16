PACKAGE_NAME :=${PACKAGE}
APP_NAME := $(if $(APP),$(APP),"APP")
OUT_NAME := $(if $(OUT),$(OUT),"out")
INPUT_NAME := $(if $(INPUT),$(INPUT),"main.go")
ASSET_VERSION := $(if $(ASSET_VERSION),$(ASSET_VERSION),"v0.0.0")

ARCH=amd64
PLATFORM_WIN=windows
PLATFORM_MAC=darwin
PLATFORM_LINUX=linux

ifeq ($(OS),Windows_NT)
 PLATFORM=${PLATFORM_WIN}
else
 ifeq ($(shell uname),Darwin)
  PLATFORM=${PLATFORM_MAC}
 else
  PLATFORM=${PLATFORM_LINUX}
 endif
endif

LDFLAGS    := -w -s
BUILD_DATE =  $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

ifneq ($(wildcard .git),)
VERSION   =  v0.0.0
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_DIRTY  =  $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")
else
VERSION=v0.0.0
GIT_COMMIT=
GIT_DIRTY=
endif

ifneq ($(PACKAGE_NAME),)
LDFLAGS += -X $(PACKAGE_NAME)/version.name=$(APP_NAME)
LDFLAGS += -X $(PACKAGE_NAME)/version.version=$(VERSION)
LDFLAGS += -X $(PACKAGE_NAME)/version.assetVersion=$(ASSET_VERSION)
LDFLAGS += -X $(PACKAGE_NAME)/version.gitCommit=$(GIT_COMMIT)
LDFLAGS += -X $(PACKAGE_NAME)/version.gitTreeState=$(GIT_DIRTY)
LDFLAGS += -X $(PACKAGE_NAME)/version.buildDate=$(BUILD_DATE)
endif
FILE_NAME=$(APP_NAME)

.PHONY: build_help
build_help:
	@echo "build options: <command> <param>"
	@echo "    build  : Build the binary of this project for current platform"
	@echo "    windows: Build the windows binary of this project"
	@echo "    linux  : Build the linux binary of this project"
	@echo "    darwin : Build the darwin binary of this project"
	@echo "build param:"
	@echo "    ARCH=?   : Set build ARCH, =amd64, =arm64, default=amd64, example: make build ARCH=amd64"

.PHONY: build
build:
	@echo "Building..."
	@go build -o ${OUT_NAME}/${FILE_NAME}_${ARCH} -ldflags "${LDFLAGS}" ${INPUT_NAME}
	@echo "Done"

.PHONY: windows
windows:
ifeq (${PLATFORM},${PLATFORM_WIN})
	@echo "Plateform windows building for windows(ARCH:${ARCH})..."
	@go build -o ${OUT_NAME}/${FILE_NAME}_win_${ARCH}.exe -ldflags "${LDFLAGS}" ${INPUT_NAME}
else
	@echo "Plateform ${PLATFORM} building for windows(ARCH:${ARCH})..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=${ARCH} go build -o ${OUT_NAME}/${FILE_NAME}_win_${ARCH}.exe -ldflags "${LDFLAGS}" ${INPUT_NAME}
endif
	@echo "Done"

.PHONY: linux
linux:
ifeq (${PLATFORM},${PLATFORM_WIN})
	@echo "Plateform windows building for linux(ARCH:${ARCH})..."
	@SET CGO_ENABLE=0
	@SET GOOS=linux
	@SET GOARCH=${ARCH}
	@go build -o ${OUT_NAME}/${FILE_NAME}_linux_${ARCH} -ldflags "${LDFLAGS}" ${INPUT_NAME}
else
ifeq (${PLATFORM},${PLATFORM_MAC})
	@echo "Plateform darwin building for linux(ARCH:${ARCH})..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build -o ${OUT_NAME}/${FILE_NAME}_linux_${ARCH} -ldflags "${LDFLAGS}" ${INPUT_NAME}
else
	@echo "Plateform linux building for linux(ARCH:${ARCH})..."
	@go build -o ${OUT_NAME}/${FILE_NAME}_linux_${ARCH} -ldflags "${LDFLAGS}" ${INPUT_NAME}
endif
endif
	@echo "Done"

.PHONY: darwin
darwin:
ifeq (${PLATFORM},${PLATFORM_WIN})
	@echo "Plateform windows building for darwin(ARCH:${ARCH})..."
	@SET CGO_ENABLE=0
	@SET GOOS=darwin
	@SET GOARCH=${ARCH}
	@go build -o ${OUT_NAME}/${FILE_NAME}_darwin_${ARCH} -ldflags "${LDFLAGS}" ${INPUT_NAME}
else
ifeq (${PLATFORM},${PLATFORM_LINUX})
	@echo "Plateform linux building for darwin(ARCH:${ARCH})..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=${ARCH} go build -o ${OUT_NAME}/${FILE_NAME}_darwin_${ARCH} -ldflags "${LDFLAGS}" ${INPUT_NAME}
else
	@echo "Plateform darwin building for darwin(ARCH:${ARCH})..."
	@go build -o ${OUT_NAME}/${FILE_NAME}_darwin_${ARCH} -ldflags "${LDFLAGS}" ${INPUT_NAME}
endif
endif
	@echo "Done"
