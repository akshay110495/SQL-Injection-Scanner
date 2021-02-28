package service

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/akshay110495/scanner/errors"
	"github.com/akshay110495/scanner/loggers"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

var (
	logger = loggers.GetLogger()
)

type scanner_service struct{}

func GetScannerService() scannerService {
	return &scanner_service{}
}

func (service *scanner_service) Crawl(web_url string) []string {
	var links []string
	u, err := url.Parse(web_url)
	if err != nil {
		loggers.WriteLog(logger, err)
		return nil
	}
	host := u.Hostname()
	collector := colly.NewCollector(
		colly.AllowedDomains(host),
	)
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		collector.Visit(e.Request.AbsoluteURL(link))
	})

	collector.OnRequest(func(r *colly.Request) {
		link := r.URL.String()
		if !contains(links, link) {
			links = append(links, link)
		}
	})

	collector.Visit(web_url)
	return links
}

func (service *scanner_service) CheckInjectability(links []string) map[string][]string {
	var vulnerableURLSList []string
	vulnerableURLS := make(map[string][]string)
	for _, url := range links {
		urlWithoutQuery, query_params, err := service.ParseURLWithoutQuery(url)
		if err != nil && urlWithoutQuery != nil {
			continue
		}
		stringified_url := urlWithoutQuery.String()
		if !contains(vulnerableURLSList, stringified_url) && len(query_params) != 0 {
			if service.IsURLVulnerable(urlWithoutQuery, query_params) {
				vulnerableURLSList = append(vulnerableURLSList, stringified_url)
				parameters := make([]string, 0, len(query_params))
				for key := range query_params {
					parameters = append(parameters, key)
				}
				vulnerableURLS[stringified_url] = parameters
			}
		}
	}
	return vulnerableURLS
}

func (service *scanner_service) ParseURLWithoutQuery(ref string) (*url.URL, url.Values, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return nil, nil, err
	}
	query_params := u.Query()
	u.RawQuery = ""

	return u, query_params, nil
}

func (service *scanner_service) IsURLVulnerable(web_url *url.URL, params url.Values) bool {
	query := web_url.Query()
	for key, values := range params {
		for _, val := range values {
			query.Set(key, val+"'")
		}
	}
	web_url.RawQuery = query.Encode()
	resp, err := http.Get(web_url.String())
	if err != nil {
		loggers.WriteLog(logger, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		loggers.WriteLog(logger, err)
	}
	for _, rgExpList := range errors.ErrorsRegexp {
		for _, regExp := range rgExpList {
			matched, err := regexp.Match(regExp, body)
			if err != nil {
				loggers.WriteLog(logger, err)
			}
			if matched == true {
				return true
			}
		}
	}
	return false
}

func (service *scanner_service) GetFormActions(links []string) []string {
	var actions []string
	var checkedURLS []string
	for _, link := range links {
		urlWithoutQuery, _, _ := service.ParseURLWithoutQuery(link)
		stringified_url := urlWithoutQuery.String()
		if !contains(checkedURLS, stringified_url) {
			checkedURLS = append(checkedURLS, stringified_url)
			resp, _ := http.Get(link)
			body, _ := ioutil.ReadAll(resp.Body)
			doc, err := html.Parse(strings.NewReader(string(body)))
			if err != nil {
				loggers.WriteLog(logger, err)
			}
			var f func(*html.Node)
			f = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "form" {
					for _, a := range n.Attr {
						if a.Key == "action" {
							actions = append(actions, urlWithoutQuery.Scheme+"://"+urlWithoutQuery.Host+"/"+a.Val)
						}
					}
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
			f(doc)
		}
	}
	return actions
}
