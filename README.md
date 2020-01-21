go-file-copies
==============

A Go program to get duplicates from specified paths.

Table of Contents
=================

+ [Images](#Images)
+ [Flags](#Flags)
+ [Install](#Install)
+ [Run](#Run)

### Images

![](https://i.imgur.com/mY3fIni.png)

### Flags

 * config - _specify the path to your config file, which has paths to directories with duplicates._
 * output - _specify the path to the output file with results._

### Install

##### Compile for yourself
Install [Go](https://golang.org/) and run [compile.sh](compile.sh) from the terminal.
Binaries will be placed in the "build" directory.

##### Use precompiled binaries
Download binary for your system from [releases](https://github.com/Dmitriy-Vas/go-file-copies/releases)

### Run

Rename [config-sample.json](config-sample.json) to config.json and add paths with duplicates to "dirs".<br>
The program will take __all files recursively__ from specified directories.<br>
You can specify paths to config and output via [flags](#Flags).<br>
Go to the directory with the program and run it like usual binary.
