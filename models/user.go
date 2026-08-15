package models

import (
	"errors"
	"fmt"
	"time"
	"sublink/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              int
	Username        string `gorm:"uniqueIndex"`
	Password        string
	Role            string
	Nickname        string
	SubscriptionID  int
	SubscriptionToken string `gorm:"uniqueIndex"`
	AllowedRegions  string
	InviteCodeID    int
	Disabled        bool
	TokenVersion    int // token版本号: 改密/登出/禁用后自增使旧JWT失效
	CreatedAt       time.Time
}

func (user *User) Create() error { // 创建用户
	if user.Username == "" || user.Password == "" {
		return errors.New("用户名或密码不能为空")
	}
	if user.Role == "" {
		user.Role = "user"
	}
	if user.Nickname == "" {
		user.Nickname = user.Username
	}
	if user.SubscriptionToken == "" {
		user.SubscriptionToken = fmt.Sprintf("sub_%d_%s", time.Now().UnixNano(), utils.RandString(24))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return DB.Create(user).Error
}
func (user *User) Set(UpdateUser *User) error { // 设置用户
	updates := map[string]any{}
	if UpdateUser.Username != "" {
		updates["username"] = UpdateUser.Username
	}
	if UpdateUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UpdateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hashedPassword)
	}
	if UpdateUser.Nickname != "" {
		updates["nickname"] = UpdateUser.Nickname
	}
	if UpdateUser.SubscriptionID != 0 {
		updates["subscription_id"] = UpdateUser.SubscriptionID
	}
	if UpdateUser.SubscriptionToken != "" {
		updates["subscription_token"] = UpdateUser.SubscriptionToken
	}
	if UpdateUser.AllowedRegions != "" {
		updates["allowed_regions"] = UpdateUser.AllowedRegions
	}
	if UpdateUser.InviteCodeID != 0 {
		updates["invite_code_id"] = UpdateUser.InviteCodeID
	}
	if len(updates) == 0 {
		return nil
	}
	return DB.Model(&User{}).Where("username = ?", user.Username).Updates(updates).Error
}
func (user *User) Verify() error { // 验证用户
	var dbUser User
	if err := DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		return err
	}
	if dbUser.Disabled {
		return errors.New("用户已被禁用")
	}
	// 仅支持bcrypt哈希密码比对(禁止明文密码回退)
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return err
	}
	*user = dbUser
	return nil
}

func (user *User) Find() error { // 查找用户
	return DB.Where("username = ? ", user.Username).First(user).Error
}

func (user *User) FindBySubscriptionToken() error { // 通过订阅token查找用户
	return DB.Where("subscription_token = ?", user.SubscriptionToken).First(user).Error
}

func (user *User) All() ([]User, error) { // 获取所有用户
	var users []User
	err := DB.Find(&users).Error
	return users, err
}

func (user *User) FindByID() error { // 通过ID查找用户
	return DB.Where("id = ?", user.ID).First(user).Error
}

func (user *User) Del() error { // 删除用户
	return DB.Delete(user).Error
}

// IncrementTokenVersion 自增用户token版本号(显式事务确保提交)
// 用于改密/登出/禁用后吊销已签发的JWT
func IncrementTokenVersion(username string) error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}
	// COALESCE 兼容存量行为NULL的旧数据(NULL+1=NULL导致自增失效)
	if _, err := tx.Exec("UPDATE users SET token_version = COALESCE(token_version, 0) + 1 WHERE username = ?", username); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
