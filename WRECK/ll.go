package main

import(
	"fmt"
	"os"
	"strings"
)

// ============ Typedefs ==============
type ParseTree struct{
	root *Node
}

type Node struct{
	data string;
	parent *Node;
	
	children []*Node;
}

type ProductionRule struct{
	lhs string;
	rhs []string;
};

type set map[string]bool;
type queue []string;
type stack []string;

// ========== Type Methods =============


// =========== Prod Rule ================
func (P ProductionRule) toString() string{
	s:= ""
	s+=P.lhs
	s+=" -> "
	for _,v:= range P.rhs{
		s+= v + " "
	}
	return s[:len(s)-1] 
}

// =========== STACK AND QUEUE ==========

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}
func (s *stack) push(v string) {
	*s = append(*s, v) 
}

func (s *stack) peek()string{
	if s.isEmpty(){
		fmt.Println("EMPTY S")
		return ""
	}
	return (*s)[len(*s)-1] 
}

func (s *stack) pop() (string) {
	if s.isEmpty() {
		return ""
	} else {
		index := len(*s)-1 
		element := (*s)[index] 
		*s = (*s)[:index] 
		return element
	}
}

func (q *queue) isEmpty() bool {
	return len(*q) == 0
}
func (q *queue) push(v string) {
	*q = append(*q, v) 
}

func (q *queue) peek()string{
	if q.isEmpty(){
		fmt.Println("EMPTY Q")
		return ""
	}
	return (*q)[0] 
}

func (q *queue) popfront() (string) {
	if q.isEmpty() {
		return ""
	} else {
		element := (*q)[0]
		*q = (*q)[1:] 
		return element
	}
}

// ========= SETS ==========

func (s set)add(v string){
	if s == nil{
		return
	}
	s[v] = true
}

// only want keys of the set/map
func (s set)getValues() []string{
	res := make([]string,0)
	for k,_ := range s{
		res = append(res, k)
	}
	return res
}

func setUnion(s1 set, s2 set)set{
	union := make(set)
	for k, _ := range s1{
		union[k] = true
	}
	for k, _ := range s2{
		union[k] = true
	}
	return union
}


// ======== PARSE TREE ==============

func makeNode(s string,parent *Node) *Node{
	n := Node{parent: parent, data:s}
	n.children = make([]*Node, 0)
	return &n
}

func addChild(t *Node, child *Node){
	t.children = append(t.children, child)
}


func printChildren(c []*Node){
	res := "children->["
	for _,v := range c{
		res+=v.data + " "
	}
	res+="]"
	fmt.Println(res) 
}

func (t Node) debug(){
	fmt.Println("------")
	if t.parent != nil{
		fmt.Println("Parent:",t.parent.data)
	}
	fmt.Println("data",t.data)
	printChildren(t.children)
}

func printTree(t *Node){
	if t == nil{
		return
	}
	
	v := *(t)
	v.debug()
	for _,x := range t.children{
		printTree(x)
	}
	return
}

// ========== File I/O =================
func readLines(path string) []string{
	f,err := os.ReadFile(path)
	if err!=nil{
		panic(err)
	}
	data := strings.ReplaceAll(string(f),"\r\n","\n")
	return strings.Split(strings.Trim(data,"\n"),"\n")
	
}

func writeToFile(path string,asm string){
	err := os.WriteFile(path,[]byte(asm),0644)
	if err != nil{
		panic(err)
	}
}
