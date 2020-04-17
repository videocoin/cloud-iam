package security

import (
	"fmt"
	"net/http"
	"strings"
)

func authFromReq(r *http.Request, expectedScheme string) (string, error) {
	val := r.Header.Get(headerAuthorize)
	if val == "" {
		return "", fmt.Errorf("Request unauthenticated with " + expectedScheme)

	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", fmt.Errorf("Bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", fmt.Errorf("Request unauthenticated with " + expectedScheme)
	}
	return splits[1], nil
}
