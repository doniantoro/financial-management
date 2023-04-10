package auth_http

import (
	"encoding/json"
	domain "multi-finance/domain/customer"
	"multi-finance/helper"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CustomerNewHttpDelivery struct {
	Usecase domain.CustomerUsecase
}

func NewCustomerHandlerDelivery(Usecase domain.CustomerUsecase) *CustomerNewHttpDelivery {
	return &CustomerNewHttpDelivery{Usecase: Usecase}
}
func (a *CustomerNewHttpDelivery) ListApplication(w http.ResponseWriter, r *http.Request) {
	var req domain.CustomerQueryParam
	id := r.Header.Get("id")
	event := "List-Application"
	limit := r.URL.Query().Get("limit")
	status := r.URL.Query().Get("status")
	page := r.URL.Query().Get("page")
	name := r.URL.Query().Get("name")
	idNumber := r.URL.Query().Get("id_number")

	req = domain.CustomerQueryParam{
		Id:       id,
		Limit:    limit,
		Page:     page,
		Status:   status,
		Name:     name,
		IdNumber: idNumber,
	}

	data, err := a.Usecase.ListApplication(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerCustomer + event + "|Error-Usecase",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     r.Method,
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
		Event:      helper.EventHandlerCustomer + event,
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

func (a *CustomerNewHttpDelivery) ApplyApplication(w http.ResponseWriter, r *http.Request) {
	var req domain.CustomerRequest
	event := "Apply-Application"
	id := r.Header.Get("id")
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerCustomer + "|Error-Decode",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     r.Method,
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
	err = helper.Validate.Struct(req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerCustomer + event + "|Error-Validate",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_INFO,
			Method:     r.Method,
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

	err = a.Usecase.ApplyApplication(&req, id)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerCustomer + event + "|Error-Usecase",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     r.Method,
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
		Event:      helper.EventHandlerCustomer + event,
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     r.Method,
		Type:       helper.LVL_INFO,
		Message:    "Success Apply Application",
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Apply Application",
	})
}

func (a *CustomerNewHttpDelivery) UpdateStatusApplication(w http.ResponseWriter, r *http.Request) {
	var req domain.CustomerRequestUpdate
	now := time.Now()
	decoder := json.NewDecoder(r.Body)
	params := mux.Vars(r)
	id := params["id"]
	accepted_user_id := r.Header.Get("id")
	event := "Update-Status-Application"

	err := decoder.Decode(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerCustomer + event + "|Error-Decode",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_ERROR,
			Method:     r.Method,
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
	req.Id = id
	req.AcceptedUserId = accepted_user_id
	err = helper.Validate.Struct(req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:      helper.EventHandlerCustomer + event + "|Validate-Payload",
			StatusCode: helper.MapHttpStatusCode(err),
			Type:       helper.LVL_INFO,
			Method:     r.Method,
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

	err = a.Usecase.UpdateStatusApplication(&req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Error-Usecase",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       r.Method,
			Message:      err.Error(),
			ClientIP:     r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:    r.UserAgent(),
			ResponseTime: time.Since(now),
		}, "service")

		helper.ResponseWithJSON(helper.ResponseHttp{
			W:          w,
			StatusCode: helper.MapHttpStatusCode(err),
			Message:    err.Error(),
		})
		return
	}

	helper.CreateLog(&helper.Log{
		Event:        helper.EventHandlerCustomer + event,
		StatusCode:   helper.MapHttpStatusCode(err),
		Method:       r.Method,
		Type:         helper.LVL_INFO,
		Message:      "Success Update Status",
		ClientIP:     r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:    r.UserAgent(),
		ResponseTime: time.Since(now),
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Update Status",
	})
}

func (a *CustomerNewHttpDelivery) FindApplication(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	now := time.Now()

	data, err := a.Usecase.FindApplication(id)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + "FindApplication|Error-Usecase",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       r.Method,
			Message:      err.Error(),
			ClientIP:     r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:    r.UserAgent(),
			ResponseTime: time.Since(now),
		}, "service")

		helper.ResponseWithJSON(helper.ResponseHttp{
			W:          w,
			StatusCode: helper.MapHttpStatusCode(err),
			Message:    err.Error(),
		})
		return
	}

	helper.CreateLog(&helper.Log{
		Event:      helper.EventHandlerCustomer + "FindApplication",
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     r.Method,
		Type:       helper.LVL_INFO,
		Message:    "Success FInd Application",
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),
		Response:   data,
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Find Application",
		Data:       data,
	})
}
