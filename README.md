# demo

只实现了注册功能，往user表中插入新数据

## 代码介绍

主要的业务逻辑处理是src文件夹中的controller包、service包、repository包

分为视图层服务层和数据层，各层尽量相互独立，视图层只调用服务层的接口，并把信息返回给客户端，服务层只调用数据层的接口，并把数据整理好返回给视图层，数据层负责直接和数据库进行交互，返回相应结果给服务层

* `configFile`文件夹

  * config.ini

    配置文件

    比如

    ```
    [server]
    HTTP_PORT = 填端口
    HTTP_HOST = 0.0.0.0 这个可以不改
    MODE = debug
    ```

    

* `public`文件夹，可以用来存视频，这个demo中没有用到

  * bear.mp4

* `sql`文件夹

  * example.sql  User表的定义创建在其中

* `src`文件夹

  * common包  定义一些公共的方法

    * jwt.go 定义ReleaseToken和ParseToken两个函数，只看注册功能的话可以基本忽略

  * config包

    * 提供了一些函数，用于读取configFile文件夹中config.ini文件中的配置信息，比如router.go中

      ```go
      httpServer.Run(config.GetServerConfig().HTTP_HOST + ":" + config.GetServerConfig().HTTP_PORT)
      ```

      把config.ini中填入的信息读出来

      或者database包中mysql.go里也有用到配置信息 dbConfig.USER这种类似的都在config.ini中填入过

      ```go
      dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s",
      		dbConfig.USER,
      		dbConfig.PASSWORD,
      		dbConfig.DB_HOST,
      		dbConfig.DB_PORT,
      		dbConfig.DB_NAME,
      		dbConfig.CHARSET,
      		dbConfig.ParseTime,
      	)
      
      	// open connection
      	MySqlDb, MySqlError = gorm.Open("mysql", dbDSN)
      ```

  * controller包

    视图层

    * common.go common包中定义Response用来复用，还有一些公用的东西
    * user.go 提供register注册接口，直接暴露给客户端，视图层调用服务层函数来完成自己的业务

  * database包

    * mysql.go mysql的配置连接都在这里完成，并且提供一个全局变量`var MySqlDb *gorm.DB`供repository包使用

  * repository包

    数据层

    * user.go 定义了user表的ORM中的对象User结构体，以及和user表相关的数据操作

    只有数据层涉及到和数据库进行交互

    Dao是数据访问对象的意思 data access object

    这里直接看代码比较好，每个函数名就是自己的用途

  * service包

    * create_user.go

    服务层，定义了一个数据流，既有输入也有返回

    create_user.go的话数据流就是CreateUserFlow结构体，主要关注这个结构体实现的方法（go中的方法对应C++类实现的函数，把结构体看成类）

    ```go
    type CreateUserFlow struct {
    	userId   int64   //用来返回
    	username string  //从上层视图层拿到的
    	password string  //从上层视图层拿到的
    }
    ```

    这部分就是在做这一件事

    ```go
    func CreateUser(username string, password string) (int64, error) { //username,password都是上层来的
    	return NewCreateUserFlow(username, password).Do() //返回的int64是新增加的用户ID
    }
    ```

    服务层作为“桥梁”，向下拿到数据层的数据，向上返回给视图层

* go.mod 不用动

* `main.go`  初始化路由引擎，将其传入`router.go`

* `router.go` 注册路由，监听端口



## 环境配置

只需要在运行这部分代码之前执行sql文件就好了

另外可能稍微需要修改一下配置文件config.ini

```
[server]
HTTP_PORT = 填监听端口
HTTP_HOST = 0.0.0.0不用改
MODE = debug 不用改

[mysql]
TYPE = mysql 不用改
USER = root 也许要改，但一般都是root
PASSWORD = 改成自己的密码

DB_HOST = 数据库所在主机的ip地址
DB_PORT = 3306 也许要改，但mysql一般都是3306
DB_NAME = user表所在database的名称
CHARSET = utf8 不用改
ParseTime = true 不用改
MaxIdleConns = 20 不用改
MaxOpenConns = 100 不用改
```

除了以上都不用改

## 编译运行命令

直接编译运行即可

```shell
go build && ./demo
```

或者

```
go run main.go router.go
```

