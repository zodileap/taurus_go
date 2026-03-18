package console

import (
	"fmt"
	"strings"
)

// Examples 用于格式化多行示例文本。
func Examples(examples ...string) string {
	for i := range examples {
		examples[i] = "  " + examples[i]
	}
	return strings.Join(examples, "\n")
}

// Module 输出一个模块标题。
func Module(format string, args ...any) {
	fmt.Printf("*** %s ***\n", fmt.Sprintf(format, args...))
}

// Step 输出一个步骤标题。
func Step(format string, args ...any) {
	fmt.Printf("-> %s\n", fmt.Sprintf(format, args...))
}

// SubStep 输出一个子步骤标题。
func SubStep(format string, args ...any) {
	fmt.Printf("   -> %s\n", fmt.Sprintf(format, args...))
}

// Skip 输出一个跳过说明。
func Skip(format string, args ...any) {
	fmt.Printf("   ! %s [skip]\n", fmt.Sprintf(format, args...))
}

// Done 输出一个完成说明。
func Done(format string, args ...any) {
	fmt.Printf("[done] %s\n", fmt.Sprintf(format, args...))
}
