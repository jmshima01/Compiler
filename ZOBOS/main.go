package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func closure(itemSet ItemSet, productionRules []ProductionRule){
	for{
		l := len(itemSet)
		for _,item := range itemSet{
			for _,p := range productionRules{
				if item.productionMarker == len(item.rhs){continue}
				if item.rhs[item.productionMarker] == p.lhs{
					newItem := makeItem(p,0)
					itemSet[newItem.toString()] = newItem
				}
			}	
		}
		if l==len(itemSet){ // unchanged
			break
		} 
	}
}


func goTo(itemSet ItemSet, productionRules []ProductionRule, symbol string)ItemSet{
	res := make(ItemSet)
	for _,item := range itemSet{
		if item.productionMarker == len(item.rhs){continue}

		if item.rhs[item.productionMarker] == symbol{
			newItem := Item{lhs: item.lhs,rhs: item.rhs,productionMarker: item.productionMarker+1}
			res[newItem.toString()] = newItem
		}
	}

	closure(res,productionRules)
	
	if len(res)==0{
		res = nil
	}
	return res
}

func main(){
	args := os.Args

	// read and clean cfg
	data := readLines(args[1])
	isComment,_ := regexp.Compile(`^(#).*`)
	grammar := make([]string,0)
	for _,line := range data{
		line = strings.TrimSpace(line)
		if isComment.MatchString(line) || len(line)==0{
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
	
	productionRules := makeProductionRules(grammar)
	// for _,v := range productionRules{
	// 	fmt.Println(v)
	// }

	// fmt.Println(terminals)
	// fmt.Println(nonTerminals)

	symbols := setUnion(terminals, nonTerminals)
	symbols.add("$")

	startState := getStartState(productionRules)

	// fmt.Println(startState)
	// fmt.Println(symbols)
	
	itemset := make(ItemSet)
	
	for _,p := range productionRules{
		if p.lhs == startState{
			item := makeItem(p,0)
			itemset[item.toString()] = item
		}
	}

	

	// for _, v:= range itemset{
	// 	fmt.Println(v)
	// }

	closure(itemset,productionRules)
	for _, v:= range itemset{
		fmt.Println(v)
	}
	newSet := goTo(itemset,productionRules,"C")
	fmt.Println("=============")
	
	for _, v:= range newSet{
		fmt.Println(v)
	}
	fmt.Println(newSet)
}