package kaucjomatapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FeatureCollection[T any] struct {
	Type     string `json:"type"`
	Features []T    `json:"features"`
}

func GetKaucjomaty[T any](path string) (T, error) {
	url := fmt.Sprintf("%s/%s", baseUrl, path)
	var kaucjomatyList T

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return kaucjomatyList, err
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return kaucjomatyList, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&kaucjomatyList)
	if err != nil {
		return kaucjomatyList, err
	}

	return kaucjomatyList, nil
}
