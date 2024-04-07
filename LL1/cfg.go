package main

import (
		"fmt"
	)

func setUnion(s1 map[string]bool, s2 map[string]bool)map[string]bool{
	s_union := map[string]bool{}
	for k, _ := range s1{
		s_union[k] = true
	}
	for k, _ := range s2{
		s_union[k] = true
	}
	return s_union
}

type CFG struct{
	P []production_rule; // prod rules
	N map[string]bool; // non terminals
	sigma map[string]bool; // terminals
	symbols map[string]bool
	S string; // start state
	nullable map[string]bool; // cache for Nonterminals w/ a lambda rule
	derivestoLambda map[string]bool; // cache for derives to lambda
	follow_sets map[string]map[string]bool;
	first_sets map[string]map[string]bool;
	predict_sets map[string]map[string]bool;
};


// ========== Methods ===========
func (cfg CFG) printCFG(){
	fmt.Println("Starting Goal:",cfg.S)
	fmt.Println("Nonterminals:",cfg.N)
	fmt.Println("Symbols:",cfg.symbols)
	fmt.Println("Rules:",cfg.P)
	fmt.Println("Nullable:",cfg.nullable)
	fmt.Println("Derv2Lambda:",cfg.derivestoLambda)
}

func (self CFG) containsTerminal(rhs []string) bool{
	for _,v := range rhs{
		_,ok := self.sigma[v]
		if ok {
			return true
		}
	}
	return false
}

func (cfg CFG) DerivesToLambda(nonterm string) bool {
	for _,p := range cfg.P{
		if p.LHS == nonterm{
			if cfg.containsTerminal(p.RHS){
				continue
			} else if p.RHS[0] == "lambda"{
				return true
			} else{
				res := true
				for _,v:=range p.RHS{
					if !cfg.nullable[v]{
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

func (cfg CFG) makeLambda()map[string]bool{
	derives_lambda_cache := make(map[string]bool) // cache all the values
	for _,p := range cfg.P{
		derives_lambda_cache[p.LHS] = cfg.DerivesToLambda(p.LHS)
	}
	return derives_lambda_cache
}

func (cfg CFG) first(n []string)map[string]bool{
	first_set := map[string]bool{}
	for _,v := range n{
		_,ok := cfg.sigma[v]
		if ok{
			first_set[v] =true
			return first_set
		}
		
	}
	
	return first_set
}

func (self CFG) follow(n []string)map[string]bool{


	return nil
}

func (cfg CFG) predict(n string)map[string]bool{
	predict := map[string]bool{}
	for _,p := range cfg.P{
		if n == p.LHS{
			if cfg.derivestoLambda[n]{
				predict = set_union(cfg.follow(p.RHS),predict)
			} else{
				predict = set_union(cfg.first(p.RHS),predict)
			}
		}
	}
	return predict
}
