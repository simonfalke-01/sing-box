package main

import (
	"net/netip"
	"testing"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-shadowsocks/shadowaead_2022"
)

func TestMuxBrutal(t *testing.T) {
	method := shadowaead_2022.List[0]
	password := mkBase64(t, 16)
	startInstance(t, option.Options{
		Inbounds: []option.Inbound{
			{
				Type: C.TypeMixed,
				Tag:  "mixed-in",
				MixedOptions: option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     option.NewListenAddress(netip.IPv4Unspecified()),
						ListenPort: clientPort,
					},
				},
			},
			{
				Type: C.TypeShadowsocks,
				ShadowsocksOptions: option.ShadowsocksInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     option.NewListenAddress(netip.IPv4Unspecified()),
						ListenPort: serverPort,
					},
					Method:   method,
					Password: password,
					Multiplex: &option.InboundMultiplexOptions{
						Enabled: true,
						Brutal: &option.BrutalOptions{
							Enabled:  true,
							UpMbps:   100,
							DownMbps: 100,
						},
					},
				},
			},
		},
		Outbounds: []option.Outbound{
			{
				Type: C.TypeDirect,
			},
			{
				Type: C.TypeShadowsocks,
				Tag:  "ss-out",
				ShadowsocksOptions: option.ShadowsocksOutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     "127.0.0.1",
						ServerPort: serverPort,
					},
					Method:   method,
					Password: password,
					Multiplex: &option.OutboundMultiplexOptions{
						Enabled:  true,
						Protocol: "smux",
						Padding:  true,
						Brutal: &option.BrutalOptions{
							Enabled:  true,
							UpMbps:   100,
							DownMbps: 100,
						},
					},
				},
			},
		},
		Route: &option.RouteOptions{
			Rules: []option.Rule{
				{
					DefaultOptions: option.DefaultRule{
						Inbound:  []string{"mixed-in"},
						Outbound: "ss-out",
					},
				},
			},
		},
	})
	testSuit(t, clientPort, testPort)
}
