include $(GOROOT)/src/Make.inc

LDFLAGS=-e
TARG=bencode

GOFILES=\
		decoder.go\
		encoder.go

include $(GOROOT)/src/Make.pkg


EXBIN=example
EXOBJ=example.$O
EXSRC=example.go

ex:
	$(GC) -o $(EXOBJ) $(EXSRC)
	$(LD) -o $(EXBIN) $(EXOBJ)

exclean:
	rm -f $(EXBIN) $(EXOBJ)

spotless: clean exclean
