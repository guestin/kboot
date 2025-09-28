package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/guestin/mob"
	"gorm.io/gorm"
)

type CreatedAt struct {
	CreatedAt   time.Time `gorm:"column:created_at" json:"-"`
	CreatedAtTs int64     `gorm:"-" json:"createdAt"`
}

func (this *CreatedAt) AfterFind(*gorm.DB) (err error) {
	this.CreatedAtTs = this.CreatedAt.Unix()
	return
}

type UpdatedAt struct {
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"-"`
	UpdatedAtTs int64     `gorm:"-" json:"updatedAt"`
}

func (this *UpdatedAt) AfterFind(*gorm.DB) (err error) {
	this.UpdatedAtTs = this.UpdatedAt.Unix()
	return
}

type DeletedAt struct {
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
	DeletedAtTs *int64         `gorm:"-" json:"deletedAt,omitempty"`
}

func (this *DeletedAt) AfterFind(*gorm.DB) (err error) {
	if this.DeletedAt.Valid {
		this.DeletedAtTs = new(int64)
		*this.DeletedAtTs = this.DeletedAt.Time.Unix()
	}
	return
}

// UuidPrimaryKey UUID主键
type UuidPrimaryKey struct {
	ID string `gorm:"column:id;primaryKey;type:varchar(32)" json:"id"`
}

func (this *UuidPrimaryKey) BeforeCreate(*gorm.DB) (err error) {
	if len(this.ID) == 0 {
		this.ID = strings.ToUpper(mob.GenRandomUUID())
	}
	return
}

// Int64PrimaryKey 自增主键
type Int64PrimaryKey struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"-"`
	//当ID数值过大时，前端js处理会丢失精度导致ID不准确，故转换为string返回
	IdStr string `gorm:"-" json:"id"`
}

//goland:noinspection ALL
func (this *Int64PrimaryKey) AfterFind(session *gorm.DB) (err error) {
	this.IdStr = fmt.Sprintf("%d", this.ID)
	return
}

// UuidPriWithCreateAtBase UUID主键 + CreatedAt
type UuidPriWithCreateAtBase struct {
	UuidPrimaryKey
	CreatedAt
}

func (this *UuidPriWithCreateAtBase) AfterFind(session *gorm.DB) (err error) {
	_ = this.CreatedAt.AfterFind(session)
	return
}

// UuidPriWithCreateDelAtBase UUID主键 + CreatedAt + DeletedAt
type UuidPriWithCreateDelAtBase struct {
	UuidPrimaryKey
	CreatedAt
	DeletedAt
}

func (this *UuidPriWithCreateDelAtBase) AfterFind(session *gorm.DB) (err error) {
	_ = this.CreatedAt.AfterFind(session)
	_ = this.DeletedAt.AfterFind(session)
	return
}

// Int64PriWithCreateAtBase 自增主键 + CreatedAt
type Int64PriWithCreateAtBase struct {
	Int64PrimaryKey
	CreatedAt
}

func (this *Int64PriWithCreateAtBase) AfterFind(session *gorm.DB) (err error) {
	_ = this.CreatedAt.AfterFind(session)
	return
}

// Int64PriWithCreateDelAtBase 自增主键 + CreatedAt + DeletedAt
type Int64PriWithCreateDelAtBase struct {
	UuidPrimaryKey
	CreatedAt
	DeletedAt
}

func (this *Int64PriWithCreateDelAtBase) AfterFind(session *gorm.DB) (err error) {
	_ = this.CreatedAt.AfterFind(session)
	_ = this.DeletedAt.AfterFind(session)
	return
}
