package tool

import "testing"

func TestEncPassword(t *testing.T) {
	password := "123456"

	result := EncryptPassword(password, RandomString(8), 50000)

	t.Logf("sum => %s", result)
}
