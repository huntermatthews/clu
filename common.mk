# common.mk
# Common Makefile settings and variables

# version check HAS to be at the top , BEFORE older versions of Make will choke on later syntax
# EL7 ships with 3.82, so this is the current minimum.
# macOS ships with an ancient version of Make (3.81), so you will need to install a newer version via Homebrew.
# 3.81 is missing the .ONESHELL feature which is critical for simple target happiness.

MINIMUM_GNU_MAKE := 3.82
ifeq "${MAKE_VERSION}" ""
  $(error This Makefile requires GNU Make $(MINIMUM_GNU_MAKE) or greater)
endif
ifneq "$(MINIMUM_GNU_MAKE)" "$(firstword $(sort $(MINIMUM_GNU_MAKE) ${MAKE_VERSION}))"
  $(error This Makefile requires GNU Make $(MINIMUM_GNU_MAKE) or greater)
endif

# Makefile settings
MAKEFLAGS += --no-builtin-rules --no-builtin-variables --no-print-directory --warn-undefined-variables
# if we want --warn-undefined-variables, we have to declare GNUMAKEFLAGS as empty to prevent a warning
GNUMAKEFLAGS ?=
SHELL := bash
.DEFAULT_GOAL := help
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c

##
##@ Help
##
# Leave the backslashes - even with .ONESHELL: they are needed for the awk script
.PHONY: help
help: ## Display this help
	@awk -v prog=$(BINARY) 'BEGIN { \
		FS = ":.*##"; \
		printf "\nUsage:\n  make <target> \033[36m\033[0m\n" \
	} \
	# make: double-dollar escapes to one dollar; backslash makes awk treat it as literal in char class \
	/^[a-zA-Z0-9_()\$$%-]+:.*?##/ { \
		# expand $(BINARY) to its value in both the target name and description \
		gsub(/\$$\(BINARY\)/, prog, $$1); gsub(/\$$\(BINARY\)/, prog, $$2); \
		printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	}' $(MAKEFILE_LIST)


.PHONY: mk-debug
mk-debug: ## Test that the Makefile is functional
	@echo "Makefile is functional."
	echo "Using Make version: ${MAKE_VERSION}"
