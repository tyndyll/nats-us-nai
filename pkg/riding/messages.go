package riding

type RideRequest struct {
	RequestID   string `json:"request_id"`
	CustomerID  string `json:"customer_id"`
	Destination string `json:"destination"`
	TimeStamp   string `json:"timestamp"`
}

type DriverInterested struct {
	ID      string       `json:"id"`
	Request *RideRequest `json:"request"`
}
