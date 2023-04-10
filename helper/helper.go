package helper

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/forgoer/openssl"
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ResponseHttp struct {
	W          http.ResponseWriter `json:",omitempty"`
	StatusCode int                 `json:"status_code,omitempty"`
	Message    string              `json:"message"`
	Data       interface{}         `json:"data,omitempty"`
}

const (
	ROLE_CUSTOMER    = "customer"
	ROLE_ADMIN       = "admin"
	ROLE_SALES       = "sales"
	STATUS_OPEN      = "open"
	STATUS_CANCELLED = "canceled"
	STATUS_ACCEPTED  = "accepted"
)
const (
	METHOD_POST     = "POST"
	METHOD_GET      = "GET"
	X_FORWARDED_FOR = "X-Forwarded-For"
)

var (
	EventHandlerAuth        = "Handler-Auth|"
	EventHandlerTransaction = "Handler-Transaction|"
	EventHandlerCustomer    = "Handler-Customer|"
	EventUsecaseAuth        = "Usecase-Auth|"
	EventMiddleware         = "Middleware|"
	EventUsecaseransaction  = "Usecase-Transaction|"
	EventUseMarketPlace     = "Usecase-MarketPlace|"
)
var (
	SqlOpen  = sql.Open
	log      = logrus.New()
	Validate = validator.New()
	validate *validator.Validate
)
var (
	ErrNotFound                     = errors.New("error not found")
	ErrPayload                      = errors.New("invalid request payload")
	ErrInternalServerError          = errors.New("error - 500 internal server error - an error occurred in the system. please try again later")
	ErrDatabase                     = errors.New("error - 500 internal server error - an error occurred in the system. please try again later")
	ErrBadRequest                   = errors.New("error - 400 bad Request - please try again later")
	ErrDuplicate                    = errors.New("the data already Exist")
	ErrForbidden                    = errors.New("error - 403 Sorry You cannot do this action")
	ErrUnAuhorized                  = errors.New("error - 401 Error Unauhorized")
	ErrUnmarshalResponse            = errors.New("error occurred while unmarshal json process")
	ErrLessThanApply                = errors.New("error Re-Apply only more than 30 Days")
	ErrThereIsApplicationProccessed = errors.New("There is application Proccessed")
	ErrThereIsApplicantApproved     = errors.New("Cannot Re-Apply,Because there is application that approved")
	ErrNoAccepted                   = errors.New("Cannot throug because there is no application that approved")
	ErrLimit                        = errors.New("Sorry,limit transaction less than product price")
)

func MapHttpStatusCode(err error) int {

	switch err {
	case nil:
		return http.StatusOK
	case ErrPayload:
		return http.StatusBadRequest
	case ErrInternalServerError, ErrUnmarshalResponse:
		return http.StatusInternalServerError
	case ErrBadRequest, ErrDuplicate:
		return http.StatusBadRequest
	case ErrNotFound, gorm.ErrRecordNotFound:
		return http.StatusNotFound
	case ErrUnAuhorized:
		return http.StatusUnauthorized
	case ErrLessThanApply, ErrThereIsApplicationProccessed,
		ErrThereIsApplicantApproved, ErrForbidden, ErrLimit, ErrNoAccepted:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}

}
func GetImageMimeType(mime string) error {
	imageAllowed := strings.Split(os.Getenv("IMAGE_ALLOWED"), ",")

	splitData := strings.Split(mime, ":")
	splitBase64 := strings.Split(splitData[1], ";")

	mimeType := splitBase64[0]

	isExists := InArray(mimeType, imageAllowed)
	if !isExists {

		return errors.New("image type not valid")
	}

	return nil
}
func InArray(val interface{}, array interface{}) (exists bool) {
	exists = false

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return
			}
		}
	}

	return
}

func ResponseWithJSON(res ResponseHttp) {
	response := ResponseHttp{
		Message: res.Message,
		Data:    res.Data,
	}
	js, _ := json.Marshal(response)

	res.W.Header().Set("Content-Type", "application/json")
	res.W.WriteHeader(res.StatusCode)
	res.W.Write(js)
}

func DynamicDir() string {
	_, b, _, _ := runtime.Caller(0)
	bStr := filepath.Dir(b)
	baseDir := strings.Replace(bStr, "helper", "", -1)

	return baseDir
}
func DecryptData(input string) (string, error) {
	keyEncrypt := []byte(os.Getenv("KEY_ENCRYPT"))

	value, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	h := md5.New()
	h.Write(keyEncrypt)
	key := []byte(string(h.Sum(nil)))

	dst, err := openssl.AesECBDecrypt(value, key, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}
func EnryptData(input string) (string, error) {
	keyEncrypt := []byte(os.Getenv("KEY_ENCRYPT"))
	value := []byte(input)

	h := md5.New()
	h.Write(keyEncrypt)
	key := []byte(string(h.Sum(nil)))

	dst, err := openssl.AesECBEncrypt(value, key, openssl.PKCS7_PADDING)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(dst)
	return encoded, nil
}
func GenerateContractNumber(code string) string {
	code = strings.ToUpper(code)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	appId := "C002"
	timeStamp := time.Now().In(loc).Format("060102150405.000")
	newTimeStamp := strings.Replace(timeStamp, ".", "", 1)

	lastDigit := code[len(code)-5:]
	changeableDigit := "5"

	contractNumber := appId + newTimeStamp + lastDigit + changeableDigit
	return contractNumber
}
