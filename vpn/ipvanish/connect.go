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
func (m *IfaceMonitor) StartOpenvpn(defaultRoute netlink.Route) {
	// prepare openvpn
	configfile := filepath.Join(configDir, m.config.FileName)
	cfg := openvpn.LoadConfig(m.iface, configfile)

	// save the route to the tunnel endpoint
	remote, err := cfg.GetRemote()
	if err != nil {
		return
	}

	err = m.saveOpenvpnRoute(remote, defaultRoute)
	if err != nil {
		return
	}

	// add command flag to openvpn
	cfg.SetLogLevel(logging.DEBUG)
	cfg.Set("ca", certFile)
	cfg.Set("dev-type", "tun")
	cfg.Set("dev", m.config.DeviceName)
	cfg.Set("auth-user-pass", authFile)
	cfg.Flag("route-noexec")
	cfg.Flag("suppress-timestamps")
	cfg.Flag("nobind")
	cfg.Flag("mute-replay-warnings")

	// event channel
	m.channel = make(chan events.OpenvpnEvent)
	m.processEvents()

	// Create the openvpn instance
	m.ctrl = cfg.Run(m.channel)
}

func (m *IfaceMonitor) saveOpenvpnRoute(remote string, defaultRoute netlink.Route) error {
	addrs, err := net.LookupIP(remote)
	if err != nil {
		return err
	}

	// should have only 1 IP address
	for _, addr := range addrs {
		if addr.To4() == nil {
			// only supports IPv4 for now
			continue
		}

		defaultRoute.Dst = &net.IPNet{
			IP:   addr,
			Mask: net.IPv4Mask(0xff, 0xff, 0xff, 0xff),
		}
		defaultRoute.Protocol = 4 // proto static

		util.Log.Noticef("adding static route to %s (%s): %+v", addr, remote, defaultRoute)
		if err := netlink.RouteAdd(&defaultRoute); err != nil {
			return err
		}
	}

	return nil
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
	// local tunnel internal ip
	ifaceAddr, err := m.ctrl.GetOpenVpnEnv("ifconfig_local")
	if err != nil {
		util.Log.Errorf("ip address not found")
		return
	}
	src := net.ParseIP(ifaceAddr)

	// remote internal ip
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

	// device attributes
	dev, err := netlink.LinkByName(m.iface)
	if err != nil {
		util.Log.Errorf("cannot get interface: %s", m.iface)
		return
	}

	// set up all the routes
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

		nlroute := netlink.Route{
			LinkIndex: dev.Attrs().Index,
			Dst:       dst,
			Src:       src,
			Gw:        gw,
			Protocol:  4, // proto: static
		}

		if err := netlink.RouteAdd(&nlroute); err != nil {
			util.Log.Errorf("cannot add route: %s", route)
			continue
		}
		util.Log.Noticef("route added: %s", route)
	}
}
