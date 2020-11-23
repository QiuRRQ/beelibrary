package models

import (
	t "city/mytoken"
	u "city/utils"
	"crypto/rand"
	_ "crypto/rand"
	"crypto/rsa"
	_ "crypto/rsa"
	"crypto/x509"
	_ "crypto/x509"
	"encoding/json"
	_ "encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	"github.com/lestrrat/go-jwx/jwa"
	_ "github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwe"
	_ "github.com/lestrrat/go-jwx/jwe"
	uuid "github.com/satori/go.uuid"
)

type Users struct {
	Id         int    `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	Address    string `json:"address"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type UserLogin struct {
	UserData  *Users         `json:"data_user"`
	TokenData UserJwtTokenVm `json:"token_data"`
}

type UserSessionVm struct {
	Session string `json:"session"`
}

type UserJwtTokenVm struct {
	Token           string `json:"token"`
	ExpTime         string `json:"exp_time"`
	RefreshToken    string `json:"refresh_token"`
	ExpRefreshToken string `json:"exp_refresh_token"`
	UserID          string `json:"user_id"`
}

func (user *Users) Login(email, password string, mydb *gorm.DB) UserLogin {

	if _, ok := user.Validate(); !ok {
		return UserLogin{
			UserData:  nil,
			TokenData: UserJwtTokenVm{},
		}
	}
	foundUser := &Users{}
	err := mydb.Table("usr").Where("email = ? AND password = ?", email, password).First(foundUser).Error
	if err != nil {
		fmt.Println(err)
		return UserLogin{
			UserData:  nil,
			TokenData: UserJwtTokenVm{},
		}
	}

	jwePayload, _ := GenerateJwePayload(strconv.Itoa(foundUser.Id))

	session, _ := UpdateSessionLogin(strconv.Itoa(foundUser.Id))

	token, refreshToken, tokenExpiredAt, refreshTokenExpiredAt, err := GenerateJwtToken(jwePayload, strconv.Itoa(foundUser.Id), session)
	userToken := UserJwtTokenVm{
		Token:           token,
		RefreshToken:    refreshToken,
		ExpTime:         tokenExpiredAt,
		ExpRefreshToken: refreshTokenExpiredAt,
		UserID:          strconv.Itoa(foundUser.Id),
	}
	res := UserLogin{
		UserData:  user,
		TokenData: userToken,
	}

	return res
}

func (user *Users) Validate() (map[string]interface{}, bool) {
	if user.Email == "" {
		return u.Message(false, "Email tidak boleh kosong!"), false
	}

	if user.Password == "" {
		return u.Message(false, "Password tidak boleh kosong!"), false
	}

	return u.Message(true, "success"), true
}

// GenerateJwtToken ...
func GenerateJwtToken(jwePayload, userid, session string) (token, refreshToken, expTokenAt, expRefreshTokenAt string, err error) {

	token, expTokenAt, err = t.GetToken(session, userid, jwePayload)

	if err != nil {
		return token, refreshToken, expTokenAt, expRefreshTokenAt, err
	}

	refreshToken, expRefreshTokenAt, err = t.GetRefreshToken(session, userid, jwePayload)
	if err != nil {
		return token, refreshToken, expTokenAt, expRefreshTokenAt, err
	}

	return token, refreshToken, expTokenAt, expRefreshTokenAt, err
}

// Generate ...
func GenerateJwePayload(id string) (res string, err error) {

	privkey, err := rsaConfigSetup("", "")
	if err != nil {
		return res, err
	}

	// Generate payload
	payload := map[string]interface{}{
		"id": id,
	}
	payloadString, err := json.Marshal(payload)
	if err != nil {
		return res, err
	}

	// Generate JWE
	jweRes, err := jwe.Encrypt([]byte(payloadString), jwa.RSA1_5, &privkey.PublicKey, jwa.A128CBC_HS256, jwa.Deflate)
	res = string(jweRes)

	return res, err
}

func rsaConfigSetup(rsaPrivateKeyLocation, rsaPrivateKeyPassword string) (*rsa.PrivateKey, error) {
	if rsaPrivateKeyLocation == "" {
		fmt.Println("No RSA Key given, generating temp one")
		return GenRSA(4096)
	}

	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		fmt.Println("No RSA private key found, generating temp one")
		return GenRSA(4096)
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {
		fmt.Println("RSA private key is of the wrong type")
	}

	if rsaPrivateKeyPassword != "" {
		privPemBytes, err = x509.DecryptPEMBlock(privPem, []byte(rsaPrivateKeyPassword))
	} else {
		privPemBytes = privPem.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			fmt.Println("Unable to parse RSA private key, generating a temp one")
			return GenRSA(4096)
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		fmt.Println("Unable to parse RSA private key, generating a temp one (2)")
		return GenRSA(4096)
	}

	return privateKey, nil
}

// GenRSA returns a new RSA key of bits length
func GenRSA(bits int) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	return key, err
}

// UpdateSessionLogin ...
func UpdateSessionLogin(ID string) (res string, err error) {
	value := uuid.NewV4().String()
	exp := os.Getenv("SESSION_EXP")
	key := "session-" + ID
	resSession := UserSessionVm{}
	resSession.Session = value

	StoreToRedistWithExpired(key, resSession, exp)

	return value, err
}

func StoreToRedistWithExpired(key string, val interface{}, duration string) error {
	var Redis *redis.Client
	dur, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = Redis.Set(key, string(b), dur).Err()

	return err
}

func (user *Users) CreatUser(mydb *gorm.DB) (map[string]interface{}, *Users) {

	err := mydb.Table("usr").Create(&user).Error
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	resp := u.Message(true, "success")
	resp["data"] = user

	return resp, user
}
