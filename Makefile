all:
	test -d bin || mkdir bin
	make -C annotate
	cp annotate/annotate bin
.PHONY: doc
doc:
	make -C doc
clean:
	make clean -C src
test:
	make test -C src
