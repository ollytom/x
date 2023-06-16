#!/bin/sh

for f in man/*
do
	name=`basename $f | sed 's/\.[0-9]$//'`
	description=`grep .Nd $f | sed 's/^\.Nd //'`
	echo $name '-' $description
done
