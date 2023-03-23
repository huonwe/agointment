package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var nowDate = time.Now().Format("2006-01-02 15")
var secret = fmt.Sprintf("%v%v", nowDate, "dF13ayA")

type MapClaims map[string]interface{}
type StrStr map[string]string

// GenerateToken 生成Token值
func GenerateToken(mapClaims jwt.MapClaims, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	return token.SignedString([]byte(key))
}

// token: "eyJhbGciO...解析token"
func ParseToken(tokenString string, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Println(claims["UserID"])
		return claims, nil
	} else {
		return nil, fmt.Errorf("token unexcepted")
	}

	// return claim.Claims.(jwt.MapClaims)["cmd"].(string), nil
}

type Response map[string]interface{}

func tokenReader(req *http.Request) map[string]interface{} {
	token, _ := req.Cookie("token")
	q, _ := ParseToken(token.Value, secret) // 解析token
	return q
}

func UIDReader(q map[string]interface{}) uint64 {
	UserID := q["UserID"].(string)
	uid, _ := strconv.Atoi(UserID)
	return uint64(uid)
}

func postReader(req *http.Request) map[string]string {
	con, _ := ioutil.ReadAll(req.Body)
	data := make(StrStr)
	_ = json.Unmarshal(con, &data)
	return data
}

func handle(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
