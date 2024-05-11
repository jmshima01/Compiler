package main

import (
	"fmt"
	"strings"
	"unicode"
)

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


