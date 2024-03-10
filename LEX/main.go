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
	startState string; // assuming should always be zero... if not bugs may arise 
}

type row struct{
	state string;
	isStartState bool;
	isAccept bool;
	transitions map[byte]string;
}

func toByteSlice(s string) []byte{
	res := make([]byte,0)
	for _,v := range s{
		res = append(res, byte(v))
	}
	return res
}

func (dfa DFA) matchSeq(seq []byte)(bool,int){
	curr_row := dfa.startState
	// fmt.Println("seq on dfa",dfa.tokenID,string(seq))
	for _,c := range seq{
		t := dfa.rows[curr_row].transitions[c] // transition
		// fmt.Println("t",t)
		// fmt.Println("row",curr_row)
		// fmt.Println("c",string(c))
		// fmt.Println()
		if t == "E"{
			return false, 0
		} else{
			curr_row = t
		}
	}
	if dfa.rows[curr_row].isAccept{return true,len(seq)} else {return false,0}
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
	return DFA{transition_chars:alphabet,rows:rows,tokenID: tokenID,tokenValueOptional: tokenData, startState: "0"}
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

func toAlphabetEncoding(s string)string{
	return ""
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

func maxMatch(m []int)(int,int){
	result := -1
	ind := 0
	for i,v := range m{
		if i==0{
			result = v
			continue
		}
		if result<v{
			ind = i
			result=v
		}

	} 

	return result,ind
}


func main(){
	args := os.Args
	
	if len(args) < 3{
		fmt.Println("Usage: go run go.main <scanningPath> <scrPath> <optional_destination>")
		os.Exit(1)
	}

	scanningPath := args[1]
	srcPath := args[2]
	// output_to_path := args[3]

	scanningData := readLines(scanningPath)
	scrData := readSrc(srcPath)

	fmt.Println(scanningData)
	fmt.Println("\n=============")
	 
	asciiAlphabet := make([]byte,0)
	dfaArray := make([]DFA, 0)
	tokenLabels := map[string]int{}
	tokenInds := map[int]string{}

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
					dfaArray = append(dfaArray, dfa)
					tokenLabels[dfa.tokenID] = i-1
					tokenInds[i-1] = dfa.tokenID
					
					// dfa.printDFA()
				case 3:
					dfa := makeDFA(ttData,asciiAlphabet,tt[1],parseAlphabetEncoding(tt[2]))
					dfaArray = append(dfaArray, dfa)
					tokenLabels[dfa.tokenID] = i-1
					tokenInds[i-1] = dfa.tokenID
					// dfa.printDFA()
				default:
					fmt.Println("error in reading tt")
					os.Exit(1)
			}
			
		}
	}
	fmt.Println(tokenLabels)
	tokenMap := map[string]DFA{}
	matchingTokenSet := map[string]bool{} // set
	dfaResults := map[string]bool{}
	longestMatch := make([]int,len(dfaArray))
	for _,v := range dfaArray{
		tokenMap[v.tokenID] = v
		matchingTokenSet[v.tokenID] = true
		dfaResults[v.tokenID] = false
	}

	
	stream := scrData
	streamStrt := 0
	streamEnd := 1
	currLine := 1
	currChar := 1
	for i := 0; i<len(scrData); i++{
		stream = scrData[streamStrt:streamEnd]
		// fmt.Println(stream)
		
		for _,v := range tokenMap{
			match,matchLen:=v.matchSeq(stream)
			
			dfaResults[v.tokenID] = match
			tokenPos:=tokenLabels[v.tokenID]
			
			fmt.Println(v.tokenID,matchLen,match)
			if match && (matchLen >= longestMatch[tokenPos]){
				longestMatch[tokenPos] = matchLen
			}
		}

		// itered through each dfa 
		atLeastOneMatch := false
		for _,x := range dfaResults{ // check all false 
			if x{
				atLeastOneMatch = true
			}
		}

		if !atLeastOneMatch{
			fmt.Println("all failed on,",string(stream),"chr:",currChar)
			
			fmt.Println("new token index:", i+1)
			fmt.Println()
		
			// maximum,maxMatch:=0,""
			fmt.Println("failing->",string(stream))
			// for s,t := range longestMatch{
			// 	fmt.Println("->",s,t)
			// }
			
			l,ind := maxMatch(longestMatch)
			fmt.Println(longestMatch)
			for s := range longestMatch{
				longestMatch[s]=0
			}
			fmt.Println("=============\nMAX:",tokenInds[ind],l,currLine,currChar-1,"\n======================")
			
			streamStrt = i
			streamEnd = i+1
			i--
			 
		} else {streamEnd++}

		currChar++

		if streamStrt == 10{ //newline
			currLine++
			currChar=1
		}
		
	}
	two := tokenMap["twosmallwords"]
	m,l := two.matchSeq(toByteSlice("rop rop "))
	fmt.Println(m,l)
}