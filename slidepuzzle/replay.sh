#!/bin/bash

cp answer.dat.tmp answer.dat

~/go/bin/6g main.go
~/go/bin/6l main.6
/usr/bin/php play.php $1 $2 $3








