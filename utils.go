package rebel

import (
	"github.com/google/uuid"
)

func GenerateNonce() string {
	nonce := uuid.New()

	return nonce.String()
}
