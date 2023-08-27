CC      	= go
PROGRAM		= cidr
prefix		= /usr

.PHONY: build clean distclean install uninstall

build: $(PROGRAM)

$(PROGRAM):
	$(CC) build

clean:
	rm -f $(PROGRAM)

distclean: clean

# https://www.gnu.org/software/make/manual/html_node/DESTDIR.html
install:
	install -D -m 0755 $(PROGRAM) $(DESTDIR)$(prefix)/bin/$(PROGRAM)

uninstall:
	-rm -f $(DESTDIR)$(prefix)/bin/$(PROGRAM)

