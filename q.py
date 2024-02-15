with open("q.txt") as f:
    d = f.read().split('\n')
    t = ""
    for i in d:
        for j in i:
            t+=j
        t+=' '
    
    t = t.split(' ')[:-1]
    print(t)
    t = list(map(int,t))
    r = []
    for i in range(1,57):
        if i not in t:
            r.append(i)
    print(r)
