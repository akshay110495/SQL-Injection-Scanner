package service

type main_service struct {
	scanner_service scannerService
	url             string
	links           []string
	injectableLinks map[string][]string
}

func GetMainService(web_url string) MainService {
	return &main_service{
		scanner_service: GetScannerService(),
		url:             web_url,
	}
}

func (service *main_service) GetLinks() []string {
	service.links = service.scanner_service.Crawl(service.url)
	service.links = append(service.links, service.scanner_service.GetFormActions(service.links)...)
	return service.links
}

func (service *main_service) GetInjectableLinks() map[string][]string {
	service.injectableLinks = service.scanner_service.CheckInjectability(service.links)
	return service.injectableLinks
}
