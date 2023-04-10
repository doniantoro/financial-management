package marketplace_shopee

type ShopeeService interface {
	FindProduct(id string) (*ReponseFindData, error)
	PostOrder(id *RequestOrder) error
}
type ShopeeUsecase interface {
	FindProduct(id string) (*ReponseFindData, error)
}
