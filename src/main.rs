use std::env;
use std::fs;



fn main() {
    
    let args: Vec<String> = env::args().collect();
    
    let file_path =  &args[1];

    let file = fs::read_to_string(&file_path).expect("File err");

    println!("content: \n{}",&file);




}
