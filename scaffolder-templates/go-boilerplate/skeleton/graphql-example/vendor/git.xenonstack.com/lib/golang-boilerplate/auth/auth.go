package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// defining structure for user defined jwt claims
type JWTClaims struct {
	Claim map[string]interface{} `json:"claim"`
	jwt.StandardClaims
}

// function for creating new signed token
func NewToken(claim map[string]interface{}) (string, error) {
	//parsing expiration time
	d, err := time.ParseDuration(os.Getenv("EXPIRATION_TIME"))
	if err != nil {
		// if any error during parsing duration
		return "", err
	}
	// setting userdefined claims in adition to standard claims
	claims := JWTClaims{
		claim,
		jwt.StandardClaims{
			// setting expirattion time of token
			ExpiresAt: time.Now().Add(d).Unix(),
		},
	}

	// creating token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// return signed token and signing with private key
	return token.SignedString([]byte(os.Getenv("PRIVATE_KEY")))
}

// function for fetching claims and validate the token
// parameter is context and it is use to fetch metadata
func ValidateToken(ctx context.Context) (map[string]interface{}, string, error) {
	// fetch metadata from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	log.Println(ok)
	if !ok {
		// if fetch failed
		return nil, "", errors.New("error in reading metadata")
	}

	// checking authorization is set in headers
	value, ok := md["authorization"]
	log.Println(ok)
	if !ok {
		return nil, "", errors.New("Please set authorization in header")
	}

	// checking authorization is of bearer type
	if !strings.HasPrefix(value[0], "Bearer ") {
		return nil, "", errors.New("Invalid Authorization Header.")
	}

	//fetching only token from whole string
	token := strings.TrimPrefix(value[0], "Bearer ")

	// parsing token and checking its validity
	rtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok = token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("PRIVATE_KEY")), nil
	})

	// if any err return nil claims
	if err != nil {
		return nil, "", err
	}

	// extracting claims from parsed token
	claims := rtoken.Claims.(jwt.MapClaims)
	zap.S().Info(claims)
	// type converion
	val, ok := claims["claim"]
	if ok {
		if reflect.TypeOf(val).String() == "map[string]interface {}" {
			filtered_claims := val.(map[string]interface{})
			// passing final claims
			return filtered_claims, token, err
		}
	}
	return claims, token, err
}

// function for fetching claims and validate the token
// parameter is token string
func ValidateTokenString(token string) (map[string]interface{}, error) {
	// checking authorization is of bearer type
	if !strings.HasPrefix(token, "Bearer ") {
		return nil, errors.New("Invalid Authorization Header.")
	}

	//fetching only token from whole string
	token = strings.TrimPrefix(token, "Bearer ")
	log.Println(token)
	// parsing token and checking its validity
	rtoken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("PRIVATE_KEY")), nil
	})
	// if any err return nil claims
	if err != nil {
		return nil, err
	}
	// extracting claims from parsed token
	claims := rtoken.Claims.(jwt.MapClaims)
	zap.S().Info(claims)
	// type converion
	val, ok := claims["claim"]
	if ok {
		if reflect.TypeOf(val).String() == "map[string]interface {}" {
			filtered_claims := val.(map[string]interface{})
			// passing final claims
			return filtered_claims, err
		}
	}
	return claims, err
}
