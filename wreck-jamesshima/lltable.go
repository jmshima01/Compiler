package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

// =============== Helpers ==========
func isNonTerminal(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func isTerminal(s string) bool {
	if s == "lambda" || s == "->" || s == "$" || s == "|" {
		return false
	}
	for _, r := range s {
		if unicode.IsUpper(r) || string(r) == " " {
			return false
		}
	}
	return true
}

func makeProductionRules(cfg []string) []ProductionRule {
	productionRules := make([]ProductionRule, 0)
	currLHS := ""
	currRHS := make([]string, 0)
	for _, line := range cfg {
		v := strings.Fields(line)
		currRHS = nil
		for i, s := range v {
			if string(s) == "->" || (string(s) == "|" && i == 0) {
				continue
			}
			if i == 0 && string(s) != "|" {
				currLHS = s
			} else if string(s) == "|" {
				productionRules = append(productionRules, ProductionRule{lhs: currLHS, rhs: currRHS})
				currRHS = nil
			} else {
				currRHS = append(currRHS, s)
			}
		}
		productionRules = append(productionRules, ProductionRule{lhs: currLHS, rhs: currRHS})
	}
	return productionRules
}

func getStartState(P []ProductionRule) string {
	seen := map[string]bool{}
	for _, v := range P {
		for _, x := range v.rhs {
			for _, y := range x {
				if string(y) == "$" {
					return v.lhs
				}
				if isNonTerminal(x) {
					seen[x] = true
				}
			}
		}
	}
	for _, v := range P {
		{
			_, ok := seen[v.lhs]
			if !ok {
				return v.lhs
			}
		}
	}
	fmt.Println("No Start State in grammar!")
	return ""
}

func containsTerminal(rhs []string) bool {
	for _, v := range rhs {
		if isTerminal(v) {
			return true
		}
	}
	return false
}

func hasLambdaRule(N string, P []ProductionRule) bool {
	for _, p := range P {
		if p.lhs == N {
			if p.rhs[0] == "lambda" {
				return true
			}
		}
	}
	return false
}

func derivesToLambda(N string, P []ProductionRule) bool {
	if !isNonTerminal(N) {
		return false
	}

	for _, p := range P {
		if p.lhs == N {
			if containsTerminal(p.rhs) {
				continue
			} else if p.rhs[0] == "lambda" {
				return true
			} else {
				res := true
				for _, v := range p.rhs {
					if !hasLambdaRule(v, P) {
						res = false
					}
				}
				if !res {
					continue
				}
				return true
			}
		}
	}
	return false
}

func first(N string, P []ProductionRule, dervLambda set, firstSet set, seen set) set {
	_, ok := seen[N]
	if ok {
		return firstSet
	}
	if N == "lambda" {
		return nil
	}
	if isTerminal(N) || N == "$" {
		s := make(set)
		s.add(N)
		return s
	}
	seen.add(N)
	for _, p := range P {
		if p.lhs == N {
			for i, v := range p.rhs {
				if i == 0 && (isTerminal(v)) {
					firstSet.add(v)
					break
				} else {
					if v == "lambda" {
						break
					}

					firstSet = setUnion(firstSet,first(v, P, dervLambda, firstSet, seen))

					if !dervLambda[v] {
						break
					}
				}
			}
		}
	}
	return firstSet
}

func needUnionFollow(dervLambda set, seq []string) bool {
	c:=0
	for _, v := range seq {
		if dervLambda[v] && isNonTerminal(v) {
			c++
		}
	}
	return c==len(seq)
}

func follow(N string, P []ProductionRule, dervLambda set, firsts map[string]set, followSet set, seen set) set{
	
	_, ok := seen[N]
	if ok {
		return followSet
	}
	seen.add(N)
	for _, p := range P {	
		for i:=0; i<len(p.rhs); i++{
			if p.rhs[i] == N {
				if i+1 != len(p.rhs){
					for _,v := range p.rhs[i+1:]{
						followSet = setUnion(followSet, firsts[v])
					}
				}
				if i+1 == len(p.rhs) || needUnionFollow(dervLambda,p.rhs[i+1:]) {
					// fmt.Println("following ",p.lhs)
					followSet = setUnion(followSet,follow(p.lhs,P,dervLambda,firsts,followSet,seen))
				}
					
			}
		}
	}
	return followSet

}


func predict(p ProductionRule, dervLambda set, firsts map[string]set, follows map[string]set) set {
	predictSet := make(set)
	flag := true
	for _, v := range p.rhs {
		if v == "lambda" {
			flag = false
			predictSet = follows[p.lhs]
			break
		}

		predictSet = setUnion(predictSet, firsts[v])

		if !dervLambda[v] {
			flag = false
			break
		}
	}
	if flag {
		predictSet = setUnion(predictSet, follows[p.lhs])
	}
	return predictSet
}

// ========== LL(1) Table driven Parser =================
func makeLLTable(grammar []string) ([][]int,string,map[int]ProductionRule,map[string]int,map[string]int){ // Produces AST given a .jack source file
	
	// ::Table gen::
	for i, v := range grammar {
		grammar[i] = strings.TrimSpace(v)
	}

	nonTerminals := make(set)
	terminals := make(set)

	for _, line := range grammar {
		v := strings.Fields(line)
		for _, s := range v {
			if isNonTerminal(s) {
				nonTerminals.add(s)
			} else if isTerminal(s) {
				terminals.add(s)
			}
		}
	}

	symbols := setUnion(terminals, nonTerminals)
	symbols.add("$")

	productionRules := makeProductionRules(grammar)

	startState := getStartState(productionRules)
	// fmt.Println(startState)

	dervLambdaCache := make(set)
	for k, _ := range symbols {
		// if isNonTerminal(k) {
		// 	fmt.Println("derv->", k, derivesToLambda(k, productionRules))
		// }
		dervLambdaCache[k] = derivesToLambda(k, productionRules)
	}

	firstCache := map[string]set{}
	for k, _ := range symbols {
		// if isNonTerminal(k) {
		// 	fmt.Println("first->", k, first(k, productionRules, dervLambdaCache, make(set), make(set)).getValues())
		// }
		firstCache[k] = first(k, productionRules, dervLambdaCache, make(set), make(set))

	}
	fmt.Println()
	followCache := map[string]set{}
	for k, _ := range nonTerminals {
		// fmt.Println("doing follow of...",k)
		// fmt.Println("follow->", k, follow(k, productionRules, dervLambdaCache, firstCache, make(set), make(set)).getValues())

		followCache[k] = follow(k, productionRules, dervLambdaCache, firstCache, make(set), make(set))

	}
	// fmt.Println()
	// for _, p := range productionRules {
	// 	fmt.Println("predict->", p, predict(p, dervLambdaCache, firstCache, followCache).getValues())
	// }

	ruleLookup := map[int]ProductionRule{}
	for i, p := range productionRules {
		ruleLookup[i+1] = p
	}

	columnValues := terminals.getValues()
	sort.Strings(columnValues)
	columnValues = append(columnValues, "$")

	columnLookup := map[string]int{}
	rowLookup := map[string]int{}

	rowValues := make([]string, 0)
	temp := make(set)
	for _, p := range productionRules {
		_, ok := temp[p.lhs]
		if !ok {
			rowValues = append(rowValues, p.lhs)
		}
		temp.add(p.lhs)
	}

	for i, v := range rowValues {
		rowLookup[v] = i
	}
	for i, v := range columnValues {
		columnLookup[v] = i
	}

	LLTable := make([][]int, len(rowLookup))
	for _, i := range rowLookup {
		LLTable[i] = make([]int, len(columnLookup))
	}

	for i, p := range productionRules {
		t := predict(p, dervLambdaCache, firstCache, followCache)
		for v, _ := range t {
			if LLTable[rowLookup[p.lhs]][columnLookup[v]] != 0 {
				fmt.Println("Grammar is not LL1 ! conflict", p.lhs, columnLookup[v], i+1, LLTable[rowLookup[p.lhs]][columnLookup[v]])
				os.Exit(1)
			}
			LLTable[rowLookup[p.lhs]][columnLookup[v]] = i + 1
		}

	}

	// fmt.Println(columnValues)
	// for _, v := range LLTable {
	// 	fmt.Println(v)
	// }
	return LLTable,startState,ruleLookup,rowLookup,columnLookup
}


func makeAST(tokenStream []token, LLTable [][]int, startState string, ruleLookup map[int]ProductionRule,rowLookup map[string]int,columnLookup map[string]int) *Node{


	S := make(stack, 0)
	Q := make(queue, 0)
	S = append(S, startState)
	
	for _, tok := range tokenStream {
		Q.push(tok)
	}

	Q.push(token{value: "$", tokenType: "$"})

	// fmt.Println(Q)
	root := makeNode("ROOT", nil, 0)
	current := root
	uniqueID := 1
	for {
		if S.isEmpty() {
			if !Q.isEmpty() {
				fmt.Println("syntax error:", Q)
				os.Exit(2)
			}

			break
		}

		// fmt.Println("S:", S)
		// fmt.Println("Q:", Q)
		s := S.peek()
		q := Q.peek()

		if s == "<*>" {

			S.pop()
			// SDT!
			switch current.data {
			
			case "CHARRNG":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {
					current.data = current.children[1].data
					current.id = current.children[1].id
					current.children = nil
					current = current.parent
				}
			case "ATOMMOD":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
				}
			case "SEQLIST":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else if len(current.children) == 1 {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = current.children[0].children
					current = current.parent
					
				} else {
					current = current.parent
					newChildren := make([]*Node, 0)
					for _, v := range current.children {
						if v.data != "SEQLIST" {
							newChildren = append(newChildren, v)
						} else {
							for _,x := range v.children{
								x.parent = current
							}
							newChildren = append(newChildren, v.children...)
						}
					}
					current.children = newChildren
				}
			case "ALTLIST":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {
					current = current.parent
					newChildren := make([]*Node, 0)
					for _, v := range current.children {
						if v.data != "ALTLIST" {
							newChildren = append(newChildren, v)
						} else {
							for _,x := range v.children{
								x.parent = current
							}
							newChildren = append(newChildren, v.children...)
						}
					}
					current.children = newChildren
				}
			case "NUCLEUS":
				if len(current.children) == 1 {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = current.children[0].children
					current = current.parent

				} else if len(current.children) == 3 {
					current.data = current.children[1].data
					current.id = current.children[1].id
					current.children = current.children[1].children
					current = current.parent

				}else if len(current.children) == 2{
					current.data = "range"
					current = current.parent
				
				}else {
					current = current.parent
				}
			case "ATOM":
				if len(current.children) == 2 {
					current.data = current.children[1].data
					current.id = current.children[1].id
					current.children = current.children[:len(current.children)-1]
					current = current.parent

				} else if len(current.children) == 1 {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = current.children[0].children
					current = current.parent

				} else {
					current = current.parent
				}
			case "ALT":
				if len(current.children) == 1{
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = current.children[0].children
					current = current.parent
				} else{
					newChildren := make([]*Node, 0)
					for _, v := range current.children {
						if v.data != "pipe" {
							newChildren = append(newChildren, v)
						} else {
							for _,x := range v.children{
								x.parent = current
							}
							newChildren = append(newChildren, v.children...)
						}
					}
					current.children = newChildren
					for _,v := range current.children{
						if v.data == "SEQ" && len(v.children)==1{
							v.data = v.children[0].data
							v.id = v.children[0].id
							v.children = v.children[0].children
						}
					}
					current = current.parent
				}
			case "RE":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = current.children[0].children
				current = current.parent
			
			default:
				current = current.parent
			}
			// current = current.parent //for no SDT
			continue
		}

		if s == "lambda" {
			S.pop()
			lambNode := makeNode("lambda", current, uniqueID)
			addChild(current, lambNode)
			uniqueID++
			continue
		}

		if isTerminal(s) || s == "$" {

			if s == q.tokenType {
				term := makeNode(q.tokenType, current, uniqueID)
				if q.tokenType == "char" {
					term.data = q.value
				}

				addChild(current, term)
				uniqueID++

				S.pop()
				Q.popfront()
				// current.debug()

			} else {
				fmt.Println("syntax error: s!=q", s, q)
				os.Exit(2)
			}

			continue
		}

		nextRule, found := ruleLookup[LLTable[rowLookup[s]][columnLookup[q.tokenType]]]
		if !found {

			fmt.Println("Parsing Error: (No such token in LL table or associated rule)")
			fmt.Println("-----")
			fmt.Println(s, q)
			fmt.Println(S)
			fmt.Println(Q)
			os.Exit(2)
		}

		// fmt.Println("fetching rule...", nextRule)
		top := S.pop()
		newNode := makeNode(top, current, uniqueID)

		addChild(current, newNode)
		// current.debug()

		current = newNode
		uniqueID++
		// add rule in reverse to stack...
		S = append(S, "<*>") // end of production symbol designation
		for i := len(nextRule.rhs) - 1; i >= 0; i-- {
			S = append(S, nextRule.rhs[i])
		}
	}
	S = nil
	Q = nil

	// printTree(current)
	ast := current.children[0]

	nodeInfo := ""
	nodeInfo = *(genNodeInfo(ast, &nodeInfo))

	edgeInfo := ""
	edgeInfo = *(genEdgeInfo(ast, &edgeInfo))

	toGraphiz := nodeInfo + "\n" + edgeInfo
	writeToFile("parsetree.txt", toGraphiz) // parsetree!
	
	return ast

}