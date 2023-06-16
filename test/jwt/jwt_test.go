package jwt_test

import (
	"testing"

	"github.com/nfwGytautas/gdev/jwt"
)

func TestGenerateToken(t *testing.T) {
	jwt.APISecret = "TEST_SECRET"
	token, err := jwt.GenerateToken(123, "TestRole")
	if err != nil {
		t.Error(err)
		return
	}

	// Try parse
	tokenInfo, err := jwt.ParseToken(token)
	if err != nil {
		t.Error(err)
		return
	}

	if tokenInfo.ID != 123 {
		println("ID doesn't match")
		t.Fail()
	}

	if tokenInfo.Role != "TestRole" {
		println("Role doesn't match")
		t.Fail()
	}
}
