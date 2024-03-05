package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// 处理 IP 地理位置
func GetIPRecord(ip string) (map[string]interface{}, error) {

	body, err := GetIPGeoJSONBytes(ip)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	if data["status"] != "success" {
		msg, _ := json.Marshal(data)
		return nil, errors.New(string(msg))
	}

	return data, nil
}

func GetIPGeoJSONBytes(ip string) ([]byte, error) {
	url := "http://ip-api.com/json/" + ip + "?lang=zh-CN"

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
