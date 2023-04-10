package shopee_usecase

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"

	domain_shopee "multi-finance/domain/market_place/shopee"
	"multi-finance/helper"
)

type MarketPlaceUsecase struct {
	service_shopee domain_shopee.ShopeeService
	redis          *redis.Client
}

func NewMarketPlaceUsecase(service_shopee domain_shopee.ShopeeService, redis *redis.Client) *MarketPlaceUsecase {
	return &MarketPlaceUsecase{service_shopee, redis}
}
func (uc *MarketPlaceUsecase) FindProduct(id string) (*domain_shopee.ReponseFindData, error) {
	// Init Variable
	now := time.Now()

	var installment domain_shopee.Installment
	product, err := uc.service_shopee.FindProduct(id)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUseMarketPlace + "|Find-Product",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      id,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	}

	installmentEnv := strings.Split(os.Getenv("INSTALLMENT"), "|")
	for _, data := range installmentEnv {

		installmentPeriod := strings.Split(data, "-")
		installmentPeriodInt, err := strconv.Atoi(installmentPeriod[0])
		if err != nil {
			log.Printf("Failed convert installmentPeriod to int %v", err)
			return nil, err
		}
		installmentAmount, err := strconv.Atoi(installmentPeriod[1])
		if err != nil {
			log.Printf("Failed convert installmentAmount to int %v", err)
			return nil, err
		}

		installment.Period = installmentPeriodInt
		installment.Interest = installmentAmount
		installment.Amount = (product.OtrPrice * installment.Amount / 100) + product.OtrPrice
		product.Installment = append(product.Installment, installment)
	}

	return product, nil
}
