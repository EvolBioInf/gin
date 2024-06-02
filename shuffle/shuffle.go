package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evolbioinf/clio"
	"github.com/evolbioinf/gin/util"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type interval struct {
	chr        string
	start, end int
}

func scan(r io.Reader, args ...interface{}) {
	n := args[0].(int)
	templates := args[1].([]*interval)
	pick := args[2].([]float64)
	ran := args[3].(*rand.Rand)
	first := args[4].(*bool)
	w := args[5].(*bufio.Writer)
	sc := bufio.NewScanner(r)
	intervals := make([]*interval, 0)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		iv := new(interval)
		s, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatalf("can't read %q", fields[1])
		}
		iv.start = s
		e, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Fatalf("can't read %q", fields[2])
		}
		iv.end = e
		intervals = append(intervals, iv)
	}
	for i := 0; i < n; i++ {
		if *first {
			*first = false
		} else {
			fmt.Fprintf(w, "\n")
		}
		for _, iv := range intervals {
			r := ran.Float64()
			i := sort.SearchFloat64s(pick, r)
			template := templates[i]
			tl := template.end - template.start + 1
			m := ran.Intn(tl)
			il := iv.end - iv.start + 1
			d := int(math.Round(float64(il) / 2.0))
			s := m - d
			e := m + d - 1
			if s < 0 {
				s = 0
			}
			if e >= tl {
				e = tl - 1
			}
			c := template.chr
			fmt.Fprintf(w, "%s\t%d\t%d\n", c, s, e)
		}
	}
	w.Flush()
}
func main() {
	util.Name("shuffle")
	u := "shuffle [-h] [option]... foo.gff [focus.txt]..."
	p := "Shuffle a set of focus intervals among the " +
		"template intervals."
	e := "shuffle -n 10000 genomic.gff f1.txt"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "print program version "+
		"and other information")
	var optN = flag.Int("n", 1, "number of iterations")
	var optS = flag.Int("s", 0, "seed for random number "+
		"generator")
	flag.Parse()
	if *optV {
		util.Version()
	}
	seed := int64(*optS)
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	ran := rand.New(rand.NewSource(seed))
	files := flag.Args()
	if len(files) < 1 {
		log.Fatal("please provide a template file")
	}
	templates := make([]*interval, 0)
	file := util.Open(files[0])
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if line[0] == '#' {
			continue
		}
		fields := strings.Fields(line)
		if fields[2] == "region" {
			template := &interval{chr: fields[0]}
			s, err := strconv.Atoi(fields[3])
			if err != nil {
				log.Fatalf("can't read %q", fields[3])
			}
			template.start = s
			e, err := strconv.Atoi(fields[4])
			if err != nil {
				log.Fatalf("can't read %q", fields[4])
			}
			template.end = e
			templates = append(templates, template)
		}
	}
	pick := make([]float64, len(templates))
	sum := 0.0
	for _, t := range templates {
		sum += float64(t.end - t.start + 1)
	}
	pick[0] = float64(templates[0].end-templates[0].start+1) / sum
	for i := 1; i < len(templates); i++ {
		f := float64(templates[i].end-
			templates[i].start+1) / sum
		pick[i] = pick[i-1] + f
	}
	files = files[1:]
	first := true
	w := bufio.NewWriter(os.Stdout)
	clio.ParseFiles(files, scan, *optN, templates,
		pick, ran, &first, w)
}
