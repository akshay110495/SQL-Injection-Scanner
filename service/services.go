package service

import "net/url"

type MainService interface {
	GetLinks() []string
	GetInjectableLinks() map[string][]string
}

type scannerService interface {
	Crawl(string) []string
	GetFormActions([]string) []string
	CheckInjectability([]string) map[string][]string
	IsURLVulnerable(*url.URL, url.Values) bool
	ParseURLWithoutQuery(string) (*url.URL, url.Values, error)
}

func contains(list []string, item string) bool {
	for _, value := range list {
		if value == item {
			return true
		}
	}
	return false
}
