#include "luther.hh"

int main(int argc, char* argv[]){

    if(argc < 2){
        perror("Usage: ./luther file.cfg");
        exit(1);
    }
    std::cout << argv[1] << std::endl;
    std::ifstream file(argv[1]);
    std::string line, raw_data;
    std::vector<std::string> lines;
    while(std::getline(file,line)){
        lines.push_back(line);
        raw_data+=line+'\n'; 
    } 
    
    file.close();

    CFG cfg = CFG(lines); 

    for(std::string s : cfg.file_data){
        std::cout << s << std:: endl;
    }
    printf("\n------------\n");
    cfg.print_cfg();
    printf("\n------------\n");
    for(auto s : cfg.non_terminals){
        std::cout << s << ": derives2lambda -> ";
        std::cout << cfg.derivesToLambda(s,std::vector<std::string>()) << std::endl;
    }

    return 0;
}