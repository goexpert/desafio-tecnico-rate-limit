package database

type IpRequests struct {
	IP  string
	Qty int
}

func NewRequest(ip string, qty int) IpRequests {
	return IpRequests{IP: ip, Qty: qty}
}
