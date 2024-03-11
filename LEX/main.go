package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

)

type DFA struct{
	transition_chars []byte;
	rows map[string]row;
	tokenID string;
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
	
	for _,c := range seq{
	
		_,ok := dfa.rows[curr_row].transitions[c]
		if !ok{
			return false,0
		}

		t := dfa.rows[curr_row].transitions[c] // transition
		
		if t == "E"{
			return false,0
		} else{
			curr_row = t
		}
	}
	if dfa.rows[curr_row].isAccept{return true,len(seq)} else {return true,0}
}


func makeDFA(tt []string, alphabet []byte, tokenID string) DFA{
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
	return DFA{transition_chars:alphabet,rows:rows,tokenID: tokenID, startState: "0"}
}

func(dfa DFA)printDFA(){
	fmt.Println("-------------- Token ID:",dfa.tokenID)
	fmt.Println("                                    ",dfa.transition_chars)
	for k,v := range dfa.rows{
		fmt.Println("ID:",k," isStart:",v.isStartState," isAcc:",v.isAccept," id:",v.state," ",v.transitions)
	}
	fmt.Println("--------------")
}

func toAlphabetEncoding(s string)string{
	result := ""
	for _,c := range s{
		re,_ := regexp.MatchString("[a-w|A-Z|y-z|0-9]",string(c))
		if re{
			result+=string(c)
		} else{
			val := strconv.FormatInt(int64(byte(c)), 16)
			if len(val) == 2{
				result+="x"+val
			} else{
				result+="x0"+val
			}
		}
	}
	return result
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



// ======================================
//              MAIN DRIVER
// ======================================
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
	 
	asciiAlphabet := make([]byte,0)
	dfaArray := make([]DFA, 0)
	tokenLabels := map[string]int{}
	tokenInds := map[int]string{}
	optionalData := map[string]string{}

	for i,v := range scanningData{
		if i == 0{
			asciiAlphabet = parseAlphabetEncoding(v)
		} else{
			tt := strings.Fields(v)
			fmt.Println(tt)
			ttPath := tt[0]
			ttData := readLines(ttPath)
			dfa := makeDFA(ttData,asciiAlphabet,tt[1])
			dfaArray = append(dfaArray, dfa)
			tokenLabels[dfa.tokenID] = i-1
			tokenInds[i-1] = dfa.tokenID
			
			
			switch len(tt){
				case 2:
				case 3:
					optionalData[tt[1]] = tt[2]
				default:
					fmt.Println("error in reading tt file row only has len",len(tt))
					os.Exit(1)
			}
			
		}
	}
	
	tokenMap := map[string]DFA{}
	
	dfaResults := map[string]bool{}
	longestMatch := make([]int,len(dfaArray))
	for _,v := range dfaArray{
		tokenMap[v.tokenID] = v
		dfaResults[v.tokenID] = false
	}

	EOF := byte(0)
	scrData = append(scrData,EOF)
	stream :=scrData
	streamStrt := 0
	streamEnd := 1
	// EOF := false
	lineNum := 1
	for streamStrt<len(scrData)-1{
		
		stream = scrData[streamStrt:streamEnd]
		
		// fmt.Println(stream)
		
		for _,v := range tokenMap{
			match,matchLen:=v.matchSeq(stream)
			
			dfaResults[v.tokenID] = match
			
			tokenPos:=tokenLabels[v.tokenID] // lookup table for token label ind
			
			// fmt.Println(v.tokenID,matchLen,match)
			if match && (matchLen >= longestMatch[tokenPos]){
				longestMatch[tokenPos] = matchLen
			}
		}
		// fmt.Println(dfaResults)
		// itered through each dfa 
		atLeastOneMatch := false
		for _,x := range dfaResults{ // check all false 
			if x{
				atLeastOneMatch = true
			}
		}

		if !atLeastOneMatch{
			
			l,ind := maxMatch(longestMatch)
			
			// if l==0{ // no more matches possible in stream (I think)
			// 	break
			// }

			fmt.Println("----------------")
			fmt.Println("all failed on,",stream)
			
			
			fmt.Println(longestMatch)

			for s := range longestMatch{
				longestMatch[s]=0
			}
			_,ok := optionalData[tokenInds[ind]]
			if ok{
				fmt.Println("=============\nMAX:",tokenInds[ind],l,optionalData[tokenInds[ind]],streamStrt,lineNum,"\n======================")
			} else{
				fmt.Println("=============\nMAX:",tokenInds[ind],l,toAlphabetEncoding(string(scrData[streamStrt:streamStrt+l])),streamStrt,lineNum,"\n======================")
			}

			streamStrt += l
			streamEnd = streamStrt+1
			 
		} else {
			if scrData[streamEnd] == 10{
				lineNum++
				fmt.Println("line found!",lineNum,stream)
			}
			streamEnd++	
		}
	}
	fmt.Println(scrData)
	fmt.Println(len(scrData))
	fmt.Println(cap(scrData))

	// two := tokenMap["twosmallwords"]
	// m,l := two.matchSeq(toByteSlice("rop rop "))
	// fmt.Println(m,l)
	// s:=toAlphabetEncoding("rop rop ")
	// fmt.Println(s)
}