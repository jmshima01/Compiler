package main

import (

	"fmt"
	"os"
	"strconv"
	"strings"
)

type DFA struct{
	transition_chars []byte;
	rows map[string]row;
	tokenID string;
	tokenValueOptional []byte;
}

type row struct{
	state string;
	isStartState bool;
	isAccept bool;
	transitions map[byte]string;
}

func makeDFA(tt []string, alphabet []byte, tokenID string, tokenData []byte) DFA{
	rows := map[string]row{}
	for _,line := range tt{
		if len(line)==0{
			continue
		}
		s := strings.Fields(line)
		trans := map[byte]string{}
		
		for i,v := range alphabet{
			j:=i+2
			trans[v]=s[j]
		}
		r := row{isStartState:s[1]=="0",isAccept:s[0]=="+",transitions: trans,state:s[1]}
		rows[s[1]]=r
	}
	return DFA{transition_chars:alphabet,rows:rows,tokenID: tokenID,tokenValueOptional: tokenData}
}

func(dfa DFA)printDFA(){
	fmt.Println("-------------- Token ID:",dfa.tokenID)
	if len(dfa.tokenValueOptional)!=0{
		fmt.Println("Optional Token Data:",dfa.tokenValueOptional)	
	}
	fmt.Println("                                    ",dfa.transition_chars)
	for k,v := range dfa.rows{
		fmt.Println("ID:",k," isStart:",v.isStartState," isAcc:",v.isAccept," id:",v.state," ",v.transitions)
	}
	fmt.Println("--------------")
}

func parseAlphabetEncoding(s string)[]byte{
	s = strings.Join(strings.Fields(s),"")
	ascii_permited := make([]byte,0)
	for j:=0; j<len(s); j++{
		if string(s[j]) == "x"{
			val,err := strconv.ParseInt(s[j+1:j+3],16,8)
			// fmt.Printf("hex %x:\n",val)
			if err != nil{
				fmt.Println("error reading parse ascii",err)
				os.Exit(1)
			}
			ascii_permited = append(ascii_permited, byte(int(val)))
			j+=2
		} else{
			// fmt.Println("non:",string(first_line[j]))
			ascii_permited = append(ascii_permited, byte(int(s[j])))
		}
	} 
	return ascii_permited
}

func readLines(path string)[]string{
	f,err := os.ReadFile(path)
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	return strings.Split(string(strings.Trim(string(f),"\n")),"\n")
}

func readSrc(path string) []byte{
	data,err := os.ReadFile(path)
	if err!=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	return data
}


func main(){
	args := os.Args

	scanningPath := args[1]
	srcPath := args[2]
	// output_to_path := args[3]

	scanningData := readLines(scanningPath)
	scrData := readSrc(srcPath)

	fmt.Println(scanningData)
	fmt.Println("\n=============")
	 
	asciiAlphabet := make([]byte,0)

	for i,v := range scanningData{
		if i == 0{
			asciiAlphabet = parseAlphabetEncoding(v)
			for _,v := range asciiAlphabet{
				fmt.Println("->",string(v))
			}
			fmt.Println(asciiAlphabet)
		} else{
			tt := strings.Fields(v)
			fmt.Println(tt)
			ttPath := tt[0]
			ttData := readLines(ttPath)

			switch len(tt){
				case 2:
					dfa := makeDFA(ttData,asciiAlphabet,tt[1],make([]byte, 0))
					dfa.printDFA()
				case 3:
					dfa := makeDFA(ttData,asciiAlphabet,tt[1],parseAlphabetEncoding(tt[2]))
					dfa.printDFA()
				default:
					fmt.Println("error in reading tt")
					os.Exit(1)
			}
			
		}
	}
	fmt.Println("\n############################")
	fmt.Println("src:",scrData)





}