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


class CFG{
    public:
        std::unordered_map<std::string,std::vector<std::string>> cfg;
        std::vector<std::string> file_data;
        std::string start_state;
        std::set<std::string> non_terminals;
        std::string raw_file;
        CFG(std::vector<std::string> file_data, std::string raw_data){
            this->raw_file = raw_data;
            this->file_data = file_data;
            this->start_state = this->_getStartState();
            this->non_terminals = this->_getNonTerminals();
            this->_makeCFG();
        }
        ~CFG();

        void print_cfg(void);
        private:
            
            bool _isNonTerminal(std::string s);
            std::string _getStartState(void);
            std::set<std::string> _getNonTerminals(void);
            void _makeCFG(void);
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


void CFG::_makeCFG(void){
    
    std::vector<std::pair<std::string,std::string>> cfg;
    std::string curr;
    for (auto s : this->file_data){
        std::vector<std::string> f = split(s,std::string(" -> "));
        if(f.size()==2){
            curr = f[0];
            std::vector<std::string> bars = split(f[1],std::string(" | "));
            for(std::string b : bars){
                cfg.push_back(std::make_pair(curr,b));
            }
            
        }
        else{
            std::string x = f[0].substr(2,f[0].length()-2);
            
            std::vector<std::string> br = split(x,std::string(" | "));
            for(std::string t : br){
                cfg.push_back(std::make_pair(curr,t));
            }
        }
    }
    for(auto i : cfg){
        std::cout << i.first << " -> " << i.second << std::endl;
    }
    

    // std::ordered_map<std::string,std::vector<std::string>> m;
    // for (std::string s : this->non_terminals) {
    //     m[s]=std::vector<std::string>();
    // }
    // int ind = 0;
    // std::string curr = this->start_state;
    // for(std::string s : this->file_data){
    //     std::vector<std::string> x = split(s,std::string(" -> "));
    //     ind = 0;
    //     if(this->non_terminals.count(x[0])){
    //         curr = x[0];
    //         ind = 1;
    //     }   
        
    //     m[curr].push_back(x[ind]);
        
    }
    
    
    // return m;


void CFG::print_cfg(void){
    std::cout << this->start_state << std::endl;
    std::cout << "Nonterminals: { ";
    
     
    for(auto s : this->non_terminals){
        std::cout << s << " ";
    } std::cout << "}" << std::endl;


    

    
    
}


