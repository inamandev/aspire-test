package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	signingKey         string
	signingKeyEnvExist bool
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	if signingKey, signingKeyEnvExist = os.LookupEnv("JWT_SIGNING_KEY"); !signingKeyEnvExist {
		log.Println("JWT_SIGNING_KEY env variable is missing")
		// panic("JWT_SIGNING_KEY env variable is missing")
	}
}

type UserTokenClaims struct {
	Id     int64  `json:"id"`
	Role   string `json:"role"`
	Status bool   `json:"status"`
	jwt.RegisteredClaims
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *UserAuth) AuthUser() (bool, error) {
	return true, nil
}

func (u *UserTokenClaims) CreateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, u)
	return token.SignedString([]byte(signingKey))
}

// validate token here
func (u *UserTokenClaims) ValidateToken(signedToken string) bool {
	var ok bool
	token, err := jwt.ParseWithClaims(signedToken, u, func(t *jwt.Token) (interface{}, error) {
		return []byte([]byte(signingKey)), nil
	})
	if err != nil {
		log.Println("error parsing jwt token", err)
		return false
	}
	u, ok = token.Claims.(*UserTokenClaims)
	if !ok {
		log.Println("not able to parse the token claims")
		return false
	}
	return u.ExpiresAt.After(time.Now().Local())
}
