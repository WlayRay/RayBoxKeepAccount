package dao

type userDao struct {
}

var (
	UserDao = userDao{}
)

func (*userDao) GetUserInfo(username string, fields ...string) (map[string]any, error) {
	// pgsqlConn := db.GetPostgresConn(db.POSTGRESQL_DB_MAIN)

	// result := make(map[string]any)
	// err := pgsqlConn.Debug().Model(&model.User{}).
	// 	Where("username = ?", username).
	// 	Select("password").
	// 	Order("id").
	// 	First(&result).Error

	result := map[string]any{
		"username": "admin",
		"password": "123456",
	}

	return result, nil
}
