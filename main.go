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

//DB mysql
var DB *sql.DB

//User 訂票人資訊
type User struct {
	ID       int    `json:"id" form:"id"`
	IDNumber string `json:"id_number" form:"id_number"`
	Mail     string `json:"mail" form:"mail"`
	Name     string `json:"name" form:"name"`
	Birth    string `json:"birthday" form:"bitrhday"`
	Status   string `json:"status" form:"status"`
}

// Detail 表演場次細節
type Detail struct {
	ID        int    `json:"id" form:"id"`
	Title     string `json:"title" form:"title"`
	Performer string `json:"performer" form:"performer"`
	Price     string `json:"ticket_price" form:"ticket_price"`
	TimeAt    string `json:"time_at" form:"time_at"`
	BookFrom  string `json:"book_from" form:"book_from"`
	EndbookAt string `json:"endbook_at" form:"endbook_at"`
	LimitSeat string `json:"limit_seat" form:"limit_seat"`
}

// Ticket 下單紀錄
type Ticket struct {
	ID       int    `json:"id" form:"id"`
	EventNum int    `json:"event_num" form:"event_num"`
	UserID   int    `json:"userid" form:"userid"`
	BookAt   string `json:"book_at" from:"book_at"`
	PayAt    string `json:"pay_at" form:"pay_at"`
	Status   string `json:"status" form:"status"`
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

	//訂票者資訊
	router.GET("/user/:id", GetOne)
	router.POST("user", AddOne)
	router.PATCH("user/:id", UpdateUser)

	//表演詳細資料
	router.GET("/detail/:id", GetOneDetail)
	router.POST("/detail", AddOneDetail)

	//已訂門票資訊
	router.GET("/tickets/:id", GetTickets)
	router.Run(":8000")
}

// GetOne 獲得一條紀錄
func GetOne(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	u := User{
		ID: id,
	}
	rs, _ := u.GetRow()
	c.JSON(http.StatusOK, gin.H{
		"result": rs,
	})
}

// AddOne 新增一位訂票者資訊
func AddOne(c *gin.Context) {
	IDNumber := c.Request.FormValue("id_number")
	mail := c.Request.FormValue("mail")
	name := c.Request.FormValue("name")
	birthday := c.Request.FormValue("birthday")
	status := c.Request.FormValue("status")

	user := User{
		IDNumber: IDNumber,
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

// UpdateUser 更改訂票人狀態 {0:註銷,1:正常}
func UpdateUser(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	u := User{
		ID: id,
	}
	status := c.Request.FormValue("status")
	user := User{
		ID:     u.ID,
		Status: status,
	}
	row := user.Update()
	msg := fmt.Sprintf("updated successful %d", row)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

// GetOneDetail 獲得詳細表演場次資訊
func GetOneDetail(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	d := Detail{
		ID: id,
	}
	rs, _ := d.GetOne()
	c.JSON(http.StatusOK, gin.H{
		"result": rs,
	})
}

// AddOneDetail 新增一筆表演場次
func AddOneDetail(c *gin.Context) {
	title := c.Request.FormValue("title")
	performer := c.Request.FormValue("performer")
	price := c.Request.FormValue("ticket_price")
	timeat := c.Request.FormValue("time_at")
	bookfrom := c.Request.FormValue("book_from")
	endbookat := c.Request.FormValue("endbook_at")
	limitseat := c.Request.FormValue("limit_seat")

	d := Detail{
		Title:     title,
		Performer: performer,
		Price:     price,
		TimeAt:    timeat,
		BookFrom:  bookfrom,
		EndbookAt: endbookat,
		LimitSeat: limitseat,
	}

	id := d.Create()
	msg := fmt.Sprintf("insert successful %d", id)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

// GetTickets 取得一筆詳細資料
func GetTickets(c *gin.Context) {
	userids := c.Param("id")
	userid, _ := strconv.Atoi(userids)
	t := Ticket{
		UserID: userid,
	}
	rs, _ := t.GetRow()
	c.JSON(http.StatusOK, gin.H{
		"list": rs,
	})
}

// GetRow 取得一名使用者資料
func (u *User) GetRow() (user User, err error) {
	user = User{}
	err = DB.QueryRow("SELECT id, id_number, mail, name, birthday, status FROM user WHERE id=?", u.ID).Scan(
		&user.ID, &user.IDNumber, &user.Mail, &user.Name, &user.Birth, &user.Status)
	return
}

// Create 新增一筆使用者資料
func (u *User) Create() int64 {
	rs, err := DB.Exec("INSERT INTO user (id_number, mail, name, birthday, status) VALUES (?, ?, ?, ?, ?);",
		u.IDNumber, u.Mail, u.Name, u.Birth, u.Status)
	if err != nil {
		log.Fatal(err)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

// Update 更新使用者狀態
func (u *User) Update() int64 {
	rs, err := DB.Exec("update user set status = ? where id = ?;", u.Status, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

// GetOne 取得一筆表演資訊
func (d *Detail) GetOne() (detail Detail, err error) {
	detail = Detail{}
	err = DB.QueryRow("SELECT id, title, performer, ticket_price, time_at, book_from, endbook_at, limit_seat from ticket_detail WHERE id = ?", d.ID).Scan(
		&detail.ID, &detail.Title, &detail.Performer, &detail.Price, &detail.TimeAt, &detail.BookFrom, &detail.EndbookAt, &detail.LimitSeat)
	return
}

// Create 新增一筆表演資訊
func (d *Detail) Create() int64 {
	rs, err := DB.Exec("INSERT INTO ticket_detail (title, performer, ticket_price, time_at, book_from, endbook_at , limit_seat) VALUES (?, ?, ?, ?, ?, ?, ?);",
		&d.Title, &d.Performer, &d.Price, &d.TimeAt, &d.BookFrom, &d.EndbookAt, &d.LimitSeat)

	if err != nil {
		log.Fatal(err)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

// GetRow 取得訂票資訊
func (t *Ticket) GetRow() (tickets []Ticket, err error) {
	rows, err := DB.Query("select id, event_num, userid, book_at, status from ticket where userid = ?", t.UserID)

	for rows.Next() {
		ticket := Ticket{}
		err := rows.Scan(&ticket.ID, &ticket.EventNum, &ticket.UserID, &ticket.BookAt, &ticket.Status)
		if err != nil {
			tickets = append(tickets, ticket)
		} else {
			log.Fatal(err)
		}
	}
	rows.Close()
	return
}
