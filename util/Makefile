NW = $(shell which noweb)

all : util

util: util.go
	go build util.go
util.go: util.org
	if [ "$(NW)" != "" ]; then\
		awk -f ../scripts/preTangle.awk util.org | bash ../scripts/org2nw | notangle -Rutil.go | gofmt > util.go;\
	fi
.PHONY: test
test: util_test.go util.go
	go test -v
util_test.go: util.org
	if [ "$(NW)" != "" ]; then\
		awk -f ../scripts/preTangle.awk util.org | bash ../scripts/org2nw | notangle -Rutil_test.go | gofmt > util_test.go;\
	fi
tangle: util.go util_test.go

clean:
	rm -f *.go
