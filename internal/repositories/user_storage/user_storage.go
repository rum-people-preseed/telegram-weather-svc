package user_storage

type UserStorage interface {
	addNewUser(userId int64, country string, city string)
}
