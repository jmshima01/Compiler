package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var count int = 0

func closure(itemSet ItemSet, productionRules []ProductionRule) {
	for {
		l := len(itemSet)
		for _, item := range itemSet {
			for _, p := range productionRules {
				if item.productionMarker == len(item.rhs) || item.rhs[0] == "lambda" {
					continue
				} // if * at end i.e S -> A B *
				if item.rhs[item.productionMarker] == p.lhs {
					newItem := makeItem(p, 0) // add a fresh start
					itemSet[newItem.toString()] = newItem
				}
			}
		}
		if l == len(itemSet) { // unchanged
			break
		}
	}
}

func goTo(itemSet ItemSet, productionRules []ProductionRule, symbol string) ItemSet {
	res := make(ItemSet)
	for _, item := range itemSet {
		if item.productionMarker == len(item.rhs) {
			continue
		}

		if item.rhs[item.productionMarker] == symbol {
			newItem := Item{lhs: item.lhs, rhs: item.rhs, productionMarker: item.productionMarker + 1}
			res[newItem.toString()] = newItem
		}
	}

	closure(res, productionRules)

	if len(res) == 0 {
		res = nil
	}
	return res
}

func makeCFSM(itemSet ItemSet, lookupCFSMID map[string]int, G CFSM, symbols set, id int, productionRules []ProductionRule, seen map[int]bool, itemsetLookup map[int]ItemSet) {
	startID := id
	numOfDone := 0
	numItems := 0

	// fmt.Println("NOW in State",id, itemSet.toString())
	for _, v := range itemSet {
		if len(v.rhs) == v.productionMarker {
			numOfDone++
			// fmt.Println("heyy",v.toString(),v.productionMarker)
		}
		numItems++
	}

	_, seenBefore := seen[id]

	if numOfDone == numItems { // no new item sets possible
		G[startID] = nil
		return
	}
	if seenBefore {
		return
	}

	G[startID] = make([]Pair, 0)

	for symbol := range symbols {
		newSet := goTo(itemSet, productionRules, symbol)
		// fmt.Println("=={",newSet.toString(),symbol,"}==")
		if newSet != nil {

			c, seenItBefore := lookupCFSMID[newSet.toString()]
			if !seenItBefore {
				count++
				lookupCFSMID[newSet.toString()] = count
				itemsetLookup[count] = newSet
				data := Pair{id: count, transitionSymbol: symbol}
				G[startID] = append(G[startID], data)

			} else {
				data := Pair{id: c, transitionSymbol: symbol}
				itemsetLookup[c] = newSet
				G[startID] = append(G[startID], data)

			}

		}
	}

	seen[startID] = true
	for _, v := range G[startID] {
		makeCFSM(itemsetLookup[v.id], lookupCFSMID, G, symbols, v.id, productionRules, seen, itemsetLookup)
		seen[v.id] = true
	}
}

