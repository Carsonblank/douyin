# 简单的抖音项目

## 配置
使用前在config文件夹下的config.go中首先设置DSN，用于连接服务器

## 默认用户和视频
首次运行时，执行Initdatabase文件夹下的addDefaultMessage.go文件，
会自动创建默认用户和默认视频信息

## 文件结构
### main.go
+ 连接数据库自动迁移表
+ 设置接口
### config文件夹
+ DSN和视频流最大视频数
### controller文件夹
+ 接口函数
+ messageType.go声明了用于接口的结构体
+ baseFunction.go定义了在数据库结构体和接口结构体之间转换的函数
### database文件夹
+ 数据库模型结构体
+ 查询操作函数
+ 删除操作函数
+ 创建操作函数
### Initdatabase文件夹
+ 第一次使用时运行，会在数据库中添加默认用户和视频信息
### public文件夹
+ 上传的视频和封面会保存到此处
+ defaultAvatar.png是默认的头像
+ defaultVideoCover.png是默认的视频封面

## 2022/5/30更新
+ 使用jwt生成token，token有效期为登陆起七天，每次登陆会更新token
+ 因为每次登陆都会更新token，所以理论上不允许多地登陆
+ 在投稿接口，赞操作，评论操作，关注操作中增加了token有效性验证
+ 使用SHA256算法对密码进行散列，注册和登陆使用
+ SQL查询时使用了Where("id=?",id)预编译，不存在SQL注入问题

## 2022/6/1更新
+ 添加了全局变量DB，定义在common文件夹的global.go
+ 增加必要的统计数据
