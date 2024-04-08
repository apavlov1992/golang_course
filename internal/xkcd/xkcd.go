package xkcd

import (
	"encoding/json"
	"github.com/apavlov1992/golang_course/internal/stemming"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	http.Client
	BaseUrl string
}

type ComicsInfo struct {
	Num         interface{}
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

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(bodyBytes, &comic)
		if err != nil {
			log.Fatal(err)
		}
		comic.Description = strings.Join(
			stemming.StemmingMain(comic.Description),
			" ",
		)

		comic.Num = strconv.Itoa(comicsNumber)
		comicsList = append(comicsList, comic)
	}
	return comicsList, nil
}
