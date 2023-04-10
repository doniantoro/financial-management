package transaction_usecase

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"

	domain_customer "multi-finance/domain/customer"
	domain_shopee "multi-finance/domain/market_place/shopee"
	domain_transaction "multi-finance/domain/transaction"

	"multi-finance/helper"
)

type TransactionUsecase struct {
	repo_transaction domain_transaction.TransactionRepository
	repo_customer    domain_customer.CustomerRepository
	service_shopee   domain_shopee.ShopeeService
	redis            *redis.Client
}

func NewTransactionUsecase(repo_transaction domain_transaction.TransactionRepository, repo_customer domain_customer.CustomerRepository, service_shopee domain_shopee.ShopeeService, redis *redis.Client) *TransactionUsecase {
	return &TransactionUsecase{repo_transaction, repo_customer, service_shopee, redis}
}

// Function to get List Transaction
func (uc *TransactionUsecase) ListTransaction(req *domain_transaction.TransactioQueryParam) (*helper.Pagination, error) {
	//Init variable
	var query string
	limit, _ := strconv.Atoi(req.Limit)
	page, _ := strconv.Atoi(req.Page)

	if req.Role == helper.ROLE_SALES {
		if query != "" {
			query = query + " AND "
		}
		query = fmt.Sprintf(query+"transactions.partner_id = '%s'", req.PartnerId)
	}
	if req.Role == helper.ROLE_CUSTOMER {
		if query != "" {
			query = query + " AND "
		}
		query = fmt.Sprintf(query+"transactions.customer_id = '%s'", req.Id)
	}
	if req.Status != "" {
		if query != "" {
			query = query + " AND "
		}

		query = fmt.Sprintf(query+"transactions.status = '%s'", req.Status)
	}

	paginate := helper.Pagination{Limit: limit,
		Page:  page,
		Query: query}
	data, err := uc.repo_transaction.Index(paginate)

	if err != nil {

		return nil, err
	}
	reflectedData := data.Rows.([]*domain_transaction.TransactionResponse)

	for _, model := range reflectedData {
		model.Installment = nil
	}
	data.Query = ""
	return data, nil
}
func (uc *TransactionUsecase) FindTransaction(idParam, idUser, role string) (*domain_transaction.TransactionResponse, error) {

	var query string

	if role == helper.ROLE_CUSTOMER {
		query = fmt.Sprintf("customer_id = '%s' AND id = '%s'", idUser, idParam)
	} else if role == helper.ROLE_ADMIN {
		query = fmt.Sprintf("id = '%s'", idParam)
	}

	data, err := uc.repo_transaction.Find(query)

	if err != nil {

		return nil, err
	} else if err == gorm.ErrRecordNotFound {
		return nil, helper.ErrNotFound
	}

	return data, nil
}

