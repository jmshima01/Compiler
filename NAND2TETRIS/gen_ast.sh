make
./Compiler jack.cfg 10/Square/Main.jack
cat parsetree.txt | ./tree-to-graphvis | dot -Tpng -o parsetree.png