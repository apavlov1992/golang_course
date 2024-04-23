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

func f(comicsNumber int) {

}
func main() {
	var configFileName string
	flag.StringVar(&configFileName, "config", "../config/config.yaml", "Specify configuration file name to use.")
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

	comicsNumber := 1
	errCount := 0
	bar := progressbar.NewOptions(-1, progressbar.OptionShowCount(), progressbar.OptionSetDescription("Getting comics..."))
	_, err = client.GetIdList()

	for {

		comicsNumber++

		bar.Add(1)
		if !xkcd.IDinDB(comicsNumber) {
			ComicsData, err := client.GetComics(comicsNumber)
			if err != nil {
				log.Fatal("Can't get comic's data: ", err)
				errCount++
			}
			if errCount > 2 {
				break
			}
			ComicsData.Description = strings.Join(stemming.StemmingString(ComicsData.Description), " ")
			comicsDataInBytes := xkcd.SerializeToMap(ComicsData)
			f, err := os.OpenFile(cfg.DBFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal("Error when opening file: ", err)
			}
			defer f.Close()
			if ComicsData.Num != 0 {
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
}
