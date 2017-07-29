package ubysseydotca

type ApiJsonResult {
	Results []Article `json:"results"`
}

type Article {
	Headline string `json:"headline"`
	FeaturedImage `json:"featured_image"`
	Snippet string `json:"snippet"`
	Url string `json:"url"`
}

type FeaturedImage {
	Url string `json:"url"`
}

func GetHttpRequestApi(url string) (ApiJsonResult, error) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client {
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	apiJsonResult := ApiJsonResult{}
	json.Unmarshal(resp, &apiJsonResult)

	return apiJsonResult, nil
}
