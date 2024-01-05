package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"golang.org/x/crypto/bcrypt"
)

func JWTAuth(w http.ResponseWriter, r *http.Request, s Database) error {
	fmt.Println("calling JWT auth middleware")

	cookie, cookerr := r.Cookie("jwtToken")
	if cookerr != nil {
		return cookerr //don't write error into JSON response for security
	}
	tokenString := cookie.Value

	token, err := validateJWT(tokenString)
	switch {
	case err != nil:
		return err
	case !token.Valid:
		return err
	}

	userName, ckerr := r.Cookie("userName")
	if ckerr != nil {
		return ckerr
	}

	acc, err := s.GetAccountByUserName(userName.Value)
	if err != nil {
		return err
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["userName"] != acc.UserName || claims["userCreated"] != acc.Created {
		return err
	}
	
	fmt.Println("JWT Authenticated")
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
		"userName":    acc.UserName,
		"userCreated": acc.Created,
	} //std jwt.Claim changes it to float64

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(jwtSecret))
}

func setCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   604800,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

func deleteCookie(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:     name,
		MaxAge:   -1,
	}
	http.SetCookie(w, &cookie)
}

func (a *Account) ValidatePassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.EncryptedPW), []byte(pw))

	return err == nil
}