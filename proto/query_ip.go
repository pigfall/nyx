package proto


// Query client vpn ip which is the server decided
type C2S_QueryIp struct{
	IpNet  string`json:"ipNet"`
}
