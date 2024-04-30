#!/bin/env bash

# prints a AST to view! takes .jack file as arg...

make
./Compiler $1
cat parsetree.txt | ./tree-to-graphvis | dot -Tpng -o parsetree.png
see parsetree.png