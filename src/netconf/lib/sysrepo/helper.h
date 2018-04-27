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

#ifndef _MAIN_H
#define _MAIN_H

#include <sysrepo.h>

#ifdef __cplusplus
extern "C" {
#endif

  sr_val_t *get_val(sr_val_t*, size_t);

  void log_cb(sr_log_level_t, const char*);

  int module_change_cb(sr_session_ctx_t*, const char*, sr_notif_event_t, void*);

  int subtree_change_cb(sr_session_ctx_t*, const char*, sr_notif_event_t, void*);

#ifdef __cplusplus
}
#endif

#endif // _MAIN_H