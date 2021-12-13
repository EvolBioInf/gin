progs = annotate ego shuffle
all:
	test -d bin || mkdir bin
	for prog in $(progs); do \
		make -C $$prog; \
		cp $$prog/$$prog bin; \
	done
.PHONY: doc
doc:
	make -C doc
clean:
	for prog in $(progs) doc; do \
		make clean -C $$prog; \
	done
test:
	for prog in $(progs); do \
		make test -C $$prog; \
	done
