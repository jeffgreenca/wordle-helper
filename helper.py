#!/usr/bin/env python3

def load(target, m=set()):
    with open(target, "r") as f:
        for line in f:
            m.add(line.strip())
    return m

def rank(target):
    d = {}
    index = 0
    with open(target, "r") as f:
        for line in f:
            d[line.strip()] = index
            index += 1
    return d

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

def _prompt(prefix, s):
    return f"{prefix} ({''.join(s)}): "

if __name__ == "__main__":
    m = load("5-unsorted.txt")
    m = load("5-ordered.txt", m)
    r = rank("5-ordered.txt")

    excluded = set()
    required = set()
    badpos = set()
    while True:
        excluded = excluded.union(set(input(_prompt("excluded", excluded))))
        print(f"> excluding {''.join(excluded)}")
        goodpos = list(input("good positions: ").strip().ljust(5, "_"))
        print(f"> must match: {''.join(goodpos)}")
        bp = input("bad positions: ").strip().ljust(5, "_")
        badpos.add(bp)
        for b in badpos:
            print(f"> cannot match: {b}")

        # compute required letters from good position and bad positions
        for c in goodpos:
            if c != "_":
                required.add(c)
        for b in badpos:
            for c in b:
                if c != "_":
                    required.add(c)
        print(f"> requiring  {''.join(required)}")
        print()

        tested = 0
        found = 0
        results = []
        for n in m:
            tested += 1
            if not check_positions(goodpos, n):
                continue
            def isbad(b):
                return not check_badpositions(list(b), n)
            if any(map(isbad, badpos)):
                continue
            nset = set(n)
            if not check_required(required, nset):
                continue
            if not check_excluded(excluded, nset):
                continue
            found += 1
            results.append(n)

        # format output
        def _ranker(n):
            if n in r:
                return r[n]
            return 50000

        cols = 14

        index = 1
        for n in sorted(results, key=_ranker, reverse=False):
            e = "\n" if index % cols == 0 else "\t"
            if n in r:
                n = n+"*"
            print(n, end=e)
            index += 1

        print("\n[%d/%d]\n" % (found, tested))
