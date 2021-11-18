package ginjwt

import (
	"time"

	"git.xenonstack.com/util/golang-boilerplate/rest-example/src/dbtypes"
	"go.uber.org/zap"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// GinJwtToken is a method to generate new token with expiry
func GinJwtToken(acs dbtypes.User) map[string]interface{} {

	// intializing middleware
	mw := MwInitializer()

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	// extracting claims in form of map
	claims := token.Claims.(jwt.MapClaims)

	// setting payload
	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(acs.Email) {
			claims[key] = value
		}
	}
	zap.S().Info(claims)

	// extracting expire time
	expire := mw.TimeFunc().Add(mw.Timeout)

	// setting claims
	claims["email"] = acs.Email
	claims["sys_role"] = acs.Role
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	mapd := map[string]interface{}{"token": "", "expire": ""}

	// signing token
	tokenString, err := token.SignedString(mw.Key)
	if err != nil {
		return mapd
	}

	// passing map eith all information
	mapd = map[string]interface{}{"error": false,
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339)}

	return mapd
}
