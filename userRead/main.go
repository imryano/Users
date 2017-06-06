package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/imryano/Users/userPackage"
	"github.com/imryano/utils/password"
	"net/http"
	"time"
)

var signingString = []byte("MYrIg08fMld1zBpb8SgddbhgbVLLIxlQS7ihHkWt")

type userClaims struct {
	Username    string
	AccessLevel int
	ID          int
	jwt.StandardClaims
}

type userCreds struct {
	Username string
	Password string
}

func main() {
	http.HandleFunc("/tryLogin", tryLogin)
	http.HandleFunc("/checkUsername", checkUsername)
	http.ListenAndServe(":8080", nil)
}

func SetCORSProps(w http.ResponseWriter, origin string) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	return w
}

func checkUsername(w http.ResponseWriter, r *http.Request) {
	preflight := (r.Method == "OPTIONS")

	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:8070"
	}

	w = SetCORSProps(w, origin)

	if !preflight {
		userExists := true
		userExistsStr := "true"
		if userExists {
			userExistsStr = "true"
		} else {
			userExistsStr = "false"
		}
		w.Write([]byte("{\"userExists\":\"" + userExistsStr + "\"}"))
	}
}

func tryLogin(w http.ResponseWriter, r *http.Request) {
	errString := ""
	redirectUrl := ""

	//CORS Origin Work Around (REMOVE FOR PROD)
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:8070"
	}

	//Check if it's the CORS preflight request
	preflight := (r.Method == "OPTIONS")
	if !preflight {
		//Check if the cookie exists
		cookie, err := r.Cookie("token")

		if err == nil {
			username := CheckToken(cookie.Value).Username
			if username != "" {
				//Cookie is working fine, can redirect
				errString += "Cookie already exists. No need to recreate."
				redirectUrl = "http://www.google.com.au"
			} else {
				//Invalid cookie. Clear and prepare for reauthentication
				cookie.MaxAge = -1
				http.SetCookie(w, cookie)
				errString += "Cookie invalid. Clearing."
			}
		} else {

			fmt.Printf("\nNew Request from: %s\n", origin)

			uc := userCreds{}

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&uc)

			if err != nil {
				errString = "Login Failed: Invalid Request\n"
			}

			fmt.Printf("Login Request Received with the following creds: \nUsername: %s \nPassword: %s\n", uc.Username, uc.Password)

			cookie = &http.Cookie{
				Name:     "token",
				Value:    CreateUserCookie(&uc),
				Expires:  time.Now().Add(time.Hour),
				MaxAge:   50000,
				Secure:   false,
				HttpOnly: true,
				Raw:      "",
				Path:     "/",
			}

			http.SetCookie(w, cookie)
			errString += CheckToken(cookie.Value).Username
		}
	}

	w = SetCORSProps(w, origin)
	w.Write([]byte("{\"error\":\"" + errString + "\", \"redirectURL\":\"" + redirectUrl + "\"}"))
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

func CheckToken(tokenString string) *userClaims {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return signingString, nil
		}
	})

	if err == nil {
		if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
			return claims
		}
	}

	return &userClaims{}
}

func CreateUserCookie(uc *userCreds) string {
	user := userPackage.User{
		Username: uc.Username,
		Password: password.HashPassword(uc.Password),
		Email:    "imryano@gmail.com",
	}

	return GetToken(&user)
}
