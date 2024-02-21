
#include <stdio.h>
#include <iostream>
#include <fstream>
#include <string>
#include "luther.hh"

int main(int argc, char* argv[]){

    if(argc < 2){
        perror("Usage: ./luther file.cfg");
        exit(1);
    }
    std::cout << argv[1] << std::endl;
    std::ifstream file(argv[1]);
    std::string line;
    std::vector<std::string> lines;
    while(std::getline(file,line)) lines.push_back(line); 
    
    file.close();

    CFG cfg = CFG(lines); 

    for(std::string s : cfg.file_data){
        std::cout << s << std:: endl;
    }

    cfg.print_cfg();


    return 0;
}