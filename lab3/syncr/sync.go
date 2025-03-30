package syncr

import (
	"ParallelProg/lab3/static"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func getLinksToArticles(client *resty.Client, url string) ([]string, error) {

	res, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode(), res.Status())
	}
	// Load the HTML page
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)
	doc.Find("a.l").Each(func(i int, s *goquery.Selection) {
		if link, exist := s.Attr("href"); exist {
			link, _ = createNormURL(link)
			links = append(links, link)
		}
	})
	return links, nil
}

func getArticleContent(client *resty.Client, article *static.Article) error {
	res, err := client.R().Get(article.Link)
	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", res.StatusCode(), res.Status())
	}
	// Load the HTML page
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		return err
	}

	content := strings.Builder{}

	if preamble := doc.Find("#preamble").Text(); preamble != "" {
		content.WriteString(preamble)
	}
	if text := doc.Find("#main-text").Text(); text != "" {
		content.WriteString(text)
	}

	article.Content = content.String()
	return nil

}

func summarizeArticle(client *resty.Client, article *static.Article) error {

	reqOR := static.ReqOpenRouter{
		Model:     static.Model,
		Prompt:    static.PromptSample + article.Summary,
		MaxTokens: static.MaxTokens,
	}

	body, err := json.Marshal(reqOR)
	if err != nil {
		return err
	}
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+static.ApiKey).
		SetBody(body).
		Post(static.AiURL)
	if err != nil {
		return fmt.Errorf("error while POST-request: %s", err.Error())
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", res.StatusCode(), res.Status())
	}

	var response map[string]interface{}
	if err = json.Unmarshal(res.Body(), &response); err != nil {
		return fmt.Errorf("error while Unmarshaling: %s", err.Error())
	}

	if choices, ok := response["choices"].([]interface{}); ok {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if text, ok := choice["text"].(string); ok {
				article.Summary = text
				//fmt.Println(text)
				return nil
			}
		}
		return fmt.Errorf("invalid response format")
	}

	return fmt.Errorf("invalid response format")
}

func createNormURL(badUrl string) (string, error) {

	// Разбираем URL на составляющие
	parsedURL, err := url.Parse(badUrl)
	if err != nil {
		return "", err
	}

	// Получаем параметры запроса
	queryParams := parsedURL.Query()

	// Закодированные параметры автоматически обновляются
	parsedURL.RawQuery = queryParams.Encode()

	// Теперь URL безопасен для запроса
	normUrl := parsedURL.String()

	return normUrl, nil
}

func writeArticle(file *os.File, article static.Article) error {

	data, err := json.Marshal(article)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	file.Write([]byte("\n"))
	return err
}

func Run() {

	start := time.Now()

	client := resty.New()

	links := make([]string, 0)
	for i := 1; i <= static.PageToParseAmount; i++ {
		url := fmt.Sprintf("%s?page=%d&query=%s&prepend=None", static.ResourceURL, i, static.Query)
		if newLinks, err := getLinksToArticles(client, url); err == nil {
			links = append(links, newLinks...)
		} else {
			panic(err)
		}
	}

	articles := make([]static.Article, len(links))
	for i := range links {

		articles[i].Link = links[i]

		if err := getArticleContent(client, &articles[i]); err != nil {
			panic(err)
		}
	}

	for i := range articles {
		if err := summarizeArticle(client, &articles[i]); err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile(static.FilenameSync, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	for _, res := range articles {
		if err := writeArticle(file, res); err != nil {
			panic(err)
		}
	}
	fmt.Println("Synchrone variant:")
	fmt.Printf("Time: %f, Articles processed: %d", time.Since(start).Seconds(), len(articles))

}
