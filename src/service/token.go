package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	Username string
	Id       int64
	jwt.StandardClaims
}

//根据用户名和Id生成token
func GetToken(username string, user_id int64) (token string, err error) {
	token, err = ReleaseToken(username, user_id)
	return
}

//验证token的有效性
func TokenValidity(tokenString string) (int64, bool) {
	token, claims, err := ParseToken(tokenString)
	if err == nil && token.Valid {
		return claims.Id, true
	} else {
		return 0, false
	}

}

//token分发
func ReleaseToken(username string, user_id int64) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) //token截至有效时间
	claims := &Claims{
		Username: username,
		Id:       user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(), //发放的时间
			Issuer:    "DouShen",         //谁发放的
			Subject:   "user token",
		},
	}
	//token分三部分，加密协议.储存的信息.前面两部分加上key再哈希得到的值
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//token解析
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
