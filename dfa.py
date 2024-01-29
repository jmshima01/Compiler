from sys import argv
from collections import deque
import copy

class DFA():
    def __init__(self, accept, non_accept, sigma):
        self.T = {"+":accept, '-':non_accept}
        self.accept = accept
        self.non_accept = non_accept
        self.sigma=sigma
    
    def print_T(self):
        print()
        print(self.sigma)
        for k,v in self.T['-'].items():
            print(f"{k} : {v}")
        print("----")
        for k,v in self.T['+'].items():
            print(f"{k} : {v}")
        print()

    def to_file(self,filename):
        with open(filename,"w+") as f: 
            for k,v in self.T['-'].items():
                f.write("- " + k + " " + " ".join(v)+"\n")
            i = 0
            for k,v in self.T['+'].items():
                i+=1
                if i==len(self.T['+'].items()):
                    f.write("+ " + k + " " + " ".join(v))
                else:
                    f.write("+ " + k + " " + " ".join(v)+"\n")
    


def update_del(s,new_s,dfa):
    for k,lis in dfa.T['-'].items():
        for i in range(len(lis)):
            if lis[i]==s:
                dfa.T['-'][k][i]=new_s
    
    for k,lis in dfa.T['+'].items():
        for i in range(len(lis)):
            if lis[i]==s:
                dfa.T['+'][k][i]=new_s
    
def merge_row(s,dfa):
    if s == []:
        print("err: s_e")
        exit(-1)
    s = s[::-1]
    for i in range(len(s)-1):
        if s[i] in dfa.accept:
            del dfa.T['+'][s[i]]
        else:
            del dfa.T['-'][s[i]]
        
        update_del(s[i],s[-1],dfa)

def partition(S,dfa,c):
    t = '-'
    if S == []:
        return []
    if S[0] in dfa.accept.keys():
        t = '+'
    res = {}
    
    for s in S:
        if dfa.T[t][s][dfa.sigma.index(c)] in res.keys():
            res[dfa.T[t][s][dfa.sigma.index(c)]].append(s)
        else:
            res[dfa.T[t][s][dfa.sigma.index(c)]] = [s]
    print(f"S: {list(res.values())}\n")
    return list(res.values())

def merge_states(dfa : DFA) -> DFA:
        L = deque()
        M = []
        n_sigma = copy.copy(dfa.sigma)
        a_sigma = copy.copy(dfa.sigma)
        L.append((list(dfa.T['+'].keys()),a_sigma))
        L.append((list(dfa.T['-'].keys()),n_sigma))
        
        while L:
            print("L:",list(L)[::-1])
            print("M:",M)
            print()
            S,C = L.pop()
            C = copy.copy(C)
            c = ''
            if len(C)!=0:
                c = C.pop(0)
            sets = partition(S,dfa,c)
            for x_i in sets:
                if(not (len(x_i)>1)):
                    continue
                if C==[]:
                    M.append(x_i)
                else:
                    L.append((x_i,C))
        print(M)
        print(L)
        for s in M:
            merge_row(s,dfa)
        
        return dfa

def read_file(file_path,sigma) -> DFA:
    data = []
    with open(file_path, "r+") as f:
        data = f.read().split("\n")

    accepting = {}
    nonacc = {}
    for i in data:
        i = i.split(" ")
        if i[0] == '+':
            accepting[i[1]] = i[2:]
        else:
            nonacc[i[1]] = i[2:]
    
    
    return DFA(accept=accepting,non_accept=nonacc,sigma=sigma)

def equal_tables(d1,d2) -> bool:
    return (len(d1.T['-'])+len(d1.T['+'])) == (len(d2.T['-'])+len(d2.T['+']))
        
if (__name__ == "__main__"):
    if len(argv) !=2:
        exit(-1) 

    sigma = ['a','b','c','q','r','s','t','u','v']
    D = read_file(argv[1],sigma)
    original = copy.deepcopy(D)
    
    while 1:
        D_ = merge_states(D)
        if equal_tables(D_ ,original):
            break
        original = copy.deepcopy(D_)
    
    D = D_

    D.print_T()
    D.to_file(argv[1].strip(".txt") + "-optimzed.txt")

