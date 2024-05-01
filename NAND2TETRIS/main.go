package main

import (
	"fmt"
	"os"
	"strings"
	"regexp"
)

func main(){
	args := os.Args
	grammar := readLines("jack.cfg")
	vmCode := ""
	jackFile,_ := regexp.Compile(`(.+\.jack)$`)
	
	// SAME DIR i.e. NO ARGS...
	if len(args) == 1{
		currPath,err := os.Getwd()
		if err!=nil{
			panic(err)
		}
		files,err := os.ReadDir(currPath)
		if err!=nil{
			panic(err)
		}
		jackFiles := make([]string,0)
		for _,f := range files{
			if jackFile.MatchString(f.Name()){
				// fmt.Println(f.Name())
				jackFiles = append(jackFiles, f.Name())

			}
		}
		if len(jackFiles) == 0{
			println("No .jack files found :^(")
		}

		for _,f := range jackFiles{
			ast := AST(grammar,f)
			
			vmCode = codeGen(ast)
			writeToFile(fmt.Sprintf("%s.vm",strings.Split(f,".")[0]),vmCode)
		}
	
	// GIVEN a <source>
	} else if len(args) == 2{
		arginfo,err := os.Stat(args[1])
		if err != nil{
			panic(err)
		}

		// DIRECTORY
		if arginfo.IsDir(){
			files,err := os.ReadDir(args[1])
			if err!=nil{
				panic(err)
			}

			jackFiles := make([]string,0)
			for _,f := range files{
				if jackFile.MatchString(f.Name()){
					// fmt.Println(f.Name())
					jackFiles = append(jackFiles, args[1]+"/"+f.Name())

				}
			}
			if len(jackFiles) == 0{
				println("No .jack files found :^(")
			}

			for _,f := range jackFiles{
				ast := AST(grammar,f)
				
				vmCode = codeGen(ast)
				writeToFile(fmt.Sprintf("%s.vm",strings.Split(f,".")[0]),vmCode)
			}	
			
		// Single File
		} else{
			if !jackFile.MatchString(args[1]){
				println("Must be .jack file!")
				os.Exit(1)
			}

			ast := AST(grammar,args[1])
			vmCode = codeGen(ast)
			writeToFile(fmt.Sprintf("%s.vm",strings.Split(args[1],".")[0]),vmCode)
			fmt.Println()
		}

	} else{
		println("USEAGE: ./JackCompiler <source>\nWhere if no source is given the current dir is searched for .jack files")
		os.Exit(1)
	}

	fmt.Println(vmCode)
	
	
}