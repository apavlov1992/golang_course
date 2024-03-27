package main

import (
	"flag"
	"fmt"
	"github.com/kljensen/snowball/english"
	"io"
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
	return slices.Contains(stopWordList, word)
}

func StemmingWords(words []string) string {
	wordsList := []string{} //to do map

	for _, word := range words {
		if len(word) > 2 && !isStopWord(word) {
			stemmedWord := english.Stem(word, true)
			wordsList = append(wordsList, stemmedWord)
		}
	}
	slices.Sort(wordsList)
	return strings.Join(slices.Compact(wordsList), " ")

}

func main() {
	text := flag.String("s", "", "Text from arguments")
	flag.Parse()
	words := strings.Fields(*text)

	err := DownloadFile(stopWordsFile, fileUrl)

	if err != nil {
		fmt.Println("Error downloading file: ", err)
		return
	}

	fileData, _ = os.ReadFile(stopWordsFile)
	stopWordList = strings.Split(string(fileData), ",")

	fmt.Println(StemmingWords(words))
}
