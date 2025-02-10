package service

const (
	READ = 1 << iota
	UPDATE
	CREATE
	DELETE
)

var (
	refreshCountKey = "refresh_token_used"
)

func OperateRefreshCount() {

}
