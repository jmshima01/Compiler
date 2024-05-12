package main

import (
	"fmt"
	"os"
)

var T [][]int
var L [][]int
var alphabet string = "ABCDg" 
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
	for i:= range L{
		L[i] = append(L[i], -1)
	}

	L = append(L, l)


	// L = append(T, make([]int, currState))
	
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
		// case "SEQ":
		// 	nodeSeq(child,src,childDest)
		// // case "ALT":
		// // 	nodeAlt(current,src,dest)
		// case "kleene":
		// 	nodeKleene(child,src,childDest)
		case "plus":
			nodePlus(child,t,childDest)
		// case "lambda":
		// default:
		// 	nodeLeaf(child,src,childDest)
		}
		t=childDest
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

}



func nodeLeaf(current *Node, src int, dest int){

}

func makeNFA(ast *Node){

	// make NFA
	

	for i,v := range alphabet{
		alphabetLookup[string(v)] = i
	}

	acceptStates := make(map[int]bool)
	acceptStates[1]=true

	addState()
	addState()
	
	fmt.Println(T)
	fmt.Println(L)
	if ast.data == "SEQ"{
		nodeSeq(ast,0,1)
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