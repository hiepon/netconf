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

package cfglxccmd

import (
	lib "netconf/app/cfg/lxc/lib"
	lxdlib "netconf/lib/lxd"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Command struct {
	verbose bool
	image   string
}

func (c *Command) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Host port")
	cmd.PersistentFlags().StringVarP(&c.image, "image", "I", lib.PROFILE_IMAGE_NANE, "Container Image.")
	return cmd
}

func (c *Command) Init() {
	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func (c *Command) Client() (*lxdlib.Client, error) {
	client := lxdlib.NewClient(c.image)
	if err := client.Connect(); err != nil {
		log.Errorf("Connect error. %s", err)
		return nil, err
	}
	return client, nil
}
