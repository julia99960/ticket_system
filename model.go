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
