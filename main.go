package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

// ============ Types ==============
type ParseTree struct{
	root *Node
}

type Node struct{
	data string;
	parent *Node;
	
	children []*Node;
}

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



type ProductionRule struct{
	lhs string;
	rhs []string;
};

func (P ProductionRule) toString() string{
	s:= ""
	s+=P.lhs
	s+=" -> "
	for _,v:= range P.rhs{
		s+= v + " "
	}
	return s[:len(s)-1] 
}

// for first, follow, & predict sets 
type set map[string]bool;

type deque []string;

func popfront(d deque)(deque,string){
	if len(d)==0{
		return d,""
	}
	
	front :=d[0] 
	d = d[1:]
	return d,front
}

func pop(d deque)(deque,string){
	if len(d)==0{
		return d,""
	}
	
	back :=d[len(d)-1] 
	d = d[:len(d)-1]
	return d,back
}


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

type Grammar struct{
	raw []string;
	P []ProductionRule;
	// Caches for memoization (maps of sets)
	dervToLamb map[string]set;
	firstSets map[string]set;
	followSets map[string]set;
	predictSets map[string]set; // N -> predict(N) 
}

// =============== Helpers ==========
func isNonTerminal(s string)bool{
    for _, r := range s {
        if unicode.IsUpper(r) && unicode.IsLetter(r) {
            return true
        }
    }
    return false
}

func isTerminal(s string)bool{
	if s == "lambda" || s == "->" || s =="$" || s=="|"{return false}
	for _, r := range s {
        if unicode.IsUpper(r) || string(r) == " "{
            return false
        }
    }
    return true
}



func makeProductionRules(cfg []string) []ProductionRule{
	productionRules := make([]ProductionRule,0)
	currLHS := ""
	currRHS := []string{}
	for _,line := range cfg{
		v := strings.Split(line, " ")
		currRHS = nil
		for i,s:= range v{
			if string(s) == "->" || (string(s) == "|" && i==0){
				continue
			}
			if i == 0 && string(s) != "|"{
				currLHS = s
			} else if string(s) == "|"{
				productionRules = append(productionRules, ProductionRule{lhs: currLHS, rhs: currRHS})
				currRHS = nil
			} else{
				currRHS = append(currRHS, s)
			}
		}
		productionRules = append(productionRules, ProductionRule{lhs: currLHS, rhs: currRHS})
	}
	return productionRules
}

func getStartState(P []ProductionRule)string{
	seen := map[string]bool{}
	for _,v := range P{
		for _,x := range v.rhs{
			for _,y := range x{
				if string(y) == "$"{
					return v.lhs
				}
				if isNonTerminal(x){
					seen[x] = true	
				}
			}
		}
	}
	for _,v := range P{{
			_,ok := seen[v.lhs]
			if !ok{
				return v.lhs
			} 
		}
	}
	fmt.Println("No Start State in grammar!")
	return ""
}


func containsTerminal(rhs []string)bool{
	for _,v := range rhs{
		if isTerminal(v){return true}
	}
	return false
}

func hasLambdaRule(N string, P []ProductionRule)bool{
	for _,p := range P{
		if p.lhs == N{
			if p.rhs[0] == "lambda"{
				return true
			}
		}
	}
	return false
}

func derivesToLambda(N string, P []ProductionRule)bool{
	if !isNonTerminal(N){return false}
	
	for _,p := range P{
		if p.lhs == N{
			if containsTerminal(p.rhs){
				continue
			} else if p.rhs[0] == "lambda"{
				return true
			} else{
				res := true
				for _,v:=range p.rhs{
					if !hasLambdaRule(v,P){
						res=false
					}
				}
				if !res{
					continue
				}
				return true
			}
		}
	}
	return false
}


func first(N string, P []ProductionRule, dervLambda set, firstSet set, seen set) set{
	_,ok := seen[N]
	if ok{
		return firstSet
	}
	if N == "lambda"{
		return nil
	}
	if isTerminal(N) || N == "$"{
		s := make(set)
		s.add(N)
		return s
	}
	seen.add(N)
	for _,p := range P{
		if p.lhs == N{
			for i,v := range p.rhs{
				if i==0 && (isTerminal(v)){
					firstSet.add(v)
					break
				} else{
					if v == "lambda" {break}
					
					firstSet = first(v,P,dervLambda,firstSet,seen)
					
					if !dervLambda[v]{
						break
					}
				}
			}
		}
	}
	return firstSet	
}

func needUnionFollow(dervLambda set, seq []string)bool{
	for _,v := range seq{
		if !dervLambda[v] && isNonTerminal(v){
			return false
		}
	}
	return true    
}


func follow(N string, P []ProductionRule, dervLambda set, firsts map[string]set, followSet set, seen set) set{
	_,ok := seen[N]
	if ok{
		return followSet
	}
	seen.add(N)
	needFollows := make(set)
	for _,p := range P{
		foundN := false
		needFollow := false
		last := p
		for i,v := range p.rhs{
			if v == N{
				foundN = true
				if i == len(p.rhs)-1{
					needFollow = true
				}
				continue
			}
			if foundN{	
				followSet = setUnion(followSet,firsts[v])
				// fmt.Println(v,dervLambda[v])
				if !dervLambda[v]{
					needFollow = false
					break	
				}	
			}
		}
		if needFollow{
			needFollows.add(last.lhs)
		}	
	}

	for s,_ := range needFollows{	
		followSet = follow(s,P,dervLambda,firsts,followSet,seen)
	}

	return followSet
}

