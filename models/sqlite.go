package models

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var isInitialized bool

func InitSqlite() {
	// 检查目录是否创建
	_, err := os.Stat("./db")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./db", os.ModePerm)
		}
	}
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("./db/sublink.db"), &gorm.Config{})
	if err != nil {
		log.Println("连接数据库失败")
	}
	DB = db
	// SQLite 单写者模型: 限制连接池为单连接，
	// 避免多连接下读事务残留导致写入被回滚/不可见(如 token_version 自增不生效、"database is locked")
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(1)
	}
	// 检查是否已经初始化
	if isInitialized {
		log.Println("数据库已经初始化，无需重复初始化")
		return
	}
	err = db.AutoMigrate(&User{}, &Subcription{}, &Node{}, &SubLogs{}, &AccessKey{}, &SubScheduler{}, &SubcriptionNode{}, &Plugin{}, &InviteCode{})
	if err != nil {
		log.Println("数据表迁移失败")
	}
	// 回填存量数据: SQLite新增列对旧行为NULL，将token_version NULL补为0
	db.Exec("UPDATE users SET token_version = 0 WHERE token_version IS NULL")
	// 初始化用户数据
	err = db.First(&User{}).Error
	if err == gorm.ErrRecordNotFound {
		admin := &User{
			Username: "admin",
			Password: "123456",
			Role:     "admin",
			Nickname: "管理员",
		}
		err = admin.Create()
		if err != nil {
			log.Println("初始化添加用户数据失败")
		}
	}
	// 设置初始化标志为 true
	isInitialized = true
	log.Println("数据库初始化成功") // 只有在没有任何错误时才会打印这个日志
}
