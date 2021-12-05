package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(text string) []string {
	wordsCountMap := make(map[string]int)

	wordsList := strings.Fields(text)

	// create map of words
	for _, word := range wordsList {
		if _, ok := wordsCountMap[word]; ok {
			wordsCountMap[word]++
			continue
		}
		wordsCountMap[word] = 1
	}
	fmt.Println(wordsCountMap)

	// find top 10
	// - create slice of words
	wordsSlice := make([]string, 0, len(wordsCountMap))
	for word := range wordsCountMap {
		wordsSlice = append(wordsSlice, word)
	}
	fmt.Println(wordsSlice, len(wordsSlice))

	// - sort slice with map data
	sort.Slice(wordsSlice, func(i, j int) bool {
		if wordsCountMap[wordsSlice[i]] != wordsCountMap[wordsSlice[j]] {
			return wordsCountMap[wordsSlice[i]] > wordsCountMap[wordsSlice[j]]
		}
		return strings.Compare(wordsSlice[i], wordsSlice[j]) < 0
	})

	// - return top 10 or len of slice
	returnSize := 10
	if len(wordsSlice) < returnSize {
		returnSize = len(wordsSlice)
	}
	return wordsSlice[:returnSize]
}
