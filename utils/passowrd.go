package utils

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

const cryptoSha512 = "sha512"
const algorithmPbkdf2 = "pbkdf2"
const algorithmScrypt = "scrypt"
const algorithmBcrypt = "bcrypt"

func GeneratePassword(input string) string {
	// Using custom options
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(input, options)
	// 密码存储: $算法$盐$加密后密码
	return fmt.Sprintf("$%s-%s$%s$%s", algorithmPbkdf2, cryptoSha512, salt, encodedPwd)
}

// CheckPassword 验证密码
func CheckPassword(input string, pwd string) bool {
	// 字符串切割成切片; 0=空;1=算法;2=盐;3=加密后密码
	pwdInfo := strings.Split(pwd, "$")
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	return password.Verify(input, pwdInfo[2], pwdInfo[3], options)
}
