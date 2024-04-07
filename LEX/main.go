package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ====== Data Structures ============

type DFA struct{ // Deterministic finite automa (aka regexp)
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

// ============ Helpers ==========

// Converts string to byte array 
func toByteSlice(s string) []byte{
	res := make([]byte,0)
	for _,v := range s{
		res = append(res, byte(v))
	}
	return res
}

// Matches a given byte stream seq to a dfa and returns if it can continue (i.e) not an E transition and the match length
func (dfa DFA) matchSeq(seq []byte)(bool,int){
	curr_row := dfa.startState
	
	for _,c := range seq{
		if c==0{return false,0}
		
		_,ok := dfa.rows[curr_row].transitions[c]
		if !ok{
			println(c,"not in alphabet!")
			os.Exit(1)
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

// DFA constuctor given an correct row from the scan.u file
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

// For debuging
func(dfa DFA)printDFA(){
	fmt.Println("-------------- Token ID:",dfa.tokenID)
	fmt.Println("                                    ",dfa.transition_chars)
	for k,v := range dfa.rows{
		fmt.Println("ID:",k," isStart:",v.isStartState," isAcc:",v.isAccept," id:",v.state," ",v.transitions)
	}
	fmt.Println("--------------")
}

// Converts a string to its coresponding "Alphabet Encoding" i.e weird hex format by keith
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

// Converts weird hex format string to byte array
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

// Read a file into list of its lines
func readLines(path string)[]string{
	f,err := os.ReadFile(path)
	if err!=nil{
		println("Can't read file",path)
		os.Exit(1)
	}
	return strings.Split(string(strings.Trim(string(f),"\n")),"\n")
}

// Write array of lines to a file
func writeLines(path string, data []string){
	result := ""
	for _,l := range data{
		result+=l+"\n"	
	}
	err := os.WriteFile(path,toByteSlice(result),0644)
	if err != nil{
		println("Could not write results to",path)
		os.Exit(1)
	}
}

// Read Src file for char stream
func readSrc(path string) []byte{
	data,err := os.ReadFile(path)
	if err!=nil{
		println("Can't read src file",err)
		os.Exit(1)
	}
	return data
}

// Finds max of int array and returns the value and its index (used for finding the max token length/index)
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

// Returns the line number and start number of a token given its index in the byte stream
func getLine(lineRange [][]int, start int)string{
	line := 1
	result := 0
	toSub := 0

	for _,v := range lineRange{
		l:=v[0]
		r:=v[1]
		if start <= r{
			line=l
			result=start-toSub
			if l==1{result++}
			break
		}
		toSub=r
	}
	return strconv.Itoa(line) + " " + strconv.Itoa(result)
}

// =========== MAIN ===============
func main(){

	args := os.Args
	if len(args) < 3{
		println("Usage: go run go.main <scanningPath> <scrPath> <outFilePath>")
		os.Exit(1)
	}

	scanningPath := args[1]
	srcPath := args[2]
	outPath := args[3]

	scanningData := readLines(scanningPath)
	srcData := readSrc(srcPath)

	asciiAlphabet := make([]byte,0) // Alphabet 
	dfaArray := make([]DFA, 0) // array of each dfa
	tokenLabels := map[string]int{} // lookup table (token -> order)
	tokenInds := map[int]string{} // lookup table (order -> token)
	optionalData := map[string]string{} // 3rd optional value per token value 

	for i,v := range scanningData{
		if i == 0{
			asciiAlphabet = parseAlphabetEncoding(v)
		} else{
			tt := strings.Fields(v)
			// fmt.Println(tt)
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
					println("error in reading tt file row only has len",len(tt))
					os.Exit(1)
			}
		}
	}
	
	tokenMap := map[string]DFA{} // tokenID -> dfa
	dfaResults := map[string]bool{} // MATCHING SET when all false then decision can be made
	longestMatch := make([]int,len(dfaArray)) // keeps track of token index and match lengths

	for _,v := range dfaArray{
		tokenMap[v.tokenID] = v
		dfaResults[v.tokenID] = false
	}

	lineRange := make([][]int,0)
	line:=1
	for i,c := range srcData{
		if c == 10{ // newline in decimal
			r:= []int{line,i}
			lineRange = append(lineRange, r)
			line++
		}
	}
	
	// fmt.Println(lineRange)

	EOF := byte(0) // assuming NUL 0x00 is not part of any alphabet provided so that can terminate loop below as all tokens will fail
	srcData = append(srcData,EOF)
	stream :=srcData
	streamStrt := 0
	streamEnd := 1
	tokens := make([]string,0) // FINAL RESULT
	for streamStrt<len(srcData)-1{
		stream = srcData[streamStrt:streamEnd]
		// fmt.Println(stream)
		for _,v := range tokenMap{
			match,matchLen:=v.matchSeq(stream)
			
			dfaResults[v.tokenID] = match
			
			tokenPos:=tokenLabels[v.tokenID] // lookup table for token label ind
			
			// fmt.Println(v.tokenID,matchLen,match)
			if match && (matchLen >= longestMatch[tokenPos]){ // set max of a token
				longestMatch[tokenPos] = matchLen
			}
		}
		
		// gone through each dfa 
		atLeastOneMatch := false
		for _,x := range dfaResults{ // check all if false 
			if x{
				atLeastOneMatch = true
			}
		}

		if !atLeastOneMatch{ // all false make token choice 
			
			l,ind := maxMatch(longestMatch) // find max length 
			
			// fmt.Println("----------------")
			// fmt.Println("all failed on,",stream)
			// fmt.Println(longestMatch)
			if l==0{
				break
			}
			
			for s := range longestMatch{ // reset match lengths
				longestMatch[s]=0
			}

			line:= getLine(lineRange,streamStrt)
			_,ok := optionalData[tokenInds[ind]]
			if ok{
				s := ""
				s += tokenInds[ind]+ " " + optionalData[tokenInds[ind]] + " " + line
				tokens = append(tokens, s)
				// fmt.Println("=============\nMAX:",tokenInds[ind],l,optionalData[tokenInds[ind]],streamStrt,lineNum,"\n======================")
			} else{
				s:= ""
				s += tokenInds[ind]+ " " + toAlphabetEncoding(string(srcData[streamStrt:streamStrt+l])) + " " + line
				tokens= append(tokens, s)
				// fmt.Println("=============\nMAX:",tokenInds[ind],l,toAlphabetEncoding(string(srcData[streamStrt:streamStrt+l])),streamStrt,lineNum,"\n======================")
			}

			streamStrt += l
			streamEnd = streamStrt+1
			 
		} else {streamEnd++} // not all false, continue stream +1 next char
	}

	// fmt.Println(tokens)
	writeLines(outPath,tokens) // OUTPUT
}