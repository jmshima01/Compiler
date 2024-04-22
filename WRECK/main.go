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

					firstSet = first(v, P, dervLambda, firstSet, seen)

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
	for _, v := range seq {
		if !dervLambda[v] && isNonTerminal(v) {
			return false
		}
	}
	return true
}

func follow(N string, P []ProductionRule, dervLambda set, firsts map[string]set, followSet set, seen set) set {
	_, ok := seen[N]
	if ok {
		return followSet
	}
	seen.add(N)
	needFollows := make(set)
	for _, p := range P {
		foundN := false
		needFollow := false
		last := p
		for i, v := range p.rhs {
			if v == N {
				foundN = true
				if i == len(p.rhs)-1 {
					needFollow = true
				}
				continue
			}
			if foundN {
				followSet = setUnion(followSet, firsts[v])
				// fmt.Println(v,dervLambda[v])
				if !dervLambda[v] {
					needFollow = false
					break
				}
			}
		}
		if needFollow {
			needFollows.add(last.lhs)
		}
	}

	for s, _ := range needFollows {
		followSet = follow(s, P, dervLambda, firsts, followSet, seen)
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

// ========== LL(1) =================
func main() {
	args := os.Args

	grammar := readLines(args[1])

	// trim whitespace
	for i, v := range grammar {
		grammar[i] = strings.TrimSpace(v)
	}

	nonTerminals := make(set)
	terminals := make(set)

	for _, line := range grammar {
		v := strings.Split(line, " ")
		for _, s := range v {
			if isNonTerminal(s) {
				nonTerminals.add(s)
			} else if isTerminal(s) {
				terminals.add(s)
			}
		}
	}
	fmt.Println(terminals)
	fmt.Println(nonTerminals)

	symbols := setUnion(terminals, nonTerminals)
	symbols.add("$")

	productionRules := makeProductionRules(grammar)
	fmt.Println("++++++++++++")

	for _, p := range productionRules {
		fmt.Println(p)
	}
	fmt.Println("++++++++++++")
	startState := getStartState(productionRules)
	fmt.Println(startState)
	fmt.Println()

	dervLambdaCache := make(set)
	for k, _ := range symbols {
		// fmt.Println("workin on",k)
		fmt.Println(k, "derv->", derivesToLambda(k, productionRules))
		dervLambdaCache[k] = derivesToLambda(k, productionRules)
	}
	fmt.Println()
	firstCache := map[string]set{}
	for k, _ := range symbols {
		firstCache[k] = first(k, productionRules, dervLambdaCache, make(set), make(set))
		if isNonTerminal(k) {
			fmt.Println("first->", k, first(k, productionRules, dervLambdaCache, make(set), make(set)).getValues())
		}

	}

	fmt.Println()

	followCache := map[string]set{}
	for k, _ := range nonTerminals {
		// if(k=="RHS"){
		fmt.Println("follow->", k, follow(k, productionRules, dervLambdaCache, firstCache, make(set), make(set)).getValues())
		followCache[k] = follow(k, productionRules, dervLambdaCache, firstCache, make(set), make(set))

	}
	fmt.Println()
	for _, p := range productionRules {
		fmt.Println("predict->", p, predict(p, dervLambdaCache, firstCache, followCache).getValues())
	}

	ruleLookup := map[int]ProductionRule{}
	for i, p := range productionRules {
		ruleLookup[i+1] = p
	}

	fmt.Println(ruleLookup)

	// LLTable := make([][]int,0)
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

	fmt.Println(rowValues)
	fmt.Println(columnValues)
	fmt.Println(columnLookup)
	fmt.Println(rowLookup)
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
	fmt.Println()
	fmt.Println(columnValues)
	for _, v := range LLTable {
		fmt.Println(v)
	}
	fmt.Println("=======================parsing============")
	test4 := "bghm$"
	// tokens_stream := getTokens()
	S := make(stack, 0)
	Q := make(queue, 0)
	S = append(S, startState)
	for _, v := range test4 {
		Q.push(string(v))
	}

	// fmt.Println(LLTable)
	

	root := makeNode("ROOT", nil)
	current := root
	current.debug()

	for {
		if S.isEmpty(){
			if !Q.isEmpty(){
				fmt.Println("syntax error:",Q)
				os.Exit(2)
			}

			break
		}

		fmt.Println("S:", S)
		fmt.Println("Q:", Q)
		
		s := S.peek()
		q := Q.peek()
		
		if s == "*"{
			S.pop()
			current = current.parent
			continue
		}

		if s == "lambda"{
			S.pop()
			lambNode:=makeNode("lambda",current)
			addChild(current,lambNode)
			current.debug()
			continue
		}
		
		if isTerminal(s) || s=="$"{
			if s==q{
				term:=makeNode(s,current)
				addChild(current,term)
				S.pop()
				Q.popfront()
				current.debug()
	
			} else{
				fmt.Println("syntax error: s!=q",s,q)
				os.Exit(2)
			}
			continue
		}

		nextRule,found := ruleLookup[LLTable[rowLookup[s]][columnLookup[q]]]
		if !found{
			fmt.Println("Parsing Error: (No such token in LL table or associated rule)",q)
			os.Exit(1)
		}
		
		fmt.Println("fetching rule...",nextRule)
		top := S.pop()
		newNode := makeNode(top, current)
		addChild(current, newNode)
		current.debug()
		
		current = newNode
		
		// add rule in reverse to stack...
		S = append(S, "*") // end of production
		for i := len(nextRule.rhs) - 1; i >= 0; i-- {
			S = append(S, nextRule.rhs[i])
		}
	}

	current.debug()

	fmt.Println("============")
	printTree(current)

}
