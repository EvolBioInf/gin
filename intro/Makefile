all: cli.sh
cli.sh: intro.org
	bash ../scripts/org2nw intro.org | notangle -Rcli > cli.sh
test: cli.sh
	cp ../data/i1.txt iv.txt
	bash cli.sh
clean:
	rm -rf cli.sh ncbi_dataset* README.md gene2go* witherrors.gaps iv.txt obsId.txt
