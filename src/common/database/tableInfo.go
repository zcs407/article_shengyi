package database

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/* 用户 table_name = user */
type User struct {
	Id           int    //用户编号
	Name         string `gorm:"size:32;unique"` //用户名
	PasswordHash string `gorm:"size:128" `      //用户密码加密的
	//Mobile        string `gorm:"size:11;unique" ` //手机号
	RollerId int       //0为管理员,1为编辑,2为普通用户
	Articles []Article //用户创建的文章
}

/* 文章表 table_name = article */
type Article struct {
	Id          int
	UserId      int
	Title       string `gorm:"not null;unique"` //文章链接
	ContentUrl  string `gorm:"size:255"`        //文章链接
	ContentPath string `gorm:"size:255"`        //文章本地路径
	CreateTime  string
	UpdateTime  string
	Images      []Image //链接地址
	Status      int     `gorm:"size:11"` //审批状态,0编辑中,1审批中,2审批通过,3审批失败
	Tags        []*Tag  `gorm:"many2many:article_tag;"`
	SpecialId   int
}
type Image struct {
	Id        int `gorm:"AUTO_INCREMENT"`
	ArticleId int
	ImageUrl  string `gorm:"type:varchar(100)"`
	ImagePath string `gorm:"type:varchar(100)"`
}
type Tag struct {
	Id      int
	TagName string
	Article []*Article
}

//专题
type Special struct {
	Id          int
	SpecialName string
	Pid         int //上级专题的id
	Articles    []*Article
	Specials    []Special `gorm:"-"`
}
