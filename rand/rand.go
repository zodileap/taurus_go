package rand

import (
	"math/rand"
	"time"
)

// 随机因子
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// 生成随机字符串
//
// 参数：
//   - length：生成字符串的长度
//   - charset：生成字符串的字符集
func StringWithCharset(length int, charset string) (string, error) {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b), nil
}
