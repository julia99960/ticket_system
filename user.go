package main

//User 訂票人資訊
type User struct {
	ID       int    `json:"id" form:"id"`
	IDNumber string `json:"id_number" form:"id_number"`
	Mail     string `json:"mail" form:"mail"`
	Name     string `json:"name" form:"name"`
	Birth    string `json:"birthday" form:"bitrhday"`
	Status   int    `json:"status" form:"status"`
}

// GetOneUser 取得一筆使用者資料
func GetOneUser(id int) (user User, err error) {
	u := User{
		ID: id,
	}
	user, err = u.GetRow()
	return
}

// AddOneUser 新增一筆使用者資料
func (u *User) AddOneUser() int64 {
	id := u.Create()
	return id
}

// UpdateUserStatus 更新使用者狀態{0:註銷,1:正常}
func UpdateUserStatus(id, status int) (row int64) {
	if status == 0 || status == 1 {
		user := User{
			ID:     id,
			Status: status,
		}
		row = user.Update()
		return
	}

	return 0
}
