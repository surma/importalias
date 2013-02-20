package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Extractor interface {
	Extract(c *http.Client) (string, error)
}

type ExtractorFunc func(c *http.Client) (string, error)

func (e ExtractorFunc) Extract(c *http.Client) (string, error) {
	return e(c)
}

func NewJSONExtractor(url string, field string) Extractor {
	return ExtractorFunc(func(c *http.Client) (string, error) {
		r, err := c.Get(url)
		if err != nil {
			return "", err
		}

		var data interface{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return "", err
		}
		switch x := data.(type) {
		case map[string]interface{}:
			data = x[field]
		default:
			return "", fmt.Errorf("Unhandled JSON type")
		}
		switch x := data.(type) {
		case float64:
			return fmt.Sprintf("%.0f", data), nil
		case string:
			return x, nil
		}
		return "", fmt.Errorf("Unsupported id type")
	})
}
