package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
