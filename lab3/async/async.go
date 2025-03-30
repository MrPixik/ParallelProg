package async

import (
	"ParallelProg/lab3/static"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func generatePageUrls(ctx context.Context) <-chan string {
	outCh := make(chan string)

	go func() {
		defer close(outCh)

		for i := 1; i <= static.PageToParseAmount; i++ {
			url := fmt.Sprintf("%s?page=%d&query=%s&prepend=None", static.ResourceURL, i, static.Query)
			outCh <- url
		}
	}()

	return outCh
}

func fanOutlinkExtractor(ctx context.Context, client *resty.Client, inCh <-chan string) []chan *static.Article {
	outChans := make([]chan *static.Article, 0)

	for url := range inCh {
		outCh := linkExtractor(ctx, client, url)
		outChans = append(outChans, outCh)
	}
	return outChans
}

func linkExtractor(ctx context.Context, client *resty.Client, url string) chan *static.Article {
	outCh := make(chan *static.Article)

	go func() {
		defer close(outCh)
		res, err := client.R().Get(url)
		if err != nil {
			outCh <- &static.Article{Link: "", Err: err}
		}
		if res.StatusCode() != http.StatusOK {
			outCh <- &static.Article{Link: "", Err: fmt.Errorf("status code error: %d %s", res.StatusCode(), res.Status())}
		}
		// Load the HTML page
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
		if err != nil {
			outCh <- &static.Article{Link: "", Err: err}
		}

		doc.Find("a.l").Each(func(i int, s *goquery.Selection) {
			if link, exist := s.Attr("href"); exist {
				link, err = createNormURL(link)
				outCh <- &static.Article{Link: link, Err: err}
			}
		})
	}()
	return outCh

}

func fanInArticle(ctx context.Context, inCh []chan *static.Article) chan *static.Article {
	outCh := make(chan *static.Article)

	var wg sync.WaitGroup

	for _, currCh := range inCh {
		in := currCh

		wg.Add(1)
		go func() {
			defer wg.Done()

			for article := range in {

				outCh <- article
			}
		}()

	}
	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func fanOutArticle(ctx context.Context, client *resty.Client, inCh <-chan *static.Article, f func(ctx context.Context, client *resty.Client, article *static.Article) chan *static.Article) []chan *static.Article {
	outChans := make([]chan *static.Article, 0)

	for article := range inCh {
		outCh := f(ctx, client, article)
		outChans = append(outChans, outCh)
	}
	return outChans
}

func contentExtractor(ctx context.Context, client *resty.Client, article *static.Article) chan *static.Article {
	outCh := make(chan *static.Article)

	go func() {
		defer close(outCh)
		if article.Err != nil {
			return
		}

		res, err := client.R().Get(article.Link)
		if err != nil {
			article.Err = err
			outCh <- article
		}

		if res.StatusCode() != http.StatusOK {
			article.Err = fmt.Errorf("status code error: %d %s", res.StatusCode(), res.Status())
			outCh <- article
		}
		// Load the HTML page
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
		if err != nil {
			article.Err = err
			outCh <- article
		}

		content := strings.Builder{}

		if preamble := doc.Find("#preamble").Text(); preamble != "" {
			content.WriteString(preamble)
		}
		if text := doc.Find("#main-text").Text(); text != "" {
			content.WriteString(text)
		}

		article.Content = content.String()
		outCh <- article
	}()
	return outCh
}

func summarizeLoader(ctx context.Context, client *resty.Client, article *static.Article) chan *static.Article {
	outCh := make(chan *static.Article)
	go func() {
		defer close(outCh)
		if article.Err != nil {
			return
		}

		reqOR := static.ReqOpenRouter{
			Model:     static.Model,
			Prompt:    static.PromptSample + article.Summary,
			MaxTokens: static.MaxTokens,
		}

		body, err := json.Marshal(reqOR)
		if err != nil {
			article.Err = err
			outCh <- article
		}
		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", "Bearer "+static.ApiKey).
			SetBody(body).
			Post(static.AiURL)
		if err != nil {
			article.Err = fmt.Errorf("error while POST-request: %s", err.Error())
			outCh <- article
		}

		if res.StatusCode() != http.StatusOK {
			article.Err = fmt.Errorf("status code error: %d %s", res.StatusCode(), res.Status())
			outCh <- article
		}

		var response map[string]interface{}
		if err = json.Unmarshal(res.Body(), &response); err != nil {
			article.Err = fmt.Errorf("error while Unmarshaling: %s", err.Error())
			outCh <- article
		}

		if choices, ok := response["choices"].([]interface{}); ok {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if text, ok := choice["text"].(string); ok {
					article.Summary = text
					outCh <- article
					return
				}
			}
			article.Err = fmt.Errorf("invalid response format")
			outCh <- article
		}

		article.Err = fmt.Errorf("invalid response format")
		outCh <- article
	}()

	return outCh
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

func writeArticle(file *os.File, article *static.Article) error {

	data, err := json.Marshal(article)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	file.Write([]byte("\n"))
	return err
}

func Run() {
	//Разрешаем использование только одного ядра для работы программы,
	//чтобы программа выполнялась асинхронно
	runtime.GOMAXPROCS(1)

	start := time.Now()

	ctx := context.Background()
	client := resty.New()

	urlsCh := generatePageUrls(ctx)

	linkedArticles := fanOutlinkExtractor(ctx, client, urlsCh)
	links := fanInArticle(ctx, linkedArticles)

	contentArticles := fanOutArticle(ctx, client, links, contentExtractor)
	contentArticle := fanInArticle(ctx, contentArticles)

	summaryArticles := fanOutArticle(ctx, client, contentArticle, summarizeLoader)
	summaryArticle := fanInArticle(ctx, summaryArticles)

	var articleNum int
	file, err := os.OpenFile(static.FilenameAsync, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	for res := range summaryArticle {
		if res.Err != nil {
			panic(res.Err)
		}
		if err := writeArticle(file, res); err != nil {
			panic(err)
		}
		articleNum++
	}
	fmt.Println("Asynchrone variant:")
	fmt.Printf("Time: %f, Articles processed: %d", time.Since(start).Seconds(), articleNum)
}
