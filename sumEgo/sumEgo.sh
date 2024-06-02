for a in $(seq 10)
do
    nohup shuffle -n 1000000 genomic.gff iv.txt |
          annotate -c genomic.gff |
          ego -o gene2go obsId.txt > ego${a}.txt &
done
wait
sumEgo ego*.txt
