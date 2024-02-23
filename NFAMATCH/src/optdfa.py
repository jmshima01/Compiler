from sys import argv
from collections import deque
import copy
import re

class DFA():
    def __init__(self, accept, non_accept, sigma):
        self.accept = accept
        self.non_accept = non_accept
        self.sigma=sigma
        self.T = {"+":self.accept, '-':self.non_accept}
    
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
    
    def reorder_keys(self):
        
        new_k = [str(i) for i in range(len(self.accept)+len(self.non_accept))]
        old_k = [i for i in self.non_accept] + [i for i in self.accept]
        
        new_keys = {old_k[i]:new_k[i] for i in range(len(new_k))}
        
        new_non_acc = {new_keys[k]:[] for k in self.non_accept}
        new_acc = {new_keys[k]:[] for k in self.accept}
        
        for k,v in self.non_accept.items():
            l = ['E']*len(v)
            for i in range(len(v)):
                if v[i] != 'E':
                    l[i] = new_keys[v[i]]
            new_non_acc[new_keys[k]] = l
        for k,v in self.accept.items():
            l = ['E']*len(v)
            for i in range(len(v)):
                if v[i] != 'E':
                    l[i] = new_keys[v[i]]
            new_acc[new_keys[k]] = l
        
        self.accept = new_acc
        self.non_accept = new_non_acc
        self.T = {"+":self.accept, '-':self.non_accept}

    def find_unreachables(self):
        print(self.accept)
        print(self.non_accept)
        combined = {**self.non_accept, **self.accept}
        print(combined)
    
        S = [combined['0']]
        seen = set('0')
        
        while(S):
            R = S.pop()
            for s in R:
                if s!='E' and s not in seen:
                    S.append(combined[s])
                    seen.add(s)
                    
            print(S)
        print(seen)
        all_states = set([i for i in combined])
        diff = list(all_states.difference(seen))
        print(diff)
        return diff

    def find_dead(self):
        dead = set()
        acc_set = set([i for i in self.accept])
        for k,v in self.non_accept.items():
            S = [v]
            seen = set(k)
            while(S):
                R = S.pop()
                for s in R:
                    if s!='E' and s not in seen:
                        if s in self.non_accept:
                            S.append(self.non_accept[s])
                        seen.add(s)

            # print(seen)
            if len(acc_set.intersection(seen))==0:
                dead.add(k)
        print(dead)
        return list(dead)
            
    def delete_states(self,states):
        print(":::::::::::::::::::::")
        S = set(states)
        new_acc = copy.deepcopy(self.accept)
        new_non_acc = copy.deepcopy(self.non_accept)
        for k,v in self.non_accept.items():
            if k in S:
                new_non_acc.pop(k)
                continue
            
            temp = ['E'] * len(v)
            for i in range(len(v)):
                if v[i] not in S or v[i]=='E':
                    temp[i] = v[i]
            
            new_non_acc[k] = temp
        for k,v in self.accept.items():
            if k in S:
                new_acc.pop(k)
                continue
            
            temp = ['E'] * len(v)
            for i in range(len(v)):
                if v[i] not in S or v[i]=='E':
                    temp[i] = v[i]
            
            new_acc[k] = temp
        print(new_non_acc)
        print(new_acc)
        self.non_accept = new_non_acc
        self.accept = new_acc
 
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
        i = re.split(r"\s+", i)
        if i[0] == '+':
            accepting[i[1]] = i[2:]
        else:
            nonacc[i[1]] = i[2:]
    
    return DFA(accepting,nonacc,sigma)

def equal_tables(d1,d2) -> bool:
    return (len(d1.T['-'])+len(d1.T['+'])) == (len(d2.T['-'])+len(d2.T['+']))


if __name__ == "__main__":
    if len(argv) != 3:
        print("usage: python3 dfa.py file sigma")
        exit(-1) 

    sigma = argv[2]
    sigma = [char for char in sigma]

    D = read_file(argv[1],sigma)
    original = copy.deepcopy(D)

    count = 0
    while 1:
        count+=1
        D_ = merge_states(D)
        if equal_tables(D_ ,original):
            break
        original = copy.deepcopy(D_)

    D = D_
    D.print_T()
    
    
    u = D.find_unreachables()
    dead = D.find_dead()
    
    D.delete_states(u+dead)
    D.reorder_keys()
    D.print_T()
    D.to_file(argv[1].strip(".txt") + "-optimzed.txt")
    print("iters:",count)