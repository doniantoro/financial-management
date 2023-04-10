package repository

import (
	domain_transaction "multi-finance/domain/transaction"
	helper "multi-finance/helper"

	"gorm.io/gorm"
)

type MysqlTransactionRepository struct {
	db *gorm.DB
}

func NewMysqlTransactionRepository(db *gorm.DB) *MysqlTransactionRepository {
	return &MysqlTransactionRepository{db}
}

func (q MysqlTransactionRepository) StoreTransactionPartner(reqPartner *domain_transaction.PartnerTransactionRequest, reqTransaction *domain_transaction.TransactionRequest) error {

	err := q.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&reqTransaction).Error; err != nil {
			return err
		}

		if err := tx.Create(&reqPartner).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (q MysqlTransactionRepository) StoreTransactionMarketPlace(req *domain_transaction.TransactionRequest) error {

	err := q.db.Create(&req)

	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (q MysqlTransactionRepository) Index(pagination helper.Pagination) (*helper.Pagination, error) {
	var transaction []*domain_transaction.TransactionResponse
	q.db.Debug().Scopes(helper.Paginate(transaction, &pagination, q.db)).Find(&transaction)

	pagination.Rows = transaction
	return &pagination, nil
}
func (q MysqlTransactionRepository) Find(condition string) (*domain_transaction.TransactionResponse, error) {
	var transaction domain_transaction.TransactionResponse
	err := q.db.Where(condition).Order("created_at desc").Preload("Installment").Last(&transaction)
	if err.Error != nil {
		return nil, err.Error
	}
	return &transaction, nil
}

func (q MysqlTransactionRepository) StoreInstallment(transaction *domain_transaction.Installment) error {
	err := q.db.Create(&transaction)
	if err != nil {
		return err.Error
	}
	return nil
}
