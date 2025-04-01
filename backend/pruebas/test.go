package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password, err := bcrypt.GenerateFromPassword([]byte("miClaveSegura"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generando el hash:", err)
		return
	}

	fmt.Println("Hash de la contraseña:", string(password))
}
