#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = pttk

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROGRAMS = $(shell ls -1 cmd)

PACKAGE = $(shell ls -1 *.go | grep -v 'version.go')

SUBPACKAGES = $(shell ls -1 */*.go)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

ifneq ($(prefix),)
        PREFIX = $(prefix)
endif

OS = $(shell uname)

EXT =
ifeq ($(OS), Windows)
	EXT = .exe
endif

DIST_FOLDERS = bin/*

build: version.go $(PROGRAMS) CITATION.cff about.md installer.sh

version.go: .FORCE
	@echo "package $(PROJECT)" >version.go
	@echo '' >>version.go
	@echo '// Version of package' >>version.go
	@echo 'const (' >>version.go
	@echo '    Version = `$(VERSION)`' >>version.go
	@echo '' >>version.go
	@echo '    LicenseText = `' >>version.go
	@cat LICENSE >>version.go
	@echo '`' >>version.go
	@echo ')' >>version.go
	@echo '' >>version.go
	@git add version.go

CITATION.cff: codemeta.json
	echo "" | pandoc --metadata title="CITATION.cff" --metadata-file=codemeta.json --from=markdown --template=codemeta-cff.tmpl >CITATION.cff

about.md: codemeta.json $(PROGRAMS)
	echo "" | pandoc --from=markdown --to=markdown --metadata title="About $(PROJECT)" --metadata-file=codemeta.json --template codemeta-md.tmpl >about.md

installer.sh: .FORCE
	echo '' | pandoc --metadata title='Installer' --metadata-file codemeta.json --template codemeta-installer.tmpl >installer.sh
	chmod 775 installer.sh
	git add -f installer.sh

$(PROGRAMS): cmd/*/*.go $(PACKAGE)
	@mkdir -p bin
	go build -o bin/$@$(EXT) cmd/$@/*.go

# NOTE: on macOS you must use "mv" instead of "cp" to avoid problems
install: build man .FORCE
	@if [ ! -d $(PREFIX)/bin ]; then mkdir -p $(PREFIX)/bin; fi
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f ./bin/$$FNAME ]; then mv -v ./bin/$$FNAME $(PREFIX)/bin/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"
	@echo ""
	@if [ ! -d $(PREFIX)/man/man1 ]; then mkdir -p $(PREFIX)/man/man1; fi
	@cp -v man/man1/pttk.1 $(PREFIX)/man/man1/
	@cp -v man/man1/pttk-blogit.1 $(PREFIX)/man/man1/
	@cp -v man/man1/pttk-phlogit.1 $(PREFIX)/man/man1/
	@cp -v man/man1/pttk-rss.1 $(PREFIX)/man/man1/
	@cp -v man/man1/pttk-ws.1 $(PREFIX)/man/man1/
	@cp -v man/man1/pttk-gs.1 $(PREFIX)/man/man1/
	@cp -v man/man1/pttk-include.1 $(PREFIX)/man/man1/
	@cp -v man/man1/pttk-frontmatter.1 $(PREFIX)/man/man1/
	@echo ""
	@echo "Make sure $(PREFIX)/man is in your MANPATH"
	@echo ""

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	-for FNAME in $(PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME ]; then rm -v $(PREFIX)/bin/$$FNAME; fi; done
	-rm $(PREFIX)/man/man1/pttk.1
	-rm $(PREFIX)/man/man1/pttk-blogit.1
	-rm $(PREFIX)/man/man1/pttk-phlogit.1
	-rm $(PREFIX)/man/man1/pttk-rss.1
	-rm $(PREFIX)/man/man1/pttk-ws.1
	-rm $(PREFIX)/man/man1/pttk-gs.1
	-rm $(PREFIX)/man/man1/pttk-include.1
	-rm $(PREFIX)/man/man1/pttk-frontmatter.1

man: man/man1/pttk.1 man/man1/pttk-blogit.1 man/man1/pttk-phlogit.1 man/man1/pttk-rss.1 man/man1/pttk-ws.1 man/man1/pttk-gs.1 man/man1/pttk-include.1 man/man1/pttk-frontmatter.1 .FORCE

man/man1/pttk.1: pttk.1.md
	mkdir -p man/man1
	pandoc pttk.1.md -s -t man -o man/man1/pttk.1

man/man1/pttk-blogit.1: pttk-blogit.1.md
	mkdir -p man/man1
	pandoc pttk-blogit.1.md -s -t man -o man/man1/pttk-blogit.1

man/man1/pttk-phlogit.1: pttk-phlogit.1.md
	mkdir -p man/man1
	pandoc pttk-phlogit.1.md -s -t man -o man/man1/pttk-phlogit.1

man/man1/pttk-rss.1: pttk-rss.1.md
	mkdir -p man/man1
	pandoc pttk-rss.1.md -s -t man -o man/man1/pttk-rss.1

man/man1/pttk-ws.1: pttk-ws.1.md
	mkdir -p man/man1
	pandoc pttk-ws.1.md -s -t man -o man/man1/pttk-ws.1

man/man1/pttk-gs.1: pttk-gs.1.md
	mkdir -p man/man1
	pandoc pttk-gs.1.md -s -t man -o man/man1/pttk-gs.1

man/man1/pttk-include.1: pttk-include.1.md
	mkdir -p man/man1
	pandoc pttk-include.1.md -s -t man -o man/man1/pttk-include.1

man/man1/pttk-frontmatter.1: pttk-frontmatter.1.md
	mkdir -p man/man1
	pandoc pttk-frontmatter.1.md -s -t man -o man/man1/pttk-frontmatter.1


#website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
#	bash mk-website.bash

gophersite: gophersite.zip

gophersite.zip: gophermap *.md
	zip gophersite.zip gophermap README.md ptty*.md CITATION.cff fix-perms.bash


check: .FORCE
	go vet *.go

test: clean build
	cd blogit && go test
	cd phlogit && go test
	cd rss && go test
	cd ws && go test
	cd frontmatter && go test
#	cd gs && go test
#	go test
#	cd help && go test
#	cd include && go test

cleanweb:
	@if [ -f index.html ]; then rm *.html; fi

clean:
	@if [ -d bin ]; then rm -fR bin; fi
	@if [ -d dist ]; then rm -fR dist; fi
	@if [ -d testout ]; then rm -fR testout; fi
	@if [ -f man/man1/pttk.1 ]; then rm man/man1/pttk.1; fi

dist/linux-amd64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env  GOOS=linux GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-linux-amd64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macos-amd64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macos-amd64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macos-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macos-arm64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/windows-amd64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-windows-amd64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/raspbian-arm7:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-rasperry-pi-os-arm7.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

distribute_docs:
	if [ -d dist ]; then rm -fR dist; fi
	mkdir -p dist
	cp -v codemeta.json dist/
	cp -v CITATION.cff dist/
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v installer.sh dist/
	cp -vR man dist/
	#cp -v INSTALL.md dist/

update_version:
	$(EDITOR) codemeta.json
	codemeta2cff codemeta.json CITATION.cff

release: CITATION.cff clean version.go distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macos-amd64 dist/macos-arm64 dist/raspbian-arm7

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

website: about.md
	make -f website.mak

publish: build website
	bash publish.bash

.FORCE:
