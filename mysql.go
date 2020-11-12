package main

import (
	"database/sql"
	"fmt"
	"log"
)

//DB mysql
var DB *sql.DB

func init() {
	var err error
	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	DB, err = sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
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
	rs, err := DB.Exec("INSERT IGNORE user (id_number, mail, name, birthday, status) VALUES (?, ?, ?, ?, ?);",
		u.IDNumber, u.Mail, u.Name, u.Birth, u.Status)
	if err != nil {
		return 0
	}
	row, err := rs.RowsAffected()
	if err == nil && row != 0 {
		return 1
	}
	return 0
}

// Update 更新使用者狀態
func (u *User) Update() int64 {
	rs, err := DB.Exec("update user set status = ? where id = ?;", u.Status, u.ID)
	if err != nil {
		return 0
	}
	row, err := rs.RowsAffected()
	if err == nil && row != 0 {
		return row
	}
	return 0
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
