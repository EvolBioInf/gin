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
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type gosym struct {
	g string
	s []string
}
type gosymSlice []*gosym
type GOterm struct {
	a, d, c      string
	o, e, pm, pp float64
}
type GOslice []*GOterm

func (g gosymSlice) Len() int {
	return len(g)
}
func (g gosymSlice) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}
func (g gosymSlice) Less(i, j int) bool {
	l1 := len(g[i].s)
	l2 := len(g[j].s)
	if l1 != l2 {
		return l1 < l2
	}
	return g[i].g < g[j].g
}
func (g GOslice) Len() int { return len(g) }
func (g GOslice) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}
func (g GOslice) Less(i, j int) bool {
	if g[i].pm != g[j].pm {
		return g[i].pm < g[j].pm
	}
	er1 := g[i].o / g[i].e
	er2 := g[j].o / g[j].e
	if er1 != er2 {
		return er1 > er2
	}
	return g[i].a < g[j].a
}
func gid2go(ids map[string]bool,
	id2go map[string]map[string]bool) map[string]int {
	og := make(map[string]int)
	for id, _ := range ids {
		gm := id2go[id]
		if gm == nil {
			continue
		}
		for g, _ := range gm {
			og[g]++
		}
	}
	return og
}
func scan(r io.Reader, args ...interface{}) {
	obsGOcounts := args[0].(map[string]int)
	id2go := args[1].(map[string]map[string]bool)
	go2descr := args[2].(map[string]string)
	go2cat := args[3].(map[string]string)
	minOcc := args[4].(int)
	raw := args[5].(bool)
	sc := bufio.NewScanner(r)
	n := 0
	ids := make(map[string]bool)
	expGOcounts := make(map[string]int)
	eCounts := make(map[string]int)
	eVals := make(map[string][]int)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) > 1 {
			log.Fatalf("can't parse %q", sc.Text())
		}
		if sc.Text() == "" {
			n++
			eg := gid2go(ids, id2go)
			for k, o := range obsGOcounts {
				e := eg[k]
				expGOcounts[k] += e
				eVals[k] = append(eVals[k], e)
				if e >= o {
					eCounts[k]++
				}
			}
			for id, _ := range ids {
				delete(ids, id)
			}
		} else {
			ids[fields[0]] = true
		}
	}
	n++
	eg := gid2go(ids, id2go)
	for k, o := range obsGOcounts {
		e := eg[k]
		expGOcounts[k] += e
		eVals[k] = append(eVals[k], e)
		if e >= o {
			eCounts[k]++
		}
	}
	GOterms := make([]*GOterm, 0)
	for a, o := range obsGOcounts {
		gt := new(GOterm)
		pm := float64(eCounts[a]) / float64(n)
		if pm == 0 {
			pm = -1
		}
		var mean, sd float64
		data := eVals[a]
		nn := len(data)
		for i := 0; i < nn; i++ {
			mean += float64(data[i])
		}
		mean /= float64(nn)
		v := 0.0
		for i := 0; i < nn; i++ {
			s := mean - float64(data[i])
			v += s * s
		}
		v /= float64(nn - 1)
		sd = math.Sqrt(v)
		x := (float64(o) - mean) / (sd * math.Sqrt(2))
		pp := 1.0 - 0.5*(1.0+math.Erf(x))
		gt.pm = pm
		gt.pp = pp
		gt.a = a
		gt.d = go2descr[a]
		gt.c = go2cat[a]
		gt.o = float64(o)
		gt.e = float64(expGOcounts[a]) / float64(n)
		GOterms = append(GOterms, gt)
	}
	sort.Sort(GOslice(GOterms))
	w := tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', 0)
	fmt.Fprintf(w, "#GO\tO\tE\tO/E\tP_m(n=%.1e)\t", float64(n))
	fmt.Fprint(w, "P_p\tCategory\tDescription\n")
	for _, g := range GOterms {
		if int(g.o) < minOcc {
			continue
		}
		fmt.Fprintf(w,
			"%s\t%d\t%.3g\t%.3g\t%.3g\t%.3g\t%s\t%s\n",
			g.a, int(g.o), g.e, g.o/g.e, g.pm, g.pp, g.c, g.d)
	}
	w.Flush()
	if raw {
		for _, term := range GOterms {
			a := term.a
			fmt.Printf("R %s", a)
			counts := eVals[a]
			for _, count := range counts {
				fmt.Printf(" %d", count)
			}
			fmt.Printf("\n")
		}
	}
}
func main() {
	util.Name("ego")
	u := "ego [-h] -o gene2go [option]... " +
		"obsID.txt [expID.txt]..."
	p := "Calculate enrichment of GO terms for " +
		"observed list of gene IDs, given a " +
		"large number of expected ID lists."
	e := "ego -o gene2go obsID.txt expID.txt"
	clio.Usage(u, p, e)
	var optV = flag.Bool("v", false, "version")
	var optO = flag.String("o", "", "file of mapping "+
		"gene IDs to GO terms")
	var optG = flag.String("g", "",
		"GFF file to print GO/sym table")
	var optM = flag.Int("m", 10, "minimum occupancy of GO term")
	var optR = flag.Bool("r", false, "print raw gene counts")
	flag.Parse()
	if *optV {
		util.Version()
	}
	if *optO == "" {
		m := "please enter a file mapping gene IDs " +
			"to GO accessions using -o"
		log.Fatal(m)
	}
	f := util.Open(*optO)
	defer f.Close()
	sc := bufio.NewScanner(f)
	id2go := make(map[string]map[string]bool)
	go2descr := make(map[string]string)
	go2cat := make(map[string]string)
	for sc.Scan() {
		if sc.Text()[0] == '#' {
			continue
		}
		fields := strings.Split(sc.Text(), "\t")
		if len(fields) != 8 {
			m := "gene2go file has %d columns instead of " +
				"the expected 8; are you using " +
				"the correct file?"
			log.Fatalf(m, len(fields))
		}
		geneId := fields[1]
		goId := fields[2]
		goDescr := fields[5]
		goCat := fields[7]
		if id2go[geneId] == nil {
			id2go[geneId] = make(map[string]bool)
		}
		id2go[geneId][goId] = true
		li := strings.LastIndex(goDescr, "[")
		if li > -1 && strings.Contains(goDescr, "GO:") {
			if li > 0 {
				goDescr = goDescr[:li-1]
			} else {
				goDescr = goDescr[:li]
			}
		}
		go2descr[goId] = goDescr
		go2cat[goId] = goCat
	}
	var id2sym map[string]string
	if *optG != "" {
		id2sym = make(map[string]string)
		f := util.Open(*optG)
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			t := sc.Text()
			if t[0] == '#' {
				continue
			}
			fields := strings.Split(t, "\t")
			if len(fields) != 9 {
				m := "expecting 9 fields in GFF file, " +
					"but you have %d"
				log.Fatalf(m, len(fields))
			}
			if fields[2] == "gene" {
				sym := ""
				gid := ""
				attributes := strings.Split(fields[8], ";")
				for _, attribute := range attributes {
					kv := strings.Split(attribute, "=")
					if kv[0] == "Dbxref" {
						ids := strings.Split(kv[1], ",")
						for _, id := range ids {
							arr := strings.Split(id, ":")
							if arr[0] == "GeneID" {
								gid = arr[1]
							}
						}
					}
					if kv[0] == "Name" {
						sym = kv[1]
					}
				}
				if gid == "" {
					fmt.Fprintf(os.Stderr,
						"couldn't find gene ID in %q\n",
						attributes)
				}
				if sym == "" {
					fmt.Fprintf(os.Stderr, "couldn't find name in %q\n",
						attributes)
				} else if gid == "" {
					log.Fatalf("found name but no gene ID in %q\n",
						attributes)
				}
				id2sym[gid] = sym
			}
		}
	}
	files := flag.Args()
	if len(files) < 1 {
		log.Fatal("please enter a file with observed IDs")
	}
	f = util.Open(files[0])
	defer f.Close()
	sc = bufio.NewScanner(f)
	ids := make(map[string]bool)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) > 1 {
			m := "gene IDs should be " +
				"in a single column"
			log.Fatal(m)
		}
		ids[fields[0]] = true
	}
	if *optG != "" {
		g2s := make(map[string][]string)
		for id, _ := range ids {
			sym := id2sym[id]
			gterms := id2go[id]
			for gterm, _ := range gterms {
				sl := g2s[gterm]
				if sl == nil {
					sl = make([]string, 0)
				}
				sl = append(sl, sym)
				g2s[gterm] = sl
			}
		}
		var gosyms []*gosym
		for g, symbols := range g2s {
			gs := new(gosym)
			gs.g = g
			sort.Strings(symbols)
			gs.s = symbols
			gosyms = append(gosyms, gs)
		}
		sort.Sort(gosymSlice(gosyms))
		w := tabwriter.NewWriter(os.Stdout, 1, 0, 2, ' ', 0)
		for i := len(gosyms) - 1; i >= 0; i-- {
			gosym := gosyms[i]
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t", gosym.g,
				len(gosym.s), go2cat[gosym.g],
				go2descr[gosym.g])
			for i, symbol := range gosym.s {
				if i > 0 {
					fmt.Fprintf(w, " ")
				}
				fmt.Fprintf(w, "%s", symbol)
			}
			fmt.Fprintf(w, "\n")
		}
		w.Flush()
		os.Exit(0)
	}
	obsGOcounts := gid2go(ids, id2go)
	files = files[1:]
	clio.ParseFiles(files, scan, obsGOcounts,
		id2go, go2descr, go2cat, *optM, *optR)
}
