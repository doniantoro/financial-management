package transaction_http

import (
	"encoding/json"
	domain "multi-finance/domain/transaction"
	"multi-finance/helper"
	"net/http"

	"github.com/gorilla/mux"
)

type TransactionNewHttpDelivery struct {
	Usecase domain.TransactionUsecase
}

func NewTransactionHandlerDelivery(Usecase domain.TransactionUsecase) *TransactionNewHttpDelivery {
	return &TransactionNewHttpDelivery{Usecase: Usecase}
}

func (a *TransactionNewHttpDelivery) ListTransaction(w http.ResponseWriter, r *http.Request) {
	var req domain.TransactioQueryParam
	req.Id = r.Header.Get("id")
	req.Role = r.Header.Get("role")
	req.Limit = r.URL.Query().Get("limit")
	req.Page = r.URL.Query().Get("page")
	req.Status = r.URL.Query().Get("status")
	req.PartnerId = r.Header.Get("partner-id")

	data, err := a.Usecase.ListTransaction(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerTransaction + "|ListTransaction",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     helper.METHOD_GET,
			Message:    err.Error(),
			Request:    req,
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
		Event:      helper.EventHandlerTransaction,
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     helper.METHOD_GET,
		Type:       helper.LVL_INFO,
		Message:    "Success Get List Transaction",
		Request:    req,
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),
		Response:   data,
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Get List Transaction",
		Data:       data,
	})
}

func (a *TransactionNewHttpDelivery) FindTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	role := r.Header.Get("role")
	idUser := r.Header.Get("id")
	idParam := params["id"]

	data, err := a.Usecase.FindTransaction(idParam, idUser, role)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerTransaction + "|Usecase-FindTransaction",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     helper.METHOD_GET,
			Message:    err.Error(),
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
		Event:      helper.EventHandlerTransaction,
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     helper.METHOD_GET,
		Type:       helper.LVL_INFO,
		Message:    "Success Get List Transaction",
		Request:    "role : " + role + " idUser : " + idUser + " idParam : " + idParam,
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),

		Response: data,
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success GET Transaction",
		Data:       data,
	})
}

func (a *TransactionNewHttpDelivery) PartnerTransaction(w http.ResponseWriter, r *http.Request) {

	event := "Partnern-Transaction"
	var req domain.PartnerTransactionRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerAuth + event + "|Error-Decode",
			StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
			Type:       helper.LVL_ERROR,
			Method:     helper.METHOD_POST,
			Message:    err.Error(),
			Request:    req,
			ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:  r.UserAgent(),
		}, "service")
		helper.ResponseWithJSON(helper.ResponseHttp{W: w,
			StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
			Message:    err.Error(),
		})
		return
	}

	req.SalesId = r.Header.Get("id")
	req.PartnerId = r.Header.Get("partner-id")

	err = helper.Validate.Struct(req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerTransaction + event + "|Validate-payload",
			StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
			Type:       helper.LVL_INFO,
			Method:     helper.METHOD_POST,
			Message:    err.Error(),
			Request:    req,
			ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:  r.UserAgent(),
		}, "service")
		helper.ResponseWithJSON(helper.ResponseHttp{W: w,
			StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
			Message:    err.Error(),
		})
		return
	}

	err = a.Usecase.PartnerTransaction(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerTransaction + event + "|Partner-Transaction-Usecase",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     helper.METHOD_POST,
			Message:    err.Error(),
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
		Event:      helper.EventHandlerTransaction + event,
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     helper.METHOD_POST,
		Type:       helper.LVL_INFO,
		Message:    "Success Submit Transaction",
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Submit Transaction",
	})
}
func (a *TransactionNewHttpDelivery) MarketPlaceTransaction(w http.ResponseWriter, r *http.Request) {
	var req domain.TransactionRequest
	event := "MarketPlace-Transaction"
	id := r.Header.Get("id")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerTransaction + event + "Error-Decode",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     helper.METHOD_POST,
			Message:    err.Error(),
			Request:    req,
			ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:  r.UserAgent(),
		}, "service")
		helper.ResponseWithJSON(helper.ResponseHttp{W: w,
			StatusCode: helper.MapHttpStatusCode(err),
			Message:    err.Error(),
		})
		return
	}
	req.CustomerId = id
	err = helper.Validate.Struct(req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerTransaction + event + "Validate-Payload",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_INFO,
			Method:     helper.METHOD_POST,
			Message:    err.Error(),
			Request:    req,
			ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:  r.UserAgent(),
		}, "service")
		helper.ResponseWithJSON(helper.ResponseHttp{W: w,
			StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
			Message:    err.Error(),
		})
		return
	}

	code, err := a.Usecase.MarketPlaceTransaction(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerTransaction + event + "Error-Usecase",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     helper.METHOD_POST,
			Message:    err.Error(),
			Request:    req,
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
		Event:      helper.EventHandlerTransaction + event,
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     helper.METHOD_POST,
		Type:       helper.LVL_INFO,
		Message:    "Success Submit Transaction",
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success ubmit Transaction",
		Data:       code,
	})
}
