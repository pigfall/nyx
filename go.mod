module github.com/pigfall/yingying

go 1.17

replace github.com/pigfall/tzzGoUtil v1.3.0 => ../tzzGoUtil

require (
	github.com/gorilla/websocket v1.4.2
	github.com/pigfall/tzzGoUtil v1.3.0
)

require (
	github.com/go-kit/kit v0.11.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/google/gopacket v1.1.17 // indirect
	github.com/pigfall/wtun-go v0.0.0-20211015084304-c7063103c2c3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8 // indirect
	github.com/vishvananda/netlink v1.1.0 // indirect
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
	golang.org/x/net v0.0.0-20210927181540-4e4d966f7476 // indirect
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	golang.zx2c4.com/wireguard v0.0.0-00010101000000-000000000000 // indirect
)

replace golang.zx2c4.com/wireguard => ../wintun-go

replace github.com/pigfall/wtun-go => ../wtun-go
