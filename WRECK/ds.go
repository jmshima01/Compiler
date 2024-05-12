package main

import(
	"fmt"
	"os"
	"strings"
	"strconv"
	"encoding/hex"
	// "regexp"
)

// James' DATA STRUCTURES
// ds.go


// ============ Typedefs ==============


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
		// fmt.Println("EMPTY Q")
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


// ALPHABET ENCODING <-->


// Converts weird hex format string to byte array
func parseAlphabetEncoding(s string)[]byte{
	s = strings.Join(strings.Fields(s),"")
	ascii_permited := make([]byte,0)
	for j:=0; j<len(s); j++{
		if string(s[j]) == "x"{
			val,err := strconv.ParseInt(s[j+1:j+3],16,8)
			// fmt.Printf("hex %x:\n",val)
			if err != nil{
				fmt.Println("error reading parse ascii",err)
				os.Exit(1)
			}
			ascii_permited = append(ascii_permited, byte(int(val)))
			j+=2
		} else{
			// fmt.Println("non:",string(first_line[j]))
			ascii_permited = append(ascii_permited, byte(int(s[j])))
		}
	} 
	return ascii_permited
}


func byteSliceToStringSlice(byteSlice []byte) []string {
    stringSlice := make([]string, len(byteSlice))
    for i, b := range byteSlice {
        stringSlice[i] = string(b)
    }
    return stringSlice
}

func alphabetEncoded(hx []byte) string{
	aEncoded := ""
	for _,c := range hx{
		x := strconv.FormatInt(int64(byte(c)), 16)
		
		if len(x) == 2{
			aEncoded += fmt.Sprintf("x%s",x)
		} else{
			aEncoded += fmt.Sprintf("x0%s",x)
		}

		// if c == byte(10) || c == byte(32) || c==byte(92) && i!=len(hexAlphabet)-1{
		// 	aEncoded+=" "
		// }
			
	}
	return aEncoded
}

func convertAlpha(s string)string{
	x := strconv.FormatInt(int64(byte(s[0])), 16)
	aEncoded:= ""	
	if len(x) == 2{
		aEncoded = fmt.Sprintf("x%s",x)
	} else{
		aEncoded = fmt.Sprintf("x0%s",x)
	}
	return aEncoded
}


func hexToASCII(hexString string) (string, error) {
    bytes, err := hex.DecodeString(hexString)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}


func convertHx(s string)string{
	x := s[1:] 
	asciiString, _ := hexToASCII(x)
	return asciiString
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
