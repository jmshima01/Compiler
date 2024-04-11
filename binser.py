
def bins(l,item):
    print(l)
    
    if l == []:
        print(l)
        return None

    if l[len(l)//2] == item:
        print("get here")
        return item

    elif l[len(l)//2] > item:
        print("left side")
        return bins(l[:len(l)//2],item)
    else:
        print("right side")
        return bins(l[len(l)//2:],item)     


l = [1,2,3,5,6,7,9,10,15,16]
print(bins(l,3))