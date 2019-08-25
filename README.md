# GoScripts

[![Build Status](https://travis-ci.org/mcai4gl2/goscripts.svg?branch=master)](https://travis-ci.org/mcai4gl2/goscripts)

### compress

Program to compress files in one folder one by one to another folder. Can use --parallel flag to run on parallel processing. In parallel mode, there is one go routine per file to be compressed. If there are too many files to be compressed, this can be slow and memory intensive.

### compressp

Similar to compress, but this version uses fixed number of go routines instead of one per file. Using --parallel to control how many go routines to use. By default, the program will create 10 go routines.

