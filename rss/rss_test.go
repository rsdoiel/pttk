// pttk is software for working with plain text content.
// Copyright (C) 2022 R. S. Doiel
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
package rss

import (
	"bytes"
	"encoding/xml"
	"net/url"
	"strings"
	"testing"
)

func TestRSS2(t *testing.T) {
	src := []byte(`<?xml version="1.0"?>
<rss version="2.0">
    <channel>
        <title>Robert&#39;s ramblings</title>
        <link>https://rsdoiel.github.io/blog</link>
        <description>Robert&#39;s ramblings and wonderigs</description>
        <pubDate>Sat, 30 Jul 2022 00:00:00 GMT</pubDate>
        <lastBuildDate>Sat, 30 Jul 2022 00:00:00 GMT</lastBuildDate>
        <generator>0.0.0</generator>
        <item>
            <title>Turbo Oberon, the dream</title>
            <link>/blog/2022/07/30/Turbo-Oberon.md</link>
            <description>Sometimes I have odd dreams and that was true last night through early this morning. The dream was set in the future. I was already retired. It was a dream about &#34;Turbo Oberon&#34;.&#xA;&#xA;&#34;Turbo Oberon&#34; was an Oberon language. The language compiler was named &#34;TO&#34; in my dream. A module&#39;s file extension was &#34;.tom&#34;, in honor of Tom Lopez (Meatball Fulton) of ZBS. There were allot of ZBS references in the dream.&#xA;&#xA;&#34;TO&#34; was very much a language in the Oberon-07 tradition with minor extensions when it came to bringing in modules. It allowed for a multi path search for module names. You could also express a Module import as a string allowing providing paths to the imported module.&#xA;&#xA;Compilation was similar to Go. Cross compilation was available out of the box by setting a few environment variables. I remember answering questions about the language and its evolution. I remember mentioning in the conversation about how I thought Go felling into the trap of complexity like Rust or C/C++ before it. The turning point for Go was generics. Complexity was the siren song to be resisted in &#34;Turbo Oberon&#34;. Complexity is seductive to language designers and implementers. I was only an implementer.&#xA;&#xA;Evolution wise &#34;TO&#34; was built initially on the Go tool chain. As a result it featured easily cross-compiled binaries and had a rich standard set of Modules like Go but also included portable libraries for implementing graphic user interfaces. &#34;Turbo Oberon&#34; evolved as a conversation between Go and the clean simplicity of Oberon-07. Two example applications &#34;shipped&#34; with the &#34;TO&#34; compiler. They were an Oberon like Operating System (stand alone and hosted) and a Turbo Pascal like IDE. The IDE was called &#34;toe&#34; for Turbo Oberon Editor. I don&#39;t remember the name of the OS implementation but it might have been &#34;toos&#34;. I remember &#34;TO&#34; caused problems for search engines and catalog systems. ...</description>
            <pubDate>Sat, 30 Jul 2022 00:00:00 GMT</pubDate>
            <guid>/blog/2022/07/30/Turbo-Oberon.md</guid>
        </item>
        <item>
            <title>Artemis Project Status, 2022</title>
            <link>/blog/2022/07/27/Artemis-Status-Summer-2022.md</link>
            <description>It&#39;s been a while since I wrote an Oberon-07 post and even longer since I&#39;ve worked on Artemis. Am I done with Oberon-07 and abandoning Artemis?  No. Life happens and free time to just hasn&#39;t been available. I don&#39;t know when that will change.&#xA;&#xA;Since I plan to continue working Artemis I need to find a way forward in much less available time. Time to understand some of my constraints.&#xA;&#xA;1. I work on a variety of machines, OBNC is the only compiler I&#39;ve consistently been able to use across all my machines&#xA;2. Porting between compilers takes energy and time, and those compilers don&#39;t work across all my machines&#xA;3. When I write Oberon-07 code I quickly hit a wall for the things I want to do, this is what original inspired Artemis, so there is still a need for a collection of modules&#xA;4. Oberon/Oberon-07 on Wirth RISC virtual machine is not sufficient for my development needs&#xA;5. A2, while very impressive, isn&#39;t working for me either (mostly because I need to work on ARM CPUs)&#xA;&#xA;These constraints imply Artemis is currently too broadly scoped. I think I need to focus on what works in OBNC for now. Once I have a clear set of modules then I can revisit portability to other compilers.&#xA;&#xA;What modules do I think I need? If I look at my person projects I tend to work allot with text, often structured text (e.g. XML, JSON, CSV). I also tend to be working with network services. Occasionally I need to interact with database (e.g. SQLite3, MySQL, Postgres).  Artemis should provide modules to make it easy to write code in Oberon-07 that works in those areas. Some of that I can do by wrapping existing C libraries. Some I can simply write from scratch in Oberon-07 (e.g. a JSON encoder/decoder). That&#39;s going to me my focus as my hobby time becomes available and then. ...</description>
            <pubDate>Wed, 27 Jul 2022 00:00:00 GMT</pubDate>
            <guid>/blog/2022/07/27/Artemis-Status-Summer-2022.md</guid>
        </item>
        <item>
            <title>Installing Golang from source on RPi-OS for arm64</title>
            <link>/blog/2022/02/18/Installing-Go-from-Source-RPiOS-arm64.md</link>
            <description>&#xA;&#xA;This are my quick notes on installing Golang from source on the Raspberry Pi OS 64 bit.&#xA;&#xA;1. Get a working compiler&#xA;&#x9;a. go to https://go.dev/dl/ and download go1.17.7.linux-arm64.tar.gz&#xA;&#x9;b. untar the tarball in your home directory (it&#39;ll unpack to $HOME/go)&#xA;&#x9;c. ` + "`" + `cd go/src` + "`" + ` and ` + "`" + `make.bash` + "`" + `&#xA;2. Move go directory to go1.17&#xA;3. Clone go from GitHub&#xA;4. Compile with the downloaded compiler&#xA;&#x9;a. ` + "`" + `cd go/src` + "`" + `&#xA;&#x9;b. ` + "`" + `env GOROOT_BOOTSTRAP=$HOME/go1.17 ./make.bash` + "`" + `&#xA;&#x9;c. Make sure ` + "`" + `$HOME/go/bin` + "`" + ` is in the path&#xA;&#x9;d. ` + "`" + `go version` + "`" + `&#xA;&#xA; ...</description>
            <pubDate>Fri, 18 Feb 2022 00:00:00 GMT</pubDate>
            <guid>/blog/2022/02/18/Installing-Go-from-Source-RPiOS-arm64.md</guid>
        </item>
        <item>
            <title>Notes on setting up a Mid-2010 Mac Mini</title>
            <link>/blog/2021/12/18/Notes-on-setting-up-a-2010-Mac-Mini.md</link>
            <description>I acquired a Mid 2010 Mac Mini. It was in good condition but lacked an OS on the hard drive.  I used a previously purchased copy of Mac OS X Snow Leopard to get an OS up and running on the bare hardware. Then it was a longer effort to get the machine into a state with the software I wanted to use on it. My goal was Mac OS X High Sierra, Xcode 10.1 and Mac Ports. The process was straight forward but very time consuming but I think worth it.  I wound up with a nice machine for experimenting with and writing blog posts.&#xA;&#xA;The setup process was as follows:&#xA;&#xA;1. Install macOS Snow Leopard on the bare disk of the Mac Mini&#xA;2. Install macOS El Capitan on the Mac Mini after manually downloading it from Apple&#39;s support site&#xA;3. Run updates indicated by El Capitan&#xA;4. Install macOS High Sierra on the Mac Mini after manually downloading it from the Apple&#39;s support site&#xA;5. Run updates indicated by High Sierra &#xA;6. Manually download and install Xcode 10.1 command line tools &#xA;7. Check and run some updates again&#xA;8. Finally install Mac Ports&#xA;&#xA;The OS installs took about 45 minutes to 90 minutes each. Installing Xcode took about 45 minutes to an hour. Installing Mac Ports was fast as was installing software via Mac Ports.&#xA;&#xA;- Apple support pages that I found helpful&#xA;    - [How to get old versions of macOS](https://support.apple.com/en-us/HT211683)&#xA;    - [How to create a bootable installer for macOS](https://support.apple.com/en-us/HT201372)&#xA;    - [macOS High Sierra - Technical Specifications](https://support.apple.com/kb/SP765?locale=en_US)&#xA;- Wikipedia page on [Xcode](https://en.wikipedia.org/wiki/Xcode) is how I sorta out what version of Xcode I needed to install&#xA;- Links to old macOS and Xcode&#xA;    - Download [Mac OS X El El Capitan](http://updates-http.cdn-apple.com/2019/cert/061-41424-20191024-218af9ec-cf50-4516-9011-228c78eda3d2/InstallMacOSX.dmg)&#xA;    - Download [Mac OX X High Sierra](https://apps.apple.com/us/app/macos-high-sierra/id1246284741?mt=12)&#xA;    - Download [Xcode 10.1](https://developer.apple.com/download/all/?q=xcode), Scroll down the list until you want it.&#xA;        - [Command Line Tools (macOS 10.13) for Xcode 10.1](https://download.developer.apple.com/Developer_Tools/Command_Line_Tools_macOS_10.13_for_Xcode_10.1/Command_Line_Tools_macOS_10.13_for_Xcode_10.1.dmg)&#xA;        - NOTE: There are two version available, you want the version for macOS 10.13 (High Sierra) NOT Mac OS 10.14. ...</description>
            <pubDate>Sat, 18 Dec 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/12/18/Notes-on-setting-up-a-2010-Mac-Mini.md</guid>
        </item>
        <item>
            <title>Setting up FreeDOS 1.3rc4 with Qemu</title>
            <link>/blog/2021/11/27/FreeDOS-1.3rc4-with-Qemu.md</link>
            <description>In this article I&#39;m going explore setting up FreeDOS with Qemu&#xA;on my venerable Dell 4319 running Raspberry Pi Desktop OS (Debian&#xA;GNU/Linux).  First step is to download FreeDOS &#34;Live CD&#34; in the&#xA;1.3 RC4 release. See http://freedos.org/download/ for that.&#xA;&#xA;I needed to install Qemu in my laptop. It runs the Raspberry Pi&#xA;Desktop OS (i.e. Debian with Raspberry Pi UI). I choose to install&#xA;the &#34;qemu-system&#34; package since I will likely use qemu for other&#xA;things besides FreeDOS. The qemu-system package contains all the&#xA;various systems I might want to emulate in other projects as well&#xA;as several qemu utilities that are handy.  Here&#39;s the full sequence&#xA;of ` + "`" + `apt` + "`" + ` commands I ran (NOTE: these included making sure my laptop&#xA;was up to date before I installed qemu-system).&#xA;&#xA;~~~&#xA;sudo apt update&#xA;sudo apt upgrade&#xA;sudo apt install qemu-system&#xA;~~~&#xA;&#xA;Now that I had the software available it was time to figure out&#xA;how to actually knit things together and run FreeDOS.&#xA;&#xA;Before I get started I create a folder in my home directory&#xA;for running everything. You can name it what you want&#xA;but I called mine ` + "`" + `FreeDOS_13` + "`" + ` and changed into that folder&#xA;for the work in this article. ...</description>
            <pubDate>Sat, 27 Nov 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/11/27/FreeDOS-1.3rc4-with-Qemu.md</guid>
        </item>
        <item>
            <title>Portable Conversions (Integers)</title>
            <link>/blog/2021/11/26/Portable-Conversions-Integers.md</link>
            <description>An area in working with Oberon-07 on a POSIX machine that has proven problematic is type conversion. In particular converting to and from INTEGER or REAL and ASCII.  None of the three compilers I am exploring provide a common way of handling this. I&#39;ve explored relying on C libraries but that approach has it&#39;s own set of problems.  I&#39;ve become convinced a better approach is a pure Oberon-07 library that handles type conversion with a minimum of assumptions about the implementation details of the Oberon compiler or hardware. I&#39;m calling my conversion module &#34;Types&#34;. The name is short and descriptive and seems an appropriate name for a module consisting of type conversion tests and transformations.  My initial implementation will focusing on converting integers to and from ASCII.&#xA;&#xA;I don&#39;t want to rely on the representation of the INTEGER value in the compiler or at the machine level. That has lead me to think in terms of an INTEGER as a signed whole number.&#xA;&#xA;The simplest case of converting to/from ASCII is the digits from zero to nine (inclusive). Going from an INTEGER to an ASCII CHAR is just looking up the offset of the character representing the &#34;digit&#34;. Like wise going from ASCII CHAR to a INTEGER is a matter of mapping in the reverse direction.  Let&#39;s call these procedures ` + "`" + `DigitToChar` + "`" + ` and  ` + "`" + `CharToDigit*` + "`" + `.&#xA;&#xA;Since INTEGER can be larger than zero through nine and CHAR can hold non-digits I&#39;m going to add two additional procedures for validating inputs -- ` + "`" + `IsIntDigit` + "`" + ` and ` + "`" + `IsCharDigit` + "`" + `. Both return TRUE if valid, FALSE if not.&#xA;&#xA;For numbers larger than one digit I can use decimal right shift to extract the ones column value or a left shift to reverse the process.  Let&#39;s called these ` + "`" + `IntShiftRight` + "`" + ` and ` + "`" + `IntShiftLeft` + "`" + `.  For shift right it&#39;d be good to capture the ones column being lost. For shift left it would be good to be able to shift in a desired digit. That way you could shift/unshift to retrieve to extract and put back values. ...</description>
            <pubDate>Fri, 26 Nov 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/11/26/Portable-Conversions-Integers.md</guid>
        </item>
        <item>
            <title>Revisiting Files</title>
            <link>/blog/2021/11/22/Revisiting-Files.md</link>
            <description>In October I had an Email exchange with Algojack regarding a buggy example in [Oberon-07 and the file system](../../../2020/05/09/Oberon-07-and-the-filesystem.html). The serious bug was extraneous non-printable characters appearing a plain text file containing the string &#34;Hello World&#34;. The trouble with the example was a result of my misreading the Oakwood guidelines and how **Files.WriteString()** is required to work. The **Files.WriteString()** procedure is supposed to write every element of a string to a file. This __includes the trailing Null character__. The problem for me is **Files.WriteString()** litters plain text files with tailing nulls. What I should have done was write my own **WriteString()** and **WriteLn()**. The program [HelloworldFile](./HelloworldFile.Mod) below is a more appropriate solution to writing strings and line endings than relying directly on **Files**. In a future post I will explorer making this more generalized in a revised &#34;Fmt&#34; module.&#xA;&#xA;~~~&#xA;MODULE HelloworldFile;&#xA;&#xA;IMPORT Files, Strings;&#xA;&#xA;CONST OberonEOL = 1; UnixEOL = 2; WindowsEOL = 3;&#xA;&#xA;VAR&#xA;  (* holds the eol marker type to use in WriteLn() *)&#xA;  eolType : INTEGER;&#xA;  (* Define a file handle *)&#xA;    f : Files.File;&#xA;  (* Define a file rider *)&#xA;    r : Files.Rider; ...</description>
            <pubDate>Mon, 22 Nov 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/11/22/Revisiting-Files.md</guid>
        </item>
        <item>
            <title>Combining Oberon-07 with C using Obc-3</title>
            <link>/blog/2021/06/14/Combining-Oberon-07-with-C-using-Obc-3.md</link>
            <description>&#xA;&#xA;This post explores integrating C code with an Oberon-07 module use&#xA;Mike Spivey&#39;s Obc-3 Oberon Compiler.  Last year I wrote a similar post&#xA;for Karl Landström&#39;s [OBNC](/blog/2020/05/01/Combining-Oberon-and-C.html).&#xA;This goal of this post is to document how I created a version of Karl&#39;s&#xA;Extension Library that would work with Mike&#39;s Obc-3 compiler.&#xA;If you want to take a shortcut you can see the results on GitHub&#xA;in my [obc-3-libext](https://github.com/rsdoiel/obc-3-libext) repository.&#xA;&#xA;From my time with OBNC I&#39;ve come to rely on three modules from Karl&#39;s&#xA;extension library. When trying to port some of my code to use with&#xA;Mike&#39;s compiler. That&#39;s where I ran into a problem with that dependency.&#xA;Karl&#39;s modules aren&#39;t available. I needed an [extArgs](http://miasap.se/obnc/obncdoc/ext/extArgs.def.html),&#xA;an [extEnv](http://miasap.se/obnc/obncdoc/ext/extEnv.def.html) and&#xA;[extConvert](http://miasap.se/obnc/obncdoc/ext/extConvert.def.html).&#xA;&#xA;Mike&#39;s own modules that ship with Obc-3 cover allot of common ground&#xA;with Karl&#39;s. They are organized differently. The trivial solution is&#xA;to implement wrapping modules using Mike&#39;s modules for implementation.&#xA;That takes case of extArgs and extEnv.&#xA;&#xA;The module extConvert is in a another category. Mike&#39;s ` + "`" + `Conv` + "`" + ` module is&#xA;significantly minimalist. To solve that case I&#39;ve create C code to perform&#xA;the needed tasks based on Karl&#39;s examples and used Mike&#39;s share library&#xA;compilation instructions to make it available inside his run time. ...</description>
            <pubDate>Mon, 14 Jun 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/06/14/Combining-Oberon-07-with-C-using-Obc-3.md</guid>
        </item>
        <item>
            <title>Combining Oberon-07 with C using Obc-3</title>
            <link>/blog/2021/06/14/Combining-Oberon-07-with-C-using-Obc-3.md</link>
            <description>&#xA;&#xA;This post explores integrating C code with an Oberon-07 module use&#xA;Mike Spivey&#39;s Obc-3 Oberon Compiler.  Last year I wrote a similar post&#xA;for Karl Landström&#39;s [OBNC](/blog/2020/05/01/Combining-Oberon-and-C.html).&#xA;This goal of this post is to document how I created a version of Karl&#39;s&#xA;Extension Library that would work with Mike&#39;s Obc-3 compiler.&#xA;If you want to take a shortcut you can see the results on GitHub&#xA;in my [obc-3-libext](https://github.com/rsdoiel/obc-3-libext) repository.&#xA;&#xA;From my time with OBNC I&#39;ve come to rely on three modules from Karl&#39;s&#xA;extension library. When trying to port some of my code to use with&#xA;Mike&#39;s compiler. That&#39;s where I ran into a problem with that dependency.&#xA;Karl&#39;s modules aren&#39;t available. I needed an [extArgs](http://miasap.se/obnc/obncdoc/ext/extArgs.def.html),&#xA;an [extEnv](http://miasap.se/obnc/obncdoc/ext/extEnv.def.html) and&#xA;[extConvert](http://miasap.se/obnc/obncdoc/ext/extConvert.def.html).&#xA;&#xA;Mike&#39;s own modules that ship with Obc-3 cover allot of common ground&#xA;with Karl&#39;s. They are organized differently. The trivial solution is&#xA;to implement wrapping modules using Mike&#39;s modules for implementation.&#xA;That takes case of extArgs and extEnv.&#xA;&#xA;The module extConvert is in a another category. Mike&#39;s ` + "`" + `Conv` + "`" + ` module is&#xA;significantly minimalist. To solve that case I&#39;ve create C code to perform&#xA;the needed tasks based on Karl&#39;s examples and used Mike&#39;s share library&#xA;compilation instructions to make it available inside his run time. ...</description>
            <pubDate>Mon, 14 Jun 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/06/14/Combining-Oberon-07-with-C-using-Obc-3.md</guid>
        </item>
        <item>
            <title>Beyond Oakwood, Modules and Aliases</title>
            <link>/blog/2021/05/16/Beyond-Oakwood-Modules-and-Aliases.md</link>
            <description>Oakwood is the name used to refer to an early Oberon language&#xA;standardization effort in the late 20th century.  It&#39;s the name&#xA;of a hotel where compiler developers and the creators of Oberon&#xA;and the Oberon System met to discuss compatibility. The lasting&#xA;influence on the 21st century Oberon-07 language can be seen&#xA;in the standard set of modules shipped with POSIX based Oberon-07&#xA;compilers like&#xA;[OBNC](https://miasap.se/obnc/), [Vishap Oberon Compiler](https://github.com/vishaps/voc) and the &#xA;[Oxford Oberon Compiler](http://spivey.oriel.ox.ac.uk/corner/Oxford_Oberon-2_compiler).&#xA;&#xA;The Oakwood guidelines described a minimum expectation for&#xA;a standard set of modules to be shipped with compilers.&#xA;The modules themselves are minimalist in implementation.&#xA;Minimalism can assist in easing the learning curve&#xA;and encouraging a deeper understanding of how things work.&#xA;&#xA;The Oberon-07 language is smaller than the original Oberon language&#xA;and the many dialects that followed.  I think of Oberon-07 as the&#xA;distillation of all previous innovation.  It embodies the&#xA;spirit of &#34;Simple but not simpler than necessary&#34;. Minimalism is&#xA;a fit description of the adaptions of the Oakwood modules for &#xA;Oberon-07 in the POSIX environment.&#xA;&#xA;Sometimes I want more than the minimalist module.  A good example&#xA;is standard [Strings](https://miasap.se/obnc/obncdoc/basic/Strings.def.html)&#xA;module.  Thankfully you can augment the standard modules with your own.&#xA;If you are creative you can even create a drop in replacement.&#xA;This is what I wound up doing with my &#34;Chars&#34; module.&#xA;&#xA;In the spirit of &#34;Simple but no simpler&#34; I originally kept Chars &#xA;very minimal. I only implemented what I missed most from Strings.&#xA;I got down to a handful of functions for testing characters,&#xA;testing prefixes and suffixes as well as trim procedures. It was&#xA;all I included in ` + "`" + `Chars` + "`" + ` was until recently. ...</description>
            <pubDate>Sun, 16 May 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/05/16/Beyond-Oakwood-Modules-and-Aliases.md</guid>
        </item>
        <item>
            <title>Ofront on Raspberry Pi OS</title>
            <link>/blog/2021/04/25/Ofront-on-Rasberry-Pi-OS.md</link>
            <description>This post is about getting Ofront[^1] up and running on Raspberry Pi OS[^2].&#xA;Ofront provides a Oberon-2 to C transpiler as well as a Oberon V4[^3]&#xA;development environment. There are additional clever tools like ` + "`" + `ocat` + "`" + `&#xA;that are helpful working with the differences in text file formats between&#xA;Oberon System 3, V4 and POSIX. The V4 implementation sits nicely on top of&#xA;POSIX with minimal compromises that distract from the Oberon experience.&#xA;&#xA;[^1]: Ofront was developed by Joseph Templ, see http://www.software-templ.com/&#xA;&#xA;[^2]: see https://www.raspberrypi.org/software/ (a 32 bit Debian based Linux for both i386 and ARM)&#xA;&#xA;[^3]: see https://ssw.jku.at/Research/Projects/Oberon.html&#xA;&#xA;I first heard of running Ofront/V4 via the ETH Oberon Mail list[^4].&#xA;What caught my eye is the reference to running on Raspberry Pi. Prof. Templ &#xA;provides two flavors of Ofront. One targets the Raspberry Pi OS on ARM&#xA;hardware the second Linux on i386. The Raspberry Pi OS for Intel is an&#xA;i386 variant. I downloaded the tar file, unpacked it and immediately ran&#xA;the &#34;oberon.bash&#34; script provided eager to see a V4 environment. It&#xA;renders but the fonts rendered terribly slowly. I should have read the&#xA;documentation first!  Prof. Templ provides man pages for the tools that&#xA;come with Ofront including the oberon application. Reading the&#xA;man page for oberon quickly addresses the point of slow font rendering.&#xA;It also discusses how to convert Oberon fonts to X Windows bitmap fonts.&#xA;If you use the X Window fonts the V4 environment is very snappy. It does&#xA;require that X Windows knows where to find the fonts used in V4. That is&#xA;done by appending the V4 converted fonts to the X Window font map. I had&#xA;installed the Ofront system in my home directory so the command was ...</description>
            <pubDate>Sun, 25 Apr 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/04/25/Ofront-on-Rasberry-Pi-OS.md</guid>
        </item>
        <item>
            <title>Updating Schema in SQLite3</title>
            <link>/blog/2021/04/16/Updating-Schema-in-SQLite3.md</link>
            <description>[SQLite3](https://sqlite.org/docs.html) is a handy little&#xA;database as single file tool.  You can interact with the file&#xA;through largely standard SQL commands and embed it easily into&#xA;applications via the C libraries that the project supports.&#xA;It is also available from various popular scripting languages&#xA;like Python, PHP, and Lua. One of the things I occasionally&#xA;need to do and always seems to forget it how to is modify a&#xA;table schema where I need to remove a column[^1]. So here are&#xA;some of the basics do I can quickly find them later and avoid&#xA;reading various articles tutorials because the search engines&#xA;doesn&#39;t return the page in the SQLite documentation.&#xA;&#xA;[^1]: The SQL ` + "`" + `ALTER TABLE table_name DROP COLUMN column_name` + "`" + ` does not work in SQLite3&#xA;&#xA;In the next sections I&#39;ll be modeling a simple person object&#xA;with a id, uname, display_name, role and updated fields.&#xA;&#xA;` + "`" + `` + "`" + `` + "`" + `sql&#xA;&#xA;CREATE TABLE IF NOT EXISTS &#34;person&#34; &#xA;        (&#34;id&#34; INTEGER NOT NULL PRIMARY KEY, &#xA;        &#34;uname&#34; VARCHAR(255) NOT NULL, &#xA;        &#34;role&#34; VARCHAR(255) NOT NULL, &#xA;        &#34;display_name&#34; VARCHAR(255) NOT NULL, &#xA;        &#34;updated&#34; INTEGER NOT NULL); ...</description>
            <pubDate>Fri, 16 Apr 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/04/16/Updating-Schema-in-SQLite3.md</guid>
        </item>
        <item>
            <title>A2 Oberon on VirtualBox 6.1</title>
            <link>/blog/2021/04/02/A2-Oberon-on-VirtualBox-6.1.md</link>
            <description>&#xA;&#xA;This is a short article documenting how I install A2 Oberon&#xA;in VirtualBox using the [FreeDOS 1.2](https://freedos.org),&#xA;the A2 [ISO](https://sourceforge.net/projects/a2oberon/files/) cd image and [VirtualBox 6.1](https://virtualbox.org).&#xA;&#xA;1. Download the ISO images for FreeDOS and A2&#xA;2. Create a new Virtual Machine&#xA;3. Install FreeDOS 1.2 (Base install) in the virtual machine&#xA;4. Install A2 from the ISO image over the FreeDOS installation&#xA;&#xA;From working with Native Oberon 2.3.7 I&#39;ve found it very helpful&#xA;to have a FreeDOS 1.2. installed in the Virtual machine first. &#xA;I suspect the reason I have had better luck taking this approach&#xA;is based on assumptions about the virtual hard disk being setup&#xA;with an existing known formatted, boot-able partition. In essence&#xA;making our Virtualbox look like a fresh out of the box vintage PC.&#xA;&#xA;You&#39;ll find FreeDOS 1.2 installation ISO image at &#xA;[FreeDos.org](http://freedos.org/download/). Download it&#xA;where you can easily find it from the VirtualBox manager. ...</description>
            <pubDate>Fri, 02 Apr 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/04/02/A2-Oberon-on-VirtualBox-6.1.md</guid>
        </item>
        <item>
            <title>ETH Oberon System 3 on VirtualBox 6.1</title>
            <link>/blog/2021/03/17/NativeOberon-VirtualBox.md</link>
            <description>In this post I am walking through installing Native Oberon 2.3.7&#xA;(aka ETH Oberon System 3) on a virtual machine running under&#xA;VirtualBox 6.1. It is a follow up to my 2019 post &#xA;[FreeDOS to Oberon System 3](/blog/2019/07/28/freedos-to-oberon-system-3.html &#34;Link to old blog post for bringing up Oberon System 3 in VirtualBox 6.0 using FreeDOS 1.2&#34;). To facilitate the install I will first prepare&#xA;my virtual machine as a FreeDOS 1.2 box. This simplifies getting the&#xA;virtual machines&#39; hard disk partitioned and formatted correctly.&#xA;When Native Oberon was released back in 1990&#39;s most Intel flavored&#xA;machines shipped with some sort Microsoft OS on them.  I believe&#xA;that is why the tools and instructions for Native Oberon assume&#xA;you&#39;re installing over or along side a DOS partition.&#xA;&#xA;1. Install VirtualBox 6.1 installed on your host computer.&#xA;2. Download and install a minimal FreeDOS 1.2 as a virtual machine&#xA;3. Downloaded a copy of Native Oberon 2.3.7 alpha from SourceForge&#xA;3. Familiarized yourself Oberon&#39;s Text User Interface&#xA;4. Boot your FreeDOS virtual machine using the Oberon0.Dsk downloaded&#xA;as part of NativeOberon_2.3.7.tar.gz&#xA;5. Mount &#34;Oberon0.Dsk&#34; and start installing Native Oberon&#xA;&#xA;Before you boot &#34;Oberon0.Dsk&#34; on your virtual machine make sure&#xA;you&#39;ve looked at some online Oberon documentation. This is important.&#xA;Oberon is very different from macOS, Windows, Linux, DOS, CP/M or&#xA;Unix. It is easy to read the instructions and miss important details &#xA;like how you use the three button mouse, particularly the selections&#xA;and execute actions of text instructions.&#xA;&#xA;VirtualBox 6.1 can be obtained from [virtualbox.org](https://www.virtualbox.org/).  This involves downloading the installer for your particular host&#xA;operating system (e.g. Linux, macOS or Windows) and follow the instructions&#xA;on the VirtualBox website to complete the installation.&#xA;&#xA;Once VirtualBox is installed, launch VirtualBox. ...</description>
            <pubDate>Wed, 17 Mar 2021 00:00:00 GMT</pubDate>
            <guid>/blog/2021/03/17/NativeOberon-VirtualBox.md</guid>
        </item>
    </channel>
</rss>`)

	r, err := Parse(src)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	results, err := r.Filter([]string{".item[].title"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(results[".item[].title"].([]string)) != len(r.ItemList) {
		t.Errorf("Expected 6 .item[].title, got %s", strings.Join(results[".item[].title"].([]string), "\t"))
		t.FailNow()
	}
	results, err = r.Filter([]string{".item[].link"})
	if err != nil {
		t.Errorf("Expected 6 .item[].link, got %+v", strings.Join(results[".item[].title"].([]string), "\t"))
		t.FailNow()
	}
	for _, link := range results[".item[].link"].([]string) {
		_, err := url.Parse(link)
		if err != nil {
			t.Errorf("expected to parse link %q into url, %s", link, err)
		}
	}
}

func TestNoTitleFeedItem(t *testing.T) {
	feed := new(RSS2)
	feed.Title = `Test titless items`
	item := new(Item)
	item.Description = "There is no title"
	feed.ItemList = append(feed.ItemList, *item)
	src, err := xml.Marshal(feed)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if bytes.Contains(src, []byte(`<title></title>`)) {
		t.Errorf("expected not title element, got\n%s\n", src)
	}
}
