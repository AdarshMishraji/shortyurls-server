package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type HTMLMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	SiteName    string `json:"site_name"`
}

func GetMetaData(link string) (*HTMLMeta, error) {
	if link == "" {
		return nil, errors.New("Invalid URL")
	}

	if _, err := url.Parse(link); err != nil {
		return nil, errors.New("Invalid URL")
	}

	resp, err := http.Get(link)
	if err != nil {
		return nil, errors.New("Unable to fetch metadata")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Unable to fetch metadata")
	}

	meta := extract(resp.Body, link)

	if meta.SiteName == "" {
		meta.SiteName = meta.Title
	}

	return meta, nil
}

func extract(resp io.Reader, link string) *HTMLMeta {
	z := html.NewTokenizer(resp)

	titleFound := false

	hm := new(HTMLMeta)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return hm
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == `body` {
				return hm
			}
			if t.Data == "title" {
				titleFound = true
			}
			if t.Data == "meta" {
				desc, ok := extractMetaProperty(t, "description")
				if ok {
					hm.Description = desc
				}

				ogTitle, ok := extractMetaProperty(t, "og:title")
				if ok {
					hm.Title = ogTitle
				}

				ogDesc, ok := extractMetaProperty(t, "og:description")
				if ok {
					hm.Description = ogDesc
				}

				ogImage, ok := extractMetaProperty(t, "og:image")
				if ok {
					hm.Image = ogImage
				}

				if hm.Image == "" {
					ogImage, ok := extractCustomProperty(t, "itemprop", "image")
					if ok {
						hm.Image = link + ogImage
					}
				}

				ogSiteName, ok := extractMetaProperty(t, "og:site_name")
				if ok {
					hm.SiteName = ogSiteName
				}
			}
		case html.TextToken:
			if titleFound {
				t := z.Token()
				hm.Title = t.Data
				titleFound = false
			}
		}
	}
}

func extractMetaProperty(t html.Token, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == "property" && attr.Val == prop {
			ok = true
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}

func extractCustomProperty(t html.Token, key string, prop string) (content string, ok bool) {
	for _, attr := range t.Attr {
		if attr.Key == key && attr.Val == prop {
			ok = true

		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return
}
