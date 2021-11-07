module github.com/pigfall/yingying

go 1.17

replace github.com/pigfall/tzzGoUtil v1.4.0 => ../tzzGoUtil

require (
	github.com/google/gopacket v1.1.17
	github.com/gorilla/websocket v1.4.2
	github.com/pigfall/tzzGoUtil v1.4.0
	github.com/pigfall/wtun-go v0.1.1
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c
	golang.zx2c4.com/wireguard/windows v0.5.0
)

require (
	github.com/go-kit/kit v0.11.0 // indirect
	github.com/go-logfmt/logfmt v0.5.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8 // indirect
	github.com/vishvananda/netlink v1.1.0 // indirect
	github.com/vishvananda/netns v0.0.0-20191106174202-0a2b9b5464df // indirect
)

// replace github.com/pigfall/wtun-go => ../wtun-go

replace golang.zx2c4.com/wireguard/windows v0.5.0 => github.com/pigfall/wireguard-windows v0.5.0

//  replace golang.zx2c4.com/wireguard/windows => ../wireguard-windows
