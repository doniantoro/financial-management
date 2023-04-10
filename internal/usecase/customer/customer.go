package customer_usecase

import (
	"fmt"
	"log"
	"multi-finance/domain/customer"
	domain_customer "multi-finance/domain/customer"
	"multi-finance/helper"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerUsecase struct {
	repo_customer domain_customer.CustomerRepository
	redis         *redis.Client
}

func NewAuthUsecase(repo_customer domain_customer.CustomerRepository, redis *redis.Client) *CustomerUsecase {
	return &CustomerUsecase{repo_customer, redis}
}

func (uc *CustomerUsecase) ListApplication(req *domain_customer.CustomerQueryParam) (*helper.Pagination, error) {
	var query string
	event := "List-Application"
	now := time.Now()
	limit, _ := strconv.Atoi(req.Limit)
	page, _ := strconv.Atoi(req.Page)

	if req.Name != "" {
		name, err := helper.EnryptData(req.Name)
		if err != nil {
			log.Print("Failed Encrypt data fullname", err)
			return nil, err
		}
		if query != "" {
			query = query + " AND "
		}
		query = fmt.Sprintf("full_name = '%s'", strings.ToLower(name))
	}
	if req.IdNumber != "" {
		idNumber, err := helper.EnryptData(req.IdNumber)
		if err != nil {
			log.Print("Failed Encrypt data fullname", err)
			return nil, err
		}
		if query != "" {
			query = query + " AND "
		}
		query = fmt.Sprintf(query+"id_number = '%s'", idNumber)
	}
	if req.Status != "" {
		if query != "" {
			query = query + " AND "
		}

		query = fmt.Sprintf(query+"status = '%s'", req.Status)
	}

	paginate := helper.Pagination{Limit: limit,
		Page:  page,
		Query: query}
	data, err := uc.repo_customer.Index(paginate)

	reflectedData := data.Rows.([]*domain_customer.Customer)
	if err != nil {

		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Get-Data-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	} else if len(reflectedData) == 0 {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Get-Data-Customer",
			StatusCode:   helper.MapHttpStatusCode(gorm.ErrRecordNotFound),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      gorm.ErrRecordNotFound.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, gorm.ErrRecordNotFound

	}

	for _, model := range reflectedData {
		uc.DecryptCustomerData(model)
	}
	data.Query = ""
	return data, nil
}

func (uc *CustomerUsecase) ApplyApplication(req *domain_customer.CustomerRequest, id string) error {
	event := "Apply-Application"
	now := time.Now()

	condition := fmt.Sprintf("user_id = '%s'", id)
	data, err := uc.repo_customer.Find(condition)

	if err != nil && err != gorm.ErrRecordNotFound {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Get-Data-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return err
	} else if data != nil {
		addOneMount := data.CreatedAt.Add(720 * time.Hour)
		isMoreThanOneMonth := now.Before(addOneMount)
		if !isMoreThanOneMonth && data.Status == helper.STATUS_CANCELLED {
			helper.CreateLog(&helper.Log{
				Event:        helper.EventHandlerCustomer + event + "|Get-Data-Customer",
				StatusCode:   helper.MapHttpStatusCode(helper.ErrLessThanApply),
				Type:         helper.LVL_INFO,
				Method:       helper.METHOD_GET,
				Message:      helper.ErrLessThanApply.Error(),
				Request:      req,
				ResponseTime: time.Since(now),
			}, "service")
			return helper.ErrLessThanApply
		} else if data.Status == helper.STATUS_OPEN {
			helper.CreateLog(&helper.Log{
				Event:        helper.EventHandlerCustomer + event + "|Get-Data-Customer",
				StatusCode:   helper.MapHttpStatusCode(helper.ErrLessThanApply),
				Type:         helper.LVL_INFO,
				Method:       helper.METHOD_GET,
				Message:      helper.ErrLessThanApply.Error(),
				Request:      req,
				ResponseTime: time.Since(now),
			}, "service")
			return helper.ErrThereIsApplicationProccessed
		} else if data.Status == helper.STATUS_ACCEPTED {
			helper.CreateLog(&helper.Log{
				Event:        helper.EventHandlerCustomer + event + "|Get-Data-Customer",
				StatusCode:   helper.MapHttpStatusCode(helper.ErrThereIsApplicantApproved),
				Type:         helper.LVL_INFO,
				Method:       helper.METHOD_GET,
				Message:      helper.ErrThereIsApplicantApproved.Error(),
				Request:      req,
				ResponseTime: time.Since(now),
			}, "service")
			return helper.ErrThereIsApplicantApproved
		}
	}

	err = helper.GetImageMimeType(req.IdCardImage)

	if err != nil {
		log.Print("Failed validate mimetype selfieImage", err)
		return helper.ErrBadRequest
	}
	err = helper.GetImageMimeType(req.SelfieImage)

	if err != nil {
		log.Print("Failed validate mimetype id card", err)
		return helper.ErrBadRequest
	}
	req.FullName, err = helper.EnryptData(strings.ToLower(req.FullName))
	if err != nil {
		log.Print("Failed Encrypt data fullname", err)
		return err
	}
	req.LegalName, err = helper.EnryptData(strings.ToLower(req.LegalName))
	if err != nil {
		log.Print("Failed Encrypt data legalName", err)
		return err
	}
	req.Msisdn, err = helper.EnryptData(req.Msisdn)
	if err != nil {
		log.Print("Failed Encrypt data msisdn", err)
		return err
	}
	req.IdNumber, err = helper.EnryptData(req.IdNumber)
	if err != nil {
		log.Print("Failed Encrypt data idNumber", err)
		return err
	}
	req.SelfieImage, err = helper.EnryptData(req.SelfieImage)
	if err != nil {
		log.Print("Failed Encrypt data selfieImage", err)
		return err
	}
	req.IdCardImage, err = helper.EnryptData(req.IdCardImage)
	if err != nil {
		log.Print("Failed Encrypt data selfieImage", err)
		return err
	}

	req.UserId = id
	req.Id = uuid.New().String()
	err = uc.repo_customer.Store(req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Store-Data-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return err
	}

	return nil
}

func (uc *CustomerUsecase) FindApplication(id string) (*domain_customer.Customer, error) {
	data := &domain_customer.Customer{}
	event := "Find-Application"
	now := time.Now()
	dataRedis, err := uc.redis.HGetAll("customer:" + id).Result()
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Store-Data-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      id,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	}

	updated_at, _ := time.Parse(time.RFC3339, dataRedis["updated_at"])
	created_at, _ := time.Parse(time.RFC3339, dataRedis["created_at"])
	sallary, _ := strconv.Atoi(dataRedis["sallary"])
	limit, _ := strconv.Atoi(dataRedis["limit"])
	remainLimit, _ := strconv.Atoi(dataRedis["remain_limit"])

	if len(dataRedis) != 0 {
		data = &domain_customer.Customer{
			Id:             dataRedis["id"],
			IdNumber:       dataRedis["id_number"],
			AcceptedUserId: dataRedis["accepted_user_id"],
			UserId:         dataRedis["user_id"],
			Msisdn:         dataRedis["msisdn"],
			FullName:       dataRedis["full_name"],
			LegalName:      dataRedis["legal_name"],
			PlaceOfBirth:   dataRedis["place_of_birth"],
			DateOfBirth:    dataRedis["date_of_birth"],
			Sallary:        sallary,
			IdCardImage:    dataRedis["id_card_image"],
			SelfieImage:    dataRedis["selfie_image"],
			Limit:          limit,
			RemainLimit:    remainLimit,
			Status:         dataRedis["status"],
			CreatedAt:      created_at,
			UpdatedAt:      updated_at,
		}
		return data, nil
	}
	condition := fmt.Sprintf("id = '%s'", id)
	data, err = uc.repo_customer.Find(condition)
	if err == gorm.ErrRecordNotFound {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(gorm.ErrRecordNotFound),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      id,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err

	} else if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      id,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	}
	uc.DecryptCustomerData(data)
	return data, nil
}

