package ubysseydotca

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ubyssey/chatbotfb/app/utils/jsonparser"
)

type ApiJsonResult struct {
	Results []Article `json:"results"`
}

type Article struct {
	Headline      string `json:"headline"`
	FeaturedImage `json:"featured_image"`
	Snippet       string `json:"snippet"`
	Url           string `json:"url"`
}

type FeaturedImage struct {
	Url string `json:"url"`
}

// Sends a GET HTTP request to a Ubyssey.com API endpoint
func GetHttpRequestApi(url string) (ApiJsonResult, error) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	apiJsonResult := ApiJsonResult{}

	resp, err := client.Get(url)
	if err != nil {
		return apiJsonResult, err
	}

	// If HTTP response is OK, then convert the response body
	// to a byte slice for unmarshalling
	if resp.StatusCode == 200 { // OK
		respBodyByteSlice, respBodyByteSliceErr := ioutil.ReadAll(resp.Body)
		if respBodyByteSliceErr != nil {
			return apiJsonResult, respBodyByteSliceErr
		}

		// Unmarshal resp body to a struct
		parsingErr := jsonparser.Parse(respBodyByteSlice, &apiJsonResult)
		if parsingErr != nil {
			return apiJsonResult, parsingErr
		}
	}
	return apiJsonResult, nil
}
