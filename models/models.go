package models

// Product schema of the product table
type Product struct {
	ID         int64  `json:"id"`
	Pname      string `json:"pname"`
	Pdesc      string `json:"pdesc"`
	Mrp        int64  `json:"mrp"`
	StBidPrice int64  `json:"stBidPrice"`
}
