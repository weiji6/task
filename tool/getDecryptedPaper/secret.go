package getDecryptedPaper

import (
	"encoding/base64"
)

// XOR 加密/解密函数，使用密钥对数据进行异或操作
func xorEncryptDecrypt(input, key string) string {
	keyLen := len(key)
	output := make([]byte, len(input))

	for i := range input {
		output[i] = input[i] ^ key[i%keyLen] // 按字节异或
	}

	return string(output)
}

// 解密用的函数,接收两个参数,第一个参数是加密的论文,第二个参数是解密用的秘钥,返回的是解密后的论文
func GetDecryptedPaper(encodedPaper, key string) (string, error) {
	// 先进行 base64 解码
	data, _ := base64.StdEncoding.DecodeString(encodedPaper)
	return xorEncryptDecrypt(string(data), key), nil
}
