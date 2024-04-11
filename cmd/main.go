package main

import "C"
import (
	"flag"
	"fmt"
	"github.com/apavlov1992/golang_course/internal/config"
	"github.com/apavlov1992/golang_course/internal/stemming"
	"github.com/apavlov1992/golang_course/internal/xkcd"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "config", "../config/config.yaml", "Specify configuration file name to use.")
	numberOfComics := flag.Int("n", 0, "Number of comics to fetch")
	outputJson := flag.Bool("o", false, "Output JSON in CLI or write to file")
	flag.Parse()

	cfg, err := config.NewConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	client := xkcd.Client{
		Client:  http.Client{Timeout: 10 * time.Second},
		BaseUrl: cfg.SourceUrl,
	}

	if *numberOfComics == 0 {
		countOfAllComics := client.GetMaxId()
		numberOfComics = &countOfAllComics
	}

	ComicsData, err := client.GetComics(*numberOfComics)
	if err != nil {
		log.Fatal(err)
	}

	for i := range ComicsData {
		ComicsData[i].Description = strings.Join(stemming.StemmingString(ComicsData[i].Description), " ")
	}

	if *outputJson == true {
		fmt.Println(ComicsData)
	} else {
		comicsDataInBytes := xkcd.SerializeToMap(ComicsData)
		err := os.WriteFile(cfg.DBFile, comicsDataInBytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
