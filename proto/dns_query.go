package proto


type DNSQuery struct{
	ID string `json:"id"`
	HostNames []string `json:"hostName"`
}
