package lib

import (
	"errors"
	"os"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

func secretKey() []byte {
	var JWT_SECRET []byte = []byte(Md5Hash(os.Getenv("PPAY_KEY")))
	return JWT_SECRET
}

type BaseInfo struct {
	IssuedAt jwt.NumericDate `json:"iat"` // Waktu token diterbitkan
}

type TokenPayload struct {
	UserId int `json:"userId"`
}

type FullToken struct {
	BaseInfo
	TokenPayload
}

func GenerateToken(userId int) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: secretKey()}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}
	tokenData := FullToken{
		BaseInfo: BaseInfo{
			IssuedAt: *jwt.NewNumericDate(time.Now()),
		},
		TokenPayload: TokenPayload{
			UserId: userId,
			// UserRole: userRole,
		},
	}
	token, err := jwt.Signed(sig).Claims(tokenData).Serialize()
	if err != nil {
		return "", nil
	}
	return token, nil
}

func VerifyToken(token string) (*TokenPayload, error) {
	tok, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.HS256})
	// fmt.Println(tok)
	if err != nil {
		return nil, err
	}

	claims := &TokenPayload{}
	err = tok.Claims(secretKey(), &claims)
	if err != nil {
		return nil, errors.New("token is invalid or expired")
	}
	return claims, nil
}
