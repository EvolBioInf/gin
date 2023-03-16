src = annotate.tex shuffle.tex ego.tex

all: $(src) ginDoc.pdf 
%.tex:	../*/%.org ginDoc.tex
	bash ../scripts/org2nw $< | noweave -n -x | sed 's/_test/\\_test/g' > $@
ginDoc.pdf: $(src) 
	latex ginDoc
	latex ginDoc
	dvips ginDoc -o -q
	ps2pdf ginDoc.ps
clean:
	rm -f $(src) *.pdf *.aux *.bbl *.blg *.dvi *.log *.toc
