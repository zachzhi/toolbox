package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type AuthClient struct {
	Username string
	Password string
	ExpDay   time.Duration
}

type TokenMid struct {
	Aeris string `json:"aeris"`
	Exp   string `json:"exp"`
}

func New(username, password string, day time.Duration) *AuthClient {
	return &AuthClient{
		Username: username,
		Password: password,
		ExpDay:   day,
	}
}

func GetUsernameFromToken(token string) ([]byte, error) {
	tarr := strings.Split(token, ".")
	if len(tarr) != 3 {
		return []byte{}, errors.New("Invalid token.")
	}

	bt, err := jwt.DecodeSegment(tarr[1])
	if err != nil {
		return []byte{}, err
	}

	return bt, nil
}

func Verify(token string, key []byte) error {
	_, err := jwt.Parse(token, func(tokenObj *jwt.Token) (interface{}, error) {
		if _, ok := tokenObj.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tokenObj.Header["alg"])
		}

		return key, nil
	})

	return err
}

func (this *AuthClient) Key(l int) []byte {
	stMap := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(stMap)
	resul := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		resul = append(resul, bytes[r.Intn(len(bytes))])
	}
	return resul
}

func (this *AuthClient) Produce() (token string, key []byte, err error) {

	jwtObj := jwt.New(jwt.SigningMethodHS256)
	jwtObj.Claims["exp"] = time.Now().Add(this.ExpDay)
	jwtObj.Claims["aeris"] = strconv.Itoa(int(time.Now().UnixNano())) + this.Username

	key = this.Key(16)
	token, err = jwtObj.SignedString(key)
	return
}
