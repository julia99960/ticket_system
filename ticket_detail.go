package main

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

// GetPerformanceDetail 獲取表演場次細節
func GetPerformanceDetail(ida int) (detail Detail, err error) {
	d := Detail{
		ID: ida,
	}
	detail, err = d.GetOne()
	return
}
