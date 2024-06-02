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
	"text/tabwriter"
)

type result struct {
	g, c, d string
	o       int
	e, f, p float64
}
type resultSlice []*result

func (r resultSlice) Len() int {
	return len(r)
}
func (r resultSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r resultSlice) Less(i, j int) bool {
	if r[i].p != r[j].p {
		return r[i].p < r[j].p
	}
	if r[i].f != r[j].f {
		return r[i].f > r[j].f
	}
	if r[i].o != r[j].o {
		return r[i].o > r[j].o
	}
	return r[i].g < r[j].g
}
func scan(r io.Reader, args ...interface{}) {
	g2e := args[0].(map[string]float64)
	g2p := args[1].(map[string]float64)
	ns := args[2].(*int)
	g2n := args[3].(map[string]int)
	nt := args[4].(*float64)
	g2o := args[5].(map[string]int)
	g2c := args[6].(map[string]string)
	g2d := args[7].(map[string]string)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		fields := strings.Fields(line)
		if line[0] == '#' {
			(*ns)++
			x := strings.Split(fields[4], "=")[1]
			x = x[0 : len(x)-1]
			n, err := strconv.ParseFloat(x, 64)
			if err != nil {
				log.Fatal(err)
			}
			(*nt) += n
		} else if len(fields) >= 8 {
			gt := fields[0]
			if g2o[gt] == 0 {
				o, err := strconv.Atoi(fields[1])
				if err != nil {
					log.Fatal(err)
				}
				g2o[gt] = o
			}
			e, err := strconv.ParseFloat(fields[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			g2e[gt] += e
			p, err := strconv.ParseFloat(fields[4], 64)
			if err != nil {
				log.Fatal(err)
			}
			if p != -1.0 {
				g2p[gt] += p
				g2n[gt]++
			} else {
				g2p[gt] += 0.0
			}
			if g2c[gt] == "" {
				g2c[gt] = fields[5]
			}
			if g2d[gt] == "" {
				d := strings.Join(fields[6:], " ")
				g2d[gt] = d
			}
		}
	}
}
func main() {
	util.Name("sumEgo")
	u := "sumEgo [-v|-h] [foo.txt]..."
	p := "Summarize output files generated with ego."
	e := "sumEgo ego*.txt"
	clio.Usage(u, p, e)
	optV := flag.Bool("v", false, "version")
	flag.Parse()
	if *optV {
		util.Version()
	}
	g2e := make(map[string]float64)
	g2p := make(map[string]float64)
	ns := 0
	g2n := make(map[string]int)
	nt := 0.
	g2o := make(map[string]int)
	g2c := make(map[string]string)
	g2d := make(map[string]string)
	files := flag.Args()
	clio.ParseFiles(files, scan, g2e, g2p,
		&ns, g2n, &nt, g2o, g2c, g2d)
	for k, _ := range g2e {
		g2e[k] /= float64(ns)
	}
	g2f := make(map[string]float64)
	for k, v := range g2o {
		f := float64(v) / g2e[k]
		g2f[k] = f
	}
	for k, v := range g2p {
		if g2n[k] > 0 {
			v /= float64(g2n[k])
		} else {
			v = -1.
		}
		g2p[k] = v
	}
	results := make([]*result, 0)
	for k, v := range g2o {
		r := new(result)
		r.g = k
		r.o = v
		r.e = g2e[k]
		r.f = g2f[k]
		r.p = g2p[k]
		r.c = g2c[k]
		r.d = g2d[k]
		results = append(results, r)
	}
	sort.Sort(resultSlice(results))
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	h := "#GO\tO\tE\tO/E\tP(n=%.1g)\t" +
		"Category\tDescription\n"
	fmt.Fprintf(w, h, float64(nt))
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%d\t%.3g\t%.3g\t%.3g\t%s\t%s\n",
			r.g, r.o, r.e, r.f, r.p, r.c, r.d)
	}
	w.Flush()
}
