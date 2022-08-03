# EntryTask
go-rpc for user management system
main function: httpServer.go
## 功能要求
实现一个用户管理系统，用户可以登录、拉取和编辑他们的profiles。

用户可以通过在Web页面输入username和password登录，backend系统负责校验用户身份。成功登录后，页面需要展示用户的相关信息；否则页面展示相关错误。

成功登录后，用户可以编辑以下内容：

1.上传profile picture

2.修改nickname（需要支持Unicode字符集，utf-8编码）

用户信息包括：

1.username（不可更改）

2.nickname

3.profile picture

需要提前将初始用户数据插入数据库用于测试。确保测试数据库中包含10,000,000条用户账号信息。

## 设计要求

- 分别实现HTTP server和TCP server，主要的功能逻辑放在TCP server实现

- Backend鉴权逻辑需要在TCP server实现

- 用户账号信息必须存储在MySQL数据库。通过MySQL Go client连接数据库

- 使用基于Auth/Session Token的鉴权机制

- TCP server需要提供RPC API，RPC机制希望自己设计实现

- Web server不允许直连MySQL。所有HTTP请求只处理API和用户输入，具体的功能逻辑和数据库操作，需要通过RPC请求TCP server完成

- 尽可能使用Go标准库

- 安全性

- 鲁棒性

- 性能

## 开发环境
OS: mocOS Monterey 12.2

MySQL: 8.0.29

Redis: 6.2.7

go: 1.18.4

## 设计简介

本项目主要由以下三部组成：
- 用户界面
- httpServer
- rpcService

**用户界面**包括用户登陆与注册界面和用户信息展示界面。

**httpServer**主要是负责处理http请求，先对数据进行预处理，然后传递给rpc服务。httpServer还负责页面的跳转，错误处理以及本地缓存的保存等。

**rpcService**主要利用rpc机制，将http的请求数据发送至tcpServer，然后tcpServer进行业务和数据操作，并将结果返回至httpServer.其主要组成为tcpC
lient和tcpServer。

## 实现流程
**整体流程图**
![整体流程图](/resource/static/flow-chart.png)

## API接口
### 1.登陆接口
|URL                        |方法 |
|:-------------------------:|:--:|
|http://localhost:8080/login|POST|

**输入参数**

|参数名  |描述  |可选|
|:-----:|:---:|:-:|
|username|用户名|否|
|password|用户名|否|

### 2.注册接口
|URL                        |方法 |
|:-------------------------:|:--:|
|http://localhost:8080/signUp|POST|

**输入参数**

|参数名  |描述  |可选|
|:-----:|:---:|:-:|
|username|用户名|否|
|password|用户名|否|

### 3.获取用户信息接口
|URL                             |方法|
|:------------------------------:|:-:|
|http://localhost:8080/getProfile|GET|

**输入参数**

|参数名  |描述  |可选|
|:-----:|:---:|:-:|
|username|用户名|否|
|token   |身份令牌|否|


### 4.修改昵称
|URL                             |方法|
|:------------------------------:|:-:|
|http://localhost:8080/updateNickName|POST|

**输入参数**

|参数名  |描述  |可选|
|:-----:|:---:|:-:|
|username|用户名  |否|
|nickname|昵称   |否|

### 5.更新头像
|URL                             |方法|
|:------------------------------:|:-:|
|http://localhost:8080/updateNickName|POST|

**输入参数**

|参数名   |描述  |可选|
|:------:|:---:|:-:|
|username|用户名  |否|
|image   |图片    |否|

### 6.注销
|URL                          |方法|
|:---------------------------:|:-:|
|http://localhost:8080/signOut|POST|

## 数据库存储
### mysql
维护一张users表
|Field|Type|Null|Key|Default|Extra|
|:---:|:--:|:--:|:-:|:-----:|:---:|
|uid  |int |NO  |PRI|NULL   |auto_increment|
|UserName|varchar(20)|NO| |NULL| |
|NickName|varchar(20)|YES| |NULL| |
|PicName|varchar(100)|YES| |NULL| |
|Password|varchar(100)|NO| |NULL| |

### redis
redis作为登陆的缓存。

1.用户账户信息
|key            |value.  |
|:-------------:|:------:|
|username+"_pwd"|password|

2.用户登陆令牌

|key            |value|
|:-------------:|:---:|
|username+"_tok"|token|

3.用户个人信息

|key            |value                                      |
|:-------------:|:-----------------------------------------:|
|username+"_inf"|{["NickName":nickname],["PicName":picname]}|

## 代码结构

```bash
usermaLX4
├── config                  //配置文件
├── dao                     //mysql
├── models                  //mysql数据结构
├── protocol                  //接口输入输出类型
├── redis                   //redis
├── resource                //页面样式，用户图片保存路径及静态图片资源
├── rpcService              //rpc实现
├── template                //页面模板
├── utils                   //辅助函数
├── httpserver              //http server
└── README.MD               //项目文档

