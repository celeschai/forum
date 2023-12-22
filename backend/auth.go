package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// func withJWTAuth(handlerFunc http.HandlerFunc, s Database) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
func JWTAuth(w http.ResponseWriter, r *http.Request, s Database) error {
		fmt.Println("calling JWT auth middleware")

		//Get returns "" if not found 
		tokenString := r.Header.Get("jwtToken") 
		token, err := validateJWT(tokenString)
		switch {
			case err != nil: 
				return permissiondenied(w, err)
			case !token.Valid: 
				return permissiondenied(w, err)
		}

		//checking if token belongs to this account
		userID, err := strconv.Atoi(r.Header.Get("userid"))
		if err != nil { return permissiondenied(w, err) }

		account, err := s.GetAccountByUserID(userID)
		if err != nil { return permissiondenied(w, err) }

		claims := token.Claims.(jwt.MapClaims)
		switch {
			case time.Now().After(claims["expiresAt"].(jwt.NumericDate).Time):
				return permissiondenied(w, err)
			case account.UserID != int(claims["userID"].(float64)): 
				return permissiondenied(w, err)
		}

		return nil
}


func validateJWT(tokenString string) (*jwt.Token, error) { //validate jwt token then the secret key
	jwtSecret := os.Getenv("JWT_SECRET")
	return jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})
}

func createJWT(acc *Account) (string, error) {
	claim := &jwt.MapClaims{
		"expiresAt": jwt.NewNumericDate(time.Now().Add(time.Hour * 672)),
		"userID":    acc.UserID} //std jwt.Claim changes it to float64

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(jwtSecret))
}

func permissiondenied(w http.ResponseWriter, err error) error {
	log.Default().Println(err) //don't show error to user for security reasons
	return WriteJSON(w, http.StatusForbidden, ApiError{Error: "permission denied"})
}
