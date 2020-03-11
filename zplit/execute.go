package zplit

import (
	"github.com/Marlinski/icannzplit/util"
	"github.com/Marlinski/icannzplit/vpn/ipvanish"

	"github.com/Marlinski/go-openvpn/events"
)

var (
	global chan events.OpenvpnEvent
)

// Execute a plan
func (rp *RoutingPlan) Execute() {
	global = make(chan events.OpenvpnEvent)

	for iface, ranges := range rp.Routes {
		mon, err := ipvanish.NewIfaceMonitor(string(iface), ranges, global)
		if err != nil {
			continue
		}

		go mon.StartOpenvpn()
	}

	for {
		select {
		case event := <-global:
			util.Log.Debugf(event.String())
		}
	}
}
