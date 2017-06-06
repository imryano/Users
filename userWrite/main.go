package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/imryano/Users/userPackage"
	"github.com/imryano/utils/password"
	"time"
)

var signingString = []byte("MYrIg08fMld1zBpb8SgddbhgbVLLIxlQS7ihHkWt")

type userClaims struct {
	Username    string
	AccessLevel int
	ID          int
	jwt.StandardClaims
}

func main() {
	history := []interface{}{
		userPackage.CreateUser{Username: "imryano", Password: password.HashPassword("password123"), Email: "imryano@gmail.com"},
	}

	aggregate := userPackage.NewUserFromHistory(history)
	fmt.Println("Before Promotion")
	fmt.Println(aggregate)

	fmt.Println(GetToken(aggregate))

	aggregate.PromoteUser()
	aggregate.PromoteUser()
	fmt.Println("After Promotion")
	fmt.Println(aggregate)

	fmt.Println(GetToken(aggregate))
	fmt.Println(CheckToken(GetToken(aggregate)))

}

func GetToken(user *userPackage.User) string {
	claims := CreateClaims(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingString)

	if err == nil {
		return ss
	} else {
		return ""
	}
}

func CreateClaims(user *userPackage.User) userClaims {
	claims := userClaims{}
	claims.Username = user.Username
	claims.AccessLevel = int(user.AccessLevel)
	claims.ID = user.ID

	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(50 * time.Duration(time.Hour)).Unix(),
		Issuer:    "imryanoUser",
	}

	return claims
}

func CheckToken(tokenString string) *userClaims {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingString, nil
	})

	if err == nil {
		if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
			return claims
		}
	}

	return &userClaims{}
}
