package zplit

import (
	"encoding/json"
	"icannzplit/util"
	"icannzplit/vpn/ipvanish"
	"math"
	"math/bits"
	"strconv"
)

// IPRange contains a list of ip address range
type IPRange []string

// RoutingPlan maps some ip ranges to an openvpn interface
type RoutingPlan struct {
	Routes map[string]IPRange `json:"routes"`
}

// BuildPlan creates the routing plan from a project configuration
func (c *Config) BuildPlan() *RoutingPlan {
	return c.filterIPVanish().splitAddressSpace()
}

func (c *Config) filterIPVanish() *RoutingPlan {
	rp := &RoutingPlan{
		Routes: make(map[string]IPRange),
	}

	ipvanish.FilterConfigs(c.IPVanishSettings)

	for _, iface := range ipvanish.ListInterface() {
		rp.Routes[iface] = make([]string, 0)
	}

	return rp
}

func (rp *RoutingPlan) splitAddressSpace() *RoutingPlan {
	nbRoutes := len(rp.Routes)

	// get a slice of all the routing table keys
	keys := make([]string, nbRoutes)
	it := 0
	for k := range rp.Routes {
		keys[it] = string(k)
		it++
	}

	// we find an address mask with at least as many network addresses as routes
	addressMask := bits.Len32(uint32(nbRoutes))
	bitshift := 32 - addressMask

	// we distribute the address space to the routes
	for i := 0; i < int(math.Pow(2, float64(addressMask))); i++ {
		intip := i << bitshift
		network := util.ITo2IP(uint32(intip))
		cidr := network.String() + "/" + strconv.Itoa(addressMask)

		iface := keys[i%nbRoutes]
		rp.Routes[iface] = append(rp.Routes[iface], cidr)
	}

	return rp
}

// Dump the routing plan
func (rp *RoutingPlan) Dump() *RoutingPlan {
	s, _ := json.MarshalIndent(rp, "", "\t")
	util.Log.Noticef("routing plan..\n" + string(s))
	return rp
}
