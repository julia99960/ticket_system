package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
