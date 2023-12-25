package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, uid int64) (User, error) {
	var res User
	//这里的 First 方法被用来执行查询，寻找满足 Where 条件的第一条记录。查询的结果会被填充到 res 变量中。
	//&res 表示传递 res 的地址，这样 GORM 可以直接在这个变量中填充数据。
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}

func (dao *UserDAO) UpdateById(ctx context.Context, entity User) error {
	// 这种写法依赖于 GORM 的零值和主键更新特性
	// Update 非零值 WHERE id = ?
	//return dao.db.WithContext(ctx).Updates(&entity).Error
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).Updates(
		map[string]any{
			"utime":    time.Now().UnixMilli(),
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

type User struct {
	AboutMe string `gorm:"type=varchar(4096)"`
	// YYYY-MM-DD
	Birthday int64
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	// 时区，UTC 0 的毫秒数
	// 创建时间
	Ctime    int64
	Nickname string `gorm:"type=varchar(128)"`
	// 更新时间
	Utime int64

	// json 存储
	//Addr string
}

//type Address struct {
//	Uid
//}
