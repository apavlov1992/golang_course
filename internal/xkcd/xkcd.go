package xkcd

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Client struct {
	http.Client
	BaseUrl string
}

type ComicsInfo struct {
	Num         int    `json:"num"`
	URL         string `json:"img"`
	Description string `json:"alt"`
}

func SerializeToMap(c interface{}) []byte {
	jsonContent, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	return jsonContent
}

func (client Client) GetMaxId() int {
	var comic ComicsInfo
	resp, err := client.Get(client.BaseUrl + "/info.0.json")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&comic)
	lastComics := comic.Num

	return lastComics
}

func (client Client) GetComics(numberOfComics int) ([]ComicsInfo, error) {
	var comicsList []ComicsInfo
	var comic ComicsInfo
	for comicsNumber := 1; comicsNumber <= numberOfComics; comicsNumber++ {
		if comicsNumber == 404 {
			continue
		}

		resp, err := client.Get(client.BaseUrl + strconv.Itoa(comicsNumber) + "/info.0.json")
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&comic)

		if err != nil {
			log.Fatal(err)
		}

		comic.Num = comicsNumber
		comicsList = append(comicsList, comic)
	}
	return comicsList, nil
}