func (uc *TransactionUsecase) PartnerTransaction(req *domain_transaction.PartnerTransactionRequest) error {
	//init variable
	var intereset int
	var reqInstallment domain_transaction.Installment
	var reqUpdateLimit domain_customer.CustomerRequestUpdate

	event := "Partner-Transaction"
	reqTransaction := domain_transaction.TransactionRequest{}
	now := time.Now()
	n := 0

	condition := fmt.Sprintf("id = '%s'", req.CustomerId)
	customer, err := uc.repo_customer.Find(condition)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|FindCustomer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return err
	}
	if customer.Status != helper.STATUS_ACCEPTED {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|FindCustomer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrNoAccepted),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_POST,
			Message:      helper.ErrLimit.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return helper.ErrNoAccepted
	} else if customer.RemainLimit < req.OtrPrice {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|FindCustomer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrLimit),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_POST,
			Message:      helper.ErrLimit.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return helper.ErrLimit
	}

	installment := strings.Split(os.Getenv("INSTALLMENT"), "|")
	for _, data := range installment {

		installmentPeriod := strings.Split(data, "-")

		installmentPeriodInt, err := strconv.Atoi(installmentPeriod[0])
		if err != nil {
			log.Printf("Failed convert installmentPeriod to int %v", err)
			return err
		}
		if req.InstallmentPeriod == installmentPeriodInt {
			installmenInterest, err := strconv.Atoi(installmentPeriod[1])
			if err != nil {
				log.Printf("Failed convert installmentPeriod to int %v", err)
				return err
			}
			intereset = installmenInterest
		}
	}
	fee, err := strconv.Atoi(os.Getenv("FEE"))
	if err != nil {
		log.Printf("Failed convert fee to int %v", err)
		return err
	}
	installAmount := (((req.OtrPrice + fee) * intereset / 100) + req.OtrPrice + fee) / req.InstallmentPeriod
	req.Id = uuid.New().String()
	reqTransaction.Id = req.TransactionId
	reqTransaction.CustomerId = req.CustomerId
	reqTransaction.ProductName = req.ProductName
	reqTransaction.InstallmentPeriod = req.InstallmentPeriod
	reqTransaction.InstallmentAmount = installAmount
	reqTransaction.OtrPrice = req.OtrPrice
	reqTransaction.Fee = fee
	reqTransaction.Interest = intereset
	reqTransaction.TransactionCome = "partner"
	reqTransaction.Status = "success"
	err = helper.GetImageMimeType(req.IdCardImage)
	if err != nil {
		log.Print("Failed validate mimetype IdCardImage", err)
		return helper.ErrBadRequest
	}
	err = helper.GetImageMimeType(req.SelfieImage)

	if err != nil {
		log.Print("Failed validate mimetype SelfieImage", err)
		return helper.ErrBadRequest
	}
	err = helper.GetImageMimeType(req.TransactionLetterImage)

	if err != nil {
		log.Print("Failed validate mimetype TransactionLetterImage", err)
		return helper.ErrBadRequest
	}

	req.SelfieImage, err = helper.EnryptData(req.SelfieImage)
	if err != nil {
		log.Print("Failed Encrypt data selfieImage", err)
	}
	req.SelfieImage, err = helper.EnryptData(req.IdCardImage)
	if err != nil {
		log.Print("Failed Encrypt data idCardImage", err)
	}
	req.TransactionLetterImage, err = helper.EnryptData(req.TransactionLetterImage)
	if err != nil {
		log.Print("Failed Encrypt data transactionLetterImage", err)
	}

	err = uc.repo_transaction.StoreTransactionPartner(req, &reqTransaction)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Store-Transaction-Partner",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrDuplicate),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_POST,
			Message:      helper.ErrDuplicate.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return helper.ErrDuplicate
	}

	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Store-Transaction-Partner",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrLimit),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return err
	}
	if now.Day() < 20 {
		now = now.AddDate(0, 1, -now.Day()+1)

	} else {
		now = now.AddDate(0, 2, -now.Day()+1)

	}

	for n < req.InstallmentPeriod {

		reqInstallment.Id = uuid.New().String()
		reqInstallment.TransactionId = req.Id
		reqInstallment.Period = now
		reqInstallment.Status = "open"
		err := uc.repo_transaction.StoreInstallment(&reqInstallment)
		if err != nil {
			helper.CreateLog(&helper.Log{
				Event:        helper.EventUsecaseransaction + event + "|Store-Installment",
				StatusCode:   helper.MapHttpStatusCode(err),
				Type:         helper.LVL_ERROR,
				Method:       helper.METHOD_POST,
				Message:      err.Error(),
				Request:      req,
				ResponseTime: time.Since(now),
			}, "service")
			log.Printf("Failed store installment %v", err)
			// return "", err
		}
		n = n + 1
	}

	reqUpdateLimit.RemainLimit = customer.RemainLimit - req.OtrPrice + 1
	reqUpdateLimit.Id = req.CustomerId
	if reqUpdateLimit.RemainLimit == 0 {
		reqUpdateLimit.RemainLimit = 1
	}

	err = uc.repo_customer.UpdateStatus(&reqUpdateLimit)

	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Update-Status|Update-Limit",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return err
	}
	return nil
}

