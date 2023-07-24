src = annotate.tex shuffle.tex ego.tex intro.tex util.tex

all: $(src) ginDoc.pdf 
%.tex:	../*/%.org ginDoc.tex
	awk -f ../scripts/preWeave.awk $< | bash ../scripts/org2nw | noweave -n -x | sed 's/_test/\\_test/g' > $@
ginDoc.pdf: $(src) ../annotate/iv.tex
	latex ginDoc
	latex ginDoc
	dvips ginDoc -o -q
	ps2pdf -dALLOWPSTRANSPARENCY ginDoc.ps
clean:
	rm -f $(src) *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc
