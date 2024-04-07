package main

import (
	"fmt"
	"os"
)

var lookupTable map[string]string = map[string]string{
	"\\s":"char 0x20\n", 
	"\\+":"char +\n",
	"\\|":"char |\n",
	"\\(":"char (\n",
	"\\)":"char )\n",
	"\\-":"char -\n",
	"\\*":"char *\n",
	"\\n":"char x0a\n",
	"\\\\":"char \\\n", 
	"+":"plus +\n",
	"(":"open (\n",
	"-":"dash +\n",
	"|":"pipe |\n",
	".":"dot .\n",
	")":"close )\n",
}

func getToken(s string) string{
	_,ok := lookupTable[s]
	if !ok{
		return fmt.Sprintf("char %s\n",s)
	} else {return lookupTable[s]}
}

func main(){
	args := os.Args
	// fmt.Println(args)

	stream := args[1]
	// fmt.Println(stream)
	
	lookahead := false
	tokens := ""

	for _,c := range stream{
		
		if string(c) == "\\" && !lookahead{
			lookahead = true
			continue
		}
		
		if !lookahead{
			tokens += getToken(string(c))
		} else{
			lookahead = false
			tokens += getToken("\\"+string(c))
		}
	}
	fmt.Println(tokens)
}