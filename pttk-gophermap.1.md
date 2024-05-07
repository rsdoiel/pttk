% pttk-gophermap(1) pttk-gophermap user manual
% R. S. Doiel
% May 7, 2024

# NAME 

pttk

# SYNOPSIS

pttk gophermap [OPTIONS] [GOPHERMAP_NAME] [FILES_TO_LIST]

# DESCRIPTIOJ

pttk gophermap provides support for generating Gophermaps, the "index" page
for directories in your Gopher Hole.

# OPTIONS

What follows are the options supported by the phlogit verb.

-help
: display gophermap help

-masthead FILENAME
: Use thie specified file contents as the "masthead" of the Gophermap

-verbose
: verbose output

# EXAMPLE

~~~shell
	pttk gophermap -masthead banner.txt gophermap \
	   README COPYING LICENSE \
	   ThePlan.txt \
	   TheBiggerPlan.txt \
	   TheWholeHole.txt
~~~


