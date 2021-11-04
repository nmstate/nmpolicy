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
//

package ast

type Meta struct {
	Position int `json:"pos"`
}

type UnaryOperator Node
type BinaryOperator [2]Node
type TernaryOperator [3]Node
type VariadicOperator []Node
type Terminal struct {
	String   *string `json:"string,omitempty"`
	Identity *string `json:"identity,omitempty"`
	Number   *int    `json:"number,omitempty"`
}

type Node struct {
	Meta
	Terminal
	EqFilter *TernaryOperator  `json:"eqfilter,omitempty"`
	Merge    *BinaryOperator   `json:"merge,omitempty"`
	Path     *VariadicOperator `json:"path,omitempty"`
	Pipe     *UnaryOperator    `json:"pipe,omitempty"`
	Replace  *TernaryOperator  `json:"replace,omitempty"`
}
