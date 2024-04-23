package xkcd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	IDMap map[int]any
)

type Client struct {
	http.Client
	BaseUrl string
	DB      string
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

func (client Client) GetComics(comicNumber int) (ComicsInfo, error) {
	var comic ComicsInfo
	if comicNumber != 404 {
		resp, err := client.Get(client.BaseUrl + strconv.Itoa(comicNumber) + "/info.0.json")

		if err != nil {
			log.Fatal("Can't get response: ", err)
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&comic)
		if err != nil {
			log.Fatal("Can't get response body: ", err)
		}

		comic.Num = comicNumber

	}
	return comic, nil
}

func IDinDB(ID int) bool {
	_, ok := IDMap[ID]
	return ok
}

func (client Client) GetIdList() (map[int]any, error) {
	var comicId ComicsInfo
	var comicsIdList []int
	content, err := os.Open(client.DB)

	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	defer content.Close()

	decoder := json.NewDecoder(content)
	for decoder.More() {
		IDMap = make(map[int]any)
		err := decoder.Decode(&comicId)
		if err != nil {
			log.Fatal("Error when opening file: ", err)
		}
		comicsIdList = append(comicsIdList, comicId.Num)
	}
	for _, ID := range comicsIdList {
		IDMap[ID] = struct{}{}
	}
	return IDMap, err
}
