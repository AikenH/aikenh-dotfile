# VsCode Setting

update 20241019, use the profile function of vscode to manager extensions and setting for different scenario.使用vscode的profile管理不同场景下的不同设置和插件。

## Profile Update. 20241019

This repo distinguishes profiles based on different langs, mainly including minimal, cpp, python, frontend for now.该repo基于不同语言对profiles进行区分，目前主要包含 minimal,cpp,python,frontend 四个配置。

- The minimal profile does not specify a language.
    其中的minimal 配置并没有针对某个语言，作为各个配置的基础;
- The Frontend & CPP profile not finish yet, I will continue to follow up on the project and provide updates as necessary. 
    前端的配置和cpp的配置还没完成，后续随着个人的使用会持续更新，并在有必要的时候进行迭代。

### How to use

Open profile setting page in vscode, 在VSCode的下述位置打开配置文件界面

![open profile](https://picture-bed-001-1310572365.cos.ap-guangzhou.myqcloud.com/mac/20241019202148.png)

Click import profile then select the profile file. 然后选中导入配置将对应的文件导入即可

![import profile](https://picture-bed-001-1310572365.cos.ap-guangzhou.myqcloud.com/mac/20241019202256.png)

Besides, you can import profile by Gist. 除了下载对应 profile 文件，也可以直接使用 Gist 链接进行导入，链接如下

- [Minimal](https://vscode.dev/profile/github/a6be778359645da3468565d9c01d463d)
- [CPP](https://vscode.dev/profile/github/bed0413d4bc4d59653e33a3c18adbfba)
- [Python](https://vscode.dev/profile/github/e6b14679818041d87ebf662e49a8d0e8)
- [Frontend](https://vscode.dev/profile/github/83e5bbb1785d082085f6e8bad9568249)

## Old Version 

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