func predict(p ProductionRule, dervLambda set, firsts map[string]set, follows map[string]set)set{
	predictSet := make(set)
	flag := true
	for _,v := range p.rhs{
		if v == "lambda"{
			flag = false
			predictSet = follows[p.lhs]
			break
		}
		
		predictSet = setUnion(predictSet,firsts[v])
		
		if !dervLambda[v]{
			flag = false
			break
		}
	}
	if flag{
		predictSet = setUnion(predictSet,follows[p.lhs])
	}
	return predictSet
}


func makeLL1Table(G Grammar){


}



func parseLL(){

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

// ========== LL(1) =================
func main(){
	args := os.Args
	fmt.Println(args)
	
	grammar := readLines(args[1])
	
	// trim whitespace
	for i,v := range grammar{
		grammar[i] = strings.TrimSpace(v)
	}

	nonTerminals := make(set)
	terminals := make(set)
	
	for _,line := range grammar{
		v := strings.Split(line, " ")
		for _,s := range v{
			if isNonTerminal(s){
				nonTerminals.add(s)
			} else if isTerminal(s){
				terminals.add(s)
			}
		}
	}
	fmt.Println(terminals)
	fmt.Println(nonTerminals)

	symbols := setUnion(terminals,nonTerminals)
	symbols.add("$")

	productionRules := makeProductionRules(grammar)
	fmt.Println(productionRules)
	startState := getStartState(productionRules)
	fmt.Println(startState)
	fmt.Println()

	dervLambdaCache := make(set)
	for k,_ := range symbols{
		fmt.Println(k,"derv->",derivesToLambda(k,productionRules))
		dervLambdaCache[k] = derivesToLambda(k,productionRules)
	}
	fmt.Println()
	firstCache := map[string]set{}
	for k,_ := range symbols{
		firstCache[k] = first(k,productionRules,dervLambdaCache,make(set),make(set))
		if isNonTerminal(k){
			fmt.Println("first->",k,first(k,productionRules,dervLambdaCache,make(set),make(set)).getValues())
		}
		
	}

	fmt.Println()

	followCache := map[string]set{}
	for k,_ := range nonTerminals{
		// if(k=="RHS"){
		fmt.Println("follow->",k,follow(k,productionRules,dervLambdaCache,firstCache,make(set),make(set)).getValues())
		followCache[k]=follow(k,productionRules,dervLambdaCache,firstCache,make(set),make(set))
		
	}
	fmt.Println()
	for _,p := range productionRules{
		fmt.Println("predict->",p,predict(p,dervLambdaCache,firstCache,followCache).getValues())
	}

	ruleLookup := map[int]ProductionRule{}
	for i,p := range productionRules{
		ruleLookup[i+1] = p
	}

	fmt.Println(ruleLookup)

	// LLTable := make([][]int,0)
	columnValues:= terminals.getValues()
	sort.Strings(columnValues)
	columnValues = append(columnValues, "$")

	columnLookup := map[string]int{}
	rowLookup := map[string]int{}

	rowValues := make([]string,0)
	temp := make(set)
	for _,p := range productionRules{
		_,ok := temp[p.lhs]
		if !ok{
			rowValues = append(rowValues, p.lhs)
		}
		temp.add(p.lhs)
	}

	for i,v:= range rowValues{
		rowLookup[v]=i
	}
	for i,v:= range columnValues{
		columnLookup[v]=i
	}

	fmt.Println(rowValues)
	fmt.Println(columnValues)
	fmt.Println(columnLookup)
	fmt.Println(rowLookup)
	LLTable := make([][]int,len(rowLookup))
	for _,i := range rowLookup{
		LLTable[i] = make([]int, len(columnLookup))
	}
	

	for i,p := range productionRules{
		t := predict(p,dervLambdaCache,firstCache,followCache)
		for v,_ :=range t{
			LLTable[rowLookup[p.lhs]][columnLookup[v]] = i+1	
		}

		
	}
	test4:= "bghm$"
	S := make(deque,0)
	Q := make(deque,0)
	S = append(S, startState)
	for _,v := range test4{
		Q = append(Q, string(v))
	} 
	fmt.Println(LLTable)
	fmt.Println(S)
	fmt.Println(Q)
	

	root := makeNode("ROOT",nil)
	current := root
	current.debug()
	
	nextRule:= ruleLookup[ LLTable[ rowLookup[S[len(S)-1]] ][ columnLookup[Q[0]] ] ]
	fmt.Println(nextRule)
	S,top := pop(S)
	
	newNode := makeNode(top,current)
	addChild(current,newNode)
	current.debug()
	current = newNode

	current.debug()

	S = append(S, "*")
	for i:= len(nextRule.rhs)-1; i>=0; i--{
		S = append(S, nextRule.rhs[i])
	}
	fmt.Println(S)
	


	
	
}