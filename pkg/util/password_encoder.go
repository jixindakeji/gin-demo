package util

import (
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
)

var options = password.Options{SaltLen: 10, Iterations: 10000, KeyLen: 50, HashFunction: md5.New}

func GetSaltAndEncodedPassword(pwd string) (string, string) {
	salt, encodedPwd := password.Encode(pwd, &options)
	return salt, encodedPwd
}

func VerifyRawPassword(rawPwd, encodedPwd, salt string) bool {
	return password.Verify(rawPwd, salt, encodedPwd, &options)
}
