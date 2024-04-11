package stemming

import (
	"fmt"
	"github.com/kljensen/snowball/english"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

var (
	stopWordsFile = "gist_stopwords.txt"
	fileUrl       = fmt.Sprintf("https://gist.githubusercontent.com/ZohebAbai/513218c3468130eacff6481f424e4e64/raw/b70776f341a148293ff277afa0d0302c8c38f7e2/" + stopWordsFile)
	fileData      []byte
	stopWordList  []string
	wordsMap      map[string]any
)

func DownloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func isStopWord(word string) bool {
	_, ok := wordsMap[word]
	return ok
}

func stemmingWords(words []string) []string {
	wordsList := []string{}

	for _, word := range words {
		if len(word) > 2 && !isStopWord(word) {
			stemmedWord := english.Stem(word, true)
			wordsList = append(wordsList, stemmedWord)
		}
	}
	slices.Sort(wordsList)

	return wordsList
}

func StemmingString(words string) []string {

	err := DownloadFile(stopWordsFile, fileUrl)
	if err != nil {
		log.Fatal(err)
	}

	fileData, err = os.ReadFile(stopWordsFile)
	if err != nil {
		log.Fatal(err)
	}

	stopWordList = strings.Split(string(fileData), ",")
	wordsMap = make(map[string]any)
	for _, word := range stopWordList {
		wordsMap[word] = struct{}{}
	}

	return stemmingWords(strings.Split(words, " "))
}
