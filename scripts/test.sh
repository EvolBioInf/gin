../annotate/annotate -s -c ../data/hsRefGene.txt ../data/i1.txt > ../data/tmpSym.txt
../shuffle/shuffle -n 10000 ../data/hsTmpl.txt ../data/i1.txt |
    ../annotate/annotate -f -s -c ../data/hsRefGene.txt |
    ../ego/ego -i ../data/Homo_sapiens.gene_info -g ../data/gene2go ../data/tmpSym.txt
rm ../data/tmpSym.txt
