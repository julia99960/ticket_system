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
	Status   int    `json:"status" form:"status"`
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
	LimitSeat int    `json:"limit_seat" form:"limit_seat"`
}

// Ticket 下單紀錄
type Ticket struct {
	ID       int    `json:"id" form:"id"`
	EventNum int    `json:"event_num" form:"event_num"`
	UserID   int    `json:"user_id" form:"user_id"`
	BookAt   string `json:"book_at" from:"book_at"`
	Status   int    `json:"status" form:"status"`
	Sum      int    `json:"sum" form:"sum"`
	Detail
}

func init() {
	var err error
	DB, err = sql.Open("mysql", "root:shupa0127@tcp(mysql:3306)/ticket?charset=utf8mb4")
	// DB, err = sql.Open("mysql", "root:demoroot@tcp(127.0.0.1:3306)/ticket?charset=utf8mb4")
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

	//表演詳細資料
	router.GET("/ticket/:id/detail", GetOneDetail)
	router.POST("/detail", AddOneDetail)

	//總計某一場次剩餘票數
	router.GET("/remain_tickets/:event_num", GetRemainTicket)

	//訂票紀錄
	router.GET("/ticket/:user_id/tickets", GetTickets)
	router.POST("/ticket/:event_num", AddTicket)
	router.PATCH("/ticket/:id/status", UpdateTicket)

	//訂票者資訊
	router.GET("/user/:id", GetOne)
	router.POST("/user", AddOne)
	router.PATCH("/user/:id/status", UpdateUser)

	router.Run(":8000")
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

// GetOneDetail 獲得詳細表演場次資訊
func GetOneDetail(c *gin.Context) {
	ids := c.Param("id")
	id, _ := strconv.Atoi(ids)
	detail, err := GetPerformanceDetail(id)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"result": detail,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": nil,
		})
	}

}

// GetPerformanceDetail 獲取表演場次細節
func GetPerformanceDetail(ida int) (detail Detail, err error) {
	d := Detail{
		ID: ida,
	}
	detail, err = d.GetOne()
	return
}

// AddOneDetail 新增一筆表演場次
func AddOneDetail(c *gin.Context) {
	title := c.Request.FormValue("title")
	performer := c.Request.FormValue("performer")
	price := c.Request.FormValue("ticket_price")
	timeat := c.Request.FormValue("time_at")
	bookfrom := c.Request.FormValue("book_from")
	endbookat := c.Request.FormValue("endbook_at")
	limitseats := c.Request.FormValue("limit_seat")
	limitseat, _ := strconv.Atoi(limitseats)

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

// GetTickets 取得使用者訂票紀錄
func GetTickets(c *gin.Context) {
	userids := c.Param("user_id")
	userid, _ := strconv.Atoi(userids)

	tickets, err := GetTicketsList(userid)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"ticket_list": tickets,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"ticket_list": nil,
		})
	}
}

// GetTicketsList 取得使用者購票紀錄
func GetTicketsList(userid int) (tickets []Ticket, err error) {
	t := Ticket{
		UserID: userid,
	}
	tickets, err = t.GetRow()
	return
}

// AddTicket 新增一筆訂票紀錄
func AddTicket(c *gin.Context) {
	eventnums := c.Param("event_num")
	userids := c.Request.FormValue("user_id")
	statuss := c.Request.FormValue("status")
	userid, _ := strconv.Atoi(userids)
	eventnum, _ := strconv.Atoi(eventnums)
	status, _ := strconv.Atoi(statuss)

	t := Ticket{
		UserID:   userid,
		EventNum: eventnum,
		Status:   status,
	}

	id := t.Create()
	msg := fmt.Sprintf("insert successful %d", id)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

// UpdateTicket 更改訂票紀錄 {0:註銷,1:正常}
func UpdateTicket(c *gin.Context) {
	ids := c.Param("id")
	statuss := c.Request.FormValue("status")
	id, _ := strconv.Atoi(ids)
	status, _ := strconv.Atoi(statuss)

	row := UpdateTicketStatus(id, status)
	msg := fmt.Sprintf("updated successful %d", row)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

// GetRemainTicket 總計某一場次剩餘票數
func GetRemainTicket(c *gin.Context) {
	eventnums := c.Param("event_num")
	eventnum, _ := strconv.Atoi(eventnums)
	row := RemainTicket(eventnum)
	c.JSON(http.StatusOK, gin.H{
		"data": row,
	})
}

// RemainTicket 剩餘票數
func RemainTicket(eventnum int) int {
	t := Ticket{
		EventNum: eventnum,
	}

	rows, err := t.GetSumTiket()

	if err != nil {
		return 0
	}

	sum := rows.Sum

	d := Detail{
		ID: eventnum,
	}

	rows1, err := d.GetLimitSeat()

	if err != nil {
		return 0
	}

	limitseat := rows1.LimitSeat
	return limitseat - sum
}

// UpdateTicketStatus 更改訂票紀錄 {0:註銷,1:正常}
func UpdateTicketStatus(id, status int) int64 {
	if id == 0 || status == 1 {
		t := Ticket{
			ID:     id,
			Status: status,
		}
		row := t.Update()
		return row
	}
	return 0
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
	rows, err := DB.Query("select * from ticket where user_id = ?", t.UserID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var ticket Ticket
		err := rows.Scan(&ticket.ID, &ticket.EventNum, &ticket.UserID, &ticket.BookAt, &ticket.Status)
		if err != nil {
			log.Fatal(err)
		}

		rs, err := GetPerformanceDetail(ticket.EventNum)
		if err != nil {
			rs = Detail{}
		}

		ticket.Detail = rs
		tickets = append(tickets, ticket)
	}
	rows.Close()
	return
}

// Create 新增一筆訂票資訊
func (t *Ticket) Create() int64 {
	rs, err := DB.Exec("INSERT INTO ticket (event_num, user_id, status) VALUES (?, ?, ?);",
		&t.EventNum, &t.UserID, &t.Status)

	if err != nil {
		log.Fatal(err)
	}
	id, err := rs.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	return id
}

// Update 更改注單狀態
func (t *Ticket) Update() int64 {
	rs, err := DB.Exec("UPDATE ticket SET status=? WHERE id=?;", t.Status, t.ID)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := rs.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

// GetSumTiket 統計同一場次票數
func (t *Ticket) GetSumTiket() (ticket Ticket, err error) {
	ticket = Ticket{}
	err = DB.QueryRow("select count(event_num) as Sum from ticket where event_num=?;", t.EventNum).Scan(
		&ticket.Sum)
	return
}

// GetLimitSeat 釋放座位數量
func (d *Detail) GetLimitSeat() (detail Detail, err error) {
	detail = Detail{}
	err = DB.QueryRow("select limit_seat as LimitSeat from ticket_detail where id=?;", d.ID).Scan(
		&detail.LimitSeat)
	return
}
