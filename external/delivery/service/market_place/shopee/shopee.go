package service_marketplace_shopee

import (
	"encoding/json"
	"log"
	domain_shopee "multi-finance/domain/market_place/shopee"
	"multi-finance/helper"
	"os"

	"github.com/doniantoro/gogix"
)

type ShopeeService struct {
}

func NewShopeeService() *ShopeeService {
	return &ShopeeService{}
}
func (service *ShopeeService) FindProduct(id string) (*domain_shopee.ReponseFindData, error) {
	baseUrl := os.Getenv("SHOPEE_FIND_DATA") + id
	var resp domain_shopee.ReponseFindData
	httpClient := gogix.NewClient(20)
	response, code, err := httpClient.Get(baseUrl, nil)
	if err != nil {
		return nil, err
	}

	if code == 200 {
		json.Unmarshal(response, &resp)
		return &resp, nil
	} else {
		return nil, helper.ErrInternalServerError
	}
}
func (service *ShopeeService) PostOrder(req *domain_shopee.RequestOrder) error {
	baseUrl := os.Getenv("SHOPEE_POST_ORDER")
	httpClient := gogix.NewClient(20)
	response, code, err := httpClient.Post(baseUrl, nil, req)
	if err != nil {
		return err
	}
	log.Print(response)
	if code == 200 {
		return nil
	} else {
		return helper.ErrInternalServerError
	}
}
