/*
 * Copyright 2021 NMPolicy Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *	  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

import "sigs.k8s.io/yaml"

func (s *NMState) UnmarshalJSON(b []byte) error {
	output, err := yaml.JSONToYAML(b)
	if err != nil {
		return err
	}
	*s = output
	return nil
}

func (s NMState) MarshalJSON() ([]byte, error) {
	return yaml.YAMLToJSON([]byte(s))
}
