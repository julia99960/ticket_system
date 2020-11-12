package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTicket 新增一筆訂票紀錄
func CreateTicket(c *gin.Context) {
	eventnum, _ := strconv.Atoi(c.Param("event_num"))
	userid, _ := strconv.Atoi(c.Request.FormValue("user_id"))
	status, _ := strconv.Atoi(c.Request.FormValue("status"))

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

// GetTickets 取得使用者訂票紀錄
func GetTickets(c *gin.Context) {
	userid, _ := strconv.Atoi(c.Param("user_id"))
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

	row := u.AddOneUser()
	if row == 1 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "unsuccess",
		})
	}
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
