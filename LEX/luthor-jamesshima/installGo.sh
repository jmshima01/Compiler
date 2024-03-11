#!/usr/bin/sh

cd
wget "https://go.dev/dl/go1.22.1.linux-amd64.tar.gz"
tar -xvf go1.22.1.linux-amd64.tar.gz -C ~/bin
export PATH="$PATH:${HOME}/bin/go/bin"
rm -rf go1.22.1.linux-amd64.tar.gz
echo "You should now have go1.22 locally on isengard!"
echo "add the export to your .bashrc if you'd like"
go version

