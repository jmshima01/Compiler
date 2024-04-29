package main

import (
	"fmt"
)


type SymbolData struct{
	kind string;
	offset int;
	name string;
	symbolType string;
}

var G map[string]SymbolData = map[string]SymbolData{}
var L map[string]SymbolData = map[string]SymbolData{}

var curClass string = ""
var counter int = 0
var numGlob int = 0


func handleParamlist(ast *Node,isMethod bool) {
	argCounter := 0
	if isMethod{
		this:= SymbolData{kind:"this",offset: argCounter,symbolType: curClass,name: "this"}
		argCounter++
		L[this.name] = this
	}
	fmt.Println("param lis", ast.children[0].data)
	symbol := SymbolData{kind: "argument", offset: argCounter}
	for i := 0; i < len(ast.children); i += 2 {
		
		symbol.symbolType = ast.children[i].data
		symbol.name = ast.children[i+1].children[0].data
		L[symbol.name] = symbol
		argCounter++
		symbol.offset = argCounter
	}
}

func handleBody(ast *Node) (int,string) {
	locals := 0
	res:=""
	for _, v := range ast.children {
		if v.data == "VarDec"{
			handleVarDec(v, &locals)
		
		
		}else{
			switch v.data {
			case "LetStatement":
				res+=handleLet(v)
			case "DoStatement":
				res+=handleDo(v)
			case "IfStatement":
				res+=handleIf(v)
			case "WhileStatement":
				res+=handleWhile(v)
			case "ReturnStatement":
				res+= handleReturn(v)
			}
		}
	}
	
	return locals,res

}



func handleVarDec(ast *Node, locals *int) {

	symbol := SymbolData{offset: *locals}
	for i, v := range ast.children {
		if i == 0 {
			symbol.symbolType = v.data
			symbol.kind = "local"
		} else {
			if v.data == "VarName" {
				*locals += 1
				symbol.name = v.children[0].data
				L[symbol.name] = symbol
				symbol.offset += 1
			}
		}
	}
}



func handleLet(ast *Node) string {
	res := ""
	
	res+=handleExpression(ast.children[1])
	res+=handleVarName(ast.children[0])
	
	res += fmt.Sprintf("pop %s %d\n", L[ast.children[0].children[0].data].kind,L[ast.children[0].children[0].data].offset)
	return res
}

func handleDo(ast *Node) string{
	return handleSubCall(ast) + "pop temp 0\n"
}

func handleWhile(ast *Node) string{
	res:=fmt.Sprintf("label While_%d\n",counter)
	res+= handleExpression(ast.children[1])+"not\n"

	res+=fmt.Sprintf("if-goto WHILE_END_%d\n",counter)

	for _, v := range ast.children {
		switch v.data {
		case "LetStatement":
			res+=handleLet(v)
		case "DoStatement":
			res+=handleDo(v)
		case "IfStatement":
			res+=handleIf(v)
		case "WhileStatement":
			res+=handleWhile(v)
		case "ReturnStatement":
			res+= handleReturn(v)
		}
	}
	res+=fmt.Sprintf("goto WHILE_%d\n",counter)
	res+=fmt.Sprintf("label WHILE_END_%d\n",counter)
	counter++
	
	return res

}

func handleReturn(ast *Node) string {
	res := ""
	if len(ast.children) > 0 {
		res =handleExpression(ast.children[0]) + "return"
	} else {
		res = "push constant 0\nreturn"
	}
	return res

}

func handleIf(ast *Node) string{
	res := handleExpression(ast.children[0])
	res += "not\n"
	res += fmt.Sprintf("if-goto IF_FALSE_%d\n",counter)
	
	for _, v := range ast.children {
		switch v.data {
		case "LetStatement":
			res+=handleLet(v)
		case "DoStatement":
			res+=handleDo(v)
		case "IfStatement":
			res+=handleIf(v)
		case "WhileStatement":
			res+=handleWhile(v)
		case "ReturnStatement":
			res+= handleReturn(v)
		}
	}
	res += fmt.Sprintf("goto IF_END_%d\n",counter)

	counter++
	return res
}

func handleExpression(ast *Node) string {
	res:=""
	fmt.Println(ast.children[0].data,"HIIIIIIIII")
	for i := 0; i < len(ast.children)-1; i++ {
		v := ast.children[i]
		if v.data == "Op" || v.data == "UnaryOp" {
			temp := ast.children[i+1] //swap for postfix notation of vm stack machine
			ast.children[i+1] = v
			ast.children[i] = temp
			i++
		}
	}

	for _, v := range ast.children {
		switch v.data {
		case "Expression":
			res+=handleExpression(v)
		case "integerconstant":
			res+=handleInt(v)
		case "stringconstant":
			res+=handleStr(v)
		case "VarName":
			res+=handleVarName(v)
		case "KeywordConstant":
			res+=handleKeyword(v)
		case "Op":
			res+=handleOp(v)
		case "SubroutineCall":
			res+=handleSubCall(v)
		}
	}
	
	return res
}

