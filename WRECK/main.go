package main

import (
	"fmt"
	"strings"
	"os"
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
				t := token{value: "0x20", tokenType: "char"}
				skip = true
				tokens = append(tokens, t)
				continue
			}else if string(s[i+1]) == `n` {
				t := token{value: "0x0a", tokenType: "char"}
				skip = true
				tokens = append(tokens, t)
				continue
			} else {
				t := token{value: string(s[i+1]), tokenType: "char"}
				t.value = convertAlpha(t.value)
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
			t.value = convertAlpha(t.value)
			tokens = append(tokens, t)

		}
	}
	return tokens
}

func chooseNewLambda(){
	lambdac := byte('A')
	for _,v := range hexAlphabet{
		if lambdac == v{
			lambdac++
		}
	}

	lambdaChar = string(lambdac) // global

}



func main() {
	args := os.Args
	
	grammar := readLines("llre.cfg")
	LLTable,startState,ruleLookup,rowLookup,columnLookup := makeLLTable(grammar)

	scan := readLines(args[1])	
	toTokenize := [][]string{}
	
	for i, line := range scan {
		vals := strings.Fields(line)
		if i == 0 {
			
			hexAlphabet = parseAlphabetEncoding(strings.TrimSpace(strings.Join(strings.Fields(line),"")))
			chooseNewLambda()
			for _,v := range hexAlphabet{
				alphabet = append(alphabet, convertAlpha(string(v)))
			}
			fmt.Println(hexAlphabet)
			continue
		}

		fmt.Println(vals)
		toTokenize = append(toTokenize, vals)
		
	}

	scanu := ""

	alphaEncoded := alphabetEncoded(hexAlphabet)
	alphaEncoded+="\n"
	scanu+=alphaEncoded

	tokenize := []string{} 
	tokenNames := []string{}
	for _,v := range toTokenize{
		tokenize = append(tokenize, v[0])
		tokenNames = append(tokenNames, v[1])
		x:=v
		x[0] = x[1]+".tt"
		scanu+=strings.Join(x," ")
		scanu+= "\n"
		
	}

	writeToFile(args[2],scanu)

	for i := range tokenize{
		tokenStream := lex(tokenize[i])
		fmt.Println("=================")
		fmt.Println(tokenStream)
		fmt.Println("================")
		ast := makeAST(tokenStream,LLTable,startState,ruleLookup,rowLookup,columnLookup)
		makeNFA(ast,tokenNames[i]+".nfa")
	}
}
