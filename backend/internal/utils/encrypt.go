package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// 使用环境变量或配置的密钥，这里使用固定密钥（生产环境应该从配置读取）
var encryptionKey = []byte("kkops-encryption-key-32-bytes!!") // 32 bytes for AES-256

// Encrypt 加密字符串
func Encrypt(plaintext string) (string, error) {
	// 如果已经是加密格式（base64），直接返回
	if len(plaintext) > 0 && plaintext[0] != '{' {
		// 简单检查：如果不是 JSON 格式，可能是已加密的
		_, err := base64.StdEncoding.DecodeString(plaintext)
		if err == nil {
			// 已经是 base64，可能是已加密的
			return plaintext, nil
		}
	}

	// 生成随机 salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// 使用 PBKDF2 派生密钥
	key := pbkdf2.Key(encryptionKey, salt, 10000, 32, sha256.New)

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 使用 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// 组合 salt 和 ciphertext
	encrypted := append(salt, ciphertext...)

	// 返回 base64 编码
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt 解密字符串
func Decrypt(ciphertext string) (string, error) {
	// 解码 base64
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		// 如果不是 base64，可能是明文，直接返回
		return ciphertext, nil
	}

	if len(encrypted) < 16 {
		// 太短，可能是明文
		return ciphertext, nil
	}

	// 提取 salt
	salt := encrypted[:16]
	ciphertextBytes := encrypted[16:]

	// 使用 PBKDF2 派生密钥
	key := pbkdf2.Key(encryptionKey, salt, 10000, 32, sha256.New)

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 使用 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 提取 nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]

	// 解密
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		// 解密失败，可能是明文，直接返回
		return ciphertext, nil
	}

	return string(plaintext), nil
}

