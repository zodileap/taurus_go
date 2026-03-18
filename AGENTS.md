# taurus_go 协作规范

- 默认使用中文。
- 先读当前包实现、测试和 README，再修改代码或文档。

## 仓库定位

- 模块路径：`github.com/zodileap/taurus_go`。
- Go 版本要求：`1.24.1`。
- 这是基础能力仓库，不承载具体业务应用逻辑。
- 常见能力集中在 `entity`、`notify/telegram`、`cache/redis`、`tlog`、`asset`、`template`、`err`、`byteutil`、`maputil`、`sliceutil`、`stringutil`、`structutil`、`geo`。

## Git 规范

- 版本号遵循三段式 `Major.Minor.Patch`。
- 提交信息默认使用中文，并遵循约定式提交格式：

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

- 当一次提交涉及多个变更类型时，允许使用多标签格式：

```text
<type1>[optional scope]: <description1>
<type2>[optional scope]: <description2>

[optional body]

[optional footer(s)]
```

- 重大不兼容变更使用 `major`，并在正文或脚注中写明 `BREAKING CHANGE: <详细描述>`。
- 常用提交类型：
  - `major`：重大更新或架构变更，版本变更 `Major (X.0.0)`。
  - `feat`：新增功能或模块，版本变更 `Minor (0.X.0)`。
  - `update`、`fix`、`perf`、`refactor`、`build`：版本变更 `Patch (0.0.X)`。
  - `docs`、`style`、`test`、`ci`、`chore`：不触发版本号变更。
  - `revert`：回退之前的提交，按实际回退内容处理。
- 提交信息要求：
  - 简洁准确反映变更内容。
  - 不添加额外签名信息。
  - 不在末尾署名。

## 文档与注释

- 文档注释中的标签均为可选。
- 注释开头先写描述，不要先写函数、方法、属性名称。
- 函数或方法注释要写清楚职责和实际行为，避免空泛描述。
- 各标签内容之间保留空行。
- 常用标签：`Params:`、`default:`、`Returns:`、`Example:`、`ExamplePath:`、`ErrCodes:`、`Verbs:`、`Extends:`。
- 结构体字段注释优先写在字段右侧。
- README 和注释里的示例代码必须基于当前真实导出 API，不能写伪代码或失效路径。

## 开发约束

- 修改任何包之前，先阅读该包当前实现和测试，确认真实对外 API。
- 新增能力优先在现有包体系中扩展；只有职责边界明确时才新增新包。
- 优先复用仓库内已有能力，例如 `err`、`asset`、`template`、`byteutil`、`maputil`、`sliceutil`、`stringutil`、`structutil`、`geo`；确认不足后再评估第三方依赖。
- 新增公共函数时优先保持清晰一致的返回风格，通常使用 `(T, error)`，但不要为了机械统一破坏已有 API 兼容性。
- 错误处理优先沿用当前风格，例如 `err.ErrCode`、各子包内的 `Err_xxx` 定义和已有包装方式；已有错误码优先复用。
- 补充错误上下文时优先包含关键调用位置、对象标识或输入边界信息，但不要泄露敏感数据。
- 与测试相关的夹具优先放 `testdata/`、`*_test.go` 或 `t.TempDir()`，不要在仓库根部散落临时文件。
- 开发后需要为对应函数、方法、模块或修复点补充测试；没有合适测试落点时，先补最接近公开行为的单元测试。
- 修改导出行为时，优先补充或更新对应测试。
- 文档类改动如果涉及示例代码，应尽量选择仓库内已有、可核验的用法模式。
- 修改 README 时，要同步核对 import path、命令、结构体名、函数名和版本信息是否真实存在。
- `README.md` 面向使用者，重点写项目定位、安装方式、核心能力和真实示例；内容应保持中文为主，必要时保留英文术语。
- 涉及 `entity`、`notify/telegram`、`cache/redis`、`tlog` 等核心子系统时，优先做局部、可验证的改动，不做无边界重构。
