package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	NotInWord       = 'n'
	WrongPosition   = '*'
	CorrectPosition = 'y'
)

const wordsListFile = "5_most_common.txt"

var words []string

func main() {
	fmt.Printf("Clever Helper v0.1\n\n")

	words = loadWords(wordsListFile)

	for {
		do(words)
	}
}

func do(words []string) {
	rules := RuleSet{}

	fmt.Printf("Reference\n")
	fmt.Printf("  %c = miss, not in word\n", NotInWord)
	fmt.Printf("  %c = correct letter, wrong position\n", WrongPosition)
	fmt.Printf("  %c = correct letter and position\n", CorrectPosition)

	for {
		attempt, restart := ui("Attempt: ")
		if restart {
			break
		}
		result, _ := ui("Result: ")

		newRules := generateRulesFromAttempt(attempt, result)
		rules = append(rules, newRules...)
		rules.Optimize()

		found := rules.FindWords(words, -1)

		fmt.Printf("\nPossible words: \n")
		printList("  ", found, 10)

		fmt.Printf("\nSuggested next guess:\n")
		guesses := suggest(truncate(found, 100), words)
		printList("  ", guesses, 5)

		fmt.Println()
	}
}

func ui(prompt string) (string, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(prompt)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	if s == "new" {
		return "", true
	}
	if len(s) != 5 {
		fmt.Println("...please enter exactly 5 characters")
		return ui(prompt)
	}
	return strings.ToLower(s), false
}

func printList(prefix string, list []string, limit int) {
	for i, s := range list {
		fmt.Printf("%s%s\n", prefix, s)
		if i == limit {
			break
		}
	}
}

func truncate(s []string, limit int) []string {
	if len(s) > limit {
		return s[:limit]
	} else {
		return s
	}
}

func score(rankedHitWords []string, newHitWords []string) int {
	// for each rankedHitWord that was eliminated (not in newHitWords), we get points equal to its inverse rank
	// but if all rankedHitWords are eliminated, we get 0 points

	score := 0

	ranked := truncate(rankedHitWords, 100)

	m := map[string]bool{}
	for _, word := range newHitWords {
		m[word] = true
	}
	allEliminated := true
	for i, word := range ranked {
		if !m[word] {
			score += len(ranked) - i
		} else {
			allEliminated = false
		}
	}
	if allEliminated {
		return 0
	}
	return score
}

