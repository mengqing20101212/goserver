package main

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"logger"
	"time"
)

func main() {
	//testLogger()
	mySigningKey := []byte("ly.1006897725")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["username"] = "nacos"
	claims["password"] = "nacos"

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("Error while signing the token:", err)
		return
	}

	fmt.Println("Generated JWT token:", base64Str(tokenString))
}

func base64Str(str string) string {
	bs := []byte(str)
	return base64.StdEncoding.EncodeToString(bs)
}

func testLogger() {
	for i := 0; i < 10; i++ {
		go func() {
			log := logger.InitNull()
			log.Debug("debug")
			log.Info("info")
			log.Warn("warn")
			log.Error("error")
		}()
	}
	time.Sleep(10 * time.Second)
}
