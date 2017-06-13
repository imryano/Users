package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/imryano/Users/userPackage"
	"github.com/imryano/utils/password"
	"github.com/imryano/utils/webservice"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

var signingString = []byte("MYrIg08fMld1zBpb8SgddbhgbVLLIxlQS7ihHkWt")

type userClaims struct {
	Username    string
	AccessLevel int
	ID          bson.ObjectId `bson:"_id,omitempty"`
	jwt.StandardClaims
}

type userCreds struct {
	Username string
	Password string
}

func main() {
	user := userPackage.User{Username: "imryano", Password: password.HashPassword("password123"), Email: "imryano@gmail.com"}
	session, err := mgo.Dial(userPackage.DbUrl)
	if err == nil {
		result := &userPackage.User{}
		col := session.DB(userPackage.DbName).C(userPackage.UserCol)
		err = col.Find(bson.M{"Username": user.Username}).One(result)
		if err == nil {
			user.ID = result.ID
			col.UpdateId(result.ID, user)
		} else {
			col.Insert(user)
		}
	}

	http.HandleFunc("/tryLogin", tryLogin)
	http.ListenAndServe(":8080", nil)
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
				errString += "Cookie already exists. Clearing. "
			} else {
				//Invalid cookie. Clear and prepare for reauthentication
				errString += "Cookie invalid. Clearing. "
			}
			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
		} else {
			fmt.Printf("\nNew Request from: %s\n", origin)

			//Check if user exists
			uc := userCreds{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&uc)

			if err != nil {
				errString += "Login Failed: Invalid Request. "
			} else {
				fmt.Printf("Login Request Received with the following creds: \nUsername: %s \nPassword: %s\n", uc.Username, uc.Password)

				session, err := mgo.Dial(userPackage.DbUrl)
				if err == nil {
					result := &userPackage.User{}
					col := session.DB(userPackage.DbName).C(userPackage.UserCol)
					err = col.Find(bson.M{"username": uc.Username}).One(result)
					if err == nil && result.Username == uc.Username {
						fmt.Printf("Request: %s, DB: %s", uc.Username, result.Username)
						if password.CheckPassword(uc.Password, result.Password) {
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
						} else {
							errString += "Authentication Failed: Username/Password combination doesn't exist. "
						}
					} else {
						errString += "Authentication Failed: Username/Password combination doesn't exist. "
					}
				} else {
					errString += "Authentication Failed: Could not connect to DB. "
				}
			}
		}
	}

	w = webservice.SetCORSProps(w, origin)
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
