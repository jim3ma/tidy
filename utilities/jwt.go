package utilities

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func UpdatePubKey(path string) {
	verifyBytes, err := ioutil.ReadFile(path)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func UpdatePriKey(path string) {
	signBytes, err := ioutil.ReadFile(path)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)
}

// JWTHandler is a Gin MinddleWare for JWT in tidy project
func JWTHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		switch c.Request.Method {
		case "GET":
			tokenString = c.DefaultQuery("auth_token", "")
		case "POST":
			tokenString = c.DefaultPostForm("auth_token", "")
		case "PUT":
			tokenString = c.DefaultPostForm("auth_token", "")
		case "DELETE":
			tokenString = c.DefaultQuery("auth_token", "")
		default:
			tokenString = c.DefaultPostForm("auth_token", "")
		}

		// for TESTING
		//log.Print(c)
		//log.Print(c.Keys)
		//log.Print(c.Params)
		//log.Print(c.Request)
		//log.Print(c.Request.Form)
		//log.Print(c.Request.PostForm)
		//tokenString, _ = NewToken(map[string]string{"uid": "tidy uid tidy-uid"})
		// for TESTING

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, "Empty Token")
			c.Abort()
			return
		}
		verified, token := VerifyToken(tokenString)
		if verified {
			appendParameter(c, token)
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, "Error Token")
		c.Abort()
		return
	}
}

func appendGetParameter(c *gin.Context, token *jwt.Token) {
	for key, val := range token.Claims {
		if str, ok := val.(string); ok {
			c.Request.URL.RawQuery += "&" + key + "=" + url.QueryEscape(str)
			continue
		}
		if stringer, ok := val.(fmt.Stringer); ok {
			c.Request.URL.RawQuery += "&" + key + "=" + url.QueryEscape(stringer.String())
			continue
		}
	}
	log.Printf("Auth parameter: %s", c.Request.URL.RawQuery)
}

func appendPostParameter(c *gin.Context, token *jwt.Token) {
	if c.Request.PostForm == nil {
		log.Print("nil postform data")
		c.Request.PostForm = url.Values{}
	}
	//c.Request.PostForm.Set("uid", token.Claims["uid"].(string))
	//c.Request.PostForm.Set("user_name", token.Claims["user_name"].(string))
	for key, val := range token.Claims {
		if str, ok := val.(string); ok {
			c.Request.PostForm.Set(key, str)
			continue
		}
		if stringer, ok := val.(fmt.Stringer); ok {
			c.Request.PostForm.Set(key, stringer.String())
			continue
		}
	}
}
func appendParameter(c *gin.Context, token *jwt.Token) {
	switch c.Request.Method {
	case "GET":
		//log.Print(c.Request.Form)
		//log.Print(c.Request)

		//if c.Request.Form == nil {
		//	log.Print("nil form data")
		//	c.Request.Form = url.Values{}
		//}

		//c.Request.Form.Set("uid", token.Claims["uid"].(string))
		// hard code
		// TBD

		//c.Request.URL.RawQuery +=
		//	"&" + "uid" + "=" + url.QueryEscape(token.Claims["uid"].(string)) +
		//		"&" + "user_name" + "=" + url.QueryEscape(token.Claims["user_name"].(string))
		appendGetParameter(c, token)
		//uid := c.DefaultQuery("uid", "none")
		//log.Print(uid)
	case "POST":
		//log.Print(c.Request.PostForm)
		//log.Print(c.Request)
		appendPostParameter(c, token)
		//uid := c.DefaultPostForm("uid", "none")
		//log.Print(uid)
	case "PUT":
		appendPostParameter(c, token)
	case "DELETE":
		appendGetParameter(c, token)
	default:
		appendGetParameter(c, token)
	}
}

// NewToken generate a jwt token
func NewToken(values map[string]string) (tokenString string, err error) {
	// create a signer for rsa 256
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	log.Printf("NewToken values: %s", values)
	for key, val := range values {
		// set our claims
		token.Claims[key] = val
	}

	// set the expire time
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	token.Claims["exp"] = time.Now().Add(time.Minute * 120).Unix()
	tokenString, err = token.SignedString(signKey)
	log.Printf("New token: %s", tokenString)
	return
}

// VerifyToken check the input token string, and return *jwt.Token for other use
func VerifyToken(tokenString string) (bool, *jwt.Token) {
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
			log.Print("Invalid Token")
			return false, nil
		}
		//log.Printf("Verified! Token:%+v\n", token)
		return true, token
	case *jwt.ValidationError: // something was wrong during the validation
		vErr := err.(*jwt.ValidationError)
		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			log.Print("Token expired")
			return false, nil
		default:
			log.Printf("ValidationError error: %+v\n", vErr.Errors)
			return false, nil
		}
	default: // something else went wrong
		log.Printf("Token parse error: %v\n", err)
		return false, nil
	}
}
