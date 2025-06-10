package helper

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lyyym/zinx-wsbase/release/mdbc_server/internal/define"
	uuid "github.com/satori/go.uuid"
)

type UserClaims struct {
	UserName    string `json:"userName"`
	AccountType uint32 `json:"accountType"`
	NickName    string `json:"nickName"`
	Ip          string `json:"ip"`
	jwt.StandardClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GetUUID
// 生成唯一码
func GetUUID() string {
	return uuid.NewV4().String()
}

// GenerateToken
// 生成 token
func GenerateToken(userName string, accountType uint32, nickname string, ip string) (string, error) {
	UserClaim := &UserClaims{
		UserName:       userName,
		AccountType:    accountType,
		NickName:       nickname,
		Ip:             ip,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(define.MyKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return define.MyKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}
