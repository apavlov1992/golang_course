package main

import (
	"flag"
	"fmt"
	"github.com/apavlov1992/golang_course/internal/config"
	"github.com/apavlov1992/golang_course/internal/xkcd"
	"log"
	"net/http"
	"os"
)

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "config", "../config/config.yaml", "Specify configuration file name to use.")
	numberOfComics := flag.Int("n", 10, "Number of comics to fetch")
	outputJson := flag.Bool("o", false, "Output JSON in CLI or write to file")
	flag.Parse()

	cfg, err := config.NewConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	client := xkcd.Client{
		Client:  http.Client{},
		BaseUrl: cfg.SourceUrl,
	}

	ComicsData, err := client.GetComics(*numberOfComics)
	if err != nil {
		log.Fatal(err)
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
