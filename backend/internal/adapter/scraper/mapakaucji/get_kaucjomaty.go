package mapakaucji

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetKaucjomaty[T any](path string) (T, error) {
	url := fmt.Sprintf("%s/%s", baseURL, path)
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
