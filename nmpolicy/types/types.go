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

package types

import "time"

type PolicySpec struct {
	Capture      map[string]string
	DesiredState []byte
}

type CachedState struct {
	Capture map[string]CaptureState
}

type GeneratedState struct {
	Cache        CachedState
	DesiredState []byte
	MetaInfo     MetaInfo
}

type CaptureState struct {
	State    []byte
	MetaInfo MetaInfo
}

type MetaInfo struct {
	Version   string
	TimeStamp time.Time
}

func NoCache() CachedState { return CachedState{} }
