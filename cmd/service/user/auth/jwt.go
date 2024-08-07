package auth

import (
	"context"
	"fmt"

	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vinicius77/ecom/config"
	"github.com/vinicius77/ecom/types"
	"github.com/vinicius77/ecom/utils"
)

type contextKey string

const UserKey contextKey = "userId"

func CreateJWT(secret []byte, userID int) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(int(userID)),
		"exp":    time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequestHeader(r)
		token, err := validateToken(tokenString)

		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// get userId from the token
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)

		userID, _ := strconv.Atoi(str)
		user, err := store.GetUserById(userID)

		if err != nil {
			log.Printf("failed to get user with id: %v", userID)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequestHeader(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(context context.Context) int {
	userID, ok := context.Value(UserKey).(int)

	if !ok {
		return -1
	}

	return userID
}
