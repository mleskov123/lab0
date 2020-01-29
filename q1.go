package cos418_hw1_1

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Find the top K most common words in a text document.
// 	path: location of the document
//	numWords: number of words to return (i.e. k)
//	charThreshold: character threshold for whether a token qualifies as a word,
//		e.g. charThreshold = 5 means "apple" is a word but "pear" is not.
// Matching is case insensitive, e.g. "Orange" and "orange" is considered the same word.
// A word comprises alphanumeric characters only. All punctuation and other characters
// are removed, e.g. "don't" becomes "dont".
// You should use `checkError` to handle potential errors.
func topWords(path string, numWords int, charThreshold int) []WordCount {
	// TODO: implement me
	// HINT: You may find the `strings.Fields` and `strings.ToLower` functions helpful
	// HINT: To keep only alphanumeric characters, use the regex "[^0-9a-zA-Z]+"

	// current idea is to keep a map of name to wordcount object. then retrieve the values put into
	// a list and send to sortWordCounts
	// another idea is to keep something like a sorted binary tree (by name) and update the numbers

	file, err := os.Open(path)
	checkError(err)
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	reg := regexp.MustCompile("[0-9a-zA-Z]+")
	myMap := make(map[string]WordCount)
	for scanner.Scan() {

		line := scanner.Text()
		words := strings.Fields(line)
		for _, word := range words {
			word = strings.ToLower(word)
			word = strings.Join(reg.FindAllString(word, -1), "")
			if len(word) < charThreshold || len(word) == 0{
				continue
			}
			if val, ok := myMap[word]; ok {
				val.Count = val.Count + 1
				// this might be unnecessary
				myMap[word] = val
			} else {
				myMap[word] = WordCount{word, 1}
			}

		}

	}
	checkError(scanner.Err())

	var counts []WordCount = make([]WordCount, 0, len(myMap))
	for _, val := range myMap {
		counts = append(counts, val)
	}
	sortWordCounts(counts)
	if len(counts) > 0{
		var endLength = numWords
		if endLength > len(counts) {
			endLength = len(counts)
		}
		return counts[0:endLength]
	}
	return counts
}

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.
// DO NOT MODIFY THIS FUNCTION!
func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}
