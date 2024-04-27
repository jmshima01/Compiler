package main

import (
	"fmt"
)




type SymbolData struct{
	name string;
	symbolType string;
	kind string;
	offset int;
}

var globalSymbolTable map[string]SymbolData = map[string]SymbolData{}
var localSymbolTable map[string]SymbolData =  map[string]SymbolData{}
var currClassName string = ""
var globalOffset int = 0
var localOffset int = 0
var counter int = 0
var numOfLocals int = 0
var currentFunctionName string = ""



func codeGen(ast *Node, code *string) *string{
	
	// printTree(ast)
	switch ast.data {
	case "Class":
		currClassName = ast.children[0].data
		fmt.Println("GOT HERE!")
		localOffset = 0
		globalOffset = 0
		
	case "SubroutineDec":
		*code += fmt.Sprintf("function %s.%s %d\n",currClassName,ast.children[2].data,numOfLocals)
		numOfLocals = 0
	
	case "ClassVarDec":
		symbol := SymbolData{offset:globalOffset}
		for i,v := range ast.children{
			if i == 0{
				if v.data == "static"{symbol.kind="static"}else{symbol.kind="this"}
				symbol.kind = v.data
			} else if i == 1{
				symbol.symbolType = v.data
			} else{
				if v.data == "ClassVarDec"{
					continue
				}
				symbol.name = v.data
				globalSymbolTable[v.data] = symbol
				globalOffset+=1
				symbol.offset+=1 

			}
		}
	case "VarDec":
		symbol := SymbolData{offset: localOffset}
		for i,v := range ast.children{
			if i == 0{
				symbol.symbolType = v.data
				symbol.kind = "local"
			} else{
				if v.data == "VarDec"{
					continue
				}
				numOfLocals +=1
				symbol.name = v.data
				localSymbolTable[v.data] = symbol
				globalOffset+=1
				symbol.offset+=1 

			}
		}
	case "SubroutineCall":
		l,inLocal := localSymbolTable[ast.children[0].data]
		
		if len(ast.children) == 2{
			if inLocal{
				*code += fmt.Sprintf("call %s.%s 0\n",l.symbolType,ast.children[1].data)
			} else{
				*code += fmt.Sprintf("call %s.%s 0\n",ast.children[0].data,ast.children[1].data)
			}
		} else if len(ast.children) == 1{
			*code += fmt.Sprintf("call %s.%s 0\n",currClassName,ast.children[0].data)
		}
	case "DoStatement":
		l,inLocal := localSymbolTable[ast.children[0].data]
		
		if len(ast.children) == 2{
			if inLocal{
				*code += fmt.Sprintf("call %s.%s 0\n",l.symbolType,ast.children[1].data)
			} else{
				*code += fmt.Sprintf("call %s.%s 0\n",ast.children[0].data,ast.children[1].data)
			}
		} else if len(ast.children) == 1{
			*code += fmt.Sprintf("call %s.%s 0\n",currClassName,ast.children[0].data)
		}
		*code += fmt.Sprintf("pop temp 0\n")	
	case "LetStatement":
		numOfLocals+=1

	case "Expression":


	case "ParameterList":
		argCounter:= 0
		symbol:= SymbolData{kind:"argument",offset: argCounter}
		for i:=0; i<len(ast.children); i+=2{
			symbol.symbolType = ast.children[i].data
			symbol.name = ast.children[i+1].data
			localSymbolTable[symbol.name] = symbol
			argCounter++
			symbol.offset=argCounter
		}


	case "ExpressionList":

	default:
	
	
	
	}

	for _,child := range ast.children{
		codeGen(child,code)
	}
	return code
	
}