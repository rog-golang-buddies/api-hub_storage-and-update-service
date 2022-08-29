package dto

import "github.com/rog-golang-buddies/api_hub_common/apispecdoc"

type ScrapingResult struct {
	IsNotifyUser bool

	ApiSpecDoc *apispecdoc.ApiSpecDoc
}
