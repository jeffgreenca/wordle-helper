#!/usr/bin/env python3

def load():
    m = set()
    with open("words_alpha.txt", "r") as f:
        for line in f:
            line = line.strip()
            if len(line) == 5:
                m.add(line)
    return m

def check_positions(p, n):
    for i in range(5):
        if p[i] == "_":
            continue
        if p[i] == n[i]:
            continue
        return False
    return True

def check_badpositions(p, n):
    for i in range(5):
        if p[i] == "_":
            continue
        if p[i] == n[i]:
            return False
    return True

def check_required(r, n):
    for c in r:
        if c not in n:
            return False
    return True

def check_excluded(e, n):
    for c in n:
        if c in e:
            return False
    return True

if __name__ == "__main__":
    m = load()

    while True:
        excluded = set(input("tried: "))
        required = set(input("required: "))
        goodpos = list(input("good positions: "))
        badpos = list(input("bad positions: "))

        tested = 0
        found = 0
        for n in m:
            tested += 1
            if not check_positions(goodpos, n):
                continue
            if not check_badpositions(badpos, n):
                continue
            nset = set(n)
            if not check_required(required, nset):
                continue
            if not check_excluded(excluded, nset):
                continue
            found += 1
            print(n)
        print("[%d/%d]" % (found, tested))
