package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var nFlag = flag.Int("n", 0, "help message for flag n")

//docker連線mysql
var (
	USERNAME = "root"
	PASSWORD = "demoroot"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = "3306"
	DATABASE = "ticket"
)

func main() {
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

	//測試是否有成功啟動
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "hellow world",
		})
	})

	//終止程式
	router.GET("/stop", func(c *gin.Context) {
		os.Exit(1)
	})

	router.Run(":8000")
}
