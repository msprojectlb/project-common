package pkg

import (
	"testing"
)

func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzExMjcxOTEsInRva2VuIjoiMTAxNCJ9.nc-i2-36CIlkhTEzxV0L3JDUXVD68ASefnl33PU9PcE"
	ParseToken(tokenString, "msproject")
}
