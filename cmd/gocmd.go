package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// GoRun 执行go run命令
//
// 执行指定目标文件的go run命令，支持传入构建标志。
//
// Params:
//
//   - target: 要运行的目标文件路径。
//   - buildFlags: 构建标志列表。
//
// Returns:
//
//   string: 命令执行的标准输出。
//   error: 如果执行过程中发生错误，返回相应的错误信息。
//
// Example:
//
//	output, err := GoRun("main.go", []string{"-ldflags", "-s -w"})
//	if err != nil {
//	    log.Printf("Run failed: %v", err)
//	    return
//	}
//	fmt.Println("Output:", output)
func GoRun(target string, buildFlags []string) (string, error) {
	args := []string{"run"}
	args = append(args, buildFlags...)
	args = append(args, target)
	
	cmd := New(append([]string{"go"}, args...)...)
	output, err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("taurus_go/cmd run error:\n%s", err)
	}
	return string(output), nil
}

// GoList 执行go list命令
//
// 执行指定目标的go list命令，支持传入构建标志。
//
// Params:
//
//   - target: 要列出的目标路径。
//   - buildFlags: 构建标志列表。
//
// Returns:
//
//   error: 如果执行过程中发生错误，返回相应的错误信息。
//
// Example:
//
//	err := GoList("./...", []string{"-json"})
//	if err != nil {
//	    log.Printf("List failed: %v", err)
//	}
func GoList(target string, buildFlags []string) error {
	args := []string{"list"}
	args = append(args, buildFlags...)
	args = append(args, target)
	
	cmd := New(append([]string{"go"}, args...)...)
	_, err := cmd.Run()
	return err
}

// IsGoModuleInitialized 检查指定目录是否已经初始化为Go模块
//
// 通过检查目录中是否存在go.mod文件来判断是否已经初始化为Go模块。
//
// Params:
//
//   - repoPath: 要检查的目录路径。
//
// Returns:
//
//   bool: 如果目录已经初始化为Go模块则返回true，否则返回false。
//
// Example:
//
//	if IsGoModuleInitialized("/path/to/repo") {
//	    fmt.Println("Directory is already a Go module")
//	} else {
//	    fmt.Println("Directory is not a Go module")
//	}
func IsGoModuleInitialized(repoPath string) bool {
	goModPath := filepath.Join(repoPath, "go.mod")
	_, err := os.Stat(goModPath)
	return err == nil
}

// GoUpgradeModules 升级所有go模块依赖到最新版本
//
// 执行 go get -u ./... 命令，将指定路径下的所有Go模块依赖升级到最新版本。
// 该操作会递归查找当前目录及其子目录中的所有go.mod文件，并升级其中的依赖。
//
// Params:
//
//   - repoPath: 要执行升级操作的仓库路径。
//
// Returns:
//
//	bool: 操作是否成功，true表示升级成功，false表示升级失败。
//	error: 如果升级过程中发生错误，返回相应的错误信息。
//
// Example:
//
//	success, err := GoUpgradeModules("/path/to/repo")
//	if err != nil {
//	    log.Printf("Upgrade failed: %v", err)
//	    return
//	}
//	fmt.Println("Upgrade success:", success)
func GoUpgradeModules(repoPath string) (bool, error) {
	output, err := New("go", "get", "-u", "./...").SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go get -u 执行失败: %s, output: %s", err.Error(), string(output))
	}
	return true, nil
}

// GoTidyModules 整理go模块依赖并下载
//
// 执行 go mod tidy 命令，整理go模块的依赖关系，移除不再使用的依赖，
// 并下载缺失的依赖包。这是Go模块管理的标准操作。
//
// Params:
//
//   - repoPath: 要执行整理操作的仓库路径。
//
// Returns:
//
//	bool: 操作是否成功，true表示整理成功，false表示整理失败。
//	error: 如果整理过程中发生错误，返回相应的错误信息。
//
// Example:
//
//	success, err := GoTidyModules("/path/to/repo")
//	if err != nil {
//	    log.Printf("Tidy failed: %v", err)
//	    return
//	}
//	fmt.Println("Tidy success:", success)
func GoTidyModules(repoPath string) (bool, error) {
	output, err := New("go", "mod", "tidy").SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go mod tidy 执行失败: %s, output: %s", err.Error(), string(output))
	}
	return true, nil
}

// GoInitModule 初始化Go模块
//
// 执行 go mod init 命令，在指定路径下初始化一个新的Go模块。
// 如果目录已经初始化为Go模块，则跳过初始化并返回成功。
// 该操作会创建go.mod文件并设置模块名称。
//
// Params:
//
//   - repoPath: 要初始化模块的路径。
//   - moduleName: 模块名称，例如 git.zodileap.com/entity/database_v1。
//
// Returns:
//
//	bool: 操作是否成功，true表示初始化成功或已经初始化，false表示初始化失败。
//	error: 如果初始化过程中发生错误，返回相应的错误信息。
//
// Example:
//
//	success, err := GoInitModule("/path/to/repo", "git.zodileap.com/entity/database_v1")
//	if err != nil {
//	    log.Printf("Init failed: %v", err)
//	    return
//	}
//	fmt.Println("Init success:", success)
func GoInitModule(repoPath string, moduleName string) (bool, error) {
	if IsGoModuleInitialized(repoPath) {
		return true, nil
	}
	
	output, err := New("go", "mod", "init", moduleName).SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go mod init 执行失败: %s, output: %s", err.Error(), string(output))
	}
	return true, nil
}

// GoGetSpecificPackage 获取指定的Go包
//
// 执行 go get 命令，下载并安装指定的Go包到指定版本。
// 常用于添加项目依赖或更新特定包到最新版本。
//
// Params:
//
//   - repoPath: 要执行操作的路径。
//   - packageName: 包名称，例如 github.com/zodileap/taurus_go/entity/cmd@latest。
//
// Returns:
//
//	bool: 操作是否成功，true表示获取成功，false表示获取失败。
//	error: 如果获取过程中发生错误，返回相应的错误信息。
//
// Example:
//
//	success, err := GoGetSpecificPackage("/path/to/repo", "github.com/zodileap/taurus_go/entity/cmd@latest")
//	if err != nil {
//	    log.Printf("Get package failed: %v", err)
//	    return
//	}
//	fmt.Println("Get package success:", success)
func GoGetSpecificPackage(repoPath string, packageName string) (bool, error) {
	output, err := New("go", "get", packageName).SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go get %s 执行失败: %s, output: %s", packageName, err.Error(), string(output))
	}
	return true, nil
}
