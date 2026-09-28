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
# Use a heredoc so the awk script can stay readable.
.PHONY: help
help: ## Display this help
	@awk -v prog='$(BINARY)' -f - $(MAKEFILE_LIST) <<'AWK'
	BEGIN {
		FS = ":.*##"
		printf "\nUsage:\n  make <target> \033[36m\033[0m\n"
	}
	/^##@/ {
		section = substr($$0, 5)
		if (!(section in seen)) {
			seen[section] = 1
			order[++count] = section
		}
		next
	}
	/^[a-zA-Z0-9_()$$%-]+:.*?##/ {
		target = $$1
		desc = $$2
		gsub(/\$$\(BINARY\)/, prog, target)
		gsub(/\$$\(BINARY\)/, prog, desc)
		if (section == "") {
			section = "Other"
			if (!(section in seen)) {
				seen[section] = 1
				order[++count] = section
			}
		}
		items[section] = items[section] sprintf("  \033[36m%-15s\033[0m %s\n", target, desc)
	}
	END {
		for (i = 1; i <= count; i++) {
			section = order[i]
			if (items[section] != "") {
				printf "\n\033[1m%s\033[0m\n%s", section, items[section]
			}
		}
	}
	AWK


.PHONY: mk-debug
mk-debug: ## Test that the Makefile is functional
	@echo "Makefile is functional."
	echo "Using Make version: ${MAKE_VERSION}"
