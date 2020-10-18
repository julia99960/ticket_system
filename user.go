package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//User 訂票人資訊
type User struct {
	ID       int    `json:"id" form:"id"`
	IDNumber string `json:"id_number" form:"id_number"`
	Mail     string `json:"mail" form:"mail"`
	Name     string `json:"name" form:"name"`
	Birth    string `json:"birthday" form:"bitrhday"`
	Status   int    `json:"status" form:"status"`
}

// GetOne 取得一筆使用者資料
func GetOne(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	rs, err := GetOneUser(id)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"result": rs,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": nil,
		})
	}
}

// GetOneUser 取得一筆使用者資料
func GetOneUser(id int) (user User, err error) {
	u := User{
		ID: id,
	}
	user, err = u.GetRow()
	return
}

// AddOne 新增一位訂票者資訊
func AddOne(c *gin.Context) {
	IDNumber := c.Request.FormValue("id_number")
	mail := c.Request.FormValue("mail")
	name := c.Request.FormValue("name")
	birthday := c.Request.FormValue("birthday")
	sta := c.Request.FormValue("status")
	status, _ := strconv.Atoi(sta)

	u := User{
		IDNumber: IDNumber,
		Mail:     mail,
		Name:     name,
		Birth:    birthday,
		Status:   status,
	}

	id := u.AddOneUser()
	msg := fmt.Sprintf("insert successful %d", id)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

// AddOneUser 新增一筆使用者資料
func (u *User) AddOneUser() int64 {
	id := u.Create()
	return id
}

// UpdateUser 更改訂票人狀態
func UpdateUser(c *gin.Context) {
	ids := c.Param("id")
	status1 := c.Request.FormValue("status")
	id, _ := strconv.Atoi(ids)
	status, _ := strconv.Atoi(status1)

	row := UpdateUserStatus(id, status)
	msg := fmt.Sprintf("updated successful %d", row)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
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
