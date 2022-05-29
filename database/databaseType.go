//该文件定义了写入数据库对象的类型

package database

import (
	"gorm.io/gorm"
)

//用户信息
type User struct {
	Id        int64  `gorm:"primaryKey;autoIncrement"`        //用户唯一标志符号
	Name      string `gorm:"type:varchar(32);not null;index"` //用户名
	Password  string `gorm:"type:varchar(32);not null"`       //用户密码
	Token     string `gorm:"not null, index"`                 //用户鉴权Token
	Avatar    string //用户头像链接Url
	Signature string //用户个性签名
}

//视频信息
type Video struct {
	Id        int64          `gorm:"primaryKey;autoIncrement"` //视频唯一标志符
	UserId    int64          `gorm:"not null"`                 //视频发布者ID
	PlayUrl   string         `gorm:"not null"`                 //视频URL
	CoverUrl  string         //视频封面URL
	Title     string         //视频标题
	CreatedAt int64          `gorm:"autoCreateTime:milli"` //视频创建时间。选择milli是因为发现本机的mysql操作是以millisecond的，即使使用nano，后边的值也只会补上0,例如使用nano：1653548764819000000
	UpdatedAt int64          `gorm:"autoCreateTime:milli"` //更新时间。
	DeletedAt gorm.DeletedAt `gorm:"index"`                //删除时间
}

//点赞信息
type Favorite struct {
	Id       int64 `gorm:"primaryKey;autoIncrement"` //点赞数据的唯一标识符
	UserId   int64 `gorm:"not null"`                 //发起点赞操作的用户ID
	ToUserId int64 `gorm:"not null"`                 //受到点赞的用户ID
	VideoId  int64 `gorm:"not null"`                 //受到点赞的视频ID
}

//评论信息
type Comment struct {
	Id        int64  `gorm:"primaryKey;autoIncrement"` //评论唯一标识符
	UserId    int64  `gorm:"not null"`                 //发起评论的用户ID
	VideoId   int64  `gorm:"not null"`                 //评论所在的视频ID
	ToUserId  int64  `gorm:"not null"`                 //视频发布者的ID
	Content   string `gorm:"not null"`                 //评论内容
	CreatedAt int64  `gorm:"autoCreateTime:milli"`     //评论创建时间
}

//关注信息
type Relation struct {
	Id       int64 `gorm:"primaryKey;autoIncrement"` //关注唯一标识符
	UserId   int64 `gorm:"not null"`                 //发起关注者ID
	ToUserId int64 `gorm:"not null"`                 //被关注者的ID
}
