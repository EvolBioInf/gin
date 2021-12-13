../annotate/annotate -s ../data/hsRefGene.txt ../data/i1.txt | awk '{print $4}' | tr ',' '\n' > ../data/obsSym.txt
../shuffle/shuffle -n 10000 ../data/hsTmpl.txt ../data/i1.txt |
    ../annotate/annotate -s ../data/hsRefGene.txt |
    awk '/^#/&&NR>1{print ""}!/^#/{print $4}' |
    tr ',' '\n' |
    ./ego -i ../data/Homo_sapiens.gene_info -g ../data/gene2go ../data/obsSym.txt
