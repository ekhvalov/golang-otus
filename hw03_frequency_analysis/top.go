package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type str string

func (s *str) getTop(size uint) []string {
	m := s.countWords()
	words := make([]string, 0)
	for w := range *m {
		words = append(words, w)
	}
	sort.Slice(words, func(i, j int) bool {
		w1, w2 := words[i], words[j]
		if (*m)[w1] == (*m)[w2] {
			return w1 < w2
		}
		return (*m)[w1] > (*m)[w2]
	})
	if len(words) < int(size) {
		return words
	}
	return words[:size]
}

var wordPattern = regexp.MustCompile(`\P{L}*(?P<word>[\p{L}-]+)\P{L}*`)

func normalizeWord(word string) string {
	matches := wordPattern.FindStringSubmatch(strings.ToLower(word))
	if len(matches) > 1 && !isHyphen(matches[1]) {
		return matches[1]
	}
	return ""
}

func isHyphen(word string) bool {
	return word == "-"
}

func (s *str) countWords() *map[string]int {
	m := make(map[string]int)
	for _, word := range strings.Fields(string(*s)) {
		word = normalizeWord(word)
		if len(word) > 0 {
			m[word]++
		}
	}
	return &m
}

func Top10(s string) []string {
	ss := str(s)
	return ss.getTop(10)
}
