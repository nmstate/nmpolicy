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

package tests

import (
	"time"

	"github.com/nmstate/nmpolicy/nmpolicy/types"
)

func resetTimeStamp(generatedState types.GeneratedState) types.GeneratedState {
	generatedState.MetaInfo.TimeStamp = time.Time{}
	for captureID, captureState := range generatedState.Cache.Capture {
		captureState.MetaInfo.TimeStamp = time.Time{}
		generatedState.Cache.Capture[captureID] = captureState
	}
	return generatedState
}
