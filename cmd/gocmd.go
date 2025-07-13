package cmd

import (
	"fmt"
)

// UpgradeGoModules 升级所有go模块依赖到最新版本
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
//	success, err := UpgradeGoModules("/path/to/repo")
//	if err != nil {
//	    log.Printf("Upgrade failed: %v", err)
//	    return
//	}
//	fmt.Println("Upgrade success:", success)
func UpgradeGoModules(repoPath string) (bool, error) {
	output, err := New("go", "get", "-u", "./...").SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go get -u 执行失败: %s, output: %s", err.Error(), string(output))
	}
	return true, nil
}

// TidyGoModules 整理go模块依赖并下载
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
//	success, err := TidyGoModules("/path/to/repo")
//	if err != nil {
//	    log.Printf("Tidy failed: %v", err)
//	    return
//	}
//	fmt.Println("Tidy success:", success)
func TidyGoModules(repoPath string) (bool, error) {
	output, err := New("go", "mod", "tidy").SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go mod tidy 执行失败: %s, output: %s", err.Error(), string(output))
	}
	return true, nil
}

// InitGoModule 初始化Go模块
//
// 执行 go mod init 命令，在指定路径下初始化一个新的Go模块。
// 该操作会创建go.mod文件并设置模块名称。
//
// Params:
//
//   - repoPath: 要初始化模块的路径。
//   - moduleName: 模块名称，例如 git.zodileap.com/entity/database_v1。
//
// Returns:
//
//	bool: 操作是否成功，true表示初始化成功，false表示初始化失败。
//	error: 如果初始化过程中发生错误，返回相应的错误信息。
//
// Example:
//
//	success, err := InitGoModule("/path/to/repo", "git.zodileap.com/entity/database_v1")
//	if err != nil {
//	    log.Printf("Init failed: %v", err)
//	    return
//	}
//	fmt.Println("Init success:", success)
func InitGoModule(repoPath string, moduleName string) (bool, error) {
	output, err := New("go", "mod", "init", moduleName).SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go mod init 执行失败: %s, output: %s", err.Error(), string(output))
	}
	return true, nil
}

// GetSpecificPackage 获取指定的Go包
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
//	success, err := GetSpecificPackage("/path/to/repo", "github.com/zodileap/taurus_go/entity/cmd@latest")
//	if err != nil {
//	    log.Printf("Get package failed: %v", err)
//	    return
//	}
//	fmt.Println("Get package success:", success)
func GetSpecificPackage(repoPath string, packageName string) (bool, error) {
	output, err := New("go", "get", packageName).SetDir(repoPath).Run()
	if err != nil {
		return false, fmt.Errorf("go get %s 执行失败: %s, output: %s", packageName, err.Error(), string(output))
	}
	return true, nil
}
