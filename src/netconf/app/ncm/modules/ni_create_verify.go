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
	"fmt"
	ncmdbm "netconf/app/ncm/dbm"
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"

	log "github.com/sirupsen/logrus"
)

type NICreateVerifyHandler struct {
	*NIAnyHandler
}

func NewNICreateVerifyHandler(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) NIChangeHandler {
	return &NICreateVerifyHandler{
		NIAnyHandler: newNIAnyHandler(ev, oper),
	}
}

func (h *NICreateVerifyHandler) Begin(name string, ni *openconfig.NetworkInstance) error {
	log.Debugf("NI/%s/%s/%s/BEGIN; %s", h.ev, h.oper, name, ni)
	h.Clear()
	return openconfig.ProcessNetworkInstance(h, false, name, ni)
}

func (h *NICreateVerifyHandler) NetworkInstanceLoopbackAddrConfig(name string, id string, index string, config *openconfig.NetworkInstanceLoopbackAddrConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, index, config)

	if err := VerifyNILoopbackAddrConfig(config); err != nil {
		log.Errorf("NI/%s/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, index, err)
		return err
	}

	return nil
}

//
// /network-instances/network-instance[name]/interfaces/interface[id]
//
func (h *NICreateVerifyHandler) NetworkInstanceInterface(name string, id string, iface *openconfig.NetworkInstanceInterface) error {
	log.Debugf("NI/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, iface)

	if err := VerifyNIInterface(iface); err != nil {
		log.Errorf("NI/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, err)
		return err
	}

	log.Debugf("NI/%s/%s/%s/%s: OK", h.ev, h.oper, name, id)
	return nil
}

//
// /network-instances/network-instance[name]/interfaces/interface[id]/config
//
func (h *NICreateVerifyHandler) NetworkInstanceInterfaceConfig(name string, id string, config *openconfig.NetworkInstanceInterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/CONF: %s", h.ev, h.oper, name, id, config)

	if err := VerifyNIInterfaceConfig(id, config); err != nil {
		log.Errorf("NI/%s/%s/%s/%s/CONF: %s", h.ev, h.oper, name, id, err)
		return err
	}

	log.Debugf("NI/%s/%s/%s/%s: OK", h.ev, h.oper, name, id)
	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/static-routes/static[rtkey]
//
func (h *NICreateVerifyHandler) StaticRoute(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, route *openconfig.StaticRoute) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s: %s", h.ev, h.oper, name, prkey, rtkey, route)

	if prkey.Ident != openconfig.INSTALL_PROTOCOL_STATIC {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s: Invalid Protocol identifier.", h.ev, h.oper, name, prkey, rtkey)
	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/static-routes/static[rtkey]/config
//
func (h *NICreateVerifyHandler) StaticRouteConfig(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, config *openconfig.StaticRouteConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/CONF %s", h.ev, h.oper, name, prkey, rtkey, config)

	if err := VerifyNIStaticRouteConfig(config); err != nil {
		log.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/CONF %s", h.ev, h.oper, name, prkey, rtkey, err)
		return err
	}

	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/CONF OK", h.ev, h.oper, name, prkey, rtkey)
	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/static-routes/static[rtkey]/next-hops/next-hop[index]
//
func (h *NICreateVerifyHandler) StaticRouteNexthop(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, nexthop *openconfig.StaticRouteNexthop) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s: %s", h.ev, h.oper, name, prkey, rtkey, index, nexthop)

	if chg := nexthop.GetChange(openconfig.OC_CONFIG_KEY); !chg {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s: do not change next-hop patrially.", h.ev, h.oper, name, prkey, rtkey, index)
	}

	_, nhType, err := nexthop.Config.GetNexthop()
	if err != nil {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s: invalid next-hop. %s", h.ev, h.oper, name, prkey, rtkey, index, err)
	}

	if nhType == openconfig.LOCAL_DEFINED_NEXT_HOP_LOCAL_LINK {
		if err := VerifyNIInterfaceRefConfig(nexthop.IfaceRef.Config); err != nil {
			return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s: %s", h.ev, h.oper, name, prkey, rtkey, index, err)
		}
	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/static-routes/static[rtkey]/next-hops/next-hop[index]/config
//
func (h *NICreateVerifyHandler) StaticRouteNexthopConfig(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, config *openconfig.StaticRouteNexthopConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF: %s", h.ev, h.oper, name, prkey, rtkey, index, config)

	if chg := config.GetChange(openconfig.STATICROUTE_NEXTHOP_KEY); !chg {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF: next-hop not specified.", h.ev, h.oper, name, prkey, rtkey, index)
	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/ospfv2
//
func (h *NICreateVerifyHandler) Ospfv2(name string, key *openconfig.NetworkInstanceProtocolKey, ospf *openconfig.Ospfv2) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s: %s", h.ev, h.oper, name, key, ospf)

	if key.Ident != openconfig.INSTALL_PROTOCOL_OSPF {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s: Invalid Protocol identifier.", h.ev, h.oper, name, key)

	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/ospfv3
//
func (h *NICreateVerifyHandler) Ospfv3(name string, key *openconfig.NetworkInstanceProtocolKey, ospf *openconfig.Ospfv3) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s: %s", h.ev, h.oper, name, key, ospf)

	if key.Ident != openconfig.INSTALL_PROTOCOL_OSPF3 {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s: Invalid Protocol identifier.", h.ev, h.oper, name, key)

	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/bgp
//
func (h *NICreateVerifyHandler) Bgp(name string, key *openconfig.NetworkInstanceProtocolKey, bgp *openconfig.Bgp) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s: %s", h.ev, h.oper, name, key, bgp)

	if key.Ident != openconfig.INSTALL_PROTOCOL_BGP {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s: Invalid Protocol identifier.", h.ev, h.oper, name, key)

	}

	if err := openconfig.ProcessBgp(h.Bgps, false, name, key, bgp); err != nil {
		log.Errorf("NI/%s/%s/%s/PROTOS/%s: Processor error. %s", h.ev, h.oper, name, key, err)
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s: Processor error. %s", h.ev, h.oper, name, key, err)
	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/bgp/neighbors/neighbor[addr]/apply-policy/config
//
func (h *NICreateVerifyHandler) BgpNeighborApplyPolicyConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.PolicyApplyConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/APPLYPOL: %s", h.ev, h.oper, name, key, addr, config)

	if config.GetChange(openconfig.POLICYAPPLY_IMPORT_KEY) {
		for _, polName := range config.ImportPolicy {
			pol, err := ncmdbm.PolicyDefinitions().Select(polName)
			if err != nil {
				return fmt.Errorf("ImportPolicy not found in Policy module. %s. %s", polName, err)
			}

			if err := openconfig.ProcessPolicyDefinition(h.Bgps, false, polName, pol); err != nil {
				return fmt.Errorf("ImportPolicy process error. %s %s", polName, err)
			}
		}
	}

	if config.GetChange(openconfig.POLICYAPPLY_EXPORT_KEY) {
		for _, polName := range config.ExportPolicy {
			pol, err := ncmdbm.PolicyDefinitions().Select(polName)
			if err != nil {
				return fmt.Errorf("ExportPolicy not found in Policy module. %s %s", polName, err)
			}

			if err := openconfig.ProcessPolicyDefinition(h.Bgps, false, polName, pol); err != nil {
				return fmt.Errorf("ExportPolicy process error. %s %s", polName, err)
			}
		}
	}

	return nil
}
