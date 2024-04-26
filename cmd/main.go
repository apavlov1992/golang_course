package main

import "C"
import (
	"flag"
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

	var configFileName string
	flag.StringVar(&configFileName, "config", "/Users/a.pavlov/GolandProjects/golang_course/config/config.yaml", "Specify configuration file name to use.")
	flag.Parse()

	cfg, err := config.NewConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	client := xkcd.Client{
		Client:  http.Client{Timeout: 10 * time.Second},
		BaseUrl: cfg.SourceUrl,
		DB:      cfg.DBFile,
	}

	_, err = client.GetIdList()

	worker := NewComicsWorker(client, 100)

	comics, err2 := worker.HandleComics()
	if err2 != nil {
		log.Fatal("Error while handling comics: ", err)
	}

	f, err := os.OpenFile(cfg.DBFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	defer f.Close()
	bar := progressbar.NewOptions(-1, progressbar.OptionShowCount(), progressbar.OptionSetDescription("Writing to DB..."))
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
}
