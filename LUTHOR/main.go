package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
	// "regexp"
)

// ============== TYPES ===============
type production_rule struct{
	LHS string;
	RHS []string;

};


// ============== Helpers =============

func isNonTerminal(s string) bool {
    for _, r := range s {
        if unicode.IsUpper(r) && unicode.IsLetter(r) {
            return true
        }
    }
    return false
}

func containsTerminal(rhs []string) bool{
	for _,v := range rhs{
		if !isNonTerminal(v) && v!= "lambda" {
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
// =============== LAMBDA/FIRST/FOLLOW ================

func derivesToLambda(nonterm string, P []production_rule) bool {
	nullable := make(map[string]bool)
	for _,p := range P{
		nullable[p.LHS] = false
	}
	for _,p := range P{
		if containsTerminal(p.RHS){
			continue;
		} 
		if p.RHS[0] == "lambda"{
			nullable[p.LHS] = true
		}
	}
	for _,p := range P{
		if p.LHS == nonterm{
			
			if containsTerminal(p.RHS){
				continue
			} else if p.RHS[0] == "lambda"{
				return true
			} else{
				res := true
				for _,v:=range p.RHS{
					if !nullable[v]{
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

func first(P []production_rule)map[string]bool{
	// first_set := map[string]bool{}

}

func follow(){
	fmt.Println("TODO")
}

func predict(){
	fmt.Println()
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
	START_SYMBOL := ""

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

	fmt.Println(NON_TERMINALS)
	fmt.Println(TERMINALS)

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
	START_SYMBOL = findStartState(P)

	fmt.Println(START_SYMBOL)
	fmt.Println()
	DERIVES_TO_LAMBDA := make(map[string]bool) // cache all the values
	for _,p := range P{
		DERIVES_TO_LAMBDA[p.LHS] = derivesToLambda(p.LHS,P)
	}
	fmt.Println(DERIVES_TO_LAMBDA)
	// fmt.Println("Rule",derivesToLambda("Rule",P))










	

	
}