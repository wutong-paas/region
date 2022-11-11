/*
Copyright 2022 The Wutong Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kube

import "testing"

func TestDownload(t *testing.T) {
	b, err := Download("https://raw.githubusercontent.com/Huawei/eSDK_K8S_Plugin/master/deploy/huawei-csi-node.yaml")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(b))
	}
}