#include <stdio.h>
#include <iostream>
#include <fstream>
#include <string>
#include <bits/stdc++.h> 
#include <set>


std::vector<std::string> split(std::string s, std::string delimiter){
    size_t pos_start = 0, pos_end, delim_len = delimiter.length();
    std::string token;
    std::vector<std::string> res;

    while ((pos_end = s.find(delimiter, pos_start)) != std::string::npos) {
        token = s.substr(pos_start, pos_end - pos_start);
        pos_start = pos_end + delim_len;
        res.push_back (token);
    }
    res.push_back(s.substr(pos_start));
    return res;
}

std::set<std::string> setUnion(std::set<std::string> a, std::set<std::string> b){
    std::set<std::string> res = a;
    res.insert(b.begin(),b.end());
    return res;
}


class CFG{
    public:
        std::unordered_map<std::string,std::vector<std::string>> cfg;
        std::vector<std::string> file_data;
        std::string start_state;
        std::set<std::string> non_terminals;
        


        std::vector<std::pair<std::string,std::string>> rules;
        
        struct xb{
            std::string X;
            std::string B;
        };
        
        CFG(std::vector<std::string> file_data){
            
            this->file_data = file_data;
            this->start_state = this->_getStartState();
            this->non_terminals = this->_getNonTerminals();
            this->rules = this->_makeProdRules();
            this->cfg = this->_makeCFG();
        }
        ~CFG();

        void print_cfg(void);
        bool derivesToLambda(std::string n, std::vector<std::string> T);
        std::pair<std::set<std::string>,std::set<std::string>> firstSet(xb XB,std::set<std::string> T);

    private:
        bool _isNonTerminal(std::string s);
        std::string _getStartState(void);
        std::set<std::string> _getNonTerminals(void);
        std::vector<std::pair<std::string,std::string>> _makeProdRules(void);
        std::unordered_map<std::string,std::vector<std::string>> _makeCFG(void);
        bool _containsTerminal(std::string s);
};

CFG::~CFG(){ 
};


bool CFG::_isNonTerminal(std::string s){
    for(char c : s){
        if (isupper(c)){
            return true;
        }
    }
    return false;
}


std::string CFG::_getStartState(void){
    std::string delimiter = " -> ";
    return this->file_data[0].substr(0, this->file_data[0].find(delimiter)); 
}

std::set<std::string> CFG::_getNonTerminals(void){
    std::set<std::string> nonterms;
    for(std::string line : this->file_data){
        std::vector<std::string> split_line = split(line,std::string(" "));
        for(std::string s : split_line){
            if(this->_isNonTerminal(s)){
                nonterms.insert(s);
            }
        }
    }
    return nonterms;
}

std::vector<std::pair<std::string,std::string>> CFG::_makeProdRules(void){
    
    std::vector<std::pair<std::string,std::string>> pr;
    std::string curr;
    for (auto s : this->file_data){
        std::vector<std::string> f = split(s,std::string(" -> "));
        if(f.size()==2){
            curr = f[0];
            std::vector<std::string> bars = split(f[1],std::string(" | "));
            for(std::string b : bars){
                pr.push_back(std::make_pair(curr,b));
            }
            
        }
        else {
            std::string x = f[0].substr(2,f[0].length()-2);
            
            std::vector<std::string> br = split(x,std::string(" | "));
            for(std::string t : br){
                pr.push_back(std::make_pair(curr,t));
            }
        }
    }
    return pr;       
}
    
void CFG::print_cfg(void){
    std::cout << this->start_state << std::endl;
    std::cout << "Nonterminals: { ";
    
     
    for(auto s : this->non_terminals){
        std::cout << s << " ";
    } std::cout << "}" << std::endl;

    for(auto i : this->rules){
        std::cout << i.first << " -> " << i.second << std::endl;
    }
    printf("\n");
    for(auto p : this->cfg){
        std::cout << p.first << ": ";
        for(auto s : p.second){
            std::cout << s << " | ";
        }std::cout << std::endl;
    }
    
}


std::unordered_map<std::string,std::vector<std::string>> CFG::_makeCFG(void){
    std::unordered_map<std::string,std::vector<std::string>> m;
    for(std::pair<std::string,std::string> p : this->rules){
        if(m.count(p.first)){
            m[p.first].push_back(p.second);
        }
        else{
            m[p.first] = std::vector<std::string>();
            m[p.first].push_back(p.second);
        }
    }
    return m;
}

bool CFG::_containsTerminal(std::string s){
    std::vector<std::string> st = split(s,std::string(" "));
    for (std::string x : st){
        if(!this->_isNonTerminal(x)){
            return true;
        }
    }
    return false;
}

bool CFG::derivesToLambda(std::string n, std::vector<std::string> T){
    for(std::string p : this->cfg.at(n)){
        if(std::find(T.begin(), T.end(), p) != T.end()){
            continue;
        }
        if(p == "lambda"){
            return true;
        }
        if(this->_containsTerminal(p)){
            continue;
        }
        bool allderivelambda = true;
        std::vector<std::string> terms = split(p,std::string(" ")); 
        for(std::string x : terms){
            T.push_back(p);
            allderivelambda = this->derivesToLambda(x,T);
            T.pop_back();
            if(!allderivelambda) break;
        }
        if(allderivelambda) return true;

    }
    return false;
}

CFG::xb makeXB(std::string s){
    CFG::xb XB = {};
    if(s.length() == 0){
        XB.X = ""; 
        XB.B = "";
    }
    
    else if(s.length() == 1){
        XB.X = s[0]; 
        XB.B = "";
    }
    else{
        XB.X = s[0];
        std::string b;
        for(char c : s){
            b+=c;
        }
        XB.B = b;
    }
    return XB;

}


std::pair<std::set<std::string>,std::set<std::string>> CFG::firstSet(CFG::xb XB, std::set<std::string> T){

    if(!this->_isNonTerminal(XB.X)){
        std::set<std::string> x;
        x.insert(XB.X);
        return {x,T};
    }

    std::set<std::string> F;
    if(T.count(XB.X)){
        T.insert(XB.X);
        for(std::string p : this->cfg.at(XB.X)){
            auto G = this->firstSet(makeXB(p),T);
            F = setUnion(F,G.first);
        }
    }

    if(this->derivesToLambda(XB.X,std::vector<std::string>()) && XB.B.length()>0){
        auto G = this->firstSet(makeXB(XB.B),T);
        F = setUnion(F,G.first);
    }
    return {F,T};
}
