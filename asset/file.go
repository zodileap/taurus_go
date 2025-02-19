package asset

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/imports"
)

// FileExists 检查文件是否存在。
//
// Params:
//
//   - filePath: 文件路径。
//
// Returns:
//
//	0: 文件是否存在。
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func DirExists(dir, name string) (bool, error) {
	// 组合完整的文件夹路径
	fullPath := filepath.Join(dir, name)

	// 获取文件信息
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件夹不存在
			return false, nil
		}
		// 其他错误
		return false, err
	}

	// 检查是否为目录
	return info.IsDir(), nil
}

// ReadFileToBuffer 读取文件内容到缓冲区。
//
// Params:
//
//   - filename: 文件名。
//
// Returns:
//
//   - 文件内容。
//   - 错误信息。
func ReadFileToBuffer(filePath string) (*bytes.Buffer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}

// FileOperator 文件操作者
//
// 这个结构体用于操作文件，比如读取文件内容，查找标记位置，插入代码等。
// 仅对文本文件有效。
type FileOperator struct {
	FilePath string
	lines    []string
}

// NewFileOperator 创建一个新的 FileOperator 实例
func NewFileOperator(filePath string) (*FileOperator, error) {
	o := &FileOperator{
		FilePath: filePath,
	}

	if !FileExists(filePath) {
		return nil, Err_0200020006.Sprintf(filePath)
	}

	if err := o.ReadFile(); err != nil {
		return nil, err
	}
	return o, nil
}

func (fo FileOperator) String() string {
	return strings.Join(fo.lines, "\n")
}

