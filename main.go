package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defer DB.Close()

	router := gin.Default()

	//總計某一場次剩餘票數
	router.GET("/remain_tickets/:event_num", GetRemainTicket)

	//訂票紀錄
	router.GET("/ticket/:user_id/tickets", GetTickets)   //查詢訂票資料
	router.POST("/ticket/:event_num/user_id", AddTicket) //一次下一張單
	router.PATCH("/ticket/:id/status", UpdateTicket)     //更改注單狀態

	//表演詳細資料
	router.GET("/ticket_detail/:id", GetOneDetail)
	router.POST("/ticket_detail", AddOneDetail)

	//訂票者資訊
	router.GET("/user/:id", GetOne)
	router.POST("/user", AddOne)
	router.PATCH("/user/:id/status", UpdateUser)

	router.Run(":8000")
}
