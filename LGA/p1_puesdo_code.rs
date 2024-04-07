fn SDTParse(Stack Queue currentNode){

    if Stack.pop() is End of Production aka `*`{
        match currentNode with{
            | case Nucleus -> currentNode = handleNucleus(currentNode)
            | case Atom -> currentNode = handleAtom(currentNode) 
            | case Seq -> currentNode = handleSeq(currentNode)
            | case SeqList -> currentNode = handleSeqList(currentNode)
            | case Alt -> currentNode =handleAlt(currentNode)
            | case AltList -> currentNode = handleAltList(currentNode)
            | case Re -> currentNode = handleRE(currentNode)
            | _ -> return currentNode
        }
    }

}

// IF ANY OF THESE FAIL OR CHECK FOR CORRECT SYNTAX, THEN ITS A SYNTAX ERROR!


fn handleNucleus(nucleus){
    // same as LGA assignment
}

fn handleAtom(atom){
    if atom.AtomMod.chr is lambda{
        return atom.Nucleus.children
    }

    if atom.AtomMod.chr is kleene{
        let head = atom.AtomMod.chr
        head.child = atom.Nucleus.children
        return head
    }

    if atom.AtomMod.chr is plus{
        let head = atom.AtomMod.chr
        head.child = atom.Nucleus.children
        return head
    }
}

fn handleSeq(Seq){
    if Seq -> lambda{
        return lambda 
    }

    else{
        return Seq.Atom
    }

}

fn handleSeqList(SeqList){
    if SeqList -> lambda{
        return SeqList.parent
    }

    else{
        return SeqList.children
    }

}

fn handleAltList(AltList){
    if AltList is just lambda{
        return AltList.parent
    }

    if AltList.children is pipe children{
        return AltList.children except pipe
    }

}

fn handleAlt(Alt){
    return Alt.children
}

fn handleRE(RE){
    return RE.children
}