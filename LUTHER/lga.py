from sys import argv
import re
def isNonterminal(s):
    for i in s:
        if i.isupper():
            return True
    return False

isSymbol = lambda s: s != "->" and s != "|"

if __name__ == "__main__":
    with open(argv[1], 'r+') as f:
        cfg = f.read().split('\n')
    
    nonTeminals = set()
    vocab = set()
    START_SYMBOL = ""
    first = True
    for row in cfg:
        for i,s in enumerate(row.split(" ")):
            if isNonterminal(s):
                if i == 0 and first:
                    START_SYMBOL = s
                    first = False
                nonTeminals.add(s)
            if isSymbol(s):
                vocab.add(s)
    
    prod_rules = {k:[] for k in nonTeminals}

    curr_symbol = START_SYMBOL
    
    for i in cfg:
        ind = 0
        i = i.split(" -> ")

        if i[0] in nonTeminals:
            curr_symbol = i[0]
            ind = 1

        prod_rules[curr_symbol].append(i[ind])
    
    for k,v in prod_rules.items():
        for i in range(len(v)):
            prod_rules[k][i] = [[i] for i in re.split(r"\s*\|\s*",v[i]) if i != '']
    
    # print()
    # print()
    # print(prod_rules)
            
    print("Grammar Non-Terminals")        
    print(nonTeminals)
    
    print("Vocab")
    print(vocab)
    print()
    
    print("Grammar Rules")
    num = 1
    for k,v in prod_rules.items():
        for i in v:
            for j in i:
                print(f"({num})  {k} -> {j}")
                num+=1
    print()
    print(f"Grammar Start Symbol or Goal: {START_SYMBOL}")
    print()
    print(prod_rules)