func (uc *TransactionUsecase) MarketPlaceTransaction(req *domain_transaction.TransactionRequest) (string, error) {
	// Init Variable
	var intereset int
	var request_order domain_shopee.RequestOrder
	var reqUpdateLimit domain_customer.CustomerRequestUpdate
	var reqInstallment domain_transaction.Installment
	var installmentConn bool
	now := time.Now()
	n := 0
	event := "MarketPlaceTransaction"

	data, err := uc.service_shopee.FindProduct(req.ProductId)

	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Find-Product",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", err
	}
	condition := fmt.Sprintf("user_id = '%s'", req.CustomerId)
	customer, err := uc.repo_customer.Find(condition)

	if customer.Status != helper.STATUS_ACCEPTED {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_POST,
			Message:      helper.ErrLimit.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", helper.ErrNoAccepted
	} else if customer.RemainLimit < data.OtrPrice {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Find-Customer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrLimit),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_POST,
			Message:      helper.ErrLimit.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", helper.ErrLimit
	}

	installment := strings.Split(os.Getenv("INSTALLMENT"), "|")
	for _, data := range installment {

		installmentPeriod := strings.Split(data, "-")
		installmentPeriodInt, err := strconv.Atoi(installmentPeriod[0])
		if err != nil {
			log.Printf("Failed convert installmentPeriod to int %v", err)
			return "", err
		}
		if req.InstallmentPeriod == installmentPeriodInt {
			installmentConn = true
			installmenInterest, err := strconv.Atoi(installmentPeriod[1])
			if err != nil {
				log.Printf("Failed convert installmentPeriod to int %v", err)
				return "", err
			}
			intereset = installmenInterest
		}
	}
	if !installmentConn {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Installment-Period",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrBadRequest),
			Type:         helper.LVL_INFO,
			Method:       helper.METHOD_POST,
			Message:      "installment not exist on our database",
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", helper.ErrBadRequest
	}

	fee, err := strconv.Atoi(os.Getenv("FEE"))
	if err != nil {
		log.Printf("Failed convert fee to int %v", err)
		return "", err
	}
	contractNumber := helper.GenerateContractNumber(req.ProductId)

	request_order.OrderId = contractNumber
	request_order.ProductId = req.ProductId
	err = uc.service_shopee.PostOrder(&request_order)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Post-Order",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", err
	}
	req.Status = "success"
	installAmount := (((req.OtrPrice + fee) * intereset / 100) + req.OtrPrice + fee) / req.InstallmentPeriod
	req.Id = contractNumber
	req.CustomerId = customer.Id
	req.ProductName = data.ProductName
	req.InstallmentAmount = installAmount
	req.OtrPrice = data.OtrPrice
	req.Fee = fee
	req.Interest = intereset
	req.TransactionCome = "market-place"
	err = uc.repo_transaction.StoreTransactionMarketPlace(req)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Store-Transaction-MarketPlacer",
			StatusCode:   helper.MapHttpStatusCode(helper.ErrDuplicate),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      helper.ErrDuplicate.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", helper.ErrDuplicate
	}

	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Store-Transaction-MarketPlacer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", err
	}

	if now.Day() < 20 {
		now = now.AddDate(0, 1, -now.Day()+1)

	} else {
		now = now.AddDate(0, 2, -now.Day()+1)

	}
	for n < req.InstallmentPeriod {

		reqInstallment.TransactionId = contractNumber
		reqInstallment.Period = now
		reqInstallment.Id = uuid.New().String()
		err := uc.repo_transaction.StoreInstallment(&reqInstallment)
		if err != nil {

			log.Printf("Failed store installment %v", err)
			// return "", err
		}
		n = n + 1
	}

	reqUpdateLimit.RemainLimit = customer.RemainLimit - req.OtrPrice + 1
	if reqUpdateLimit.RemainLimit == 0 {
		reqUpdateLimit.RemainLimit = 1
	}
	reqUpdateLimit.Id = req.CustomerId
	err = uc.repo_customer.UpdateStatus(&reqUpdateLimit)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseransaction + event + "|Store-Transaction-MarketPlacer",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return "", err
	}
	return contractNumber, nil
}
