
# Installation

The plain text toolkit is experiemental. It's probably break and surely has bugs.

## Quick install via curl

This project can be installed via curl and sh for POSIX compatible systems.
Enter the following in your shell.

~~~
curl https://rsdoiel.github.io/pttk/installer.sh | sh
~~~

## Install from source

### Requirements

- Golang >= 1.21.4
- Pandoc >= 3.1
- GNU Make
- Git

### Steps

1. clone this repository
2. change into the cloned directory
3. Run `make`, `make test` and `make install`

My default the program(s) are installed in `$HOME/bin`. You will
need to make sure that `$HOME/bin` is in your path.

Here's what I run on a bare machine to install from source.

```
git clone https://github.com/rsdoiel/pttk src/github.com/rsdoiel/pttk
cd src/github.com/rsdoiel/pttk
make
make test
make install
```



