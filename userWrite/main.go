package main

import (
	"fmt"
	"github.com/imryano/Users/userPackage"
	"github.com/imryano/utils/password"
)

func main() {
	history := []interface{}{
		userPackage.CreateUser{Username: "imryano", Password: password.HashPassword("password123"), Email: "imryano@gmail.com"},
	}

	aggregate := userPackage.NewUserFromHistory(history)
	fmt.Println("Before Promotion")
	fmt.Println(aggregate)

	aggregate.PromoteUser()
	aggregate.PromoteUser()
	fmt.Println("After Promotion")
	fmt.Println(aggregate)
}
