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

package ncgobgp

type Statement Entries

func NewStatements(i interface{}) []Statement {
	stmts := []Statement{}

	switch i.(type) {
	case nil:
	default:
		for _, s := range i.([]interface{}) {
			stmts = append(stmts, NewStatement(s))
		}
	}

	return stmts
}

func NewStatement(i interface{}) Statement {
	return Statement(NewEntries(i))
}

func RawStatements(stmts []Statement) interface{} {
	list := make([]interface{}, len(stmts))
	for index, stmt := range stmts {
		list[index] = Entries(stmt).Raw()
	}
	return list
}
