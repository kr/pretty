include $(GOROOT)/src/Make.inc

TARG=pretty
GOFILES=\
	formatter.go\
	pretty.go\

include $(GOROOT)/src/Make.pkg
