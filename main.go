package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var nFlag = flag.Int("n", 1, "help message for flag n")

//本地測試資料庫連線資訊
const (
	USERNAME = "root"
	PASSWORD = "shupa0215"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = "3306"
	DATABASE = "ticket"
)

func main() {

	flag.Parse()

	//server資料庫連線資訊
	if *nFlag == 0 {
		const (
			USERNAME = "root"
			PASSWORD = "shupa0127"
			NETWORK  = "tcp"
			SERVER   = "34.80.173.177"
			PORT     = "3306"
			DATABASE = "ticket"
		)
	}

	defer DB.Close()

	router := gin.Default()

	//訂票紀錄
	router.POST("/ticket/:event_num/user_id", CreateTicket)   //一次下一張單
	router.GET("/ticket/:user_id/tickets", GetTickets)        //查詢訂票資料
	router.PATCH("/ticket/:id/status", UpdateTicket)          //更改注單狀態
	router.GET("/remain_tickets/:event_num", GetRemainTicket) //總計某一場次剩餘票數

	//表演詳細資料
	router.POST("/ticket_detail", AddOneDetail)
	router.GET("/ticket_detail/:id", GetOneDetail)

	//訂票者資訊
	router.POST("/user", AddOne)
	router.GET("/user/:id", GetOne)
	router.PATCH("/user/:id/status", UpdateUser)

	router.Run(":8000")
}
