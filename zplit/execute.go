package zplit

import (
	"net"
	"os"

	"github.com/Marlinski/icannzplit/util"
	"github.com/Marlinski/icannzplit/vpn/ipvanish"
	"github.com/vishvananda/netlink"

	"github.com/Marlinski/go-openvpn/events"
)

var (
	global chan events.OpenvpnEvent
)

// Execute a plan
func (rp *RoutingPlan) Execute() {
	// we select the default route as the one routing toward google dns server
	routes, err := netlink.RouteGet(net.ParseIP("8.8.8.8"))
	if err != nil {
		util.Log.Errorf("%+v", err)
		os.Exit(1)
	}

	if len(routes) == 0 {
		util.Log.Errorf("no default route found")
		os.Exit(1)
	}

	util.Log.Noticef("using default route via: %+v", routes[0].Gw)
	for iface, ranges := range rp.Routes {
		mon, err := ipvanish.NewIfaceMonitor(string(iface), ranges, global)
		if err != nil {
			continue
		}

		go mon.StartOpenvpn(routes[0])
	}
}
