package main

import (
	"log"
	serviceShopee "multi-finance/external/delivery/service/market_place/shopee"
	"multi-finance/helper"
	repoCustomer "multi-finance/internal/repository/mysql"
	repoUser "multi-finance/internal/repository/mysql"

	"fmt"
	"multi-finance/config"
	auth_handler "multi-finance/internal/delivery/handler/auth"

	productHandler "multi-finance/external/delivery/handler/product"
	product_usecase "multi-finance/external/usecase/market_place"
	customerHandler "multi-finance/internal/delivery/handler/customer"
	transactionHandler "multi-finance/internal/delivery/handler/transaction"
	auth_usecase "multi-finance/internal/usecase/auth"
	customerUsecase "multi-finance/internal/usecase/customer"
	transactionUsecase "multi-finance/internal/usecase/transaction"
	"multi-finance/middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	baseDir := helper.DynamicDir()
	fmt.Println("baseDir", baseDir)
	envPath := fmt.Sprintf("%s.env", baseDir)

	err := godotenv.Load(envPath)
	if err != nil {
		log.Print("Error Env", err.Error())
	}
	cache, err := config.RedisConnection()
	log.Println("redis accecpted connection", cache.Ping())
	failOnError(err, "Error Redis : ")
	dbMysqlGorm := config.MysqlGorm()

	repoMysqlUser := repoUser.NewMysqlUserRepository(dbMysqlGorm)
	authUsecase := auth_usecase.NewAuthUsecase(repoMysqlUser, cache)
	authHttp := auth_handler.NewAuthHandlerDelivery(authUsecase)

	repoMysqlCustomer := repoCustomer.NewMysqlCustomerRepository(dbMysqlGorm)
	customerUsecase := customerUsecase.NewAuthUsecase(repoMysqlCustomer, cache)
	customerHttp := customerHandler.NewCustomerHandlerDelivery(customerUsecase)

	serviceShopee := serviceShopee.NewShopeeService()
	repoMysqlTransaction := repoCustomer.NewMysqlTransactionRepository(dbMysqlGorm)
	transactionUsecase := transactionUsecase.NewTransactionUsecase(repoMysqlTransaction, repoMysqlCustomer, serviceShopee, cache)
	transactionHttp := transactionHandler.NewTransactionHandlerDelivery(transactionUsecase)

	product_usecase := product_usecase.NewMarketPlaceUsecase(serviceShopee, cache)
	productHandler := productHandler.NewProductHandlerDelivery(product_usecase)
	// define mux router
	r := mux.NewRouter()
	// define router

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/login", authHttp.Login).Methods(http.MethodPost)
	api.Handle("/logout", middleware.ValidationToken(http.HandlerFunc(authHttp.Logout), cache)).Methods(http.MethodPost)
	api.Handle("/refresh-token", middleware.ValidationToken(http.HandlerFunc(authHttp.Refresh), cache)).Methods(http.MethodPost)

	api.Handle("/product/{id}", middleware.ValidationToken(http.HandlerFunc(productHandler.FindProduct), cache)).Methods(http.MethodGet)

	customer := api.PathPrefix("/customer/application").Subrouter()
	customer.Handle("/{id}", middleware.ValidationToken(http.HandlerFunc(customerHttp.FindApplication), cache)).Methods(http.MethodGet)
	customer.Handle("", middleware.ValidationToken(http.HandlerFunc(customerHttp.ApplyApplication), cache)).Methods(http.MethodPost)
	customer.Handle("", middleware.ValidationToken(http.HandlerFunc(customerHttp.ListApplication), cache)).Methods(http.MethodGet)
	customer.Handle("/{id}", middleware.ValidationToken(http.HandlerFunc(customerHttp.UpdateStatusApplication), cache)).Methods(http.MethodPut)

	transaction := api.PathPrefix("/transaction").Subrouter()

	transaction.Handle("", middleware.ValidationToken(http.HandlerFunc(transactionHttp.ListTransaction), cache)).Methods(http.MethodGet)
	transaction.Handle("/{id}", middleware.ValidationToken(http.HandlerFunc(transactionHttp.FindTransaction), cache)).Methods(http.MethodGet)
	transaction.Handle("/partner", middleware.ValidationToken(http.HandlerFunc(transactionHttp.PartnerTransaction), cache)).Methods(http.MethodPost)
	transaction.Handle("/market-place", middleware.ValidationToken(http.HandlerFunc(transactionHttp.MarketPlaceTransaction), cache)).Methods(http.MethodPost)

	if os.Getenv("USE_SSL") == "true" {

		log.Print("Run on port :", os.Getenv("PORT"), " With SSL:", os.Getenv("USE_SSL"))
		err = http.ListenAndServeTLS(":"+os.Getenv("PORT"), baseDir+"server.crt", baseDir+"server.key", r)
		fmt.Println(err)

	} else {

		log.Print("Run on port :", os.Getenv("PORT"), " With SSL:", os.Getenv("USE_SSL"))
		err = http.ListenAndServe(":"+os.Getenv("PORT"), r)
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