func (uc *CustomerUsecase) UpdateStatusApplication(req *domain_customer.CustomerRequestUpdate) error {
	event := "Update-Status-Application"
	now := time.Now()
	condition := fmt.Sprintf("id = '%s'", req.Id)

	data, err := uc.repo_customer.Find(condition)

	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return err
	}
	if data.AcceptedUserId != "" && data.AcceptedUserId != req.AcceptedUserId {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrForbidden),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      "Only User Accpted that can update status",
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return helper.ErrForbidden
	}

	if req.Status == helper.STATUS_CANCELLED && req.Limit != 0 || req.Status == helper.STATUS_OPEN && req.Limit != 0 {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrBadRequest),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      "cancelled or open status must have limit 0",
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return helper.ErrBadRequest
	} else if req.Status == helper.STATUS_ACCEPTED && req.Limit <= 0 {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrBadRequest),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      "accepted status must have limit more than 0",
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return helper.ErrBadRequest
	} else if req.Status != helper.STATUS_OPEN && req.Status != helper.STATUS_CANCELLED && req.Status != helper.STATUS_ACCEPTED {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrBadRequest),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      "status must be open, cancelled or accepted",
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return helper.ErrBadRequest
	}
	req.RemainLimit = req.Limit
	req.UpdatedAt = time.Now()
	err = uc.repo_customer.UpdateStatus(req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventHandlerCustomer + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_GET,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return err
	}
	isExist, err := uc.redis.HExists("customer:"+req.Id, "limit").Result()
	if isExist {
		var m = make(map[string]interface{})
		m["limit"] = req.Limit
		m["remain_limit"] = req.RemainLimit
		m["status"] = req.Status
		m["updated_at"] = req.UpdatedAt
		m["accepted_user_id"] = req.AcceptedUserId

		ttl, err := strconv.Atoi(os.Getenv("EXPIRE_REDIS_FIND_CUSTOMER"))
		if err != nil {
			log.Println("error convert expire redis find customer", err)
			return err
		}
		err = uc.redis.HMSet("customer:"+req.Id, m).Err()
		if err != nil {
			log.Println("error HMSet", err)
			return err
		}
		err = uc.redis.Do("EXPIRE", "customer:"+req.Id, ttl).Err()
		if err != nil {
			log.Println("error set expire", err)
			return err
		}
	}
	return nil

}

func (uc *CustomerUsecase) DecryptCustomerData(model *customer.Customer) error {
	idCardImage, err := helper.DecryptData(model.IdCardImage)
	if err != nil {
		log.Print("Failed Decrypt data IdCardImage", err)
	}
	model.SelfieImage, err = helper.DecryptData(model.SelfieImage)
	if err != nil {
		log.Print("Failed Decrypt data IdCardImage", err)
	}
	//decrypted msisdn
	model.Msisdn, err = helper.DecryptData(model.Msisdn)
	if err != nil {
		log.Print("Failed Decrypt data msisdn", err)
	}

	//decrypt legal name
	model.LegalName, err = helper.DecryptData(model.LegalName)
	if err != nil {
		log.Print("Failed Decrypt data LegalName", err)
	}

	//decrypt id number
	model.IdNumber, err = helper.DecryptData(model.IdNumber)
	if err != nil {
		log.Print("Failed Decrypt data IdNumber", err)
	}

	//decrypt full name
	model.FullName, err = helper.DecryptData(model.FullName)
	if err != nil {
		log.Print("Failed Decrypt data FullName", err)
	}

	model.IdCardImage = idCardImage
	return nil
}
