
[![Project Status: Concept – Minimal or no implementation has been done yet, or the repository is only intended to be a limited example, demo, or proof-of-concept.](https://www.repostatus.org/badges/latest/concept.svg)](https://www.repostatus.org/#concept)

pdtk
====

A light weight pre-processor for Pandoc. Target use case is JSON object
documents rendered via Pandoc templates.

This is a proof-of-concept Go package which makes it easy to extend
your Go application to incorporate Pandoc template processing. It includes
the Go package and an example command line tool (cli).

A motivation for pdtk is to deconstruct a content management system
into easily scriptable components and drive website assembly via
a simple Makefile. Pandoc is relied for processing input and
rendering content but pdtk makes it easy to integrate with other
activity such as maintain a blog stuctured document tree,
rendering RSS and sitemap.xml.

cli
---

The command line implementation is a proof of concept Pandoc
pre-processor in Go. An example usage would be to process
[example.json](example.json) JSON document using a Pandoc template
called [example.tpml](example.tmpl).

```shell
    pdtk -- --template example.tmpl < example.json > example.html
```

Go package
----------

Here's some simple use examples of the three functions supplied
in the pdtk package.

Given a JSON Object document  as a slice of bytes render formated
output based on the Pandoc template `example.tmpl`

```go
    src, err := ioutil.ReadFile("example.json")
    if err != nil {
        // ... handle error
    }
    // options passed to Pandoc
    opt := []string{"--template", "example.tmpl"}
    src, err = pdtk.Apply(src, opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
```

Using an `io.Reader` to retrieve the JSON content, process with the
`example.tmpl` template and write standard output

```go
    f, err := Open("example.json")
    if err != nil {
        // ... handle error
    }
    defer f.Close()
    // options passed to Pandoc
    opt := []string{"--template", "example.tmpl"}
    src, err := pdtk.ReadAll(f, opt)
    if err != nil {
        // ... handle error
    }
    fmt.Fprintf(os.Stdout, "%s", src)
```

Using an `io.Reader` and `io.Writer` to read JSON source from standard
input and write the processed pandoc templated standard output.

```go
    // options passed to Pandoc
    opt := []string{"--template", "example.tmpl"}
    err := pdtk.ApplyIO(os.Stdin, os.Stdout, opt)
    if err != nil {
        // ... handle error
    }
```

Requirements
------------

- [Pandoc](https://pandoc.org) 2.18 or better
- [Go](https://golang.org) 1.18.4 or better to compile from source
- [GNU Make](https://www.gnu.org/software/make/) (optional) to automated compilation
- [Git](https://git-scm.com/) or other Git client to retrieve this repository

Installation
------------

1. Clone https://github.com/rsdoiel/pdtk to your local machine
2. Change directory into the git repository (i.e. `pdtk`
3. Compile using `go build`
4. Install using `go install`

```shell
    git clone https://github.com/rsdoiel/pdtk
    cd pdtk
    git fetch origin
    git pull origin main
    go build -o bin/pdtk cmd/pdtk/pdtk.go
    go install
```

NOTE: This recipe assumes' you are familar with setting up a
Go development environment (e.g. You've set GOPATH environment
appropriately). See the [go website](https://golang.org) for
details about setting up and compiler programs.

License
-------

BSD 3-Clause License

Copyright (c) 2022, R. S. Doiel
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

