package oauthprovider

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type OAuthProvider interface {
	GetUserEmail(idToken string) (string, error)
	getJsonWebKeys() (string, error)
	validateSignature(idToken string, signingKey interface{}) (jwt.Claims, error)
}

type oauthProviderImpl struct{}

type JsonWebKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func NewOAuthProvider() OAuthProvider {
	return &oauthProviderImpl{}
}

func (provider *oauthProviderImpl) GetUserEmail(idToken string) (string, error) {
	keyString, err := provider.getJsonWebKeys()
	if err != nil {
		return "", err
	}

	jsonWebKeys, err := provider.parseKeys(keyString)
	if err != nil {
		return "", err
	}

	kid, err := getKidFromToken(idToken)
	if err != nil {
		return "", err
	}

	var selectedKey JsonWebKey
	found := false
	for _, jsonWebKey := range jsonWebKeys {
		if kid == jsonWebKey.Kid {
			selectedKey = jsonWebKey
			found = true
			break
		}
	}

	if !found {
		return "", errors.New("key not found")
	}

	claims, err := provider.validateSignature(idToken, getRSAPublicKey(selectedKey))
	if err != nil {
		return "", err
	}

	userEmail, ok := claims.(jwt.MapClaims)["email"].(string)
	if !ok {
		return "", errors.New("email claim not found")
	}

	return userEmail, nil
}

func getKidFromToken(idToken string) (string, error) {
	header, err := getHeader(idToken)
	if err != nil {
		return "", err
	}

	decodedHeader, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return "", err
	}

	var headerJSON map[string]interface{}
	if err := json.Unmarshal(decodedHeader, &headerJSON); err != nil {
		return "", err
	}

	kid, ok := headerJSON["kid"].(string)
	if !ok {
		return "", errors.New("kid not found in token header")
	}

	return kid, nil
}

func getHeader(idToken string) (string, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid token format")
	}
	return parts[0], nil
}

func getRSAPublicKey(jwk JsonWebKey) interface{} {
	nBytes, _ := base64.RawURLEncoding.DecodeString(jwk.N)
	eBytes, _ := base64.RawURLEncoding.DecodeString(jwk.E)

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}
}

func (provider *oauthProviderImpl) parseKeys(keyString string) ([]JsonWebKey, error) {
	var keys struct {
		Keys []JsonWebKey `json:"keys"`
	}

	if err := json.Unmarshal([]byte(keyString), &keys); err != nil {
		return nil, err
	}

	return keys.Keys, nil
}

func (provider *oauthProviderImpl) getJsonWebKeys() (string, error) {
	// 여기서 실제로 키를 가져오는 구현을 해야 합니다.
	// 예: HTTP 요청으로 JWK를 가져옴
	// resp, err := http.Get("https://example.com/.well-known/jwks.json")
	// if err != nil {
	//     return "", err
	// }
	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	//     return "", err
	// }
	// return string(body), nil

	// 예제 JSON Web Key 문자열 반환 (실제 구현에서는 위의 HTTP 요청을 사용)
	return `{"keys":[{"kid":"example-kid","kty":"RSA","alg":"RS256","use":"sig","n":"...","e":"..."}]}`, nil
}

func (provider *oauthProviderImpl) validateSignature(idToken string, signingKey interface{}) (jwt.Claims, error) {
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
