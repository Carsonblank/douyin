# 配置需求
+ 本代码基于`go version go1.18.1 linux/amd64`开发
+ mysql数据库，版本 `Ver 14.14`

# 快速开始
+ 修改配置文件
+ `go run main.go`

# 代码介绍
## config文件夹
+ config包
    配置文件
    config.init
    + [server]gin配置
        + HTTP_PORT 端口号
        + HTTP_HOST  ip地址
        + MODE 开启模式，debug/release
    + [mysql]数据库配置，目前只支持mysql数据库
        + USER 用户名
        + PASSWORD 密码
        + DB_HOST 数据库IP地址
        + DB_PORT 数据库端口
        + DB_NAME 数据库名称
        + CHARSET 编码方式
        + ParseTime 时间设置
        + Loc 时区
    + [otherInfo] 其他参数，目前只有一个控制视频流最大视频数量
        + MAX_FEED_VIDEO_NUMS feed视频流最大视频数量

## public文件夹
    保存视频文件

## src文件夹 
代码实现
+ database包
    + init.go:数据库的初始化和连接
    + databaseType.go:定义了和数据库交互的GORM结构体
    + init.go:提供了全局变量`var MySqlDb *gorm.DB`供repository包使用
+ repository包
数据层，使用Gorm直接与数据库进行交互
    + user.go 定义了User表的相关数据操作
    + video.go 定义了Video表的相关数据操作
    + comment.go 定义了Comment表的相关数据操作
    + favorite.go 定义了Favorite表的相关数据操作
    + relation.go 定义了Relation表的相关数据操作
+ service包
服务层，为视图层提供服务
    + encrption.go
        用户密码加密
    + inputeValid.go
        验证端口接受到的参数的有效性
    + token.go
        实现鉴权token的生成和分发
    + messageType.go
        定义了直接用于发包的json结构体
    + transformation.go
        定义了数据库GORM结构体和JSON结构体之间的变换
+ controller包
表示层，直接为douyin.apk端口提供服务
    + 基础接口
        + feed.go
            /douyin/feed:视频流接口实现
        + userRegister.go
            /douyin/user/register/:用户注册接口实现
        + userLogin.go
            /douyin/user/login/:用户登录接口实现
        + userInfo.go
            /douyin/user/:用户信息接口实现
        + publishAction.go
            /douyin/publish/action/:用户投稿接口实现
        + publishList.go
            /douyin/publish/list/:用户视频发布列表实现
    + 扩展接口-I
        + favoriteAction.go
            /douyin/favorite/action/:点赞操作实现
        + favoriteList.go
            /douyin/favorite/list/:点赞列表实现
        + commentAction.go
            /douyin/comment/action/:评论操作实现
        + commentList.go
            /douyin/comment/list/:评论列表实现
    + 扩展接口-II
        + relationAction.go
            /douyin/relation/action/:关注操作实现
        + relationFollowList.go
            /douyin/relation/follow/list/:关注列表实现
        + relationFollowerList.go
            /douyin/relation/follower/list/:粉丝列表实现
## log文件夹
生成的log会输出到这里

# 飞书项目文档
https://u2jpafxkyk.feishu.cn/docs/doccnrOfLvAQ5xCGgSvrtr8AKph#GrGL7k