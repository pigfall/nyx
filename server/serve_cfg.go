package server

type ServeCfg struct{
	Port int `json:"port"`
	Mode string `json:"mode"`// udp or ws
}
