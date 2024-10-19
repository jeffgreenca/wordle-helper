# wordle helper

## newer (go) usage
Clone this repo. Make sure you have Go 1.20 or better.

```
$ go run ./clever.go
Clever Helper v0.1

2024/10/19 01:08:04 loaded 39933 words from 5_most_common.txt
Reference
  n = miss, not in word
  * = correct letter, wrong position
  y = correct letter and position
Attempt: bling
Result : nyn**

Possible words:      
  gluon
  glean
  glsen
  algan
  rlngs
  algun
  elgon
  glahn
  algen
  glnpo
  elgan

Suggested next guess:
  asean
  asean
  ahern
  again
  asian
  women

Attempt:
```

# older python version
## usage
```
$ python3 helper.py 
excluded (): ar
> excluding ar
good positions: w
> must match: w____
bad positions: __t
> cannot match: __t__
> requiring  tw

white*  westy   wlity   weent   wilts   whute   whewt   whets   whits   wiste   wists   wests   width   weety
wicht   wecht   weste   welts   whift   weets   wiyot   wevet   wisht   wysty   whity   wheft   whist   whipt
wefty   weest   wonts   whoot   wefts   wight   wootz
[35/15956]
```

## help

- excluded: letters not allowed (additive each loop)
- good positions: known correct letter positions, using underscore `_` as placeholder (replaced each loop)
- bad positions: required letters in wrong positions, using `_` placeholder (additive each loop)
- words marked with `*` are common words, presented in frequency ranked order
