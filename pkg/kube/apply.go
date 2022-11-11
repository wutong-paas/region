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

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlserializer "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	sigsyaml "sigs.k8s.io/yaml"
)

type ApplyExector struct {
	KubeConfig    *rest.Config
	KubeClient    kubernetes.Interface
	DynamicClient dynamic.Interface
	Content       string
	Namespace     string
}

func NewApplyExector(content, namespace string, restConfig *rest.Config) *ApplyExector {
	return &ApplyExector{
		KubeConfig:    restConfig,
		KubeClient:    kubernetes.NewForConfigOrDie(restConfig),
		DynamicClient: dynamic.NewForConfigOrDie(restConfig),
		Content:       content,
		Namespace:     namespace,
	}
}

func (exector *ApplyExector) GetGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	grs, err := restmapper.GetAPIGroupResources(exector.KubeClient.Discovery())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(grs)

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return mapping.Resource, nil
}

func (exector *ApplyExector) Apply() error {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer([]byte(exector.Content)), 4096)

	var err error
	for {
		var rawObj runtime.RawExtension
		err = decoder.Decode(&rawObj)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("decode failed: %s", err.Error())
		}

		obj, _, err := yamlserializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return fmt.Errorf("new yaml decoding serializer failed: %s", err.Error())
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return fmt.Errorf("convert to unstructured object faild: %s", err.Error())
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, err := exector.GetGVR(unstructuredObj.GroupVersionKind())
		if err != nil {
			return fmt.Errorf("get gvr failed: %s", err.Error())
		}

		unstructuredContent, err := sigsyaml.Marshal(unstructuredObj)
		if err != nil {
			return fmt.Errorf("marshal unstructured object failed: %s", err.Error())
		}

		objNamespace := exector.Namespace
		if unstructuredObj.GetNamespace() != "" {
			objNamespace = unstructuredObj.GetNamespace()
		}

		_, err = exector.DynamicClient.Resource(gvr).Namespace(objNamespace).Get(context.Background(), unstructuredObj.GetName(), metav1.GetOptions{})
		if err != nil {
			_, err = exector.DynamicClient.Resource(gvr).Namespace(objNamespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("create resources failed: %s", err.Error())
			}
		}

		force := true

		_, err = exector.DynamicClient.Resource(gvr).Namespace(objNamespace).Patch(context.Background(), unstructuredObj.GetName(), types.ApplyPatchType, unstructuredContent, metav1.PatchOptions{
			FieldManager: unstructuredObj.GetName(),
			Force:        &force,
		})
		if err != nil {
			return fmt.Errorf("apply resource failed: %s", err.Error())
		}
	}
	return nil
}

func (exector *ApplyExector) Delete() error {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer([]byte(exector.Content)), 4096)

	var err error
	for {
		var rawObj runtime.RawExtension
		err = decoder.Decode(&rawObj)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("decode failed: %s", err.Error())
		}

		obj, _, err := yamlserializer.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return fmt.Errorf("new yaml decoding serializer failed: %s", err.Error())
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return fmt.Errorf("convert to unstructured object faild: %s", err.Error())
		}
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, err := exector.GetGVR(unstructuredObj.GroupVersionKind())
		if err != nil {
			return fmt.Errorf("get gvr failed: %s", err.Error())
		}

		objNamespace := exector.Namespace
		if unstructuredObj.GetNamespace() != "" {
			objNamespace = unstructuredObj.GetNamespace()
		}

		_, err = exector.DynamicClient.Resource(gvr).Namespace(objNamespace).Get(context.Background(), unstructuredObj.GetName(), metav1.GetOptions{})
		if err != nil {
			_, err = exector.DynamicClient.Resource(gvr).Namespace(objNamespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("create resources failed: %s", err.Error())
			}
		}

		err = exector.DynamicClient.Resource(gvr).Namespace(objNamespace).Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("delete resource falied: %s", err.Error())
		}
	}
	return nil
}
