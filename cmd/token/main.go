package main

import (
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"notify/app"
	"time"
)

func main() {
	userId := flag.String("userId", "1234", "User ID, string")
	channel := flag.String("channel", "test", "Channel to subscribe to, string")
	expires := flag.Int64("exp", time.Now().Unix()+60*60*24, "Expiration time, default now+24hours")
	secret := flag.String("secret", "secret", "secret")

	flag.Parse()

	session := &app.Session{
		UserId:         *userId,
		Channel:        *channel,
		StandardClaims: jwt.StandardClaims{ExpiresAt: *expires},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	ss, err := token.SignedString([]byte(*secret))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v", ss)
	}

}
