package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const textURLFmt = "https://%s/translate_a/single?client=gtx&dt=t&dt=bd&dt=md&dt=ex&sl=%s&tl=%s&q=%s"

// DefaultHost is the Google Translate host used when none is given.
const DefaultHost = "translate.google.com"

// DefaultUserAgent is the User-Agent header used when none is given.
const DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type Translation struct {
	Text string `json:"text"`          // translated text
	POS  string `json:"pos,omitempty"` // part of speech breakdown
	Def  string `json:"def,omitempty"` // definitions
}

func Translate(srcLangCode, dstLangCode, message, proxyURL, host, userAgent string) (*Translation, error) {
	translation := new(Translation)

	if host == "" {
		host = DefaultHost
	}
	if userAgent == "" {
		userAgent = DefaultUserAgent
	}
	urlStr := fmt.Sprintf(textURLFmt, host, srcLangCode, dstLangCode, url.QueryEscape(message))

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

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
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
