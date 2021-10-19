package client

import(
	"github.com/google/gopacket"	
	"github.com/google/gopacket/layers"	
)


type dnsToQuery struct{
	packet gopacket.Packet
	layerDNS *layers.DNS
}
