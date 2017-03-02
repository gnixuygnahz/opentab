![Opentab Logo](.github/opentab-logo.png)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Gitter](https://badges.gitter.im/Join Chat.svg)](https://gitter.im/opentab-server/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link)
开源的轻应用后端（Open Tiny App Backend），轻量，高效，易部署。

## 入门
Opentab 是一个开源的轻量级的通用应用后端，具备应用存储、用户管理、权限控制、文件管理等基础功能，为移动应用开发提供强有力的后端支持。
## 文档
[Opentab 部署指南](https://leanote.com/blog/post/58aedb41ab6441490b001595)

[Opentab REST API使用指南](http://leanote.com/blog/post/58a86b5aab644109c3000377)

[Opentab Master开发文档](http://leanote.com/blog/post/58a697676761373561000000)

## 依赖
 - Golang 推荐1.7以上
 - PostgreSQL 9.5
 - gin （内置组件）
 - gopsutil （内置组件）
 - lib/pq （内置组件）
 - qiniu sdk （内置组件）
 - Unknwon/goconfig （内置组件）
 - pmylund/go-cache （内置组件）

## 部署
### 直接部署
安装PostgreSQL数据库。配置项目根目录`config.ini`相关参数，设置环境变量`GIN_MODE=release`。在windows环境下，请将编译源代码后的二进制文件与config.ini放置在同一目录下，直接运行二进制文件即可；在linux环境下，需要将二进制文件与config.ini放置在/usr/bin/下，直接运行即可。
### Docker部署
根目录中包含Docker部署相关的文件

 - `Dockerfile`
 - `Dockerfile.package` - 安全镜像打包文件
 - `docker-compose.yml` - 应用编排配置文件
 - `daocloud.yml` - daocloud 安全镜像配置文件

### 初次使用
项目运行后，会在日志中打印出 AppId，AppKey，MasterKey 。具体使用请参照 [Opentab REST API使用指南](http://leanote.com/blog/post/58a86b5aab644109c3000377)

## 管理
个人桌面级管理端正在开发，敬请期待。

## 支持
目前项目处于测试阶段，为了使作者更好地完善项目，欢迎大家提出issue，作者会及时做出回应。如有需求或者意见，也可以直接发送邮件至936269579@qq.com。

## 版权说明
```
Copyright (c) 2017-present, Zhang Yuxing.
All rights reserved.

This source code is licensed under the Apache License version 2.0
found in the LICENSE file in the root directory of this source tree.
An additional grant of patent rights can be found in the PATENTS
file in the same directory.
```