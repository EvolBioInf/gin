VERSION = $(shell bash ../scripts/getVersion.sh)
DATE := $(shell bash ../scripts/getDate.sh)

EXE = sumEgo

all: $(EXE) sumEgo.sh

$(EXE):
	go build -ldflags "-X github.com/evolbioinf/gin/util.version=$(VERSION) -X github.com/evolbioinf/gin/util.date=$(DATE)" $(EXE).go
$(EXE).go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE).go | gofmt > $(EXE).go
test: $(EXE) $(EXE)_test.go sumEgo.sh
	go test -v
$(EXE)_test.go: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -R$(EXE)_test.go | gofmt > $(EXE)_test.go
tangle: $(EXE).go $(EXE)_test.go sumEgo.sh
sumEgo.sh: $(EXE).org
	bash ../scripts/org2nw $(EXE).org | notangle -RsumEgo.sh > sumEgo.sh
clean:
	rm -f $(EXE) *.go sumEgo.sh
