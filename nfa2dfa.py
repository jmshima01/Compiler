from copy import copy,deepcopy
from sys import argv
from collections import deque


def get_transitions(c,state,N): # returns a set of trans_states where c is the trans char for a given state
    trans_set = set()
    for trans_state,t in N[state].items():
        for i in t:
            if i == c:
                trans_set.add(trans_state)
    print(f"trans_set:{list(trans_set)}")
    return trans_set


def follow_char(c,states,N): # returns a set of transitions from all the states that have transition c
    F = set()
    for t in states:
        for q in get_transitions(c,t,N):
            F.add(q)
    return F

def follow_lambda(LAMBDA, states, N): # returns condensed set of combined states to remove lambda transitions using a DFS
    S = states
    M = deque()
    
    for t in S:
        M.append(t)
    print()
    print(f"M:{list(M)}")

    while(M):
        t = M.pop()
        for q in get_transitions(LAMBDA,t,N):
            if q not in S:
                S.add(q)
                M.append(q)
    print(f"S after follow_lambda() -> {S}")
    return S


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
    
    print(f"header:{header}")
    print(f"nfa:{nfa}")
    
    # parse header:
    NUM_STATES = int(header[0])
    LAMBDA = header[1]
    
    print(f"Lambda={LAMBDA}")
    ALPHABET = header[2:]
    print(f"Alphabet:{ALPHABET}")
    START_STATE = '0' # can assume 0 is always the start state...
    
    A=set() # set of accepting states
    for row in nfa:
        if row[0] == '+':
            A.add(row[1])
    
    print(f"A:{A}")
    N = {str(i):{} for i in range(NUM_STATES)} #NFA file in nice dict format

    for row in nfa:
        state = row[1]
        transition_state = row[2]
        if len(row) < 4: # No transition char on this line of the file i.e dead(if not finish state) or finish state
            transition_chars = []
        else:
            transition_chars = row[3:]

        N[state][transition_state] = transition_chars
    
    print(f"\n\nN:{N}")


    T={} # Transition Function aka DFA
    L = deque()

    B = tuple(follow_lambda(LAMBDA,set(START_STATE),N))
    DFA_START_STATE = B
    DFA_ACCEPT_STATES = set() # lookup for accept states of Transition function
    T[B] = {}
    
    if len(A.intersection(set(B))) != 0:
        DFA_ACCEPT_STATES.add(B)
    print(f"accp: {DFA_ACCEPT_STATES}")

    L.append(B)
    while(L):
        S = L.pop()
        for c in ALPHABET:
            R = tuple(follow_lambda(LAMBDA,follow_char(c,set(S),N),N))
            if len(R)!=0: # save space in T (ignoring Save states)
                T[S][c] = R
            if (len(R)>0) and (R not in T.keys()):
                T[R] = {}
                if len(A.intersection(set(R))) != 0:
                    DFA_ACCEPT_STATES.add(R)
                L.append(R)
    
    print()
    print(T)

    with open("nfa2dfa.txt","w+") as f:
        for k,v in T.items():
            pass
    print(DFA_ACCEPT_STATES)

    new_keys = {}
    for i,k in enumerate(T.keys()):
        new_keys[k] = str(i)
    print(new_keys)

    
    def convert_old_keys(new_keys, row):
        n_row = {}
        for k,v in row.items():
            n_row[k] = new_keys[v]
        
        return n_row


    T_final = {k:{} for k in new_keys.values()}
    # print()
    # print(T_final)
    for k,v in T.items():
        # print("val:",v)
        # print("k:",k)
        T_final[new_keys[k]] = convert_old_keys(new_keys,v)
    print()
    print()
    print(T)
    print("===============")
    print(T_final)

    new_accept = set()
    for i in DFA_ACCEPT_STATES:
        new_accept.add(new_keys[i])
    print(new_accept)

    def convert_row_to_dfa_list(row,alphabet):
        row_to_write = {i:'E' for i in ALPHABET}
        for k,v in row.items():
            row_to_write[k] = v
        # print(row_to_write)
        return " ".join(row_to_write.values())


    print("++++++++++")
    a = " ".join(ALPHABET)
    print("    " + a)
    with open("nfa2dfa.txt","w+") as f:
        line = ""
        for k,v in T_final.items():
            line = ""
            if k in new_accept:
                line+= "+ " + k + " "
            else:
                line+="- " + k + " "
            print(line + convert_row_to_dfa_list(v,ALPHABET))
            if k != list(T_final.keys())[-1]:
                f.write(line + convert_row_to_dfa_list(v,ALPHABET) + "\n")
            else:
                f.write(line + convert_row_to_dfa_list(v,ALPHABET))

    
    