func suggest(hitWords, allWords []string) []string {
	results := []string{}

	// brute force
	// step 1 - finding and ranking rules
	// for each hitword,
	//   for each letter,
	//     if we had a rule that said "this letter is not at this position",
	//       and we applied that rule to all the hitwords,
	//       what would be the score?
	type ruleScore struct {
		Rule  Rule
		Score int
	}
	scores := []ruleScore{}
	for _, hitWord := range hitWords {
		for i, c := range hitWord {
			rule := Rule{Type: NotAtPosition, Position: i, Character: byte(c)}
			/* this check is unnecessary since we are doing the NotAtPosition, if we knew it was NotAtPosition, we wouldn't see it in result word
			if currentRuleSet.Contains(rule) {
				continue
			}
			*/
			// only score each unique rule once
			for _, tried := range scores {
				if rule.Equal(tried.Rule) {
					continue
				}
			}
			rs := RuleSet{rule}
			newHitWords := rs.FindWords(hitWords, -1)
			scores = append(scores, ruleScore{Rule: rule, Score: score(hitWords, newHitWords)})
		}
	}
	log.Printf("found %d possible new rules", len(scores))

	// sort the rules based on score so we can try the best ones first
	rankedRules := []Rule{}
	for _, score := range scores {
		rankedRules = append(rankedRules, score.Rule)
	}
	sort.Slice(rankedRules, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if len(rankedRules) == 0 {
		log.Printf("can't suggest a word, no helpful rules found")
		return results
	}
	log.Printf("found %d helpful rules from %d possible rules", len(rankedRules), len(scores))

	// step 2 - finding the best word to try, which would maximize our chance
	// of learning positional data that is most helpful
	// we are going to invert the search. we want a negative match. so if we have the following rules
	// - a not at position 0
	// - b not at position 1
	// then we want a word that has a at position 0 and b at position 1

	// we start with adding our highest ranked rule to our ruleset
	filterRules := map[int]Rule{}
	filterRules[rankedRules[0].Position] = rankedRules[0]

	toRuleSet := func(m map[int]Rule) RuleSet {
		var rs RuleSet
		for _, rule := range m {
			inverted := Rule{Type: RequiredAtPosition, Position: rule.Position, Character: rule.Character}
			rs = append(rs, inverted)
		}
		return rs
	}

	// now we have two dimensions we could search in
	// 1. we could try to add more rules in other positions
	// 2. we could try to replace a rule in a current position with a different rule
	for {
		valid := toRuleSet(filterRules).FindWords(allWords, 2)
		if len(valid) > 0 {
			// insert this rule at head of results
			results = append(valid, results...)
		} else {
			// no valid words, so we can give up or try something else
			break
		}

		if len(filterRules) == 5 {
			// can't have more than 5 rules because only 5 positions
			break
		}

		// try to add next most valuable rule at a different position
		for _, rule := range rankedRules {
			if _, ok := filterRules[rule.Position]; ok {
				continue
			}
			filterRules[rule.Position] = rule
			break
		}
	}

	return results
}

func (rs RuleSet) Contains(r Rule) bool {
	for _, rule := range rs {
		if rule.Equal(r) {
			return true
		}
	}
	return false
}

func (r Rule) Equal(other Rule) bool {
	return r.Type == other.Type && r.Position == other.Position && r.Character == other.Character
}

func (rs RuleSet) Equal(other RuleSet) bool {
	if len(rs) != len(other) {
		return false
	}
	for i, r := range rs {
		if !r.Equal(other[i]) {
			return false
		}
	}
	return true
}

func loadWords(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	log.Printf("loaded %d words from %s", len(words), filename)

	return words
}

func generateRulesFromAttempt(attempt string, result string) RuleSet {
	var rules RuleSet

	repeatedLetters := map[byte]bool{}
	for _, c := range attempt {
		if strings.Count(attempt, string(c)) > 1 {
			repeatedLetters[byte(c)] = true
		}
	}

	for i, c := range result {
		switch c {
		case NotInWord:
			if repeatedLetters[byte(attempt[i])] {
				// defer adding rule for this character, as it might be too restrictive
				continue
			} else {
				rules = append(rules, Rule{Type: NotAnywhere, Character: attempt[i]})
			}
		case WrongPosition:
			rules = append(rules, Rule{Type: NotAtPosition, Position: i, Character: attempt[i]})
			rules = append(rules, Rule{Type: RequiredAnywhere, Character: attempt[i]})
		case CorrectPosition:
			rules = append(rules, Rule{Type: RequiredAtPosition, Position: i, Character: attempt[i]})
		}
	}

	for c := range repeatedLetters {
		firstPos := strings.IndexByte(attempt, c)
		lastPos := strings.LastIndexByte(attempt, c)
		firstResult := result[firstPos]
		lastResult := result[lastPos]

		if firstResult == NotInWord && lastResult == NotInWord {
			// if both come back with negative result, its not anywhere and we are done here
			rules = append(rules, Rule{Type: NotAnywhere, Character: c})
		} else {
			// could be one or both
			if (firstResult != NotInWord) && (lastResult != NotInWord) {
				// both, so AppearsTwice
				rules = append(rules, Rule{Type: Repeated, Character: c})
			} else {
				// one, so ExactlyOnce
				rules = append(rules, Rule{Type: ExactlyOnce, Character: c})
			}
			// the position specific rules were already added above, so we are done
		}
	}

	return rules
}

// rules:
// - ch required anywhere
// - ch required at position
// - ch not at position
// - ch not anywhere

type RuleType int

const (
	RequiredAnywhere RuleType = iota
	RequiredAtPosition
	NotAtPosition
	NotAnywhere
	Repeated // special case for words like `foods`
	ExactlyOnce
)

type Rule struct {
	Type      RuleType
	Position  int
	Character byte
}

// validate a rule against a string
func (r *Rule) validate(s string) bool {
	switch r.Type {
	case RequiredAnywhere:
		return strings.IndexByte(s, r.Character) != -1
	case RequiredAtPosition:
		return s[r.Position] == r.Character
	case NotAtPosition:
		return s[r.Position] != r.Character
	case NotAnywhere:
		return strings.IndexByte(s, r.Character) == -1
	case Repeated:
		return strings.Count(s, string(r.Character)) > 1
	case ExactlyOnce:
		return strings.Count(s, string(r.Character)) == 1
	}
	return false
}

type RuleSet []Rule

func (rs RuleSet) Validate(s string) bool {
	for _, r := range rs {
		if !r.validate(s) {
			return false
		}
	}
	return true
}

func (rs RuleSet) Optimize() {
	sort.Slice(rs, func(i, j int) bool {
		order := map[RuleType]int{
			RequiredAtPosition: 0,
			NotAtPosition:      1,
			RequiredAnywhere:   2,
			NotAnywhere:        3,
			Repeated:           4,
			ExactlyOnce:        5,
		}
		return order[rs[i].Type] < order[rs[j].Type]
	})
}

func (rs RuleSet) FindWords(words []string, maxHits int) []string {
	var result []string
	for _, word := range words {
		if rs.Validate(word) {
			result = append(result, word)
		}
		if len(result) == maxHits {
			break
		}
	}
	return result
}
