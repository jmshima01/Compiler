use std::env;
use std::fs;
use std::collections::HashMap;
use regex::Regex;
use std::collections::VecDeque;
fn main() { // DFA OPTIMIZE...
    let mut T: HashMap<&str, VecDeque<&str>> = HashMap::new();

    let args: Vec<String> = env::args().collect();
    
    let file_path : &String =  &args[1];

    let file : String = fs::read_to_string(&file_path).expect("File err");
    let l: Vec<&str> = file.lines().collect();
    
    dbg!("content: \n{}",&l);
    
    let re = Regex::new(r" +").unwrap();
    for i in l{
        let v: Vec<&str> = re.split(i).collect();
        let mut q: VecDeque<&str> = VecDeque::from(v);
        let key = q[1];
        q.pop_front();
        q.pop_front();
        T.insert(key, q);
    }
    
    
    dbg!("{}",&T);
    


    
}