func handleInt(ast *Node) string {
	return fmt.Sprintf("push constant %s\n", ast.children[0].data)
}

func handleStr(ast *Node) string{
	res := fmt.Sprintf("push constant %d\n",len(ast.children[0].data))
	res += "call String.new 1\n"
	for _,v := range ast.children[0].data{
		res += fmt.Sprintf("push constant %v\ncall String.appendChar 2\n",v)

	}
	return res
}

func handleKeyword(ast *Node) string {
	res := ""
	switch ast.children[0].data {
	case "null":
		res = "push constant 0\n"
	case "false":
		res = "push constant 0\n"
	case "this":
		res = "push pointer 0\n"
	case "true":
		res = "push constant 0\nnot\n"
	}
	return res
}

func handleOp(ast *Node) string {
	res := ""
	switch ast.children[0].data {
	case "+":
		res = "add\n"
	case "-":
		res = "sub\n"
	case "*":
		res = "Math.multiply 2\n"
	case "/":
		res = "Math.divide 2\n"
	case "pipe":
		res = "or\n"
	case "&":
		res = "and\n"
	case "=":
		res = "eq\n"
	case "<":
		res = "gt\n"
	case ">":
		res = "lt\n"
	}
	return res
}

func handleVarName(ast *Node) string {
	res := ""
	g ,inGlobal:= G[ast.children[0].data]
	l, inLocal := L[ast.children[0].data]
	if inLocal {
		res = fmt.Sprintf("push %s %d\n", l.kind, l.offset)
	} else if inGlobal{
		res = fmt.Sprintf("push %s %d\n", g.kind, g.offset)
	}
	return res
}

func handleSubCall(ast *Node) string {
	l, inLocal := L[ast.children[0].data]
	res:=""
	
	v:= ast.children[len(ast.children)-1]
	if v.data == "ExpressionList"{
		temp,exprCount:=handleExprList(v)
		if len(ast.children)==2{
			res = fmt.Sprintf("push pointer 0\n%scall %s.%s %d\n",temp,curClass,ast.children[0].data,exprCount)
		} else{
			if inLocal{
				res= fmt.Sprintf("%spush %s %d\ncall %s.%s %d\n",temp,l.kind,l.offset,curClass,ast.children[0].data,exprCount)
			} else{
				res+= fmt.Sprintf("%scall %s.%s %d\n",temp,ast.children[0].data,ast.children[1].data,exprCount)

			}
		}
	
	} else {
		if len(ast.children) == 1 {
			res= fmt.Sprintf("call %s.%s 1\n", curClass, ast.children[0].data)
		
		}else if inLocal {
			res = fmt.Sprintf("call %s.%s 1\n", l.symbolType, ast.children[1].data)
		} else {
			res= fmt.Sprintf("call %s.%s 0\n", ast.children[0].data, ast.children[1].data)
		}
	}
	
	return res
}

func handleExprList(ast *Node)(string,int){
	res:=""
	
	exprCount:=0
	for _,v:= range ast.children{
		exprCount++
		res+=handleExpression(v)
		
	}
	return res,exprCount
}

func handleClassVar(ast *Node){

	symbol := SymbolData{offset: numGlob}
		
		for i, v := range ast.children {
			if i == 0 {
				if v.data == "static" {
					symbol.kind = "static"
					
				} else {
					symbol.kind = "field"
					
				}
				symbol.kind = v.data
			} else if i == 1 {
				symbol.symbolType = v.data
			} else {
				if v.data == "ClassVarDec" {
					handleClassVar(v)
				}else{
					symbol.name = v.children[0].data
					G[symbol.name] = symbol
					symbol.offset++
					numGlob++
				}
			}
		}
}


func handleSubrDec(ast *Node)string{
	res:=""
	// clear(L) // reset local symbol table NOTE:go1.22 clear()

	numLocals := 0
	res,r,s:="","",""
	for _, v := range ast.children {
		if v.data == "SubroutineDec"{
			s=handleSubrDec(v)
		}
		if v.data == "ParameterList" {
			handleParamlist(v,false)

		} else if v.data == "SubroutineBody" {
			numLocals,r = handleBody(v)
		}
	}
	res += fmt.Sprintf("function %s.%s %d\n%s\n%s", curClass, ast.children[2].data, numLocals,r,s)
	
	

	if ast.children[0].data == "constructor"{
		res += fmt.Sprintf("push constant %d\ncall Memory.alloc 1\npop pointer 0\n%s",numGlob,r)
	}
	
	return res
	
}

func codeGen(ast *Node) string {
	code := ""
	switch ast.data {

	case "Class":
		numGlob = 0
		for i,v := range ast.children{
			if i == 0{
				curClass = v.data
				continue
			}
			if v.data == "ClassVarDec"{
				handleClassVar(v)
			} 
			if v.data == "SubroutineDec"{
				fmt.Println("HADKELL,",v.children[2].data)
				code += handleSubrDec(v)
								
			}
		}
	}
	fmt.Println(G)
	fmt.Println(L)
	fmt.Println("ENDED",numGlob)

	return code
}
