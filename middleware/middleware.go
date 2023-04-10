package middleware

import (
	"multi-finance/helper"
	"net/http"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func ValidationToken(next http.Handler, cache *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			defer helper.CreateLog(&helper.Log{
				Event:      helper.EventMiddleware + "|Check-Token",
				StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
				Type:       helper.LVL_INFO,
				Method:     r.Method,
				Message:    "Token cannot be empty empty",
				ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
				UserAgent:  r.UserAgent(),
			}, "service")
			helper.ResponseWithJSON(helper.ResponseHttp{
				W:          w,
				StatusCode: helper.MapHttpStatusCode(helper.ErrBadRequest),
				Message:    helper.ErrBadRequest.Error(),
			})
			return
		}
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, helper.ErrUnAuhorized
			}
			return SECRET_KEY, nil
		})

		if err != nil {
			defer helper.CreateLog(&helper.Log{
				Event:      helper.EventMiddleware + "|Jwt-parse",
				StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
				Type:       helper.LVL_ERROR,
				Method:     helper.METHOD_GET,
				Message:    "Token not valid",
				ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
				UserAgent:  r.UserAgent(),
			}, "service")
			helper.ResponseWithJSON(helper.ResponseHttp{
				W:          w,
				StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
				Message:    helper.ErrUnAuhorized.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if r.URL.String() == "/api/refresh-token" {
				checkRefreshTokenToken := cache.Get("authorization-" + claims["id"].(string))
				refreshToken, _ := checkRefreshTokenToken.Result()
				if refreshToken == "" {
					defer helper.CreateLog(&helper.Log{
						Event:      helper.EventMiddleware + "|Check-Redis",
						StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
						Type:       helper.LVL_ERROR,
						Method:     r.Method,
						Message:    "Refresh Token There is no in cache",
						ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
						UserAgent:  r.UserAgent(),
					}, "service")
					helper.ResponseWithJSON(helper.ResponseHttp{
						W:          w,
						StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
						Message:    helper.ErrUnAuhorized.Error(),
					})
					return
				}

			}
			checkToken := cache.Get("authorization-" + claims["id"].(string))
			accesToken, _ := checkToken.Result()
			if accesToken == "" {
				defer helper.CreateLog(&helper.Log{
					Event:      helper.EventMiddleware + "|Check-Redis",
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Type:       helper.LVL_ERROR,
					Method:     r.Method,
					Message:    "Token There is no in cache",
					ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
					UserAgent:  r.UserAgent(),
				}, "service")
				helper.ResponseWithJSON(helper.ResponseHttp{
					W:          w,
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Message:    helper.ErrUnAuhorized.Error(),
				})
				return
			}

			params := mux.Vars(r)
			id := params["id"]
			if r.URL.String() == "/api/customer/application" && r.Method == "POST" && claims["role"].(string) != helper.ROLE_CUSTOMER {
				defer helper.CreateLog(&helper.Log{
					Event:      helper.EventMiddleware + "|Check-Permission",
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Type:       helper.LVL_INFO,
					Method:     r.Method,
					Message:    "Only " + helper.ROLE_CUSTOMER + " allowed access " + r.URL.String(),
					ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
					UserAgent:  r.UserAgent(),
				}, "service")
				helper.ResponseWithJSON(helper.ResponseHttp{
					W:          w,
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Message:    helper.ErrUnAuhorized.Error(),
				})
				return
			} else if r.URL.String() == "/api/customer/application/"+id && r.Method == "PUT" && claims["role"].(string) != helper.ROLE_ADMIN {
				defer helper.CreateLog(&helper.Log{
					Event:      helper.EventMiddleware + "|Check-Permission",
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Type:       helper.LVL_INFO,
					Method:     r.Method,
					Message:    "Only " + helper.ROLE_ADMIN + " allowed access " + r.URL.String(),
					ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
					UserAgent:  r.UserAgent(),
				}, "service")
				helper.ResponseWithJSON(helper.ResponseHttp{
					W:          w,
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Message:    helper.ErrUnAuhorized.Error(),
				})
				return
			} else if r.URL.String() == "/api/customer/application/"+id && r.Method == "GET" && claims["role"].(string) == helper.ROLE_SALES {
				defer helper.CreateLog(&helper.Log{
					Event:      helper.EventMiddleware + "|Check-Permission",
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Type:       helper.LVL_INFO,
					Method:     r.Method,
					Message:    "Role " + helper.ROLE_SALES + " not allowed access " + r.URL.String(),
					ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
					UserAgent:  r.UserAgent(),
				}, "service")
				helper.ResponseWithJSON(helper.ResponseHttp{
					W:          w,
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Message:    helper.ErrUnAuhorized.Error(),
				})
				return
			} else if r.URL.String() == "/api/transaction/partner/" && r.Method == "POST" && claims["role"].(string) != helper.ROLE_SALES {
				defer helper.CreateLog(&helper.Log{
					Event:      helper.EventMiddleware + "|Check-Permission",
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Type:       helper.LVL_INFO,
					Method:     r.Method,
					Message:    "Only " + helper.ROLE_SALES + "  allowed access " + r.URL.String(),
					ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
					UserAgent:  r.UserAgent(),
				}, "service")

				helper.ResponseWithJSON(helper.ResponseHttp{
					W:          w,
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Message:    helper.ErrUnAuhorized.Error(),
				})
				return
			} else if r.URL.String() == "/api/transaction/" && r.Method == "GET" && claims["role"].(string) != helper.ROLE_CUSTOMER {
				defer helper.CreateLog(&helper.Log{
					Event:      helper.EventMiddleware + "|Check-Permission",
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Type:       helper.LVL_INFO,
					Method:     r.Method,
					Message:    "Only " + helper.ROLE_CUSTOMER + "  allowed access " + r.URL.String(),
					ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
					UserAgent:  r.UserAgent(),
				}, "service")

				helper.ResponseWithJSON(helper.ResponseHttp{
					W:          w,
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Message:    helper.ErrUnAuhorized.Error(),
				})
				return
			} else if r.URL.String() == "/api/transaction/"+id && r.Method == "GET" && claims["role"].(string) == helper.ROLE_SALES {
				defer helper.CreateLog(&helper.Log{
					Event:      helper.EventMiddleware + "|Check-Permission",
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Type:       helper.LVL_INFO,
					Method:     r.Method,
					Message:    "Role " + helper.ROLE_SALES + " not  allowed access " + r.URL.String(),
					ClientIP:   r.Header.Get(helper.X_FORWARDED_FOR),
					UserAgent:  r.UserAgent(),
				}, "service")
				helper.ResponseWithJSON(helper.ResponseHttp{
					W:          w,
					StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
					Message:    helper.ErrUnAuhorized.Error(),
				})
				return
			} else if claims["role"].(string) == helper.ROLE_SALES {
				r.Header.Set("partner-id", claims["partner_id"].(string))

			}

			r.Header.Set("id", claims["id"].(string))
			r.Header.Set("role", claims["role"].(string))

			next.ServeHTTP(w, r)
			return

		} else {
			helper.ResponseWithJSON(helper.ResponseHttp{
				W:          w,
				StatusCode: helper.MapHttpStatusCode(helper.ErrUnAuhorized),
				Message:    helper.ErrUnAuhorized.Error(),
			})
			return
		}

	})
}
