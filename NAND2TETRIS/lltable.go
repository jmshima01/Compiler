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

// ========== LL(1) Table driven Parser =================
func AST(grammar []string, tokFilepath string) *Node{
	
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
	fmt.Println()
	fmt.Println("ruleLookup := map[int]ProductionRule{")
	for k,v := range ruleLookup{
		r := "string{"
		for i,s := range v.rhs{
			r+= fmt.Sprintf("\"%s\"",s)
			if i == len(v.rhs)-1{
				continue
			}
			r+=","
		}
		r+="}"
		
		x:=fmt.Sprintf("%d: ProductionRule{lhs:\"%s\", rhs:%s)",k,ruleLookup[k].lhs,r)
		fmt.Println(x)
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
	fmt.Println("=======================LL table driven parsing============")
	// test4 := "bghm$"
	tokenStream := readTokens(tokFilepath)

	for _, v := range tokenStream {
		fmt.Println(v)
	}

	S := make(stack, 0)
	Q := make(queue, 0)
	S = append(S, startState)
	for i, v := range tokenStream {
		if v.tokenType == "identifier" {
			if tokenStream[i+1].value == "(" {
				Q.push(token{value: v.value, tokenType: "subroutinename"})
			} else if tokenStream[i+1].value == "." {
				Q.push(token{value: v.value, tokenType: "objectname"})
			} else if tokenStream[i+1].value == "[" {
				Q.push(token{value: v.value, tokenType: "array"})
			} else {
				Q.push(v)
			}

		} else if v.tokenType == "stringConst" {
			Q.push(token{value: v.value, tokenType: "stringconstant"})
		} else if v.tokenType == "integerConst" {
			Q.push(token{value: v.value, tokenType: "integerconstant"})
		} else {
			Q.push(v)
		}

	}
	Q.push(token{value: "$", tokenType: "$"})
	fmt.Println("---------------")
	for _, v := range Q {
		fmt.Println(v)
	}

	root := makeNode("ROOT", nil, 0)
	current := root
	current.debug()
	uniqueID := 1
	for {
		if S.isEmpty() {
			if !Q.isEmpty() {
				fmt.Println("syntax error:", Q)
				os.Exit(2)
			}

			break
		}

		fmt.Println("S:", S)
		fmt.Println("Q:", Q)

		s := S.peek()
		q := ""
		t := Q.peek()

		if t.tokenType == "keyword" || t.tokenType == "symbol" {
			q = t.value
		} else {
			q = t.tokenType
		}

		if s == "<*>" {

			S.pop()

			// SDT Ast conversion:
			switch current.data {
			case "ArrayName":
				current = current.parent
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data != "ArrayName" {
						newChildren = append(newChildren, v)
					} else {
						for _,x := range v.children{
							x.parent = current
						}
						newChildren = append(newChildren, v.children...)
					}
				}
				current.children = newChildren

			case "ClassName":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = nil
				current = current.parent
			
			case "SubroutineCallName":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = nil
				current = current.parent
			case "SubroutineName":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = nil
				current = current.parent
			case "Type":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = nil
				current = current.parent
			case "ClassVarDecSF":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = nil
				current = current.parent
			case "SubroutineDecCFM":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = nil
				current = current.parent

			case "SubroutineDecType":
				current.data = current.children[0].data
				current.id = current.children[0].id
				current.children = nil
				current = current.parent
			// case "KeywordConstant":
			// 	current.data = current.children[0].data
			// 	current.id = current.children[0].id
			// 	current.children = nil
			// 	current = current.parent
			
			case "Term":
				current = current.parent
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data != "Term" {
						newChildren = append(newChildren, v)
					} else {
						for _,x := range v.children{
							x.parent = current
						}
						newChildren = append(newChildren, v.children...)
					}
				}
				current.children = newChildren

			case "Expression":
				newChildren := make([]*Node, 0)
				for _, v := range current.children {

					if v.data != ")" && v.data != "(" {
						newChildren = append(newChildren, v)
					}
				}
				current.children = newChildren
				current = current.parent

			case "Statement":
				current = current.parent
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data != "Statement" {
						newChildren = append(newChildren, v)
					} else {
						for _,x := range v.children{
							x.parent = current
						}
						newChildren = append(newChildren, v.children...)
					}
				}
				current.children = newChildren

			case "ExpressionTerms":
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
						if v.data != "ExpressionTerms" {
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

			case "SubroutineBodyVarDec":
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
						if v.data != "SubroutineBodyVarDec" {
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
			case "ExtraVarExt":
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
						if v.data != "ExtraVarExt" {
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
			case "Statements":
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
						if v.data != "Statements" {
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
			case "LetExpressionCheck":
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
						if v.data != "LetExpressionCheck" {
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

			case "ReturnExpressionCheck":
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
						if v.data != "ReturnExpressionCheck" {
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
			case "ExpressionList":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {

					current = current.parent
				}
			case "ParameterList":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {

					current = current.parent
				}
			case "ClassVarDec":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {
					newChildren := make([]*Node, 0)
					for _, v := range current.children {
						if v.data == "," || v.data == ";" {
							continue
						}
						newChildren = append(newChildren, v)

					}
					current.children = newChildren
					current = current.parent

				}
			case "SubroutineDec":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {
					newChildren := make([]*Node, 0)
					for _, v := range current.children {
						if v.data == "(" || v.data == ")" {
							continue
						}
						newChildren = append(newChildren, v)

					}
					current.children = newChildren
					current = current.parent
				}
			case "VarDecExt":
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
						if v.data != "VarDecExt" {
							newChildren = append(newChildren, v)
						} else {
							for _, x := range v.children {
								x.parent = current
								if x.data != "," {
									newChildren = append(newChildren, x)
								}
							}

						}
					}
					current.children = newChildren
				}

			case "ExpressionListExt":
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
						if v.data != "ExpressionListExt" {
							newChildren = append(newChildren, v)
						} else {
							for _, x := range v.children {
								x.parent = current
								if x.data != "," {
									newChildren = append(newChildren, x)
								}
							}

						}
					}
					current.children = newChildren
				}
			case "ParameterListExt":
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
						if v.data != "ParameterListExt" {
							newChildren = append(newChildren, v)
						} else {
							for _, x := range v.children {
								x.parent = current
								if x.data != "," {
									newChildren = append(newChildren, x)
								}
							}

						}
					}
					current.children = newChildren
				}

			case "ElseStatment":
				if current.children[0].data == "lambda" {
					current.data = current.children[0].data
					current.id = current.children[0].id
					current.children = nil
					current = current.parent
					current.children = current.children[:len(current.children)-1]
				} else {
					newChildren := make([]*Node, 0)
					for _, v := range current.children {
						if v.data == "{" || v.data == "}" || v.data == "else" {
							continue
						}
						newChildren = append(newChildren, v)

					}
					current.children = newChildren
					current = current.parent
				}

			case "DoStatement":
				current = current.parent
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data != "DoStatement" {
						newChildren = append(newChildren, v)
					} else {
						for _, x := range v.children {
							x.parent = current
							if x.data != "do" && x.data != ";" {
								newChildren = append(newChildren, x)
							}
						}

					}
				}
				current.children = newChildren
				current.children[0].data = "DoStatement"

			case "LetStatement":
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data == "let" || v.data == ";" || v.data == "=" {
						continue
					}
					newChildren = append(newChildren, v)

				}
				current.children = newChildren
				current = current.parent

			case "WhileStatement":
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data == "while" || v.data == "(" || v.data == ")" || v.data == "{" || v.data == "}" {
						continue
					}
					newChildren = append(newChildren, v)

				}
				current.children = newChildren
				current = current.parent
			
			case "VarDec":
				current.children = current.children[1 : len(current.children)-1]
				current = current.parent
			
			case "ReturnStatement":
				if len(current.children) == 2 {
					current.children = nil
				} else {
					current.children = current.children[1:2]
				}
				current = current.parent
			
			case "SubroutineBody":
				current.children = current.children[1 : len(current.children)-1]
				current = current.parent
			
			case "SubroutineCall":
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data == "." || v.data == "(" || v.data == ")" {
						continue
					}
					newChildren = append(newChildren, v)

				}
				current.children = newChildren
				current = current.parent

			case "Class":
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data == "{" || v.data == "}" || v.data == "class" || v.data == "$" {
						continue
					}
					newChildren = append(newChildren, v)

				}
				current.children = newChildren
				current = current.parent
			
			case "IfStatement":
				newChildren := make([]*Node, 0)
				for _, v := range current.children {
					if v.data == "{" || v.data == "}" || v.data == "(" || v.data == ")" || v.data == "if" {
						continue
					}
					newChildren = append(newChildren, v)

				}
				current.children = newChildren
				current = current.parent

			
			default:
				current = current.parent
			}

			continue
		}

		if s == "lambda" {
			S.pop()
			lambNode := makeNode("lambda", current, uniqueID)
			addChild(current, lambNode)
			current.debug()
			uniqueID++
			continue
		}

		if isTerminal(s) || s == "$" {
			if s == q {
				
				if t.tokenType == "stringconstant" || t.tokenType == "integerconstant"{
					v := makeNode(t.value, current, uniqueID)
					uniqueID++
					term := makeNode(t.tokenType, current, uniqueID)
					addChild(term,v)
					uniqueID++
					addChild(current,term)
				} else{
					term := makeNode(t.value, current, uniqueID)
					addChild(current, term)
					uniqueID++
				}
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

			fmt.Println("Parsing Error: (No such token in LL table or associated rule)", s, q, Q)
			fmt.Println("-----")
			fmt.Println(s, S)
			os.Exit(2)
		}

		fmt.Println("fetching rule...", nextRule)
		top := S.pop()
		newNode := makeNode(top, current, uniqueID)
		addChild(current, newNode)
		current.debug()

		current = newNode
		uniqueID++
		// add rule in reverse to stack...
		S = append(S, "<*>") // end of production
		for i := len(nextRule.rhs) - 1; i >= 0; i-- {
			S = append(S, nextRule.rhs[i])
		}
	}
	S = nil
	Q = nil
	current.debug()
	fmt.Println("============")
	// printTree(current)
	ast := current.children[0]
	// g := ""
	// graphiz:=*(toGraphiz(current,&g))

	nodeInfo := ""
	nodeInfo = *(genNodeInfo(ast, &nodeInfo))

	edgeInfo := ""
	edgeInfo = *(genEdgeInfo(ast, &edgeInfo))

	toGraphiz := nodeInfo + "\n" + edgeInfo
	writeToFile("parsetree.txt", toGraphiz)
	fmt.Println(toGraphiz)
	printTree(ast)

	return ast
}