package utils

import (
	"encoding/base64"
	"strings"
)

// DecodeBase64 解码 Base64 字符串（支持标准和 URL 安全编码）
func DecodeBase64(s string) (string, error) {
	// 去除空白字符
	s = strings.TrimSpace(s)

	// 补齐 padding
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}

	// 尝试标准 Base64 解码
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		return string(decoded), nil
	}

	// 尝试 URL 安全 Base64 解码
	decoded, err = base64.URLEncoding.DecodeString(s)
	if err == nil {
		return string(decoded), nil
	}

	// 尝试 RawStdEncoding（无 padding）
	s = strings.TrimRight(s, "=")
	decoded, err = base64.RawStdEncoding.DecodeString(s)
	if err == nil {
		return string(decoded), nil
	}

	// 尝试 RawURLEncoding（无 padding）
	decoded, err = base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

// EncodeBase64 编码为 Base64 字符串
func EncodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// IsBase64 判断字符串是否为 Base64 编码
func IsBase64(s string) bool {
	// 去除空白字符
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return false
	}

	// Base64 字符集
	base64Chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	urlSafeChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_="

	isStandard := true
	isURLSafe := true

	for _, c := range s {
		if !strings.ContainsRune(base64Chars, c) {
			isStandard = false
		}
		if !strings.ContainsRune(urlSafeChars, c) {
			isURLSafe = false
		}
	}

	return isStandard || isURLSafe
}
