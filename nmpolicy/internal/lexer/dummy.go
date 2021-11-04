// Copyright 2021 The NMPolicy Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lexer

var (
	DefaultGwTokens = []Token{
		{0, IDENTITY, "routes"},
		{6, DOT, "."},
		{7, IDENTITY, "running"},
		{14, DOT, "."},
		{15, IDENTITY, "destination"},
		{26, EQFILTER, "=="},
		{28, STRING, "0.0.0.0/0"},
		{38, EOF, ""},
	}
	BaseIfaceRoutesTokens = []Token{
		{0, IDENTITY, "routes"},
		{6, DOT, "."},
		{7, IDENTITY, "running"},
		{14, DOT, "."},
		{15, IDENTITY, "next-hop-interface"},
		{33, EQFILTER, "=="},
		{35, IDENTITY, "capture"},
		{43, DOT, "."},
		{44, IDENTITY, "default-gw"},
		{54, DOT, "."},
		{55, IDENTITY, "routes"},
		{61, DOT, "."},
		{62, IDENTITY, "running"},
		{69, DOT, "."},
		{70, NUMBER, "0"},
		{71, DOT, "."},
		{72, IDENTITY, "next-hop-interface"},
		{89, EOF, ""},
	}
	BaseIfaceTokens = []Token{
		{0, IDENTITY, "interfaces"},
		{10, DOT, "."},
		{11, IDENTITY, "name"},
		{15, EQFILTER, "=="},
		{17, IDENTITY, "capture"},
		{25, DOT, "."},
		{26, IDENTITY, "default-gw"},
		{36, DOT, "."},
		{37, IDENTITY, "routes"},
		{43, DOT, "."},
		{44, IDENTITY, "running"},
		{51, DOT, "."},
		{52, NUMBER, "0"},
		{53, DOT, "."},
		{54, IDENTITY, "next-hop-interface"},
		{71, EOF, ""},
	}
	BridgeRoutesTokens = []Token{
		{0, IDENTITY, "capture"},
		{8, DOT, "."},
		{9, IDENTITY, "base-iface-routes"},
		{27, PIPE, "|"},
		{29, IDENTITY, "routes"},
		{35, DOT, "."},
		{36, IDENTITY, "running"},
		{43, DOT, "."},
		{44, IDENTITY, "next-hop-interface"},
		{62, REPLACE, ":="},
		{64, STRING, "br1"},
		{68, EOF, ""},
	}
	DeleteBaseIfaceRoutesTokens = []Token{
		{0, IDENTITY, "capture"},
		{8, DOT, "."},
		{9, IDENTITY, "base-iface-route"},
		{26, PIPE, "|"},
		{28, IDENTITY, "routes"},
		{34, DOT, "."},
		{35, IDENTITY, "running"},
		{42, DOT, "."},
		{43, IDENTITY, "state"},
		{48, REPLACE, ":="},
		{50, STRING, "absent"},
		{57, EOF, ""},
	}
	BridgeRoutesTakeoverTokens = []Token{
		{0, IDENTITY, "capture"},
		{8, DOT, "."},
		{9, IDENTITY, "delete-primary-nic-routes"},
		{34, DOT, "."},
		{35, IDENTITY, "routes"},
		{41, DOT, "."},
		{42, IDENTITY, "running"},
		{50, MERGE, "+"},
		{52, IDENTITY, "capture"},
		{60, DOT, "."},
		{61, IDENTITY, "bridge-routes"},
		{74, DOT, "."},
		{75, IDENTITY, "routes"},
		{81, DOT, "."},
		{82, IDENTITY, "running"},
		{88, EOF, ""},
	}
	dummy = map[string][]Token{
		`routes.running.destination=="0.0.0.0/0"`:                                                   DefaultGwTokens,
		`routes.running.next-hop-interface==capture.default-gw.routes.running.0.next-hop-interface`: BaseIfaceRoutesTokens,
		`interfaces.name==capture.default-gw.routes.running.0.next-hop-interface`:                   BaseIfaceTokens,
		`capture.base-iface-routes | routes.running.next-hop-interface:="br1"`:                      BridgeRoutesTokens,
		`capture.base-iface-route | routes.running.state:="absent"`:                                 DeleteBaseIfaceRoutesTokens,
		`capture.delete-base-iface-routes.routes.running + capture.bridge-routes.routes.running`:    BridgeRoutesTakeoverTokens,
	}
)
