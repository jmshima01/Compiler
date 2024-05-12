package main

import (
	"fmt"
	"os"
)

var T [][]int
var L [][]int
var alphabet string = "" 
var numStates int = -1
var alphabetLookup map[string]int = map[string]int{}


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
			nodeAlt(child,t,dest)
		case "kleene":
			nodeKleene(child,src,childDest)
		case "plus":
			nodePlus(child,t,childDest)
		// case "lambda":
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
			nodeAlt(child,t,dest)
		case "kleene":
			nodeKleene(child,src,childDest)
		case "plus":
			nodePlus(child,t,childDest)
		// case "lambda":
		default:
			nodeLeaf(child,t,childDest)
		}
		
		addLambda(childDest,dest)

	}
	

}



func nodePlus(current *Node, src int, dest int){
	data := current.children[0].data
	nxt := addState()
	addEdge(data,src,nxt)
	addEdge(data,nxt,nxt)
	fmt.Println(dest,nxt)
	addLambda(nxt,dest)
}


func charRange(start, end rune) []string {
    
	if start > end {
        fmt.Println("Invalid range found: start character must be less than or equal to end character",string(start),string(end))
        os.Exit(2)
    }
    result := make([]string,0)
    for ch := start; ch <= end; ch++ {
        result = append(result, string(ch))
    }
    return result
}

func nodeRange(current *Node, src int, dest int){
	start := rune(current.children[0].data[0])
	end := rune(current.children[1].data[0])
	r := charRange(start,end)
	fmt.Println(r)
	for _,v := range r{
		fmt.Println(dest)
		addEdge(v,src,dest)
	}
}

func nodeDot(current *Node, src int, dest int){
	for _,v := range alphabet{
		addEdge(string(v),src,dest)
	}
}

func nodeKleene(current *Node, src int, dest int){
	fir := addState()
	sec := addState()
	addLambda(src,fir)
	addLambda(fir,sec)
	addLambda(sec,dest)
	addLambda(sec,fir)
	addEdge(current.children[0].data,fir,sec)

}



func nodeLeaf(current *Node, src int, dest int){
	addEdge(current.data,src,dest)
}

func makeNFA(ast *Node){

	// make NFA
	

	for i,v := range alphabet{
		alphabetLookup[string(v)] = i
	}

	acceptStates := make(map[int]bool)
	acceptStates[1]=true

	// init start and goal
	addState()
	addState()
	
	fmt.Println(T)
	fmt.Println(L)
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
	// case "lambda":
	default:
		nodeLeaf(ast,0,1)

	}
	fmt.Println("======  T ======")
	for _,v:= range T{
		fmt.Println(v)
	}
	fmt.Println("===== L ======")
	for _,v:= range L{
		fmt.Println(v)
	}


}