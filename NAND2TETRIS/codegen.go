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

var numGlob int = 0
// var numStatic int = 0
var curSubDecType string = ""
var counterw,counterif = -1,-1

func handleParamlist(ast *Node) string{
	argCounter:=0
	if curSubDecType == "method"{
		argCounter++
	}
	
	
	res:= ""
	symbol := SymbolData{kind: "argument", offset: argCounter}
	for i := 0; i < len(ast.children); i += 2 {
		
		symbol.symbolType = ast.children[i].data
		symbol.name = ast.children[i+1].children[0].data
		L[symbol.name] = symbol
		
		// res+=fmt.Sprintf("%s %d",symbol.kind,argCounter)
		
		argCounter++
		symbol.offset = argCounter
	}
	return res
}

func handleBody(ast *Node) (int,string) {
	locals := 0
	res:=""
	// counterw = 0
	// counterif = 0
	for _, v := range ast.children {
		if v.data == "VarDec"{
			handleVarDec(v, &locals)
		
		} else{
			switch v.data {
			case "LetStatement":
				res+=handleLet(v)
			case "DoStatement":
				res+=handleDo(v)
			case "IfStatement":
				res+=handleIf(v)
			case "ElseStatement":
				res+=handleElse(v)
			case "WhileStatement":
				
				res+=handleWhile(v)
				counterw--
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
	if ast.children[0].children[0].data == "ArrayName"{ // ARRAY from AST --> a [ Expr ] 
		
		// res+="TESTARR11!\n"
		res+=handleExpression(ast.children[0].children[0].children[1])
		res += handleExprVarName(ast.children[0].children[0])
		res+= "add\n"
		
		res+= handleExpression(ast.children[1])
		res+="pop temp 0\n"
		res+="pop pointer 1\npush temp 0\npop that 0\n"
										                    


	} else{
		l,inLocal := L[ast.children[0].children[0].data]
		g := G[ast.children[0].children[0].data]
		res+=handleExpression(ast.children[1])
		
		if inLocal{
			res += fmt.Sprintf("pop %s %d\n", l.kind,l.offset)
		} else { // must be global then!
			if g.kind != "static"{
				switch curSubDecType{
				case "method":
					res += fmt.Sprintf("pop %s %d\n","this",g.offset)
				case "constructor":
					res += fmt.Sprintf("pop %s %d\n","this",g.offset)
				default: // "function --> field"
					res += fmt.Sprintf("pop %s %d\n", g.kind,g.offset)

				}
			} else{res += fmt.Sprintf("pop %s %d\n", g.kind,g.offset)}
		}
	}
	return res
}

func handleDo(ast *Node) string{
	return handleSubCall(ast) + "pop temp 0\n"
}

func handleWhile(ast *Node ) string{
	counterw++
	res := ""
	res=fmt.Sprintf("label WHILE_EXP%d\n",counterw)
	res+= handleExpression(ast.children[0])+"not\n"

	res+=fmt.Sprintf("if-goto WHILE_END%d\n",counterw)

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
	res+=fmt.Sprintf("goto WHILE_EXP%d\n",counterw)
	res+=fmt.Sprintf("label WHILE_END%d\n",counterw)
	
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
	
	counterif++
	res := handleExpression(ast.children[0])
	res += fmt.Sprintf("if-goto IF_TRUE%d\n",counterif)
	res += fmt.Sprintf("goto IF_FALSE%d\n",counterif)
	res += fmt.Sprintf("label IF_TRUE%d\n",counterif)
	for _,v := range ast.children {
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

	res += fmt.Sprintf("label IF_FALSE%d\n",counterif)
	
	return res
}


func handleElse(ast *Node)string{
	
	res:=fmt.Sprintf("goto IF_END%d\n",counterif)
	for _,v := range ast.children {
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
	res += fmt.Sprintf("label IF_END%d\n",counterif)
	return res
}



func handleExpression(ast *Node) string {
	res:=""
	// fmt.Println(ast.children[0].data,"HIIIIIIIII")
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
			res+=handleExprVarName(v)
		case "KeywordConstant":
			res+=handleKeyword(v)
		case "Op":
			res+=handleOp(v)
		case "SubroutineCall":
			res+=handleSubCall(v)
		case "UnaryOp":
			res+=handleUnaryOp(v)
		}
	}
	
	return res
}

func handleExprVarName(ast *Node)string{
	res:=""
	if ast.children[0].data == "ArrayName"{
		
		res+=handleExpression(ast.children[0].children[1])
		l,inLocal := L[ast.children[0].children[0].data]
		glob := G[ast.children[0].children[0].data]
		if inLocal{
			res += fmt.Sprintf("push %s %d\n",l.kind,l.offset)
		} else{
			if glob.kind!="static"{
				res += fmt.Sprintf("push this %d\n",glob.offset)
			} else{
				res += fmt.Sprintf("push static %d\n",glob.offset)
			}
		}
		res+= "add\npop pointer 1\npush that 0\n"

	} else{
		l,inLocal := L[ast.children[0].data]
		glob := G[ast.children[0].data]
		if inLocal{
			res = fmt.Sprintf("push %s %d\n",l.kind,l.offset)
		} else{
			if glob.kind!="static"{
				res = fmt.Sprintf("push this %d\n",glob.offset)
			} else{
				res = fmt.Sprintf("push static %d\n",glob.offset)
			}
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
		// _,inLocal := L["this"]
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
		res = "call Math.multiply 2\n"
	case "/":
		res = "call Math.divide 2\n"
	case "pipe":
		res = "or\n"
	case "&":
		res = "and\n"
	case "=":
		res = "eq\n"
	case "<":
		res = "lt\n"
	case ">":
		res = "gt\n"
	}
	return res
}

func handleUnaryOp(ast *Node)string{
	if ast.children[0].data == "-"{
		return "neg\n"
	} else{return "not\n"}
}


func handleSubCall(ast *Node) string {
	l, inLocal := L[ast.children[0].data]
	
	res:=""
	
	v:= ast.children[len(ast.children)-1]
	
	temp,exprCount:=handleExprList(v)
	if len(ast.children)==2{
		res += fmt.Sprintf("push pointer 0\n%scall %s.%s %d\n",temp,curClass,ast.children[0].data,exprCount+1)
	} else if inLocal{
		
		res+= fmt.Sprintf("%spush %s %d\ncall %s.%s %d\n",temp,l.kind,l.offset,l.symbolType,ast.children[1].data,exprCount+1)
	} else{
		glob,inGlob := G[ast.children[0].data]
		if inGlob{
			res+= fmt.Sprintf("%spush this %d\ncall %s.%s %d\n",temp,glob.offset,glob.symbolType,ast.children[1].data,exprCount+1)
		}else{
			res+= fmt.Sprintf("%scall %s.%s %d\n",temp,ast.children[0].data,ast.children[1].data,exprCount)
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
	
	// clear(L) // reset local symbol table NOTE:go1.22 clear()
	numLocals := 0
	res,r,s,p:="","","",""
	switch ast.children[0].data{
	case "constructor":
		curSubDecType = "constructor"
		for _, v := range ast.children {
			if v.data == "SubroutineDec"{
				s=handleSubrDec(v)
			}
			if v.data == "ParameterList" {
				p=handleParamlist(v)

			} else if v.data == "SubroutineBody" {
				numLocals,r = handleBody(v)
			}
		}
		res += fmt.Sprintf("function %s.%s %d\npush constant %d\ncall Memory.alloc 1\npop pointer 0\n%s%s\n%s", curClass, ast.children[2].data, numLocals,numGlob,p,r,s)
		return res
	
	case "method":
		curSubDecType = "method"
		L["this"]=SymbolData{name: "this",symbolType: curClass,offset: 0,kind:"argument"}
		for _, v := range ast.children {
			if v.data == "SubroutineDec"{
				s=handleSubrDec(v)
			}
			if v.data == "ParameterList" {
				p=handleParamlist(v)

			} else if v.data == "SubroutineBody" {
				numLocals,r = handleBody(v)
			}
		}
		res += fmt.Sprintf("function %s.%s %d\npush argument 0\npop pointer 0\n%s%s\n%s", curClass, ast.children[2].data, numLocals,p,r,s)
		return res
	
	default: // class function
		curSubDecType = "function"
		for _, v := range ast.children {
			if v.data == "SubroutineDec"{
				s=handleSubrDec(v)
			}
			if v.data == "ParameterList" {
				p=handleParamlist(v)

			} else if v.data == "SubroutineBody" {
				numLocals,r = handleBody(v)
			}
		}
		res += fmt.Sprintf("function %s.%s %d\n%s%s\n%s", curClass, ast.children[2].data, numLocals,p,r,s)
		return res
	}
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
				// fmt.Println("HADKELL,",v.children[2].data)
				code += handleSubrDec(v)
								
			}
		}
	}
	fmt.Println(G)
	fmt.Println(L)
	fmt.Println("ENDED",numGlob)

	return code
}
