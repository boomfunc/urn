# Import helper makefiles.
# All commands will be available in sub-modules.
include .make/golang.mk


.PHONY: lint
lint: golang-fmt


.PHONY: test
test: golang-test
