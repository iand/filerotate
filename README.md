filerotate
==========

Go package to rotate writes to a file based on a time interval. Canonical use case is logging.

The NewFileRotater function creates a FileRotater object which implements tha standard Write interface

NewFileRotater takes 3 parameters:

* filename is the base name of the file.
* format is a date format string that will be applied to the rotation time and appended to filename as an extension. Use an empty string to append unixtime seconds.
* rotationtime is the interval between rotations, starting at midnight of current day

If you specify 86400 for the rotationtime then the file will rotate 86400 seconds after midnight on the current day and every 86400 seconds thereafter. If the 
rotationtime were 5 minutes then the file will rotate at the nearest multiple of 5 minutes from midnight and every 5 minutes after that.

EXAMPLE
=======

This example rotates a file every 30 seconds while writing to it


	package main

	import (
		"github.com/iand/filerotate"
		"time"
	)

	func main() {
		fr, _ := filerotate.NewFileRotater("/tmp/filerotatetest", "", 30*time.Second)

		for {
			fr.Write([]byte("hello\n"))

			time.Sleep(5 * time.Second)
		}

	}



INSTALLATION
============

Simply run

	go get github.com/iand/filerotate

Documentation is at [http://go.pkgdoc.org/github.com/iand/filerotate](http://go.pkgdoc.org/github.com/iand/filerotate)

LICENSE
=======
This code and associated documentation is in the public domain.

To the extent possible under law, Ian Davis has waived all copyright
and related or neighboring rights to this file. This work is published 
from the United Kingdom. 

TIP
===
If you like this code and want to show your appreciation, I accept bitcoin tips at 1NMjYDmQq9X2m8oSSieGh6J6tmJY11K47X

