package main

import (
	"fmt"
	"os"
	"strings"
	"regexp"
)


func (t token)toString()string{
	return fmt.Sprintf("<%s> %s </%s>",t.tokenType,t.value,t.tokenType)
}

func lex(data []string)[]token{
	keywords := set{"class": true, "constructor": true, "function": true, "method": true, "field": true, "static": true, "var": true, "int": true, "char": true, "boolean": true, "void": true, "true": true, "false": true, "null": true, "this": true, "return": true, "let": true, "do": true, "if": true, "else": true, "while": true}
	// keyword,_:= regexp.Compile("class|constructor|function|method|field|static|var|int|char|boolean|void|true|false|null|this|return|let|do|if|else|while")
	integerConst,_ := regexp.Compile(`[1-3][0-2][0-7][0-6][0-7]|^(\d{1,4}$)`)
	stringConst,_ := regexp.Compile("(^(\")(.*\"))")
	identifier,_ := regexp.Compile("^([a-zA-Z_])([a-zA-Z_0-9])*$")
	isAlpha,_ := regexp.Compile("[a-zA-Z]")
	tokens := make([]token,0)
	// matchSet := map[string]int{}
	for _,v := range data{
		s:=""
		l:=len(v)
		i:=0
		c:=""
		for i<l{
			c = string(v[i])			
			if c=="{"||c== "}"||c=="("||c==")"|| c=="["||c=="]"||c=="+"||c=="-"||c=="&"||c=="|"||c=="."||c==","||c==";"||c=="*"||c=="/"||c=="~"||c=="="||c=="<"||c==">"{
				if c == "|"{ // to deal with cfg
					tokens = append(tokens, token{value:"pipe",tokenType:"symbol"})
				} else{
					tokens = append(tokens, token{value:c,tokenType:"symbol"})
				}
				
			} else if isAlpha.MatchString(c){
				for i<l{
					if !identifier.MatchString(s+c){
						break
					}
					s+=c
					i++
					c=string(v[i])
				}
				_,isKeyword := keywords[s]
				if isKeyword{
					tokens = append(tokens, token{value:s,tokenType:"keyword"})
					fmt.Println("found keyword",s)
				} else if identifier.MatchString(s){
					tokens = append(tokens, token{value:s,tokenType:"identifier"})
					fmt.Println("found ident",s)
				}
				s=""
				i--
			} else if c=="0"||c=="1"||c=="2"||c=="3"||c=="4"||c=="5"||c=="6"||c=="7"||c=="8"||c=="9"{
				for i<l{
					if !integerConst.MatchString(s+c){
						break
					}
					s+=c
					i++
					c=string(v[i])
				}
				if integerConst.MatchString(s){
					tokens = append(tokens, token{value:s,tokenType:"integerConst"})
				} else{
					fmt.Println("Lex error int",s)
					os.Exit(1)
				}
				s=""
				i--
			} else if c == "\""{
				for i<l{
					s+=c
					i++
					c=string(v[i])
					if c == "\""{
						s+=c
						break
					}
				}
				if stringConst.MatchString(s){
					tokens = append(tokens, token{value:s[1:len(s)-1],tokenType:"stringConst"})
				} else{
					fmt.Println("Lex error str",s)
					os.Exit(1)
				}
				s=""
			}
			i++
		}
	}
	return tokens
}

func readTokens() []token{
	args := os.Args
	lines := readLines(args[2])

	// remove comments from .jack file:
	clean := make([]string,0)
	comment,_:= regexp.Compile(`^(//).*|^(/\*).*\*/`)
	sameLineComment, _ := regexp.Compile("//")
	multiLine, _ := regexp.Compile(`^(/\*).*`)
	multiLineEnd, _ := regexp.Compile(`\*/`)
	multiLineFlag := false
	for _,v := range lines{
		line := strings.TrimSpace(v)
		if line == ""{
			continue
		}
		if comment.MatchString(v){
			continue
		}
		
		if multiLine.MatchString(v){
			multiLineFlag = true
			continue
		}
		if multiLineEnd.MatchString(v){
			multiLineFlag = false
			continue
		}

		if !comment.MatchString(line) && !multiLineFlag{
			if sameLineComment.MatchString(line){
				line = strings.Split(line,"//")[0]
				clean = append(clean, line)
			} else {clean = append(clean, line)}
		}
	}

	// debug...
	for _,v := range clean{
		fmt.Println(v)
	}
	tokens := lex(clean)
	fmt.Println(tokens)
	return tokens
}