# EntryTask
go-rpc for user management system
main function: httpServer.go
****
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
- 用户登陆界面
- httpServer
- rpcService

**httpServer**主要是负责处理http请求，先对数据进行预处理，然后传递给rpc服务。httpServer还负责页面的跳转，错误处理以及本地缓存的保存等。

**rpcService**主要利用rpc机制，将http的请求数据发送至tcpServer，然后tcpServer进行业务和数据操作，并将结果返回至httpServer.其主要组成为tcpC
lient和tcpServer。

