package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type User struct {
	Id       int    `json:"id" form:"id"`
	IdNumber string `json:"id_number" form:"id_number"`
	Mail     string `json:"mail" form:"mail"`
	Name     string `json:"name" form:"name"`
	Birth    string `json:"birthday" form:"bitrhday"`
	Status   string `json:"status" form:"status"`
}

type Detail struct {
	EventNum  int    `json:"event_num" form:"event_num"`
	Title     string `json:"title" form:"title"`
	Performer string `json:"performer" form:"performer"`
	Price     string `json:"ticket_price" form:"ticket_price"`
	TimeAt    string `json:"time_at" form:"time_at"`
	BookFrom  string `json:"book_from" form:"book_from"`
	EndbookAt string `json:"endbook_at" form:"endbook_at"`
	LimitSeat string `json:"limit_seat" form:"limit_seat"`
}

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:demoroot@tcp(127.0.0.1:3306)/ticket?charset=utf8mb4")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	defer DB.Close()

	router := gin.Default()
	router.GET("/user/:id", GetOne)
	router.POST("user", AddOne)
	router.PATCH("user/:id", UpdateUser)

	router.GET("/detail/:id", GetOneDetail)
	router.Run(":8000")
}

//獲得一條紀錄
func GetOne(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	u := User{
		Id: id,
	}
	rs, _ := u.GetRow()
	c.JSON(http.StatusOK, gin.H{
		"result": rs,
	})
}

//新增一位訂票者資訊
func AddOne(c *gin.Context) {
	id_number := c.Request.FormValue("id_number")
	mail := c.Request.FormValue("mail")
	name := c.Request.FormValue("name")
	birthday := c.Request.FormValue("birthday")
	status := c.Request.FormValue("status")

	user := User{
		IdNumber: id_number,
		Mail:     mail,
		Name:     name,
		Birth:    birthday,
		Status:   status,
	}

	id := user.Create()
	msg := fmt.Sprintf("insert successful %d", id)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

//更改訂票人狀態 {0:註銷,1:正常}
func UpdateUser(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	u := User{
		Id: id,
	}
	status := c.Request.FormValue("status")
	user := User{
		Id:     u.Id,
		Status: status,
	}
	row := user.Update()
	msg := fmt.Sprintf("updated successful %d", row)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}
func (d *Detail) Get() (detail Detail, err error) {
	detail = Detail{}
	err = DB.QueryRow("SELECT * from ticket_detail WHERE event_num = ?", detail.EventNum).Scan(
		&detail.EventNum, &detail.Title, &detail.Performer, &detail.Price, &detail.TimeAt, &detail.BookFrom, &detail.EventNum, &detail.LimitSeat)
	return
}

//獲得詳細表演資訊
func GetOneDetail(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	d := User{
		Id: id,
	}
	rs, _ := d.Get()
	c.JSON(http.StatusOK, gin.H{
		"result": rs,
	})
}

func (u *User) GetRow() (user User, err error) {
	user = User{}
	err = DB.QueryRow("SELECT id_number, mail, name, birthday, status FROM user WHERE id=?", u.Id).Scan(
		&user.IdNumber, &user.Mail, &user.Name, &user.Birth, &user.Status)
	return
}

func (u *User) Create() int64 {
	rs, err := DB.Exec("INSERT INTO user (id_number, mail, name, birthday, status) VALUES (?, ?, ?, ?, ?);",
		u.IdNumber, u.Mail, u.Name, u.Birth, u.Status)
	if err != nil {
		log.Fatal(err)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func (u *User) Update() int64 {
	rs, err := DB.Exec("update user set status = ? where id = ?;", u.Status, u.Id)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}
