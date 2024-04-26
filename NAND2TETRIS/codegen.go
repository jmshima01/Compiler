package main

import (
	// "fmt"
)


type SymbolData struct{
	name string;
	symbolType string;
	kind string;
	index int;
}

var globalSymbolTable map[string]SymbolData = map[string]SymbolData{}
var localSymbolTable map[string]SymbolData =  map[string]SymbolData{}
var currClassName string = ""
var currOffset int = 0
var counter int = 0


func codeGen(ast *Node){
	if ast.data == "Class"{
		
	}
	
}