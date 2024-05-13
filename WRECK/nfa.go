package main

import (
	"fmt"
	"os"
)

var T [][]int
var L [][]int
var alphabet []string = []string{} 
var numStates int = -1
var alphabetLookup map[string]int = map[string]int{}
var hexAlphabet []byte = []byte{}
var lambdaChar string = "Z"


func addState()int{
	s:=make([]int, len(alphabet))
	for i := range s{
		s[i]=-1
	}
	T = append(T, s)
	numStates++


	l:=make([]int, numStates)
	for i := range l{
		l[i]=-1
	}
	L = append(L, l)
	for i:= range L{
		L[i] = append(L[i], -1)
	}

	
	
	return numStates
}

func addLambda(src int, dest int){
	
	L[src][dest] = 0 // *
}

func addEdge(c string, src int, dest int){
	// fmt.Println("ADDINNG",alphabetLookup[c],c,src,dest)
	// fmt.Println(alphabetLookup)
	// _,x := alphabetLookup[c]
	// fmt.Println(x)
	T[src][alphabetLookup[c]] = dest 
}

func nodeSeq(current *Node, src int, dest int){
	t := src
	childDest:= 0
	for _,child := range current.children{
		childDest = addState()
		
		switch child.data{
		case "dot":
			nodeDot(child,t,childDest)
		case "range":
			nodeRange(child,t,childDest)
		case "SEQ":
			nodeSeq(child,t,childDest)
		case "ALT":
			nodeAlt(child,t,childDest)
		case "kleene":
			nodeKleene(child,t,childDest)
		case "plus":
			nodePlus(child,t,childDest)
		case "lambda":
			nodeLambda(child,t,childDest)
		default:
			nodeLeaf(child,t,childDest)
		}
		t=childDest
	}
	
	addLambda(childDest,dest)
}


func nodeAlt(current *Node, src int, dest int){
	t:=src
	for _,child := range current.children{
		t = addState()
		addLambda(src,t)
		childDest := addState()
		switch child.data{
		case "dot":
			nodeDot(child,t,childDest)
		case "range":
			nodeRange(child,t,childDest)
		case "SEQ":
			nodeSeq(child,t,childDest)
		case "ALT":
			nodeAlt(child,t,childDest)
		case "kleene":
			nodeKleene(child,t,childDest)
		case "plus":
			nodePlus(child,t,childDest)
		case "lambda":
			nodeLambda(child,t,childDest)
		default:
			// fmt.Println("DEBUG",child.data,t,childDest)
			nodeLeaf(child,t,childDest)
		}
		
		addLambda(childDest,dest)

	}
	

}

func nodeLambda(current *Node, src int, dest int){
	
	addLambda(src,dest)
}


func nodePlus(current *Node, src int, dest int){
	data := current.children[0].data
	child := current.children[0]

	t := addState()
	switch data{
	case "dot":
		nodeDot(child,src,t)
	case "range":
		nodeRange(child,src,t)
	case "SEQ":
		nodeSeq(child,src,t)
	case "ALT":
		nodeAlt(child,src,t)
	case "kleene":
		nodeKleene(child,src,t)
	case "plus":
		nodePlus(child,src,t)
	case "lambda":
		nodeLambda(child,src,t)
	default:
		nodeLeaf(child,src,t)
	}
	
	
	addLambda(t,src)
	addLambda(t,dest)
}


func charRange(start, end rune) []string {
    
	if start > end {
        fmt.Println("Semantic Error: start character must be less than or equal to end character",string(start),string(end))
        os.Exit(3)
    }
    result := make([]string,0)
    for ch := start; ch <= end; ch++ {
        result = append(result, string(ch))
    }
    return result
}

func nodeRange(current *Node, src int, dest int){
	
	start := rune(convertHx(current.children[0].data)[0])
	end := rune(convertHx(current.children[1].data)[0])
	r := charRange(start,end)
	// fmt.Println(r)
	for _,v := range r{
		// fmt.Println(dest)
		addEdge(convertAlpha(v),src,dest)
	}

}

func nodeDot(current *Node, src int, dest int){
	for _,v := range alphabet{
		addEdge(string(v),src,dest)
	}
}

func nodeKleene(current *Node, src int, dest int){
	first := addState()
	addLambda(src,first)
	out := addState()
	addLambda(out,dest)
	data := current.children[0].data
	child := current.children[0]

	
	switch data{
	case "dot":
		nodeDot(child,first,out)
	case "range":
		nodeRange(child,first,out)
	case "SEQ":
		nodeSeq(child,first,out)
	case "ALT":
		nodeAlt(child,first,out)
	case "kleene":
		nodeKleene(child,first,out)
	case "plus":
		nodePlus(child,first,out)
	case "lambda":
		nodeLambda(child,first,out)
	default:
		nodeLeaf(child,first,out)
	}
	
	addLambda(first,out)
	addLambda(out,first)
	
}



func nodeLeaf(current *Node, src int, dest int){
	addEdge(current.data,src,dest)
}

func makeNFA(ast *Node, filename string){
	T = nil // clear globals
	L = nil
	
	// make NFA
	numStates = -1 // reset
	// fmt.Println(alphabet)
	alphaIndLookup := make(map[int]string)
	for i,v := range alphabet{
		alphabetLookup[v] = i
		alphaIndLookup[i] = v
	}

	// fmt.Println(alphaIndLookup)
	// fmt.Println(alphabetLookup)
	acceptStates := make(map[int]bool)
	acceptStates[1]=true

	// init start and goal
	addState()
	addState()
	
	// fmt.Println(T)
	// fmt.Println(L)
	switch ast.data{
	case "dot":
		nodeDot(ast,0,1)
	case "range":
		nodeRange(ast,0,1)
	case "SEQ":
		nodeSeq(ast,0,1)
	case "ALT":
		nodeAlt(ast,0,1)
	case "kleene":
		nodeKleene(ast,0,1)
	case "plus":
		nodePlus(ast,0,1)
	case "lambda":
		nodeLambda(ast,0,1)
	default:
		nodeLeaf(ast,0,1)

	}
	// fmt.Println("======  T ======")
	// for _,v:= range T{
	// 	fmt.Println(v)
	// }
	// fmt.Println("===== L ======")
	// for _,v:= range L{
	// 	fmt.Println(v)
	// }

	toNFA := ""

	alphaHeader := alphabetEncoded(hexAlphabet)
	a := []string{}
	for i:=0; i<len(alphaHeader); i+=3{
		a = append(a, alphaHeader[i:i+3])
	}
	alphaHeader = ""
	for i,v:= range a{
		alphaHeader+=v
		if i !=len(a)-1{
			alphaHeader+=" "
		}
	}

	header := fmt.Sprintf("%d %s %s\n",numStates+1,lambdaChar,alphaHeader)
	toNFA+=header

	for i:=0; i<len(T); i++{
		for j:=0; j<len(T[0]); j++{
			if T[i][j]==-1{
				continue
			}
			toNFA += fmt.Sprintf("- %d %d %s\n",i,T[i][j],alphaIndLookup[j])
		}
	}
	for i:=0; i<len(L); i++{
		for j:=0; j<len(L[0]); j++{
			if L[i][j]==0{
				toNFA += fmt.Sprintf("- %d %d %s\n",i,j,lambdaChar)
			}
		}
	}
	toNFA += "+ 1 1"
	// fmt.Println(alphabetLookup)
	// fmt.Println(toNFA)

	writeToFile(filename,toNFA)

}