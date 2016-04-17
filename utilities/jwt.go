package utilities

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// location of the files used for signing and verification
const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

// keys are held in global variables
// i havn't seen a memory corruption/info leakage in go yet
// but maybe it's a better idea, just to store the public key in ram?
// and load the signKey on every signing request? depends on  your usage i guess
var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// read the key files before starting http handlers
func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func JWTHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		switch c.Request.Method {
		case "GET":
			token = c.DefaultQuery("token", "")
		case "POST":
			token = c.DefaultPostForm("token", "")
		default:
			token = ""
		}
		if token == "" {
			//c.JSON(http.StatusUnauthorized, "Empty Token")
			token, _ = NewToken()
			VerifyToken(token)
		}
		c.Next()
		return
	}
}

func NewToken() (tokenString string, err error) {
	// create a signer for rsa 256
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	// set our claims
	token.Claims["AccessToken"] = "level1"

	// set the expire time
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	token.Claims["exp"] = time.Now().Add(time.Minute * 120).Unix()
	tokenString, err = token.SignedString(signKey)
	log.Printf("New token: %s", tokenString)
	return
}

func VerifyToken(tokenString string) bool {
	// validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})

	// branch out into the possible error from signing
	switch err.(type) {
	case nil: // no error
		if !token.Valid { // but may still be invalid
			fmt.Println("WHAT? Invalid Token? F*** off!")
			return false
		}
		// see stdout and watch for the CustomUserInfo, nicely unmarshalled
		log.Printf("Someone accessed resricted area! Token:%+v\n", token)
		return true
	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)
		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			fmt.Println("Token Expired, get a new one.")
			return false
		default:
			log.Printf("ValidationError error: %+v\n", vErr.Errors)
			return false
		}
	default: // something else went wrong
		log.Printf("Token parse error: %v\n", err)
		return false
	}
}
