package service

type accountService struct {
}

var (
	AccountService = accountService{}
)

func (*accountService) GetAccountInfo(username string) (map[string]any, error) {
	return map[string]any{
		"username": "admin",
		"password": "123456",
	}, nil
}
