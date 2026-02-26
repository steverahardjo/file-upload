package middleware

import(
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	uuid "github.com/google/uuid"
	"time"
)

type JwtClaim struct{
	TokenID uuid.UUID
	CreatedAt time.Time
	ExpireOn time.Time
	Freq int
}
b

func GenerateToken()
