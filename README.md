# api 文档生成和同步工具

## 使用方法

假如 golang 项目已经按照 swag 要求的注释方式做好了注释

### 创建 apidoc 配置文件

在 golang 项目的根目录创建配置文件 .apidoc.yaml，填入 yapi 和 metersphere 的配置信息（注意不要将此配置文件提交到代码库，在 .gitignore 文件中添加忽略），示例格式如下：

```
yapi:
host: <your yapi host>
# 项目
project: <your project name>
# 项目token
token: <your token>

ms:
host: <your ms host>
# 工作空间
workspace: <your workspace>
# 项目
project: <your project name>
# 应用名称
application: <your application name>
accessKey: <your accessKey>
secretKey: <your secretKey>
```

### 生成api文档

在项目根目录或者项目的 app 目录执行如下命令（使用此命令生成 api 文档的前提是安装了 swag 工具，这个命令底层调用了 swag 命令来实现）

```
apidoc build
```

执行此命令后，会在当前目录下生成 docs 文件夹，此文件下面会有 json 和 yaml 格式的接口文档。

### 上传到 yapi 或metersphere

在项目的根目录执行如下命令同时上传到 yapi 和 metersphere

```
apidoc upload
```

上传到 yapi

```
apidoc upload yapi
```

上传到 meterspher

```
apidoc upload ms
```