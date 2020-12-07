package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

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

// fileName 文件檔案名稱
var fileName string = fmt.Sprintf("%s_errorlog", time.Now().Format("20060102"))

// ErrLogPath 文件檔案路徑
var ErrLogPath string = fmt.Sprintf("logs/%s.txt", fileName)

func init() {
	// 建立每日紀錄log的檔案
	if !FileExists(ErrLogPath) {
		file, err := os.Create(ErrLogPath)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	}
}

// WriteErrorLog 紀錄 error log
func WriteErrorLog(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY, 777)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

// FileExists 判斷檔案是否存在
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

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
