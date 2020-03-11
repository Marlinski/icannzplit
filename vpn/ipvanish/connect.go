package ipvanish

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/Marlinski/icannzplit/util"

	"github.com/op/go-logging"

	"github.com/Marlinski/go-openvpn"
	"github.com/Marlinski/go-openvpn/events"
	"github.com/vishvananda/netlink"
)

// IfaceMonitor manage a specific interface
type IfaceMonitor struct {
	iface   string
	config  TunnelConfig
	routes  []string
	ctrl    openvpn.Controller
	channel chan events.OpenvpnEvent
}

// NewIfaceMonitor opens up a tunnel, sets up the routes  and monitor its status
func NewIfaceMonitor(iface string, routes []string, global chan events.OpenvpnEvent) (*IfaceMonitor, error) {
	// find config for this interface
	var config TunnelConfig
	for _, cfg := range configs {
		if cfg.DeviceName == iface {
			config = cfg
			break
		}
	}

	if config == (TunnelConfig{}) {
		return nil, fmt.Errorf("%s: config not found", iface)
	}

	// start it
	mon := IfaceMonitor{
		iface:  iface,
		config: config,
		routes: routes,
	}

	return &mon, nil
}

// StartOpenvpn fires up the tunnel
func (m *IfaceMonitor) StartOpenvpn() {
	// prepare openvpn
	configfile := filepath.Join(configDir, m.config.FileName)
	openvpnCfg := openvpn.LoadConfig(m.iface, configfile)
	openvpnCfg.SetLogLevel(logging.DEBUG)
	//openvpnCfg.SetLogStd(true)
	openvpnCfg.Set("ca", certFile)
	openvpnCfg.Set("dev-type", "tun")
	openvpnCfg.Set("dev", m.config.DeviceName)
	openvpnCfg.Set("auth-user-pass", authFile)
	openvpnCfg.Flag("route-noexec")
	openvpnCfg.Flag("suppress-timestamps")
	openvpnCfg.Flag("nobind")
	openvpnCfg.Flag("mute-replay-warnings")

	// event channel
	m.channel = make(chan events.OpenvpnEvent)
	m.processEvents()

	// Create the openvpn instance
	m.ctrl = openvpnCfg.Run(m.channel)
}

func (m *IfaceMonitor) processEvents() {
	go func() {
		for {
			e := <-m.channel
			util.Log.Noticef("event: %s", e.String())
			if e.Code() == events.OpenvpnEventUp {
				m.setUpRoute()
			}
		}
	}()
}

func (m *IfaceMonitor) setUpRoute() {
	vpnGateway, err := m.ctrl.GetOpenVpnEnv("route_vpn_gateway")
	if err != nil {
		util.Log.Errorf("vpn gateway not found")
		return
	}

	gw := net.ParseIP(vpnGateway)
	if gw == nil {
		util.Log.Errorf("wrong vpn gateway ip %s: ", gw)
		return
	}

	for _, route := range m.routes {
		_, ipnet, err := net.ParseCIDR(route)
		if err != nil {
			util.Log.Errorf("cannot parse route: %s", route)
			continue
		}

		dst := &net.IPNet{
			IP:   ipnet.IP,
			Mask: ipnet.Mask,
		}

		dev, err := netlink.LinkByName(m.iface)
		if err != nil {
			util.Log.Errorf("cannot get interface: %s", m.iface)
			continue
		}

		nlroute := netlink.Route{
			LinkIndex: dev.Attrs().Index,
			Dst:       dst,
			Gw:        gw,
		}

		if err := netlink.RouteAdd(&nlroute); err != nil {
			util.Log.Errorf("cannot add route: %s", route)
			continue
		}
		util.Log.Noticef("route added: %s", route)
	}
}
