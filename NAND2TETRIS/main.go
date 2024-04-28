package main

import (
	"fmt"
	"os"
)

func main(){

	args := os.Args
	grammar := readLines("jack.cfg")
	vmCode := ""
	
	if len(args) == 1{
		currPath,err := os.Getwd()
		if err!=nil{
			panic(err)
		}
		files,err := os.ReadDir(currPath)
		if err!=nil{
			panic(err)
		}
		for _,f := range files{
			fmt.Println(f.Name())
		}
	
	} else if len(args) == 2{
		arginfo,err := os.Stat(args[1])
		if err != nil{
			panic(err)
		}
		if arginfo.IsDir(){	
			fmt.Println("isDir")
		} else{
			ast := AST(grammar,args[1])
			fmt.Println(ast)
			vmCode = codeGen(ast)
			fmt.Println(vmCode)
		}

	} else{
		println("USEAGE: ./JackCompiler source")
		os.Exit(1)
	}
	fmt.Println("\nVMcode:")
	fmt.Println(vmCode)
}