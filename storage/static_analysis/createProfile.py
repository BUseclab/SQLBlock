#!/usr/bin/python3
import sys
import json 
profiledic = {}
classes = set()
numchecked = 0
FP = 0
logcheck = False
logprofile = False
f = open("log","w")
fpCallstacks = set()
totalCallstacks = set()
foundset = set()
notfoundset = set()
def processCallstack(line):
    lines = line.strip('\n').split("#")
    if len(lines) == 1:
        return None
    callstack = line.strip('\n').strip(' ').split("#")[-1]
    calls = callstack.split("@")
    if len(calls) == 1:
        return None
    stack = ""
    for call in calls:
        # append mysqli_query to the beginning
        if "mysqli_query" in call:
            continue
        # append classes 
        if "::" in call:
            procls = call.strip(" ").strip("\n").split("::")[0]
            if procls in classes:
                continue
        elif call.strip(" ") in classes:
                continue
        stack = call.strip(" ").strip("\n")
        break
    if logprofile == True:
        f.write("========================\n")
#        f.write("classes are %s\n"%(classes))
        f.write("the callstack is: %s\n"%(calls))
        f.write("the function is: %s\n"%(stack))
    
    if stack == "": ## if stack is empty this means the last database method communicate with db
        stack = calls[-1].strip(" ").strip("\n")
    if callstack == "":
        return None
#    print("[%s] changes to [%s]" %(callstack, stack))
    return stack.strip(" ").strip("\n")

def processfuncs(line):
    funcList = set()
    items = line.strip('\n').split("#")
    andcond = 0
    orcond = 0
    cond = -1 # variable
    ## funcs
    for i in range(len(items)):
        if items[i] == "FIELD":
            continue
        ## processing a function
        if "FUNC:" in items[i]:
            funcName = items[i][5:]
            numargs = int(items[i+1])
            leftarg = 1
            rightarg = 2
            if "FIELD" in items[i+2]: ## first arg to the func
                leftarg = 1
            else:
                leftarg = 0
#            for j in range(numargs-1):
#                print(len(items),i+j+3)
#                print(str(items))
            if len(items) <= (i+numargs+2):
                break
            if "FIELD" in items[i+3:i+numargs+2] and "LITERAL" in items[i+3:i+numargs+2]:
                rightarg = 2 # var
            elif "LITERAL" in items[i+3:i+numargs+2]:
                rightarg = 0 # literal
            elif "FIELD" in items[i+3:i+numargs+2]:
                 rightarg = 1 # field
                
#                if len(items) <= (i+j+3):
#                    break
#                if "FIELD" not in items[i+j+3]: ## starting from the second arg to the func
#                    rightarg = 0
#                    break
            funcList.add((funcName,leftarg,rightarg))
            i = i+ numargs + 2 ## point to after funcName,numarg and args of function
            continue
        ## conds
        ## THIS IS WRONG
        if "COND:" in items[i]:
            if "and" == (items[i][5:]).lower():
                andcond = andcond + 1
            if "or" == (items[i][5:]).lower():
                orcond = orcond + 1
    if andcond >= 1 and orcond == 0: ## if all conds are AND
        cond = 1
    if orcond >= 1 and andcond == 0: ## if all conds are OR
        cond = 0
            ## o.w. we can't conclude anything from conds so we let it to be -1
    return (funcList, cond)
 
def processtab(line):
    tabacc = line.strip('\n').split("##")
    if len(tabacc) != 2:
        return (None,None)
    return(tabacc[0],tabacc[1])

def processEntry(lines):
    try:
#        print(lines)
        callstack = processCallstack(lines[0])
        tabaccList = []
        funcList = set()
        cond = -1
#        print(callstack)
        if callstack == None:
            return
        if len(lines) == 1:
            return
        else:
            (funcList, cond) = processfuncs(lines[1])
#            print ("func, cond: %s,%s" %(funcList,cond))
            for i in range(2,len(lines)):
                (table,action) = processtab(lines[i])
                if table != None and action != None:
                    tabaccList.append((action,table))
#                    print ("table, action : %s,%s" %(table,action))

    except:
        print("something happened while extracting data::ignoring them")
        return
    if logprofile == True:
        f.write("==========================\n")
        f.write("query is [%s]\n"%lines[0].strip("\n"))
        f.write("the function is %s\n"%(callstack))
        f.write("adding tables %s\n"%str((tabaccList)))
        f.write("adding func %s\n"%str(list(funcList)))
        f.flush()
        
    item = [list(tabaccList), cond, list(funcList)]
    if callstack in profiledic.keys():
        items = profiledic[callstack]
        if item not in profiledic[callstack]:
           profiledic[callstack].append(item)
        else:
            if logprofile == True:
                f.write("it was already in the list\n")
                f.flush()
    else:
        profiledic[callstack] = [item]



# read db-related functions and classes for specific web app
with open(sys.argv[3]) as dbfile:
    alllines = dbfile.readlines()
    for line in alllines:
        if line.strip("\n") == "":
            continue
        classes.add(line.strip('\n'))
#    for item in classes:
#        print(item)

# read infoFile
with open(sys.argv[1],'rb') as profile:
    lines = []
#    alllines = profile.readlines()
#    del alllines[0]
    for line in profile:
        try:
#    for line in alllines:
            if "####" not in str(line):
                if line.decode('utf8', errors='ignore') != "":
                    lines.append(line.decode('utf8', errors='ignore'))
            else:
                if len(lines) != 0 :
#                    print("========================")
#                    print("lines to process")
#                    print(lines)
                    processEntry(lines)
                    lines.clear()
                else:
                    continue
        except UnicodeDecodeError as e :
            continue
#            print("ERROR IS HERE: ["+str(line).strip('\n')+"]")
#  print("error")

with open(sys.argv[2], 'w') as outFile:
#    print(profiledic)
    for key in profiledic:
        items = profiledic[key]
        for item in items:
            outFile.write(key+'\n')
            # write funcList
            outFile.write("##".join("@@".join(str(t) for t in x) for x in item[0])+"##\n")
            outFile.write(str(item[1])+"\n")
            outFile.write("##".join("@@".join(str(t) for t in x) for x in item[2])+"##\n")
            outFile.write("####\n")

if logprofile == True:
    f.close()
