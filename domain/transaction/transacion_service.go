package transaction

import "multi-finance/helper"

type TransactionUsecase interface {
	PartnerTransaction(req *PartnerTransactionRequest) error
	MarketPlaceTransaction(req *TransactionRequest) (string, error)
	ListTransaction(req *TransactioQueryParam) (*helper.Pagination, error)
	FindTransaction(idParam, idUser, role string) (*TransactionResponse, error)
}

type TransactionRepository interface {
	StoreTransactionPartner(reqPartner *PartnerTransactionRequest, reqTransaction *TransactionRequest) error
	StoreTransactionMarketPlace(req *TransactionRequest) error
	Index(req helper.Pagination) (*helper.Pagination, error)
	Find(condition string) (*TransactionResponse, error)
	StoreInstallment(transaction *Installment) error
}
