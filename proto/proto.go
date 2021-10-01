package proto

import(
		"fmt"
)


type Msg struct {
	Id int `json:"id"`
	ErrReason string  `json:"errReason"`
	ErrMsg string `json:"errMsg"`
	Body []byte `json:"body"`
}

func (this *Msg) Err()error{
	if this.ErrReason == ""{
		return nil
	}

	return this
}

func (this *Msg) Error() string{
	return fmt.Sprintf("%s %s",this.ErrReason,this.ErrMsg)
}



type S2C_ClientVPNIpNet struct{
	// format 127.0.0.1/8
	IpNet string `json:"ipNet"`
}
