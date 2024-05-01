## 1) 
#### James Shima, Keith Hellman

## 2) 
Weeks, Didn't turn in Project 10 as used a different parsing algorithm and token layout but def worth it and the late penalty

## 3) 
# LL(1) Predict Set Table Parser/Compiler 
I'm currently in the Compiler Course `CSCI425` and learned how to make a frontend compiler toolchain from scratch so I used my knowladge from there to make this complier for Jack. It starts by parsing a CFG `jack.cfg` into a data structure of Production Rules and calculates the `first` `follow` and `predict` sets of all the symbols of the Grammar to then make a Parse Table.

From there, I do a table parse and make a parse tree and preform Syntax Directed Translation to make the AST.
If you would like to view this parse/AST tree, Prof. Hellman provided me a nice python script to do so `tree-to-graphvis`

### AST Visualization script (Credit Prof. Hellman):
`cat parsetree.txt | ./tree-to-graphvis | dot -Tpng -o parsetree.png` after running my bin

After that my code goes through and reads the AST in a recursive manner based on the Nonterminals and produces the `vm` code.

## 4) 
See `LANGINFO`, used most recent version of Golang `Go.1.22.2`

# Requirements/Dependencies
`jack.cfg` must be in the same dir as the binary as it is the CFG used for my algorithms and code base

all `.go` files must be in the same dir 

## Running
`make`

`./JackCompiler <source>`

- Followed directions this time, sorry about last time. No args will search current dir for `.jack` files and if given one argument will either read the dir or single `.jack` file it was given.

#### Thanks for a great semester and your hard work this was an enjoyment!