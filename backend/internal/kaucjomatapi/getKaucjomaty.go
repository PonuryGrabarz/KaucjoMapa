package kaucjomatapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetKaucjomaty(path string) ([]Kaucjomat, error) {
	url := fmt.Sprintf("%s/%s", baseUrl, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Kaucjomat{}, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []Kaucjomat{}, err
	}
	defer res.Body.Close()

	var kaucjomatyList []Kaucjomat

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&kaucjomatyList)
	if err != nil {
		return []Kaucjomat{}, nil
	}

	return kaucjomatyList, nil
}
