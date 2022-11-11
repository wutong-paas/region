package kube

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	//"k8s.io/apimachinery/pkg/types"
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	controllerruntime "sigs.k8s.io/controller-runtime"
	sigyaml "sigs.k8s.io/yaml"
)

var s = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpbin-test
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: httpbin
      version: v1
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: httpbin
        version: v1
    spec:
      containers:
      - image: docker.io/kennethreitz/httpbin
        imagePullPolicy: IfNotPresent
        name: httpbin
        ports:
        - containerPort: 80
          protocol: TCP
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: httpbin
      serviceAccountName: httpbin
      terminationGracePeriodSeconds: 30
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: httpbin
    service: httpbin
  name: httpbin-test
  namespace: default
spec:
  ports:
  - name: http
    port: 8000
    protocol: TCP
    targetPort: 80
  selector:
    app: httpbin
  sessionAffinity: None
  type: ClusterIP
`

type ExecuteYaml struct {
	applyYaml string
	namespace string
}

func NewYaml(applyYaml, namespace string) *ExecuteYaml {
	return &ExecuteYaml{
		applyYaml: applyYaml,
		namespace: namespace,
	}
}

func (y *ExecuteYaml) GtGVR(gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {

	config := controllerruntime.GetConfigOrDie()

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)

	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return mapping.Resource, nil
}

func (y *ExecuteYaml) UpdateFromYaml() error {
	config := controllerruntime.GetConfigOrDie()
	dynameicclient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	d := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(y.applyYaml), 4096)

	for {
		var rawObj runtime.RawExtension
		err = d.Decode(&rawObj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decode is err %v", err)
		}

		obj, _, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return fmt.Errorf("rawobj is err%v", err)
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return fmt.Errorf("tounstructured is err %v", err)
		}

		unstructureObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, err := y.GtGVR(unstructureObj.GroupVersionKind())
		if err != nil {
			return err
		}
		unstructuredYaml, err := sigyaml.Marshal(unstructureObj)
		if err != nil {
			return fmt.Errorf("unable to marshal resource as yaml: %w", err)
		}
		_, getErr := dynameicclient.Resource(gvr).Namespace(y.namespace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		if getErr != nil {
			_, createErr := dynameicclient.Resource(gvr).Namespace(y.namespace).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		}

		force := true
		if y.namespace == unstructureObj.GetNamespace() {

			_, err = dynameicclient.Resource(gvr).
				Namespace(y.namespace).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						FieldManager: unstructureObj.GetName(),
						Force:        &force,
					})

			if err != nil {
				return fmt.Errorf("unable to patch resource: %w", err)
			}

		} else {

			_, err = dynameicclient.Resource(gvr).
				Patch(context.Background(),
					unstructureObj.GetName(),
					types.ApplyPatchType,
					unstructuredYaml, metav1.PatchOptions{
						Force:        &force,
						FieldManager: unstructureObj.GetName(),
					})
			if err != nil {
				return fmt.Errorf("ns is nil unable to patch resource: %w", err)
			}

		}

	}
	return nil

}

func TestApplyResources(t *testing.T) {
	//fmt.Println(s)
	ey := NewYaml(s, "default")
	err := ey.UpdateFromYaml()
	if err != nil {
		t.Error(err)
	}
	t.Log("done")
}
