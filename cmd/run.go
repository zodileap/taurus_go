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

type Cmd struct {
	cmd *exec.Cmd
}

var (
	baseCmd   []string = []string{}
	outerrFmt string   = " 命令: %s,\n 运行错误: %s, \n 标准输出: %s, \n 标准错误: %s"
)

// SetBaseCmd 设置基础命令，这个会在每个命令的前面添加。
// 例如sudo
func SetBaseCmd(cmd ...string) {
	baseCmd = cmd
}

// NewCmd 执行 cmd 命令,这个用来调用那种不常用的命令
// [cmd_1] [cmd_2] ... [cmd_n].
func NewCmd(cmd ...string) *Cmd {
	args := appendArgs(cmd...)
	// 构建并运行 cmd 命令
	c := &Cmd{
		cmd: exec.Command(args[0], args[1:]...),
	}
	return c
}

// NewBash 创建一个新的命令，用来执行一些简单的命令。
//
// 参数：
//  - sh 用来执行的shell,例如bash,zsh,csh,ksh,/bin/sh
//  - [cmd] 用来执行的命令。

func NewSh(sh string, cmd string) *Cmd {
	return NewCmd(sh, "-c", cmd)
}

// NewUser 创建命令，用于切换的别的用户。
// sudo -i -u user.
// 不受baseCmd影响
func NewUser(user string) *Cmd {
	cmd := &Cmd{
		cmd: exec.Command("sudo", "-i", "-u", user),
	}
	return cmd
}

// NewCurl 创建一个新的curl命令，用来下载文件。
// curl [url] | sudo bash.
// 不受baseCmd影响
func NewCurl(url string) *Cmd {
	// 构建并运行 curl 命令
	cmd := &Cmd{
		cmd: exec.Command("sh", "-c", fmt.Sprintf("curl %s | sudo bash", url)),
	}
	return cmd
}

// NewInstall 创建新的安装文件命令。
// apt-get install [path].
func NewInstall(path string) *Cmd {
	args := appendArgs("apt-get", "-y", "install", path)
	// 构建并运行 apt-get 命令
	cmd := &Cmd{
		cmd: exec.Command(args[0], args[1:]...),
	}
	return cmd
}

// NewGrep 创建 grep 命令用来查找文件中是否存在某一行的。
// grep -qF -- [line] [file]
func NewGrep(file string, line string) *Cmd {
	args := appendArgs("grep", "-qF", "--", line, file)
	cmd := &Cmd{
		cmd: exec.Command(args[0], args[1:]...),
	}
	return cmd
}

// runInsertText 创建 sed 命令。
//
// flag:
//   - [$ a ] 在文件的最后一行添加内容。
//   - [s/] 替换匹配行。
//
// sudo sed -i [flag][line] [file]
func NewSed(file string, flag string, line string) *Cmd {
	args := appendArgs("sed", "-i", fmt.Sprintf(`%s%s`, flag, line), file)
	cmd := &Cmd{
		cmd: exec.Command(args[0], args[1:]...),
	}
	return cmd
}

// split 用来分割命令。
// 例如：sudo -i -u user会被分割成sudo,-i,-u,user
func Split(commandStr string) []string {
	return strings.Fields(commandStr)
}

func (c *Cmd) SetDir(dir string) *Cmd {
	c.cmd.Dir = dir
	return c
}

func (c *Cmd) SetEnv(env []string) *Cmd {
	c.cmd.Env = env
	return c
}

// Run 运行命令，返回运行信息，以及出现的error，不会对error进行处理
func (c *Cmd) Run() ([]byte, error) {
	var stdout, stderr bytes.Buffer
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
// 但是对于错误为空，但是标准错误不为空的情况，会输出错误信息，但是不会停止运行。
func (c *Cmd) Must() {
	var stdout, stderr bytes.Buffer
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
// 但是对于错误为空，但是标准错误不为空的情况，会输出错误信息，但是不会停止运行。
func (c *Cmd) MustReturn() []byte {
	var stdout, stderr bytes.Buffer
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
// 例如原本的命令是sudo -i -u user，添加psql -c,
// 现在会变成 sudo -i -u user psql -c
func (c *Cmd) AddCmd(cmd ...string) *Cmd {
	c.cmd.Args = append(c.cmd.Args, cmd...)
	return c
}

func (c *Cmd) AddPipe(cmd ...string) *Cmd {
	cmd = append([]string{" | "}, cmd...)
	c.cmd.Args = append(c.cmd.Args, cmd...)
	return c
}

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

// outputErr 用来输出错误
func outputErr(format string, args ...interface{}) error {
	err := errors.Errorf(format, args...)
	fmt.Println(err)
	return err
}
