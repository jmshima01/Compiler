use std::env;
use std::fs;
use std::fmt;


fn main() {
    
    let args: Vec<String> = env::args().collect();
    
    let file_path = &args[1];
    println!("file_path {}",file_path);

    let file = fs::read_to_string(file_path).expect("File err");

    println!("content: {}",file)

}
