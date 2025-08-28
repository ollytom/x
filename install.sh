#!/bin/sh

mkdir -p $HOME/bin
cp bin/* $HOME/bin

mkdir -p $HOME/.config/git
cp lib/git $HOME/.config/git/config

mkdir -p $HOME/lib
cp lib/plumbing $HOME/lib

if uname | grep OpenBSD
then
	mkdir -p $HOME/.config/gtk-3.0 $HOME/.config/gtk-4.0
	cp lib/gtk.ini $HOME/.config/gtk-3.0/settings.ini
fi
