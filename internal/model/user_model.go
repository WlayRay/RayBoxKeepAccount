package model

type User struct {
	ID        int64  `gorm:"column:id"`         // 主键
	Username  string `gorm:"column:username"`   // 用户名
	Password  string `gorm:"column:password"`   // 密码
	Phone     string `gorm:"column:phone"`      // 手机号
	Email     string `gorm:"column:email"`      // 邮箱
	Avatar    string `gorm:"column:avatar"`     // 头像地址
	CreatedAt int64  `gorm:"column:created_at"` //创建时间，时间戳，单位秒
	UpdatedAt int64  `gorm:"column:updated_at"` //更新时间，时间戳，单位秒
	State     int    `gorm:"column:state"`      // 状态
}

func (User) TableName() string {
	return "work.user"
}

var UserColums = struct {
	ID        string
	Username  string
	Password  string
	Phone     string
	Email     string
	Avatar    string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Username:  "username",
	Password:  "password",
	Phone:     "phone",
	Email:     "email",
	Avatar:    "avatar",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