func main() {
	args := os.Args

	// read and clean cfg
	data := readLines(args[1])
	isComment, _ := regexp.Compile(`^(#).*`)
	grammar := make([]string, 0)
	for _, line := range data {
		line = strings.TrimSpace(line)
		if isComment.MatchString(line) || len(line) == 0 {
			continue
		}
		grammar = append(grammar, line)
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

	// first and follow sets
	dervLambdaCache := make(set)
	for k := range symbols {
		dervLambdaCache[k] = derivesToLambda(k, productionRules)
	}
	firstCache := map[string]set{}
	for k := range symbols {
		firstCache[k] = first(k, productionRules, dervLambdaCache, make(set), make(set))
	}
	followCache := map[string]set{}
	for k := range nonTerminals {
		followCache[k] = follow(k, productionRules, dervLambdaCache, firstCache, make(set), make(set))
	}

	itemsetZero := make(ItemSet)

	for _, p := range productionRules {
		if p.lhs == startState {
			item := makeItem(p, 0)
			itemsetZero[item.toString()] = item
		}
	}

	closure(itemsetZero, productionRules)
	for _, v := range itemsetZero {
		fmt.Println(v)
	}

	fmt.Println("=============")
	lookupCFSMID := make(map[string]int)
	lookupCFSMID[itemsetZero.toString()] = 0

	cfsm := make(CFSM)
	seen := make(map[int]bool)

	itemsetLookup := make(map[int]ItemSet)
	itemsetLookup[0] = itemsetZero

	makeCFSM(itemsetZero, lookupCFSMID, cfsm, symbols, 0, productionRules, seen, itemsetLookup)

	fmt.Println(cfsm)
	fmt.Println("++++++++++++++++")
	for k, v := range itemsetLookup {
		fmt.Println(k, "->", v.toString())
	}

	SLRTable := make([][]string, len(cfsm))

	for i := range cfsm {
		SLRTable[i] = make([]string, len(symbols)-1)
		for j := range SLRTable[i] {
			SLRTable[i][j] = "__"
		}
	}

	columnHeader := make([]string, 0)
	for v := range terminals {
		columnHeader = append(columnHeader, v)
	}
	sort.Strings(columnHeader)
	columnHeader = append(columnHeader, "$")
	n := make([]string,0)
	for v := range nonTerminals{
		if v != startState {
			n = append(n, v)
		}
	}
	sort.Strings(n)
	columnHeader = append(columnHeader, n...)
	
	SLRRowLookup := make(map[string]int)
	for i, v := range columnHeader {
		SLRRowLookup[v] = i
	}
	fmt.Println(columnHeader)
	fmt.Println(SLRRowLookup)

	// for k, v := range followCache {
	// 	fmt.Println(k, "->", v)
	// }

	productionRuleLookup := make(map[string]int)
	for i,p := range productionRules{
		productionRuleLookup[p.toString()] = i+1
	}

	for id, itemLis := range cfsm {
		if len(itemLis) == 0{
			for _,v := range itemsetLookup[id]{
				if v.lhs == startState{
					SLRTable[id][0] = fmt.Sprintf("RD-%d",productionRuleLookup[ProductionRule{rhs: v.rhs,lhs: v.lhs}.toString()])
				} else{
					for f := range followCache[v.lhs]{
						SLRTable[id][SLRRowLookup[f]] = fmt.Sprintf("rd-%d",productionRuleLookup[ProductionRule{rhs: v.rhs,lhs: v.lhs}.toString()])
					}
				}
			}
		}
		for _, x := range itemLis {
			SLRTable[id][SLRRowLookup[x.transitionSymbol]] = fmt.Sprintf("sh-%d", x.id)

			for _, y := range itemsetLookup[x.id]{
				if y.rhs[0] == "lambda"{
					for f := range followCache[y.lhs]{
						SLRTable[id][SLRRowLookup[f]] = fmt.Sprintf("rd-l%d",productionRuleLookup[ProductionRule{rhs: y.rhs,lhs: y.lhs}.toString()])
					}
				}
			}
				// if len(y.rhs) == y.productionMarker || y.rhs[0] == "lambda" {

				// 	fmt.Println(y.lhs,"->",y.rhs,y.productionMarker)

				// 	f := followCache[y.lhs]
				// 	for z := range f{
				// 		// if SLRTable[id][SLRRowLookup[z]] != "__"{
				// 		// 	println("CONFLICT",SLRTable[id][SLRRowLookup[x.transitionSymbol]])
				// 		// 	os.Exit(2)
				// 		// }
				// 		SLRTable[id][SLRRowLookup[z]] = fmt.Sprintf("rd-p")
				// 	}
				// }
				
			// }
		}
		// fmt.Println(id,itemsetLookup[x.id].toString(),x.transitionSymbol)
	}

	fmt.Println(columnHeader)
	for _, v := range SLRTable {
		fmt.Println(v)
	}

	// let us parse...
	// tokenStream := []string{"y","x","x","h","x","p","y","$"}
	// S := make(stack,0)
	// Q := make(queue,0)

}
