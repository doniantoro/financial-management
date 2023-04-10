package user_usecase

import (
	"log"
	domain_user "multi-finance/domain/user"
	"multi-finance/helper"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
)

type AuthUsecase struct {
	repo_user domain_user.MysqlUserRepository
	Redis     *redis.Client
}

func NewAuthUsecase(repo_user domain_user.MysqlUserRepository, Redis *redis.Client) *AuthUsecase {
	return &AuthUsecase{repo_user, Redis}
}
func (uc AuthUsecase) Get(req *domain_user.UserRequest) (*domain_user.UserResponse, error) {
	now := time.Now()
	encryptedPassword, err := helper.EnryptData(req.Password)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseAuth + "|Get|Error-Encrypt",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	}
	req.Password = encryptedPassword
	data, err := uc.repo_user.Find(req)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseAuth + "|Get|Error-Find-User",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	}
	signedToken, refreshoken, err := uc.GetToken(data)

	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseAuth + "|Get|Error-Generate-Token",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      req,
			ResponseTime: time.Since(now),
		}, "service")

		return nil, err
	}

	data.AccesToken = *signedToken
	data.RefreshToken = *refreshoken
	data.Customer = nil
	return data, nil
}

func (uc AuthUsecase) Logout(id string) error {

	err := uc.Redis.Del("authorization-" + id)
	if err := err.Err(); err != nil {
		return err
	}
	err = uc.Redis.Del("authorization-refresh-" + id)
	if err := err.Err(); err != nil {
		return err
	}
	return nil
}
func (uc AuthUsecase) Refresh(id string) (*domain_user.UserResponse, error) {
	now := time.Now()

	data, err := uc.repo_user.FindById(id)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseAuth + "Refresh|Find-User",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      id,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	}
	err = uc.Logout(id)
	if err != nil {
		helper.CreateLog(&helper.Log{
			Event:        helper.EventUsecaseAuth + "Refresh|Logout",
			StatusCode:   helper.MapHttpStatusCode(err),
			Type:         helper.LVL_ERROR,
			Method:       helper.METHOD_POST,
			Message:      err.Error(),
			Request:      id,
			ResponseTime: time.Since(now),
		}, "service")
		return nil, err
	}
	signedToken, refreshoken, err := uc.GetToken(data)
	if err != nil {
		return nil, err
	}
	data.AccesToken = *signedToken
	data.RefreshToken = *refreshoken
	return data, nil
}
func (uc AuthUsecase) GetToken(data *domain_user.UserResponse) (*string, *string, error) {
	ttl, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {

		return nil, nil, err
	}
	ttlRefresh, err := time.ParseDuration(os.Getenv("REFRESH_JWT_DURATION"))
	if err != nil {

		return nil, nil, err
	}
	fullName, _ := helper.DecryptData(data.FullName)
	exp := time.Now().UTC().Add(ttl)
	claim := jwt.MapClaims{}
	claim["id"] = data.ID
	claim["email"] = data.Email
	claim["role"] = data.Role
	claim["full_name"] = fullName
	claim["partner_id"] = data.PartnerId
	claim["exp"] = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {

		return nil, nil, err
	}
	claim["exp"] = time.Now().UTC().Add(ttlRefresh).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedrefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {

		return nil, nil, err
	}

	setNewToken := uc.Redis.Set("authorization-"+data.ID, signedToken, ttl)
	if err := setNewToken.Err(); err != nil {
		log.Printf("unable to SET data. error: %v", err)
	}
	setRefreshToken := uc.Redis.Set("authorization-refresh-"+data.ID, signedrefreshToken, ttlRefresh)
	if err := setRefreshToken.Err(); err != nil {
		log.Printf("unable to SET data. error: %v", err)
	}
	return &signedToken, &signedrefreshToken, nil
}
