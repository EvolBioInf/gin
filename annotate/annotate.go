package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/gin/util"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type interval struct {
	start, end int
	id         string
}
type ivSlice []*interval

func scan(r io.Reader, args ...interface{}) {
	genes := args[0].(map[string][]*interval)
	maxGeneLengths := args[1].(map[string]int)
	w := args[2].(*bufio.Writer)
	header := args[3].(string)
	col := args[4].(bool)
	sc := bufio.NewScanner(r)
	ids := make([]string, 0)
	if !col {
		fmt.Fprint(w, header)
	}
	for sc.Scan() {
		if len(sc.Text()) == 0 {
			fmt.Fprint(w, header)
			continue
		}
		fields := strings.Fields(sc.Text())
		chr := fields[0]
		s, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal(err.Error())
		}
		e, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Fatal(err.Error())
		}
		g := genes[chr]
		i := sort.Search(len(g), func(i int) bool {
			return g[i].end >= s
		})
		mgl := maxGeneLengths[chr]
		ids = ids[:0]
		for i < len(g) {
			d := g[i].start - e + 1
			if d > mgl {
				break
			}
			if g[i].start <= e {
				ids = append(ids, g[i].id)
			}
			i++
		}
		if len(ids) > 0 {
			if !col {
				fmt.Fprintf(w, "%s\t%d\t%d\t%s", chr, s, e, ids[0])
				for i := 1; i < len(ids); i++ {
					fmt.Fprintf(w, ",%s", ids[i])
				}
				fmt.Fprint(w, "\n")
			} else {
				for _, id := range ids {
					fmt.Fprintf(w, "%s\n", id)
				}
			}
		}
	}
	w.Flush()
}
func (s ivSlice) Len() int { return len(s) }
func (s ivSlice) Less(i, j int) bool {
	return s[i].end < s[j].end
}
func (s ivSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func main() {
	util.Name("annotate")
	u := "annotate [-h] [option]... foo.gff [intervals.txt]..."
	p := "Annotate genome intervals with intersecting genes."
	e := "annotate genomic.gff iv.txt"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optT = flag.Bool("t", false, "intersect transcript "+
		"instead of promoter")
	var optL = flag.Int("l", 2000, "promoter length")
	var optS = flag.Bool("s", false, "gene symbols instead of "+
		"gene IDs")
	var optC = flag.Bool("c", false, "print gene IDs "+
		"or gene symbols in a single column")
	flag.Parse()
	if *optV {
		util.Version()
	}
	files := flag.Args()
	if len(files) < 1 {
		log.Fatal("please supply a GFF file")
	}
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
	}
	genes := make(map[string][]*interval)
	f := util.Open(files[0])
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		t := sc.Text()
		if t[0] == '#' {
			continue
		}
		fields := strings.Split(sc.Text(), "\t")
		if len(fields) > 2 && fields[2] == "gene" {
			gene := new(interval)
			chr := fields[0]
			tss, err := strconv.Atoi(fields[3])
			if err != nil {
				log.Fatal(err)
			}
			tes, err := strconv.Atoi(fields[4])
			if err != nil {
				log.Fatal(err)
			}
			strand := fields[6]
			attributes := strings.Split(fields[8], ";")
			var gid, name string
			for _, attribute := range attributes {
				arr := strings.Split(attribute, "=")
				if arr[0] == "Name" {
					name = arr[1]
				}
				if arr[0] == "Dbxref" {
					ids := strings.Split(arr[1], ",")
					for _, id := range ids {
						arr := strings.Split(id, ":")
						if arr[0] == "GeneID" {
							gid = arr[1]
							break
						}
					}
					if gid == "" {
						log.Fatal("couldn't find GeneID")
					}
				}
			}
			gene.id = gid
			if *optS {
				gene.id = name
			}
			if *optT {
				gene.start = tss
				gene.end = tes
			} else {
				if strand == "+" {
					gene.start = tss - *optL + 1
					gene.end = tss
				} else {
					gene.start = tes
					gene.end = tes + *optL - 1
				}
			}
			if gene.end-gene.start >= 0 {
				if genes[chr] == nil {
					genes[chr] = make([]*interval, 0)
				}
				genes[chr] = append(genes[chr], gene)
			}
		}
	}
	f.Close()
	for _, v := range genes {
		is := ivSlice(v)
		sort.Sort(is)
	}
	maxGeneLengths := make(map[string]int)
	for chr, gs := range genes {
		maxGeneLengths[chr] = -1
		for _, g := range gs {
			l := g.end - g.start + 1
			if maxGeneLengths[chr] < l {
				maxGeneLengths[chr] = l
			}
		}
	}
	files = files[1:]
	w := bufio.NewWriter(os.Stdout)
	header := ""
	header = "\n"
	if !*optC {
		header = "#Chr\tStart\tEnd\t"
		if *optS {
			header += "Sym...\n"
		} else {
			header += "ID...\n"
		}
	}
	clio.ParseFiles(files, scan, genes, maxGeneLengths,
		w, header, *optC)
}
