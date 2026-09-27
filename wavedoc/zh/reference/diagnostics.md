---
translation_set_id: diagnostics
path: reference/diagnostics
locale: zh
group: reference
group_order: 5
order: 2
title: 故障排除：从安装到执行
summary: 隔离失败的步骤并通过可重现的信息缩小原因范围。
---

## 一、区分故障阶段

|观察到的现象|先检查一下|下一步行动|
| --- | --- | --- |
|wavec 未找到命令|PATH 和可执行文件位置|使用绝对路径运行并设置PATH|
|找不到运行所需的文件|安装文件夹中是否缺少任何文件？|再次解压并安装整个包|
|std import 失败| `wavec print std-path` |对应：安装std或指定`--std-root`|
|源位置和类型错误输出| `wavec check main.wave` |修复第一个错误并再次检查|
|其他 OS·CPU 目标构建失败|指定target和目标环境|[交叉构建设置](/docs/zh/whale/build-link-targets) 确认|
|构建成功后执行失败|退出代码、输入、工作目录|执行环境及API错误检查|

## 一个小诊断示例

这是整个程序，这是故意错误的：

```wave
fun main() {
    var count: i32 = 1;
    println("{}", missing);
}
```

`wavec check main.wave` 必须指向未声明的名称missing。将变量名称更改为count，然后再次检查并运行。重点关注文件、位置和原因，而不是整个诊断文本。后续错误可能是最初错误的结果。

## 当可执行文件失败时

在Linux/macOS shell中，执行后立即检查退出代码为`echo $?`，在PowerShell中，它是`$LASTEXITCODE`。输入错误和显式`return 1`不是同一原因。班次计数或实际转换的无效运行时值可能会导致trap。查看[操作规则](/docs/zh/language/expressions-and-operators)。

相对文件路径受可执行文件工作目录而不是源文件位置的影响。不要将文件读取失败视为字符串长度为0，首先检查返回错误。通过地址查找、服务器等待、权限和超时来检查网络连接故障。

## 报告问题所需的信息

1. `wavec --version` 输出和执行的确切命令。
2. target. 与主机分开指定 OS·架构
3. 与选定的 std 路径一起使用的编译器源。
4. 重现问题所需的最少源、输入和文件。
5. 预期结果、实际结果、诊断和退出代码。

密码、令牌和个人文件内容将被删除。如果当您减少最小示例时问题就消失了，那么您删除的最后一个元素就是线索。当工具收集诊断信息时，`--error-format=json` 可用。

[安装](/docs/zh/getting-started/install) · [编译器命令](/docs/zh/getting-started/compiler) · [目标和链接](/docs/zh/whale/build-link-targets)
