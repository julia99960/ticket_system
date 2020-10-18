package main

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

// GetTicketsList 取得使用者購票紀錄
func GetTicketsList(userid int) (tickets []Ticket, err error) {
	t := Ticket{
		UserID: userid,
	}
	tickets, err = t.GetRow()
	return
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
