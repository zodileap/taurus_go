package byteutil

// ReplaceAll 替换字节切片中的特定字节
//
// Params:
//
//   - slice: 原始字节切片。
//   - old: 要替换的字节。
//   - new: 替换后的字节。
func ReplaceAll(slice []byte, old byte, new []byte) []byte {
	// 创建一个新的字节切片来存储结果
	result := make([]byte, len(slice))
	copy(result, slice) // 复制原始数据到新切片

	// 遍历所有字节，替换符合条件的字节
	for i := 0; i < len(result); i++ {
		if result[i] == old {
			// 如果找到了要替换的字节，将新字节插入到原始字节切片中
			result = append(result[:i], append(new, result[i+1:]...)...)
		}
	}

	return result
}
