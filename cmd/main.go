package main

import (
	"flag"
	"fmt"
	"github.com/apavlov1992/golang_course/cmd/services"
	"github.com/apavlov1992/golang_course/internal/config"
	"github.com/apavlov1992/golang_course/internal/stemming"
	"github.com/apavlov1992/golang_course/internal/xkcd"
	"github.com/schollz/progressbar/v3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var words string
	var configFileName string
	var index bool
	flag.StringVar(&words, "s", "", "Text from arguments")
	flag.StringVar(&configFileName, "config", "../config/config.yaml", "Specify configuration file name to use.")
	flag.BoolVar(&index, "i", false, "Find in index file")
	flag.Parse()

	cfg, err := config.NewConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	client := xkcd.Client{
		Client:    http.Client{Timeout: 10 * time.Second},
		BaseUrl:   cfg.SourceUrl,
		DB:        cfg.DBFile,
		IndexFile: cfg.IndexFile,
	}

	_, err = client.GetIdList()

	worker := services.NewComicsWorker(client, 100)

	comics, err2 := worker.HandleComics()
	if err2 != nil {
		log.Fatal("Error while handling comics: ", err)
	}

	f, err := os.OpenFile(cfg.DBFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	i, err := os.OpenFile(cfg.IndexFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	defer f.Close()
	bar := progressbar.NewOptions(-1, progressbar.OptionShowCount(), progressbar.OptionSetDescription("Writing to DB..."))
	stemmingFromCmd := stemming.StemmingString(words)
	fmt.Println(stemmingFromCmd)
	for _, c := range comics {
		bar.Add(1)
		if c.Num != 0 {
			c.Description = strings.Join(stemming.StemmingString(c.Description), " ")
			comicsDataInBytes := xkcd.SerializeToMap(c)
			_, err = f.Write(comicsDataInBytes)
			if err != nil {
				log.Fatal("Error when writing to file: ", err)
			}
			_, err = f.Write([]byte("\n"))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	IndexId := map[string]interface{}{}
	if !index {
		for _, c := range stemmingFromCmd {
			var a []int
			a, err = client.GetWordIdList(c)
			IndexId[c] = a
		}
		b := xkcd.SerializeToMap(IndexId)
		err = i.Truncate(0)
		if _, err := i.Write(b); err != nil {
			log.Fatal("Error when writing to file: ", err)
		}
	} else {
		for _, c := range stemmingFromCmd {
			var a interface{}
			a, err = client.GetIdFromIndex(c)
			IndexId[c] = a
		}
		fmt.Println(IndexId)
	}
}
