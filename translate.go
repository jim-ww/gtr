package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const textURL = "https://translate.googleapis.com/translate_a/single?client=gtx&dt=t&dt=bd&dt=md&dt=ex&sl=%s&tl=%s&q=%s"

type Translation struct {
	Text string `json:"text"`          // translated text
	POS  string `json:"pos,omitempty"` // part of speech breakdown
	Def  string `json:"def,omitempty"` // definitions
}

func Translate(srcLangCode, dstLangCode, message, proxyURL string) (*Translation, error) {
	translation := new(Translation)

	urlStr := fmt.Sprintf(textURL, srcLangCode, dstLangCode, url.QueryEscape(message))

	client := http.DefaultClient
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy url: %w", err)
		}
		client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
		}
	}

	res, err := client.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() //nolint:errcheck

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data []any
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if len(data) <= 0 {
		return nil, errors.New("translation not found")
	}

	// translated text = data[0]
	for _, line := range data[0].([]any) {
		translatedLine := line.([]any)[0]
		translation.Text += translatedLine.(string)
	}

	// part of speech = data[1]
	if len(data) > 1 && data[1] != nil {
		for _, partOfSpeeches := range data[1].([]any) {
			partOfSpeeches := partOfSpeeches.([]any)
			pos := partOfSpeeches[0]
			translation.POS += fmt.Sprintf("[%v]\n", pos)
			for _, words := range partOfSpeeches[2].([]any) {
				words := words.([]any)
				dstWord := words[0]
				translation.POS += fmt.Sprintf("\t%v:", dstWord)
				firstWord := true
				for _, word := range words[1].([]any) {
					if firstWord {
						translation.POS += fmt.Sprintf(" %v", word)
						firstWord = false
					} else {
						translation.POS += fmt.Sprintf(", %v", word)
					}
				}
				translation.POS += "\n"
			}
		}
	}

	// definitions = data[12]
	if len(data) >= 13 && data[12] != nil {
		for _, definitions := range data[12].([]any) {
			definitions := definitions.([]any)
			pos := definitions[0]
			translation.Def += fmt.Sprintf("[%v]\n", pos)
			for _, sentences := range definitions[1].([]any) {
				sentences := sentences.([]any)
				def := sentences[0]
				translation.Def += fmt.Sprintf("\t- %v\n", def)
				if len(sentences) >= 3 && sentences[2] != nil {
					example := sentences[2]
					translation.Def += fmt.Sprintf("\t\t\"%v\"\n", example)
				}
			}
		}
	}

	return translation, nil
}
