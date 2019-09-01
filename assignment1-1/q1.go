package cos418_hw1_1

import (
	"strings"
	"fmt"
	"sort"
	"bufio"
	"os"
	"log"
	"regexp"
)

// Find the top K most common words in a text document.
// 	path: location of the document
//	numWords: number of words to return (i.e. k)
//	charThreshold: character threshold for whether a token qualifies as a word,
//		e.g. charThreshold = 5 means "apple" is a word but "pear" is not.
// Matching is case insensitive, e.g. "Orange" and "orange" is considered the same word.
// A word comprises alphanumeric characters only. All punctuations and other characters
// are removed, e.g. "don't" becomes "dont".
// You should use `checkError` to handle potential errors.
func topWords(path string, numWords int, charThreshold int) []WordCount {
	// TODO: implement me
	// HINT: You may find the `strings.Fields` and `strings.ToLower` functions helpful
	// HINT: To keep only alphanumeric characters, use the regex "[^0-9a-zA-Z]+"
	//return nil
	file, err := os.Open(path)   	//Open the file for reading
    if err != nil {					//Check if error happened when opening a file
        log.Fatal(err)
    }
    defer file.Close()

	scanner := bufio.NewScanner(file)

	var y string  					//Append each line to y
	for scanner.Scan() {             
        y=y+scanner.Text()  		
    }

	y = strings.Replace(y,"."," ",-1)   //Corner case, word1.Word2 was considered as 1 big word
	m := make(map[string]int)
	
	re := regexp.MustCompile("[^0-9a-zA-Z]+")

    for _, word := range strings.Fields(y){
        word = strings.ToLower(word)
		word = re.ReplaceAllLiteralString(word,"")
		if(len(word)>=charThreshold){					//Checking if word is longer or equal than th
			m[word] = m[word]+1
		}     
	}
	var words [] WordCount
	for k,v := range m {								//Transform map to WordCount Array
		words =append(words, WordCount{k,v})
	}

	sortWordCounts(words)

	return words[0:numWords]							//Return slice of k words
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
