package auth

import (
	"encoding/json"
	"strings"

	jwt "github.com/videocoin/jwt-go"
)

// note: the original method contains a lot of logic that is not required.
func parseHeader(tokenStr string) (*jwt.Token, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, jwt.NewValidationError("token contains an invalid number of segments", jwt.ValidationErrorMalformed)
	}

	token := new(jwt.Token)

	// parse Header
	var (
		headerBytes []byte
		err         error
	)
	if headerBytes, err = jwt.DecodeSegment(parts[0]); err != nil {
		if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
			return nil, jwt.NewValidationError("tokenstring should not contain 'bearer '", jwt.ValidationErrorMalformed)
		}
		return nil, &jwt.ValidationError{Inner: err, Errors: jwt.ValidationErrorMalformed}
	}
	if err = json.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, &jwt.ValidationError{Inner: err, Errors: jwt.ValidationErrorMalformed}
	}

	return token, nil
}
