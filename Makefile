include $(GOROOT)/src/Make.inc

TARG=github.com/kr/pretty.go
GOFILES=\
	formatter.go\
	pretty.go\

include $(GOROOT)/src/Make.pkg
