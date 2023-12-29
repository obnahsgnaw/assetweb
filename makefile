export PACKAGE=github.com/obnahsgnaw/assetweb
export INPUT=cmd/main.go
export OUT=out
export APP=asset-web


.PHONY: help
help:base_help build_help test_help version_help

.PHONY: base_help
base_help:
	@echo "usage: make <option> <params>"
	@echo "options and effects:"
	@echo "    help   : Show help"
include ./build/build/makefile
include ./build/test/makefile
include ./build/version/makefile