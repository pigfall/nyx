package proto


type DNSQuery struct{
	ID string `json:"id"`
	HostNames []string `json:"hostName"`
}

type DNSQueryRes struct{
	ID string
	Answers []*DNSQueryAnswers
}

type DNSQueryAnswers struct{
	HostName string
	Ips []string
}
