import copy
from sys import argv
from collections import deque


class Row():

    def __init__(self,S,isStart,isAccept):
        self.alphabet = {}
        self.S = S
        self.isStart = isStart
        self.isAccept = isAccept
        
    def add_alphabet(self,c,state):
        self.alphabet[c] = state
    
    def __str__(self):
        return f"isStart:{self.isStart} isAcc:{self.isAccept} State:{self.S} Alphabet:{self.alphabet}"
        

def char_transitions(state,transitions,c):
    res = set()
    for k,v in transitions[state].items():
        for i in v:
            if i == c:
                res.add(k)
    print(f"char trans: {state} -> {list(res)}")
    return list(res)


def follow_char(S,transitions,c):
    F = set()
    for t in S:
        for q in char_transitions(t,transitions,c):
            F.add(q)
    return list(F)

# returns list of lambda transitions for a state or [] if None
def lambda_transitions(l,state,transitions):
    res = set()
    for k,v in transitions[state].items():
        for i in v:
            if i == l:
                res.add(k)
    print(f"trans: {state} -> {list(res)}")
    return list(res)




def follow_lambda(l, S, transitions):
    M = deque()
    for t in S:
        M.append(t)
    
    print(M)
    
    while(M):
        t = M.pop()
        print(M)
        for q in lambda_transitions(l,t,transitions):
            print(f"q:{q}")
            if q not in S:
                S.append(q)
                M.append(q)
        
    return S

def intersectA(A, B):
    for i in B:
        if i in A:
            return True
    return False  


if __name__ == "__main__":

    if len(argv) < 2:
        print("Usage: python3 nfa2dfa.py filepath -optional:outfile")
        exit(-1)
    
    filepath = argv[1]

    nfa = ""
    whitespace = lambda w: w.split(" ")
    with open(filepath,"r+") as f:
        nfa = list(map(whitespace, f.read().split('\n')))
    
        header = nfa[0]
        nfa = nfa[1:]
    
    print(header)

    
    NUM_STATES = int(header[0])
    LAMBDA = header[1]
    print(NUM_STATES)
    print(LAMBDA)
    SIGMA = header [2:]
    print(SIGMA)
    START_STATE = '0' # can assume 0 is always the start state...
    
    for i in nfa:
        if len(i) == 3:
            i.append(None)
    print(nfa)
    transitions = {str(i):{} for i in range(NUM_STATES)}
    
    accept = lambda x: True if x == "+" else False
    isAccept = {i[1]: accept(i[0]) for i in nfa}
    
    print(isAccept)
    print()
    for t in nfa:

        transitions[t[1]][t[2]] = t[3:]
    print(transitions)
    print()
    print()

    T = []

    L = deque()
    A = list(set([i for i in isAccept if isAccept[i]]))
    print(A)


    B = follow_lambda(LAMBDA,[START_STATE],transitions)
    
    firstRow = Row(B,True,intersectA(A,B))
    T.append(firstRow)
    print(firstRow)
    L.append(B)


    # C = follow_char([START_STATE],transitions,'/')
    # print(C)
    # while(L):
    #     S = L.pop()
    #     for c in SIGMA:




    




    
    
