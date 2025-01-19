package service

import (
	"ray_box/internal/dao"
)

type accountService struct {
}

var (
	AccountService = accountService{}
)

func (*accountService) GetAccountInfo(username string) (map[string]any, error) {
	return dao.UserDao.GetUserInfo(username, "password")
}