// ReadFile 读取文件内容
func (fo *FileOperator) ReadFile() error {
	file, err := os.Open(fo.FilePath)
	if err != nil {
		return Err_0200020003.Sprintf(fo.FilePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fo.lines = append(fo.lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return Err_0200020007.Sprintf(fo.FilePath, err)
	}

	return nil
}

// GetLine 获取指定行的内容
// 如果行号超出范围，则返回空字符串。
func (fo *FileOperator) GetLine(pos int) string {
	if pos < 0 || pos >= len(fo.lines) {
		return ""
	}
	return fo.lines[pos-1]
}

func (fo *FileOperator) GetLines() []string {
	return fo.lines
}

// Find 查找文件中的标记位置，
// 如果找到则返回行号，行号从1开始；否则返回-1。
//
// Params:
//
//   - marker: 要查找的标记。
//   - isLike: 是否使用模糊匹配。
//   - startPos: 查找的起始位置，从0开始。如果小于0，则从第一行开始查找。
//   - endPos: 查找的结束位置，如果小于0或大于文件行数，则到文件末尾。
//
// Returns:
//
//	0: 标记所在行号。如果找不到，则返回-1。
//	1: 错误信息。
//
// Example:
//
//	// 全文查询
//	pos, err := fo.Find("func main() {", true, -1, -1)
//
//	// 从第10行开始查询
//	pos, err := fo.Find("func main() {", true, 10, -1)
func (fo *FileOperator) Find(marker string, isLike bool, startPos, endPos int) (int, error) {
	// 确保 startPos 和 endPos 在有效范围内
	if startPos < 0 {
		startPos = 0
	}
	if endPos < 0 || endPos > len(fo.lines) {
		endPos = len(fo.lines)
	}

	// 确保 startPos 不大于 endPos
	if startPos > endPos {
		startPos = endPos
	}
	lines := fo.lines[startPos:endPos]
	for i, line := range lines {
		if isLike {
			if strings.Contains(line, marker) {
				return i + 1, nil
			}
		} else {
			if strings.TrimSpace(line) == marker {
				return i + 1, nil
			}
		}
	}
	return -1, nil
}

// FindRange 查找文件中的两个标记之间的位置。
//
// Params:
//
//   - start: 起始标记。
//   - end: 结束标记。
//   - isLike: 是否使用模糊匹配。
//   - startPos: 查找的起始位置，从0开始。如果小于0，则从第一行开始查找。
//   - endPos: 查找的结束位置，如果小于0或大于文件行数，则到文件末尾。
//
// Returns:
//
//	0: 起始行号。如果找不到，则返回-1，同时结束行号也为-1。
//	1: 结束行号。如果找不到，则返回-1。存在有起始行号但找不到结束行号的情况，此时起始行号大于0，结束行号为-1。
//	2: 错误信息。
func (fo *FileOperator) FindRange(startMarker, endMarker string, isLike bool, startPos, endPos int) (int, int, error) {
	startPos, err := fo.Find(startMarker, isLike, startPos, endPos)
	if err != nil {
		return -1, -1, err
	}
	if startPos == -1 {
		return -1, -1, nil
	}

	if startPos < len(fo.lines) {
		for p := startPos + 1; p <= len(fo.lines); p++ {
			if isLike {
				if strings.Contains(fo.lines[p-1], endMarker) {
					return startPos, p, nil
				}
			} else {
				if strings.TrimSpace(fo.lines[p-1]) == endMarker {
					return startPos, p, nil
				}
			}
		}
	}

	return startPos, -1, nil
}

// Insert 在指定位置插入内容
//
// Params:
//
//   - pos: 插入位置，会在这个位置插入，从1开始。比如 pos = 1，会在第一行后面插入，原本的第一行会变成第二行。
//   - content: 要插入的内容。
//
// Returns:
//
//	0: 插入内容后一行的行号。
//	1: 错误信息。
func (fo *FileOperator) Insert(pos int, content string) (nextPos int, err error) {
	if pos < 1 {
		return -1, Err_0200020008.Sprintf(1)
	}

	newLines := make([]string, 0, len(fo.lines))
	if pos > len(fo.lines) {
		newLines = make([]string, 0, pos)
		newLines = append(newLines, fo.lines...)
		for i := len(fo.lines); i < pos-1; i++ {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, content)
		fo.lines = newLines
		return pos + 1, nil
	} else {
		newLines = append(newLines, fo.lines[:pos-1]...)
		newLines = append(newLines, content)
		newLines = append(newLines, fo.lines[pos-1:]...)
		fo.lines = newLines
		return pos + 1, nil
	}
}

// BatchInsert 批量插入内容
//
// Params:
//
//   - index: 插入位置，会在这个位置插入，从1开始。比如 pos = 1，会在第一行后面插入，原本的第一行会变成第二行。
//   - contents: 要插入的内容列表。
//
// Returns:
//
//	0: 插入内容后一行的行号。
//	1: 错误信息。
func (fo *FileOperator) BatchInsert(pos int, contents []string) (nextPos int, err error) {
	nextPos = pos
	for _, content := range contents {
		nextPos, err = fo.Insert(nextPos, content)
		if err != nil {
			return -1, err
		}
	}
	return nextPos, nil
}

// Push 在文件末尾追加内容
func (fo *FileOperator) Append(content string) {
	fo.lines = append(fo.lines, content)
}

// Replace 替换文件中的指定内容
//
// Params:
//
//   - oldContent: 要替换的原内容。
//   - newContent: 新的内容。
//
// Returns:
//
//   - 是否找到并替换了内容。
//   - 错误信息。
//
// Example:
//
//	oldResponse := `export class WasmApiResponse {
//	  free(): void;
//	  readonly code: number;
//	  readonly data: any;
//	  readonly message: string;
//	}`
//	newResponse := `export class WasmApiResponse<T> {
//	  free(): void;
//	  readonly code: number;
//	  readonly data: T;
//	  readonly message: string;
//	}`
//	replaced, err := fo.Replace(oldResponse, newResponse)
func (fo *FileOperator) Replace(oldContent, newContent string) (bool, error) {
	// 将整个文件内容转换为字符串
	content := fo.getContent()

	// 检查原内容是否存在
	if !strings.Contains(content, oldContent) {
		return false, nil
	}

	// 替换内容
	newFileContent := strings.ReplaceAll(content, oldContent, newContent)

	// 按行分割新内容
	fo.lines = strings.Split(strings.TrimRight(newFileContent, "\n"), "\n")

	return true, nil
}

// ReplaceWithRegexp 使用正则表达式和回调函数替换文件内容
//
// Params:
//
//   - re: 正则表达式。
//   - replacer: 替换回调函数，接收匹配的字符串，返回替换后的字符串。
//
// Returns:
//
//   - 是否找到并替换了内容。
//   - 错误信息。
//
// Example:
//
//	// 将所有的数字替换为其两倍
//	re := regexp.MustCompile(`\d+`)
//	result, err := fo.ReplaceWithRegexp(re, func(match string) string {
//		num, _ := strconv.Atoi(match)
//		return strconv.Itoa(num * 2)
//	})
func (fo *FileOperator) ReplaceWithRegexp(re *regexp.Regexp, replacer func(string) string) (bool, error) {
	// 获取当前内容
	content := fo.getContent()

	// 找到所有匹配项
	matches := re.FindAllString(content, -1)
	if len(matches) == 0 {
		return false, nil
	}

	// 使用正则表达式和回调函数替换内容
	newContent := re.ReplaceAllStringFunc(content, replacer)

	// 如果内容没有变化，返回 false
	if newContent == content {
		return false, nil
	}

	// 更新行内容
	fo.lines = strings.Split(strings.TrimRight(newContent, "\n"), "\n")

	return true, nil
}

// Save 保存修改后的内容到文件
func (fo *FileOperator) Save() error {
	// 将 lines 合并成一个字符串
	content := fo.getContent()

	// 如果是 Go 文件，先进行格式化
	if filepath.Ext(fo.FilePath) == ".go" {
		src, err := imports.Process(fo.FilePath, []byte(content), nil)
		if err != nil {
			return Err_0200020002.Sprintf(fo.FilePath, err)
		}
		if err != nil {
			return err
		}
		content = string(src)
	}

	// 写入文件
	err := os.WriteFile(fo.FilePath, []byte(content), 0644)
	if err != nil {
		return Err_0200020001.Sprintf(fo.FilePath, err)
	}

	return nil
}

// getContent 将 lines 合并成一个字符串
func (fo *FileOperator) getContent() string {
	var buffer bytes.Buffer
	for _, line := range fo.lines {
		buffer.WriteString(line)
		buffer.WriteString("\n")
	}
	return buffer.String()
}
