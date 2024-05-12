package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func lex(s string) []token {
	tokens := make([]token, 0)
	skip := false
	for i, v := range s {
		if skip {
			skip = false
			continue
		}
		if string(v) == `\` {
			if string(s[i+1]) == "s" {
				t := token{value: "SP", tokenType: "char"}
				skip = true
				tokens = append(tokens, t)
				continue
			} else {
				t := token{value: string(v), tokenType: "char"}
				skip = true
				tokens = append(tokens, t)
				continue

			}
		}
		switch string(v) {
		case "(":
			t := token{value: string(v), tokenType: "open"}
			tokens = append(tokens, t)
		case ")":
			t := token{value: string(v), tokenType: "close"}
			tokens = append(tokens, t)
		case "|":
			t := token{value: string(v), tokenType: "pipe"}
			tokens = append(tokens, t)
		case "*":
			t := token{value: string(v), tokenType: "kleene"}
			tokens = append(tokens, t)
		case "+":
			t := token{value: string(v), tokenType: "plus"}
			tokens = append(tokens, t)
		case "-":
			t := token{value: string(v), tokenType: "dash"}
			tokens = append(tokens, t)
		case ".":
			t := token{value: string(v), tokenType: "dot"}
			tokens = append(tokens, t)
		case "lambda":
			t := token{value: string(v), tokenType: "lambda"}
			tokens = append(tokens, t)
		default:
			t := token{value: string(v), tokenType: "char"}
			tokens = append(tokens, t)

		}
	}
	return tokens
}

func main() {
	// args := os.Args
	grammar := readLines("llre.cfg")
	fmt.Println(grammar)

	scan := readLines("scan.lut")
	// fmt.Println(scan)

	toTokenize := []string{}
	for i, line := range scan {
		vals := strings.Fields(line)
		if i == 0 {
			continue
		}

		fmt.Println(vals)
		toTokenize = append(toTokenize, vals[0])
	}
	tokenStream := []token{}
	for _, v := range toTokenize {
		tokenStream = lex(v)
		break
	}

	

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
	fmt.Println(startState)

	dervLambdaCache := make(set)
	for k, _ := range symbols {
		if isNonTerminal(k) {
			fmt.Println("derv->", k, derivesToLambda(k, productionRules))
		}
		dervLambdaCache[k] = derivesToLambda(k, productionRules)
	}
	
	firstCache := map[string]set{}
	for k, _ := range symbols {
		if isNonTerminal(k) {
			fmt.Println("first->", k, first(k, productionRules, dervLambdaCache, make(set), make(set)).getValues())
		}
		firstCache[k] = first(k, productionRules, dervLambdaCache, make(set), make(set))

	}
	fmt.Println("------------------")
	followCache := map[string]set{}
	for k, _ := range nonTerminals {
		// fmt.Println("doing follow of...",k)
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

	// fmt.Println(rowValues)
	// fmt.Println(columnValues)
	// fmt.Println(columnLookup)
	// fmt.Println(rowLookup)
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
	// fmt.Println()
	fmt.Println(columnValues)
	for _, v := range LLTable {
		fmt.Println(v)
	}

	S := make(stack, 0)
	Q := make(queue, 0)
	S = append(S, startState)
	for _, v := range tokenStream {
		Q.push(v.tokenType)
	}

	Q.push("$")
	fmt.Println(Q)
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

		fmt.Println("S:",S)
		fmt.Println("Q:",Q)
		s := S.peek()
		q := Q.peek()

		if s == "<*>" {

			S.pop()
			current = current.parent 
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
			if s == q {

				term := makeNode(q, current, uniqueID)
				addChild(current, term)
				uniqueID++

				S.pop()
				Q.popfront()
				current.debug()

			} else {
				fmt.Println("syntax error: s!=q", s, q)
				os.Exit(2)
			}

			continue
		}

		nextRule, found := ruleLookup[LLTable[rowLookup[s]][columnLookup[q]]]
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
		current.debug()

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
	writeToFile("parsetree.txt", toGraphiz)
}
