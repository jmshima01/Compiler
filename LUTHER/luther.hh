#include <bits/stdc++.h> 
#include <set>

std::vector<std::string> split(std::string s, std::string delimiter){

    size_t pos_start = 0, pos_end, delim_len = delimiter.length();
    std::string token;
    std::vector<std::string> res;

    while ((pos_end = s.find(delimiter, pos_start)) != std::string::npos) {
        token = s.substr (pos_start, pos_end - pos_start);
        pos_start = pos_end + delim_len;
        res.push_back (token);
    }

    res.push_back (s.substr (pos_start));
    return res;
}




class CFG{
    public:
        std::unordered_map<std::string,std::vector<std::string>>* cfg;
        std::vector<std::string> file_data;
        CFG(std::vector<std::string> file_data){
            this->cfg = __makeCFG();
            this->file_data = file_data;
        }
        ~CFG();

        void print_cfg(void);
        std::string getStartState(void);
        std::set<std::string> getNonTerminals(void);
        
        private:
            std::unordered_map<std::string,std::vector<std::string>>* __makeCFG(void);
            bool __isNonTerminal(std::string s);
};

CFG::~CFG(){ delete this->cfg;};

bool CFG::__isNonTerminal(std::string s){
    for(char c : s){
        if (isupper(c)){
            return true;
        }
    }
    return false;
}


std::string CFG::getStartState(void){
    std::string delimiter = " -> ";
    return this->file_data[0].substr(0, this->file_data[0].find(delimiter)); 
}

std::set<std::string> CFG::getNonTerminals(void){
    std::set<std::string> nonterms;
    for(std::string line : this->file_data){
        std::vector<std::string> split_line = split(line,std::string(" "));
        for(std::string s : split_line){
            if(this->__isNonTerminal(s)){
                nonterms.insert(s);
            }
        }
    }
}


std::unordered_map<std::string,std::vector<std::string>>* CFG::__makeCFG(void){
    
    std::unordered_map<std::string,std::vector<std::string>>* m = new std::unordered_map<std::string,std::vector<std::string>>();
    
    // for(int i = 0; i<this->file_data.size(); i++){
    //     m[]
        
    // }
    
    
    
    
    return m;

}

void CFG::print_cfg(void){
    
}


