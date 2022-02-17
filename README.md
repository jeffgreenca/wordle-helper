# wordle helper

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
