package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Cmd 执行命令的结构体。
type Cmd struct {
	cmd *exec.Cmd
}

var (
	baseCmd   []string = []string{}
	outerrFmt string   = " 命令: %s,\n 运行错误: %s, \n 标准输出: %s, \n 标准错误: %s"
	outerr    bool     = false
)

// SetBaseCmd 设置基础命令，这个会在每个命令的前面添加。
//
// Params:
//
//   - cmd: 基础命令。
//
// Example:
//
//	cmd.SetBaseCmd("sudo")
//	c := cmd.New("mkdir", "dir").String()
//	fmt.Println(c)
//
// ExamplePath: taurus_go_demo/cmd/cmd_test.go - TestSetBaseCmd
func SetBaseCmd(cmds ...string) {
	baseCmd = cmds
}

func SetOuterr(b bool) {
	outerr = b
}

// New 用于创建一个新的命令。
// 注意，在终端中执行的命令，比如 mrkdir dir，需要拆分成两个字符串”mkdir","dir"传入，而不是直接传入“mrkdir dir”。
//
// Params:
//
//   - cmd: 用来执行的命令。
//
// Returns:
//
// Example:
//
//	cmd.New("mkdir", "dir").Must()
func New(cmds ...string) *Cmd {
	args := appendArgs(cmds...)
	// 构建并运行 cmd 命令
	c := &Cmd{
		cmd: exec.Command(args[0], args[1:]...),
	}
	return c
}

// NewSh 自定义命令行解释器并执行命令。例如，sh,bash,zsh等。
//
// Params:
//
//   - sh: 自定义命令行解释器。
//   - cmd: 用来执行的命令。
//
// Example:
//
//	cmd.NewSh("bash","mkdir dir").Must()
func NewSh(sh string, cmd string) *Cmd {
	c := &Cmd{
		cmd: exec.Command(sh, "-c", cmd),
	}
	return c
}

// NewUser 创建命令，用于切换的别的用户，通过sudo执行，所以需要有root权限。
// sudo -i -u user,不受基础命令影响。
// 可以通过 cat /etc/passwd 命令查看用户列表。
//
// Params:
//
//   - user: 需要切换到的用户.
//
// Example:
//
//	cmd.NewUser("root").AddCmd("mkdir", "dir").Must()
//
// ExamplePath: taurus_go_demo/cmd/cmd_test.go - TestNewUser
func NewUser(user string) *Cmd {
	c := &Cmd{
		cmd: exec.Command("sudo", "-i", "-u", user),
	}
	return c
}

// split 用来分割命令。
// 例如："sudo -i -u user"会被分割成"sudo","-i","-u","user"。
//
// Params:
//
//   - commandStr: 需要分割的命令。
//
// Returns:
//
//	0: 分割后的命令。
//
// Example:
//
// c := cmd.New(cmd.Split("mkdir dir")...)
//
// ExamplePath:
//
//   - taurus_go_demo/cmd/cmd_test.go - TestSplit
func Split(commandStr string) []string {
	return strings.Fields(commandStr)
}

// SetDir 设置命令运行的目录。类似于cd 命令。
//
// Params:
//
//   - dir: 需要切换的目录。
//
// Example:
//
//	cmd.New("mkdir", "dir").SetDir("/home").Must()
func (c *Cmd) SetDir(dir string) *Cmd {
	c.cmd.Dir = dir
	return c
}

// SetEnv 设置命令运行的环境变量。
//
// Params:
//
//   - env: 需要设置的环境变量。
//
// Example:
//
//	cmd.New("go", "build").SetEnv(append(os.Environ(), "GOOS=linux", "GOARCH=amd64")).Must()
func (c *Cmd) SetEnv(env []string) *Cmd {
	c.cmd.Env = env
	return c
}

// Run 运行命令，返回运行信息，以及出现的error，不会对error进行处理。
//
// Returns:
//
//	0: 运行信息。
//	1: 错误信息。
//
// Example:
//
//	r, err := cmd.New("cat", "/etc/passwd").Run()
func (c *Cmd) Run() ([]byte, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	c.cmd.Stdout = &stdout
	c.cmd.Stderr = &stderr
	err := c.cmd.Run()
	if err != nil {
		err := errors.Errorf(outerrFmt, c.cmd.String(), err.Error(), stdout.String(), stderr.String())
		return stdout.Bytes(), err
	}
	return stdout.Bytes(), nil
}

// Must 运行指令，并判断是否出现错误，
// 如果出现错误，抛出panic，并停止运行，不会返回执行信息。
// 对于错误为空，但是标准错误不为空的情况，如果设置了SetOuterr(true)会输出错误信息，但是不会停止运行。
//
// Example:
//
//	cmd.New("mkdir", "dir").Must()
func (c *Cmd) Must() {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	c.cmd.Stdout = &stdout
	c.cmd.Stderr = &stderr
	err := c.cmd.Run()
	if err != nil {
		outputErr(outerrFmt, c.cmd.String(), err.Error(), stdout.String(), stderr.String())
		os.Exit(1)
	} else if stderr.Len() > 0 {
		outputErr(outerrFmt, c.cmd.String(), "nil", stdout.String(), stderr.String())
	}
}

// MustReturn 运行指令，并判断是否出现错误，
// 如果出现错误，抛出panic，并停止运行，会返回执行信息。
// 对于错误为空，但是标准错误不为空的情况，SetOuterr(true)会输出错误信息，但是不会停止运行。
//
// Example:
//
//	r := cmd.New("cat", "/etc/passwd").MustReturn()
func (c *Cmd) MustReturn() []byte {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	c.cmd.Stdout = &stdout
	c.cmd.Stderr = &stderr
	err := c.cmd.Run()
	if err != nil {
		outputErr(outerrFmt, c.cmd.String(), err.Error(), stdout.String(), stderr.String())
		os.Exit(1)
	} else if stderr.Len() > 0 {
		outputErr(outerrFmt, c.cmd.String(), "nil", stdout.String(), stderr.String())
	}
	return stdout.Bytes()
}

// AddCmd 在原本的命令的末尾添加新的命令。
// 例如原本的命令是“sudo -i -u user”，添加“psql -c”,
// 现在会变成“sudo -i -u user psql -c”。
//
// Params:
//
//   - cmd: 需要添加的命令。
//
// Returns:
//
//	0: 添加后的命令。
//
// Example:
//
//	c := cmd.NewUser("root").AddCmd("mkdir", "dir")
func (c *Cmd) AddCmd(cmd ...string) *Cmd {
	c.cmd.Args = append(c.cmd.Args, cmd...)
	return c
}

// String 返回命令的字符串形式。
//
// Returns:
//
//	0: 命令的字符串形式。
//
// Example:
//
//	c := cmd.New("mkdir", "dir").String()
func (c *Cmd) String() string {
	return c.cmd.String()
}

func (c *Cmd) SetSrdin(stdin io.Reader) *Cmd {
	c.cmd.Stdin = stdin
	return c
}

// appendArgs 用来在命令前面添加基础命令。
// 例如：sudo -i -u user
func appendArgs(args ...string) []string {
	return append(baseCmd, args...)
}

// outputErr 用来输出错误。
func outputErr(format string, args ...interface{}) error {
	err := errors.Errorf(format, args...)
	if outerr {
		fmt.Println(err)
	}
	return err
}
