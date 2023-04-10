package product_http

import (
	"multi-finance/helper"

	domain "multi-finance/domain/market_place/shopee"
	"net/http"

	"github.com/gorilla/mux"
)

type ProductNewHttpDelivery struct {
	Usecase domain.ShopeeUsecase
}

func NewProductHandlerDelivery(Usecase domain.ShopeeUsecase) *ProductNewHttpDelivery {
	return &ProductNewHttpDelivery{Usecase: Usecase}
}
func (a *ProductNewHttpDelivery) FindProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	data, err := a.Usecase.FindProduct(id)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerCustomer + "|Find-Product|Error-Usecase",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     r.Method,
			Message:    err.Error(),
			Request:    id,
			ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:  r.UserAgent(),
		}, "service")

		helper.ResponseWithJSON(helper.ResponseHttp{
			W:          w,
			StatusCode: helper.MapHttpStatusCode(err),
			Message:    err.Error(),
		})
		return
	}

	helper.CreateLog(&helper.Log{
		Event:      helper.EventHandlerCustomer + "|Find-Product",
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     r.Method,
		Type:       helper.LVL_INFO,
		Response:   data,
		Message:    "Success Get All Application",
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Get All Application",
		Data:       data,
	})
}
