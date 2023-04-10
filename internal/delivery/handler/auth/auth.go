package auth_http

import (
	"encoding/json"
	domain "multi-finance/domain/user"
	"multi-finance/helper"
	"net/http"
	"time"
)

type AuthNewHttpDelivery struct {
	Usecase domain.UserUsecase
}

func NewAuthHandlerDelivery(Usecase domain.UserUsecase) *AuthNewHttpDelivery {
	return &AuthNewHttpDelivery{Usecase: Usecase}
}

func (a *AuthNewHttpDelivery) Login(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	var req domain.UserRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&req)
	if err != nil {
		defer helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerAuth + "Login|Error-Decode",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       r.Method,
			Message:      err.Error(),
			Request:      req,
			ClientIP:     r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:    r.UserAgent(),
			ResponseTime: time.Since(now),
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
			Event:        helper.EventHandlerAuth + "Login|Vailed-Payload",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_INFO,
			Method:       r.Method,
			Message:      err.Error(),
			Request:      req,
			ClientIP:     r.Header.Get(helper.X_FORWARDED_FOR),
			UserAgent:    r.UserAgent(),
			ResponseTime: time.Since(now),
		}, "service")
		helper.ResponseWithJSON(helper.ResponseHttp{W: w,
			StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
			Message:    err.Error(),
		})
		return
	}

	data, err := a.Usecase.Get(&req)
	if err != nil {
		defer helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerAuth + "|Do-Usecase",
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

	defer helper.CreateLog(&helper.Log{
		Event:      helper.EventHandlerAuth + "Login|Success",
		StatusCode: helper.MapHttpStatusCode(err),
		Method:     r.Method,
		Type:       helper.LVL_INFO,
		Message:    "Success Login",
		Response:   data,
		ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:  r.UserAgent(),
	}, "service")

	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Login",
		Data:       data,
	})
}

func (a *AuthNewHttpDelivery) Logout(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	now := time.Now()
	err := a.Usecase.Logout(id)
	if err != nil {
		defer helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerAuth + "Logout|Do-Usecase",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       r.Method,
			Message:      err.Error(),
			Request:      id,
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
	defer helper.CreateLog(&helper.Log{
		Event:        helper.EventHandlerAuth + "Logout|Do-Usecase",
		StatusCode:   helper.MapHttpStatusCode(err),
		Type:         helper.LVL_INFO,
		Method:       r.Method,
		Message:      "Success Logout",
		Request:      id,
		ClientIP:     r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:    r.UserAgent(),
		ResponseTime: time.Since(now),
	}, "service")
	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Logout",
	})
}
func (a *AuthNewHttpDelivery) Refresh(w http.ResponseWriter, r *http.Request) {
	id := r.Header.Get("id")
	now := time.Now()

	data, err := a.Usecase.Refresh(id)
	if err != nil {
		defer helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerAuth + "Refresh|Do-Usecase",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       r.Method,
			Message:      err.Error(),
			Request:      id,
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
	defer helper.CreateLog(&helper.Log{
		Event:        helper.EventHandlerAuth + "Refresh",
		StatusCode:   helper.MapHttpStatusCode(err),
		Type:         helper.LVL_INFO,
		Method:       r.Method,
		Message:      "Success Referesh Token",
		Request:      id,
		ClientIP:     r.Header.Get(helper.X_FORWARDED_FOR),
		UserAgent:    r.UserAgent(),
		ResponseTime: time.Since(now),
	}, "service")
	helper.ResponseWithJSON(helper.ResponseHttp{W: w,
		StatusCode: helper.MapHttpStatusCode(err),
		Message:    "Success Referesh Token",
		Data:       data,
	})
}
