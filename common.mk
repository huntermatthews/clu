# common.mk
# Common Makefile settings and variables

# version check HAS to be at the top , BEFORE older versions of Make will choke on later syntax
# EL7 ships with 3.82, so this is the current minimum.
# macOS ships with an ancient version of Make (3.81), so sysadmins will need to install a newer version via Homebrew.
MINIMUM_GNU_MAKE := 3.81
ifeq "${MAKE_VERSION}" ""
  $(error This Makefile requires GNU Make $(MINIMUM_GNU_MAKE) or greater)
endif
ifneq "$(MINIMUM_GNU_MAKE)" "$(firstword $(sort $(MINIMUM_GNU_MAKE) ${MAKE_VERSION}))"
  $(error This Makefile requires GNU Make $(MINIMUM_GNU_MAKE) or greater)
endif

# Makefile settings
# Note: --warn-undefined-variables is too noisy with recursive make (warns about GNUMAKEFLAGS)
MAKEFLAGS += --no-builtin-rules --no-builtin-variables --no-print-directory
SHELL := bash
.DEFAULT_GOAL := help
.ONESHELL:

# VERSION determination
# If VERSION is not set in the environment, try to get it from git
ifndef VERSION
  # Check if we're in a git repository
  IS_GIT_REPO := $(shell git rev-parse --is-inside-work-tree 2>/dev/null)
  ifeq "$(IS_GIT_REPO)" "true"
    VERSION := $(shell git describe --dirty --always --match "v[0-9]*")
  else
    $(error VERSION not set and build not in a git repository. Please set a VERSION environment variable.)
  endif
endif
# Strip leading 'v' from VERSION if present
VERSION := $(patsubst v%,%,$(VERSION))

# TODO: Validate VERSION format
# SEMVER_REGEX := ^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$$
# VALID_VERSION := $(shell echo "$(VERSION)" | grep -E '$(SEMVER_REGEX)' >/dev/null && echo yes)
# ifneq ($(VALID_VERSION),yes)
#   $(error VERSION "$(VERSION)" is not valid semantic versioning format)
# endif


##
##@ Help
##
.PHONY: help
help: ## Display this help
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "\nUsage:\n  make <target> \033[36m\033[0m\n" \
	} \
	/^[a-zA-Z0-9_%-]+:.*?##/ { \
		printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
	}' $(MAKEFILE_LIST)


.PHONY: mk-debug
mk-debug: ## Test that the Makefile is functional
	@echo "Makefile is functional."
	echo "Using Make version: ${MAKE_VERSION}"
	echo "VERSION: ${VERSION}"
