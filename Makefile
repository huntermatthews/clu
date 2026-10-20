# Clu project Makefile

# Variables
BINARY := clu
MANPAGE := pkg/subcmd/$(BINARY).1
SRCS := $(shell find cmd pkg -name '*.go') $(BINARY).1.md $(MANPAGE) go.mod go.sum Makefile go-common.mk
GO_PKGS := ./cmd/... ./pkg/...

include go-common.mk


##
##@ Build
##

dist/$(BINARY): $(MANPAGE)


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

$(MANPAGE): $(BINARY).1.md dist/go-md2man
	dist/go-md2man -in $< -out $@

dist/go-md2man:
	@mkdir -p $(@D)
	go install -o $@ github.com/cpuguy83/go-md2man/v2@v2.0.7
