package main

import(
	"fmt"
	"os"
	"strings"
)

// James' DATA STRUCTURES
// ds.go


// ============ Typedefs ==============

type ItemSet map[string]Item;

type Item struct{
	productionMarker int;
	lhs string;
	rhs []string;
}

func makeItem(rule ProductionRule, marker int)Item{
	return Item{lhs: rule.lhs,rhs: rule.rhs, productionMarker: marker}
}


type AdjList []Pair;

type Pair struct{
	itemsetID int;
	transitionSymbol string;
}

type CFSM map[int]AdjList

func (i Item)toString()string{
	
	return fmt.Sprint(i)
}



type SymbolData struct{
	hasType string;
	isConst bool;
	isUsed bool;
	initialized bool;
	lineNum int;
	columnNum int;
	name string;
}

type SymbolTable map[string]SymbolData;

type Node struct{
	data string;
	parent *Node;
	id int;
	children []*Node;
}

type ProductionRule struct{
	lhs string;
	rhs []string;
};



type token struct{
	value string;
	tokenType string;
}


type set map[string]bool;
type queue []token;
type stack []string;

// Type Methods: 
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

// ======== STACK AND QUEUE ==========

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
func (q *queue) push(v token) {
	*q = append(*q, v) 
}

func (q *queue) peek()token{
	if q.isEmpty(){
		fmt.Println("EMPTY Q")
		return token{}
	}
	return (*q)[0] 
}

func (q *queue) popfront()token {
	if q.isEmpty() {
		return token{}
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

func makeNode(s string,parent *Node, id int) *Node{
	n := Node{parent: parent, data:s}
	n.children = make([]*Node, 0)
	n.id = id
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
	fmt.Println("data",t.data,t.id)
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

func genNodeInfo(root *Node, s *string) *string{
	if root == nil{
		return s
	}

	r := *(root)
	if r.data == "lambda"{ // avoid implicit label for python script
		*s+=fmt.Sprintf("%d %s\n",r.id,"Lambda")
	// } else if r.data == "identifier" || r.data == "objectname" || r.data == "array" || r.data == "subroutinename" || r.data == "integerconstant" || r.data == "stringconstant"{
	// 	*s+= fmt.Sprintf("%d %s\n",r.id,r.identInfo)
	} else{
		*s+=fmt.Sprintf("%d %s\n",r.id,r.data)
	}
	
	for _,child := range root.children{
		genNodeInfo(child,s)
	}
	return s
}

func genEdgeInfo(root *Node, dfs *string) *string{
	if root == nil{
		return dfs
	}
	v := *(root)
	*dfs+= fmt.Sprintf("%d",v.id)
	for _,child := range v.children{
		*dfs+=fmt.Sprintf(" %d",child.id)
		
	}
	*dfs+="\n"
	for _,x := range root.children{
		child := *(x)
		if len(child.children) != 0{ // if not leaf
			genEdgeInfo(x,dfs)
		}
	}

	return dfs
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

func writeToFile(path string,raw string){
	err := os.WriteFile(path,[]byte(raw),0644)
	if err != nil{
		panic(err)
	}
}
