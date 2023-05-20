# VsCode Setting

> 对于我本人来说只需要借助官方的服务，使用Github账号进行账号同步即可，这里主要是介绍如何导入我的配置。

该文件夹中的两个文件分别是：

- extensions.txt 保存了我安装的所有插件列表
- .aiken.code-profile 使用`code --list-extensions` 导出的插件列表

插件列表可以通过如下命令进行安装:

Powershell:

```powershell
cat extensions.txt |% { code --install-extension $_}
```

Bash:

```bash
cat extensions.txt | xargs code --list-extensions {}
```

FI.
