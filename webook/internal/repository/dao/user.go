package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u DaoUser) error
	FindByEmail(ctx context.Context, email string) (DaoUser, error)
	UpdateById(ctx context.Context, entity DaoUser) error
	FindById(ctx context.Context, uid int64) (DaoUser, error)
	FindByPhone(ctx context.Context, phone string) (DaoUser, error)
}

type GormUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}

func (dao *GormUserDAO) Insert(ctx context.Context, u DaoUser) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (DaoUser, error) {
	var u DaoUser
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *GormUserDAO) UpdateById(ctx context.Context, entity DaoUser) error {

	// 这种写法依赖于 GORM 的零值和主键更新特性
	// Update 非零值 WHERE id = ?
	//return dao.db.WithContext(ctx).Updates(&entity).Error
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.Id).
		Updates(map[string]any{
			"utime":    time.Now().UnixMilli(),
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

func (dao *GormUserDAO) FindById(ctx context.Context, uid int64) (DaoUser, error) {
	var res DaoUser
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&res).Error
	return res, err
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (DaoUser, error) {
	var res DaoUser
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&res).Error
	return res, err
}

type DaoUser struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 代表这是一个可以为 NULL 的列
	//Email    *string
	Email    sql.NullString `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(128)"`
	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	// 代表这是一个可以为 NULL 的列
	Phone sql.NullString `gorm:"unique"`

	// 时区，UTC 0 的毫秒数
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64

	// json 存储
	//Addr string
}

//type Address struct {
//	Uid
//}
