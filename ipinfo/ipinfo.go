package ipinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetInfo() (*Info, error) {
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(string(responseBody))
	}

	info := new(Info)

	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
