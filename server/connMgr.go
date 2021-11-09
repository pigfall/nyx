package server

import(
	yy "github.com/pigfall/yingying"
)

type connMgr struct {
	conns      map[string]yy.Transport
	ipPoolIfce ipPoolIfce
}

func newConnMgr(ipPoolIfce ipPoolIfce) *connMgr{
	return &connMgr{
		conns:make(map[string]yy.Transport),
		ipPoolIfce:ipPoolIfce,
	}
}

func (this *connMgr) HasTransport(remoteAddr string)bool{
	_,ok := this.conns[remoteAddr]
	return ok
}
