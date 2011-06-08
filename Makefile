GOROOT ?= $(shell printf 't:;@echo $$(GOROOT)\n' | gomake -f -)
include $(GOROOT)/src/Make.inc

TARG=github.com/kr/pretty.go
GOFILES=\
	diff.go\
	formatter.go\
	pretty.go\

include $(GOROOT)/src/Make.pkg
