package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"


)

// ============ Types ============== 
type ProductionRule struct{
	lhs string;
	rhs []string;
};

// for first, follow, & predict sets 
type set struct{
	items map[string]bool;
};

func makeSet() set{
	var s set
	s.items = make(map[string]bool,0)
	return s
}

func (s set)add(v string){
	if s.items == nil{
		return
	}
	s.items[v] = true

}

// only want keys of the set/map
func (s set)getValues() []string{
	res := make([]string,0)
	for k,_ := range s.items{
		res = append(res, k)
	}
	return res
}

func setUnion(s1 set, s2 set)set{
	union := makeSet()
	for k, _ := range s1.items{
		union.items[k] = true
	}
	for k, _ := range s2.items{
		union.items[k] = true
	}
	return union
}

type Grammar struct{
	raw []string;
	P []ProductionRule;
	N set; // Nonterminals
	sigma set; // terminals
	symbols set;
	// Caches for memoization (maps of sets)
	dervToLamb set;
	firstSets map[string]set;
	followSets map[string]set;
	predictSets map[string]set; // N -> predict(N) 
}
func makeGrammar(cfg []string) Grammar{
	P := makeProductionRules(cfg)
	
	nonTerminals := makeSet()
	terminals := makeSet()
	for _,line := range cfg{
		v := strings.Split(line, " ")
		for _,s := range v{
			if isNonTerminal(s){
				nonTerminals.add(s)
			} else if isTerminal(s){
				terminals.add(s)
			}
		}
	}
	symbols := setUnion(terminals,nonTerminals)
	symbols.add("$")
	d := makeDerivesToLambda(P)
	return Grammar{raw:cfg, P:P, N:nonTerminals, symbols:symbols, sigma:terminals,dervToLamb:d, firstSets: make(map[string]set,0),followSets: make(map[string]set,0),predictSets: make(map[string]set,0)}
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


func (G Grammar) first(N string, firstSet set, seen set) set{
	_,ok := seen.items[N]
	if ok{
		return firstSet
	}
	if N == "lambda"{
		return set{}
	}
	if isTerminal(N){
		return set{items:map[string]bool{N:true}}
	}
	seen.items[N]=true
	for _,p := range G.P{
		if p.lhs == N{
			for i,v := range p.rhs{
				if i==0 && isTerminal(v){
					firstSet.add(v)
					break
				} else{
					if v == "lambda" {break}
					
					firstSet = G.first(v,firstSet,seen)
					
					if !G.dervToLamb.items[v]{
						break
					}
				}
			}
		}
	}
	return firstSet	
}

func (G Grammar) follow(N string, followSet set, seen set) set{
	_,ok := seen.items[N]
	if ok{
		return followSet
	}
	seen.add(N)
	for _,p := range G.P{
		flag := false
		last := p.lhs
		for _,v := range p.rhs{
			
			if v == N{
				flag = true
				continue
			}
			if flag{
				followSet = setUnion(followSet,G.firstSets[v])
				if !G.dervToLamb.items[v]{
					flag = false
					break
				}
			}
		}
		if flag{
			N=last
			followSet = G.follow(N,followSet,seen)
		}
	}
	return followSet
}

func (G Grammar) predict(p ProductionRule)set{
	predictSet := set{}
	flag := true
	for _,v := range p.rhs{
		if v == "lambda"{
			flag = false
			predictSet = G.followSets[p.lhs]
			break
		}
		
		predictSet = setUnion(predictSet,G.firstSets[v])
		
		if !G.dervToLamb.items[v]{
			flag = false
			break
		}
	}
	if flag{
		predictSet = setUnion(predictSet,G.followSets[p.lhs])
	}
	return predictSet


}

func makeDerivesToLambda(P []ProductionRule)set{
	lookup := makeSet()
	for _,p := range P{
		lookup.items[p.lhs] = derivesToLambda(p.lhs,P)
	}
	return lookup
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


	cfg := makeGrammar(grammar)

	P := cfg.P
	fmt.Println(P)
	
	







}