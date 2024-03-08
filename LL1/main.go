package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	
)

// ============== Helpers =============

func isNonTerminal(s string) bool {
    for _, r := range s {
        if unicode.IsUpper(r) && unicode.IsLetter(r) {
            return true
        }
    }
    return false
}



func findStartState(P []production_rule) string{
	seen := map[string]bool{}
	for _,v := range P{
		for _,x := range v.RHS{
			for _,y := range x{
				if string(y) == "$"{
					return v.LHS
				}
				if isNonTerminal(x){
					seen[x] = true	
				}
			}
		}
	}
	for _,v := range P{{
			_,ok := seen[v.LHS]
			if !ok{
				return v.LHS
			} 
		}
	}
	fmt.Println("No Start State in grammar!")
	return ""
}



func main() {

	args := os.Args

	cfg_filepath := args[1]
	data,err := os.ReadFile(cfg_filepath)

	if err!=nil{
		fmt.Println("Error reading file: ",err);
	}

	lines := strings.Split(string(data),"\n");

	// clean grammar input whitespace
	for i,v := range lines{
		lines[i] = strings.TrimSpace(v)
		fmt.Println(lines[i])
	}
	fmt.Println()
	
	P := make([]production_rule,0)
	curr_rhs := ""
	
	TERMINALS := map[string]bool{} // set in golang
	NON_TERMINALS := map[string]bool{}
	GRAMMAR_SYMBOLS := map[string]bool{}

	// get terminals and non-terminals
	for _,v := range lines{
		line := strings.Split(v," ")
		for _,t := range line{
			if isNonTerminal(t) && t != "|" && t!= "->"{
				NON_TERMINALS[t] = true
			} else if !isNonTerminal(t) && t != "|" && t!= "->" && t != "lambda" && t!="$"{
				TERMINALS[t] = true
			} 
			if !isNonTerminal(t) && t != "|" && t!= "->"{
				GRAMMAR_SYMBOLS[t] = true
			}

		}	

	}

	// fmt.Println(NON_TERMINALS)
	// fmt.Println(TERMINALS)

	for _,v := range lines{
		s := strings.Split(v," -> ")
		if len(s) == 1{
			bars := strings.Split(v[2:]," | ")
			for _,b := range bars{
				P = append(P,production_rule{LHS:curr_rhs,RHS:strings.Split(b," ")})
			}
		} else{
			bars := strings.Split(s[1]," | ")
			for _,b := range bars{
				P = append(P,production_rule{LHS:s[0],RHS:strings.Split(b," ")})
			}
			curr_rhs = s[0]
		}
	}
	
	cfg := CFG{sigma:TERMINALS, S:findStartState(P), P:P, N:NON_TERMINALS, derivestoLambda: nil, symbols:GRAMMAR_SYMBOLS, nullable:nil}
	cfg.nullable = make(map[string]bool)

	for _,p := range P{
		_,ok := cfg.nullable[p.LHS]
		if !ok{
			cfg.nullable[p.LHS]=false
		}
		if p.RHS[0] == "lambda"{
			cfg.nullable[p.LHS] = true
		}
	}
	
	cfg.derivestoLambda = cfg.genLambda()

	cfg.printCFG()
	// fmt.Println("Rule",derivesToLambda("Rule",P))
	
}