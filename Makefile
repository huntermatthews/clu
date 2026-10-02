# Clu project Makefile

# Variables
BINARY := clu
MANPAGE := pkg/subcmd/$(BINARY).1
SRCS := $(shell find cmd pkg -name '*.go') $(BINARY).1.md $(MANPAGE) go.mod go.sum Makefile common.mk go-common.mk

include common.mk
include go-common.mk


##
##@ Build
##

dist/$(BINARY): $(MANPAGE)

.PHONY: clean
clean: clean-default
	@rm -f  .go-md2man-installed


.PHONY: install
install: install-default man ## Install $(BINARY) binary, manpage, and documentation
	install -d $(PREFIX)/share/man/man1
	install -d $(PREFIX)/share/doc/$(BINARY)
	install -m 644 $(MANPAGE) $(PREFIX)/share/man/man1/$(BINARY).1
	install -m 644 README.md $(PREFIX)/share/doc/$(BINARY)/README.md


##
##@ Documentation
##

.PHONY: man
man: $(MANPAGE) ## Generate the embedded man page from markdown using go-md2man

$(MANPAGE): $(BINARY).1.md .go-md2man-installed
	go-md2man -in $< -out $@

.go-md2man-installed:
	go install github.com/cpuguy83/go-md2man/v2@latest
	@touch .go-md2man-installed
