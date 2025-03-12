package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	wordsCounter := map[string]int{}
	var words []string

	for _, v := range strings.Fields(s) {
		if _, ok := wordsCounter[v]; ok {
			wordsCounter[v]++
		} else {
			words = append(words, v)
			wordsCounter[v] = 1
		}
	}

	if len(words) < 2 {
		return words
	}

	sort.Slice(words, func(i, j int) bool {
		a, b := wordsCounter[words[i]], wordsCounter[words[j]]
		if a == b {
			return words[i] < words[j]
		}

		return a > b
	})

	if l := len(words); l > 10 {
		return words[:10]
	}

	return words
}
