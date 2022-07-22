#
# Makefile for running pandoc on all Markdown docs ending in .md
#
PROJECT = pdtk

MD_PAGES = $(shell ls -1 *.md)

HTML_PAGES = $(shell ls -1 *.md | sed -E 's/.md/.html/g')

build: $(HTML_PAGES) $(MD_PAGES) LICENSE.html

$(HTML_PAGES): $(MD_PAGES) .FORCE
	pandoc -s --to html5 $(basename $@).md -o $(basename $@).html \
	    --template=page.tmpl
	@if [ $@ = "README.html" ]; then mv README.html index.html; fi

LICENSE.html: LICENSE
	pandoc -s --from Markdown --to html5 LICENSE -o license.html \
	    --template=page.tmpl

clean:
	@if [ -f index.html ]; then rm *.html; fi
	#@if [ -f docs/index.html ]; then rm docs/*.html; fi

.FORCE:
