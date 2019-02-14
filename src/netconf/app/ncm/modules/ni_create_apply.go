// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ncm

import (
	ncmdbm "netconf/app/ncm/dbm"
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"

	log "github.com/sirupsen/logrus"
)

type NICreateApplyHandler struct {
	*NIAnyHandler
}

func NewNICreateApplyHandler(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) NIChangeHandler {
	return &NICreateApplyHandler{
		NIAnyHandler: newNIAnyHandler(ev, oper),
	}
}

func (h *NICreateApplyHandler) Begin(name string, ni *openconfig.NetworkInstance) error {
	log.Debugf("NI/%s/%s/%s/BEGIN: %s", h.ev, h.oper, name, ni)

	h.Clear()
	if err := openconfig.ProcessNetworkInstance(h, false, name, ni); err != nil {
		h.Clear()
		return err
	}

	return h.DoCmds()
}

func (h *NICreateApplyHandler) NetworkInstance(name string, ni *openconfig.NetworkInstance) error {
	log.Debugf("NI/%s/%s/%s: %s", h.ev, h.oper, name, ni)

	AddNIContainerCmd(h, name, true)

	return nil
}

func (h *NICreateApplyHandler) NetworkInstanceConfig(name string, config *openconfig.NetworkInstanceConfig) error {
	log.Debugf("NI/%s/%s/%s/CONF: %s", h.ev, h.oper, name, config)

	if config.GetChange(openconfig.OC_TYPE_KEY) {
		AddNIContainterInitCmd(h, name, config)
	}

	if config.GetChange(openconfig.NETWORKINSTANCE_ROUTERID_KEY) {
		AddNIRouterIdCmd(h, name, config.RouterId.String(), true)
	}

	if config.OneOfChange(openconfig.NETWORKINSTANCE_RD_KEY, openconfig.NETWORKINSTANCE_RT_KEY) {
		AddNIVrfCmd(h, name, config.RD.String(), config.RT.String(), true)
	}

	return nil
}

func (h *NICreateApplyHandler) NetworkInstanceLoopbackAddrConfig(name string, id string, index string, config *openconfig.NetworkInstanceLoopbackAddrConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, index, config)

	AddNIVtyInterfaceCmd(h, name, "lo", "ip address", config.IFAddr().IPNet(), true)

	return nil
}

func (h *NICreateApplyHandler) NetworkInstanceInterface(name string, id string, iface *openconfig.NetworkInstanceInterface) error {
	log.Debugf("NI/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, iface)
	return nil
}

func (h *NICreateApplyHandler) NetworkInstanceInterfaceConfig(name string, id string, config *openconfig.NetworkInstanceInterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/CONF: %s", h.ev, h.oper, name, id, config)

	subif, device, _ := ncmdbm.Subinterfaces().SelectById(id)

	if config.GetChange(openconfig.OC_ID_KEY) {
		if subif.Index == 0 {
			iface, _ := ncmdbm.Interfaces().Select(id)
			hwaddr := iface.Ethernet.Config.MacAddr.String()
			AddNIContainerInterfaceCmd(h, name, id, hwaddr, true)
		}

		AddNIInterfaceNetworkCmd(h, name, device, subif, true)
		AddNIInterfaceSysctlCmd(h, name, id, true)

	} else if config.GetChanges(openconfig.SUBINTERFACE_KEY) {
		AddNIInterfaceNetworkCmd(h, name, device, subif, true)
	}

	if config.GetChanges(openconfig.INTERFACE_KEY, openconfig.SUBINTERFACE_KEY) {
		AddNIVtyInterfaceCmd(h, name, id, "shutdown", "", false)

		for _, ifv4 := range subif.IPv4.Addresses {
			AddNIVtyInterfaceCmd(h, name, id, "ip address", ifv4.Config.IFAddr(), true)
		}
	}

	log.Debugf("NI/%s/%s/%s/%s: OK", h.ev, h.oper, name, id)
	return nil
}

func (h *NICreateApplyHandler) Ospfv2GlobalConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.Ospfv2GlobalConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/CONF: %s", h.ev, h.oper, name, key, config)

	if config.GetChange(openconfig.OSPFV2_ROUTERID_KEY) {
		AddNIOspfRouterCmd(h, name, "router-id", config.RouterId, true)
	}

	return nil
}

func (h *NICreateApplyHandler) Ospfv2InterfaceConfig(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, config *openconfig.Ospfv2InterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF: %s", h.ev, h.oper, name, key, areaId, ifaceId, config)

	if config.GetChange(openconfig.OSPFV2_METRIC_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf cost", config.Metric, true)
	}

	if config.GetChange(openconfig.OSPFV2_PASSIVE_KEY) {
		AddNIOspfRouterCmd(h, name, "passive-interface", ifaceId, config.Passive)
	}

	if config.GetChange(openconfig.OSPFV2_PRIORITY_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf priority", config.Priority, true)
	}

	if config.GetChange(openconfig.OSPFV2_NETWORK_TYPE_KEY) {
		n, _ := getOspfNetworkType(config)
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf network", n, true)
	}

	return nil
}

func (h *NICreateApplyHandler) Ospfv2InterfaceRefConfig(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, config *openconfig.InterfaceRefConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/IFREF: %s", h.ev, h.oper, name, key, areaId, ifaceId, config)

	if config.GetChanges(openconfig.INTERFACE_KEY, openconfig.SUBINTERFACE_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf area", areaId, true)
	}

	return nil
}

func (h *NICreateApplyHandler) Ospfv2InterfaceTimers(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, timers *openconfig.Ospfv2InterfaceTimers) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/TIMERS: %s", h.ev, h.oper, name, key, areaId, ifaceId, timers)

	if timers.GetChange(openconfig.OSPFV2_DEAD_INTERVAL_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf dead-interval", timers.DeadInterval, true)
	}

	if timers.GetChange(openconfig.OSPFV2_HELLO_INTERVAL_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf hello-interval", timers.HelloInterval, true)
	}

	return nil
}

func (h *NICreateApplyHandler) MplsInterfaceAttrRefConfig(name string, ifaceId string, config *openconfig.InterfaceRefConfig) error {
	log.Debugf("NI/%s/%s/%s/MPLS/%s/REF: %s", h.ev, h.oper, name, ifaceId, config)

	AddNIMplsInterfaceSysctlCmd(h, name, ifaceId, true)

	return nil
}

func (h *NICreateApplyHandler) MplsLdpConfig(name string, config *openconfig.MplsLdpConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/CONF: %s", h.ev, h.oper, name, config)

	if config.GetChange(openconfig.MPLS_LDP_ROUTERID_KEY) {
		AddNIMplsLdpCmd(h, name, "router-id", config.LsrId, true)
	}

	return nil
}

func (h *NICreateApplyHandler) MplsLdpAddressFamilyV4Config(name string, config *openconfig.MplsLdpAddressFamilyV4Config) error {
	log.Debugf("NI/%s/%s/%s/LDP/ADDR: %s", h.ev, h.oper, name, config)

	if config.GetChange(openconfig.MPLS_LDP_AF_TARNSADDR_KEY) {
		AddNIMplsLdpIPv4Cmd(h, name, "discovery transport-address", config.TransportAddr, true)
	}

	if config.GetChange(openconfig.MPLS_LDP_AF_SESSION_HOLDTIME_KEY) {
		AddNIMplsLdpIPv4Cmd(h, name, "session holdtime", config.SessionHoldTime, true)
	}

	explicitNll := config.LabelPolicy.Advertise.EngressExplicitNull
	if explicitNll.GetChange(openconfig.MPLS_LDP_LABELPOLICY_ENABLE_KEY) {
		AddNIMplsLdpIPv4Cmd(h, name, "label local advertise explicit-null", "", explicitNll.Enable)
	}

	return nil
}

func (h *NICreateApplyHandler) MplsLdpDiscovInterfacesConfig(name string, config *openconfig.MplsLdpDiscovInterfacesConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/IFCONF: %s", h.ev, h.oper, name, config)

	if config.GetChange(openconfig.MPLS_LDP_HELLO_HOLDTIME) {
		AddNIMplsLdpCmd(h, name, "discovery hello holdtime", config.HelloHoldTime, true)
	}

	if config.GetChange(openconfig.MPLS_LDP_HELLO_INTERVAL) {
		AddNIMplsLdpCmd(h, name, "discovery hello interval", config.HelloInterval, true)
	}

	return nil
}

func (h *NICreateApplyHandler) MplsLdpInterface(name string, ifaceId string, iface *openconfig.MplsLdpInterface) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/%s: %s", h.ev, h.oper, name, ifaceId, iface)
	AddNIMplsLdpIPv4Cmd(h, name, "interface", ifaceId, true)
	return nil
}

func (h *NICreateApplyHandler) MplsLdpInterfaceConfig(name string, ifaceId string, config *openconfig.MplsLdpInterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/%s/CONF: %s", h.ev, h.oper, name, ifaceId, config)

	if config.GetChange(openconfig.MPLS_LDP_HELLO_HOLDTIME) {
		AddNIMplsLdpIPv4IfaceCmd(h, name, ifaceId, "discovery hello holdtime", config.HelloHoldTime, true)
	}

	if config.GetChange(openconfig.MPLS_LDP_HELLO_INTERVAL) {
		AddNIMplsLdpIPv4IfaceCmd(h, name, ifaceId, "discovery hello interval", config.HelloInterval, true)
	}

	return nil
}

func (h *NICreateApplyHandler) StaticRouteNexthop(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, nexthop *openconfig.StaticRouteNexthop) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s: %s", h.ev, h.oper, name, prkey, rtkey, index, nexthop)

	ip, nhtype, _ := nexthop.Config.GetNexthop()
	switch nhtype {
	case openconfig.LOCAL_DEFINED_NEXT_HOP_LOCAL_LINK:
		AddNIStaticRouteCmd(h, name, rtkey.String(), nexthop.IfaceRef.Config.IFName(), true)
	case openconfig.LOCAL_DEFINED_NEXT_HOP_DROP:
		AddNIStaticRouteCmd(h, name, rtkey.String(), "null0", true)
	default:
		AddNIStaticRouteCmd(h, name, rtkey.String(), ip.String(), true)
	}

	return nil
}

func (h *NICreateApplyHandler) Bgp(name string, key *openconfig.NetworkInstanceProtocolKey, bgp *openconfig.Bgp) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s: %s", h.ev, h.oper, name, key, bgp)

	restart := bgp.Global.Config.OneOfChange(openconfig.BGP_AS_KEY, openconfig.BGP_ROUTERID_KEY)

	openconfig.ProcessBgp(h.Bgps, false, name, key, bgp)
	AddNIBgpConfigCmd(h, name, h.Bgps.Bytes(), restart, true)

	return nil
}

func (h *NICreateApplyHandler) BgpNeighborApplyPolicyConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.PolicyApplyConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/APPLYPOL: %s", h.ev, h.oper, name, key, addr, config)

	if config.GetChange(openconfig.POLICYAPPLY_IMPORT_KEY) {
		for _, polName := range config.ImportPolicy {
			pol, _ := ncmdbm.PolicyDefinitions().Select(polName)
			openconfig.ProcessPolicyDefinition(h.Bgps, false, polName, pol)
		}
	}

	if config.GetChange(openconfig.POLICYAPPLY_EXPORT_KEY) {
		for _, polName := range config.ExportPolicy {
			pol, _ := ncmdbm.PolicyDefinitions().Select(polName)
			openconfig.ProcessPolicyDefinition(h.Bgps, false, polName, pol)
		}
	}

	AddNIBgpConfigCmd(h, name, h.Bgps.Bytes(), false, true)

	return nil
}
