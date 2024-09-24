package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	testCases := []struct {
		name      string
		obj       Cryptor
		clearTest string
	}{
		{
			name:      "len 16",
			obj:       NewCrypto("testqwer12345678"),
			clearTest: "hello world美国",
		},
		{
			name:      "len 24",
			obj:       NewCrypto("testqwer12345678testqwer"),
			clearTest: "hello 「」{}[]」world!@#$%^&*()_+_",
		},
		{
			name:      "len 32",
			obj:       NewCrypto("testqwer12345678testqwer12345678"),
			clearTest: "hello world!@#$%^&*()_+_<>?;'\"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypt := tc.obj.Encrypt(tc.clearTest)
			decrypt := tc.obj.Decrypt(encrypt)
			assert.Equal(t, tc.clearTest, decrypt)
		})
	}
}

func TestMd5(t *testing.T) {
	testCases := []struct {
		name string
		data string
		want string
	}{
		{
			name: "empty",
			data: "",
			want: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name: "hello",
			data: "hello",
			want: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name: "Chinese",
			data: "美国",
			want: "1efc978309b2a7a32b3c8db1bcc5cf58",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := Md5(tc.data)
			assert.Equal(t, tc.want, res)
		})
	}
}
