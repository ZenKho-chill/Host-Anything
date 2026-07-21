// Copyright 2026 Host Anything Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/host-anything/hostanything/pkg/types"
)

// Adapter implements types.RuntimeAdapter for Kubernetes.
type Adapter struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewAdapter creates a new Kubernetes Adapter.
func NewAdapter() (*Adapter, error) {
	// Try in-cluster config first, then fallback to ~/.kube/config
	// For simplicity in M5, we just use the default fallback mechanism.
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// Attempt in-cluster config if out-of-cluster fails
		// rest.InClusterConfig() can be used here, but we will omit it for brevity
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return &Adapter{
		clientset: clientset,
		namespace: "default", // Using default namespace for M5
	}, nil
}

func (a *Adapter) Deploy(ctx context.Context, spec *types.ServiceSpec) error {
	name := "ha-" + spec.ServiceID
	labels := map[string]string{
		"app":                     name,
		"sh.hostanything.managed": "true",
	}

	// 1. Prepare Environment Variables
	var envVars []corev1.EnvVar
	for k, v := range spec.ResolvedEnv {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	// 2. Prepare Ports
	var containerPorts []corev1.ContainerPort
	var servicePorts []corev1.ServicePort
	hasExternalPort := false

	for i, netCfg := range spec.Template.Network {
		proto := corev1.ProtocolTCP
		if netCfg.Protocol == "udp" {
			proto = corev1.ProtocolUDP
		}

		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: int32(netCfg.InternalPort),
			Protocol:      proto,
		})

		svcPort := corev1.ServicePort{
			Name:       fmt.Sprintf("port-%d", i),
			Port:       int32(netCfg.InternalPort),
			TargetPort: intstr.FromInt(netCfg.InternalPort),
			Protocol:   proto,
		}

		if netCfg.ExternalPort > 0 {
			svcPort.NodePort = int32(netCfg.ExternalPort)
			hasExternalPort = true
		}
		servicePorts = append(servicePorts, svcPort)
	}

	// 3. Prepare Volumes
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount
	for _, vol := range spec.Template.Volumes {
		vName := "vol-" + vol.Name
		volumes = append(volumes, corev1.Volume{
			Name: vName,
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/tmp/ha-data/" + spec.ServiceID + "/" + vol.Name, // Simplistic hostPath mapping for M5
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      vName,
			MountPath: vol.MountPath,
		})
	}

	// 4. Create Deployment
	replicas := int32(1)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:         name,
							Image:        spec.Template.Runtime.Image,
							Command:      spec.Template.Runtime.Command,
							Env:          envVars,
							Ports:        containerPorts,
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}

	// Apply Deployment (Create or Update)
	deploymentsClient := a.clientset.AppsV1().Deployments(a.namespace)
	_, err := deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		// If it exists, update it
		_, err = deploymentsClient.Update(ctx, deployment, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("kubernetes: failed to create/update deployment: %w", err)
		}
	}

	// 5. Create Service
	if len(servicePorts) > 0 {
		svcType := corev1.ServiceTypeClusterIP
		if hasExternalPort {
			svcType = corev1.ServiceTypeNodePort
		}

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name,
				Labels: labels,
			},
			Spec: corev1.ServiceSpec{
				Selector: labels,
				Ports:    servicePorts,
				Type:     svcType,
			},
		}

		servicesClient := a.clientset.CoreV1().Services(a.namespace)
		_, err := servicesClient.Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			// Update if exists. Services need ResourceVersion to update,
			// so typically we get it first, but for M5 we'll just ignore if it exists.
		}
	}

	return nil
}

func (a *Adapter) Stop(ctx context.Context, serviceID string) error {
	// Stopping in k8s means scaling to 0
	name := "ha-" + serviceID
	scale, err := a.clientset.AppsV1().Deployments(a.namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes: get scale: %w", err)
	}

	scale.Spec.Replicas = 0
	_, err = a.clientset.AppsV1().Deployments(a.namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes: update scale: %w", err)
	}
	return nil
}

func (a *Adapter) Start(ctx context.Context, serviceID string) error {
	// Starting in k8s means scaling to 1
	name := "ha-" + serviceID
	scale, err := a.clientset.AppsV1().Deployments(a.namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes: get scale: %w", err)
	}

	scale.Spec.Replicas = 1
	_, err = a.clientset.AppsV1().Deployments(a.namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes: update scale: %w", err)
	}
	return nil
}

func (a *Adapter) Remove(ctx context.Context, serviceID string) error {
	name := "ha-" + serviceID

	// Delete Service
	_ = a.clientset.CoreV1().Services(a.namespace).Delete(ctx, name, metav1.DeleteOptions{})

	// Delete Deployment
	err := a.clientset.AppsV1().Deployments(a.namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("kubernetes: delete deployment: %w", err)
	}

	return nil
}

func (a *Adapter) Status(ctx context.Context, serviceID string) (*types.ServiceStatus, error) {
	name := "ha-" + serviceID

	deploy, err := a.clientset.AppsV1().Deployments(a.namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("kubernetes: get deployment: %w", err)
	}

	state := types.ServiceStateStopped
	if deploy.Status.AvailableReplicas > 0 {
		state = types.ServiceStateRunning
	} else if deploy.Status.UnavailableReplicas > 0 {
		state = types.ServiceStateError
	}

	// Fetch port mappings from service
	mappings := make(map[int]int)
	svc, err := a.clientset.CoreV1().Services(a.namespace).Get(ctx, name, metav1.GetOptions{})
	if err == nil {
		for _, port := range svc.Spec.Ports {
			if port.NodePort > 0 {
				mappings[int(port.Port)] = int(port.NodePort)
			}
		}
	}

	return &types.ServiceStatus{
		State:        state,
		PortMappings: mappings,
		RuntimeID:    string(deploy.UID),
	}, nil
}

func (a *Adapter) Logs(ctx context.Context, serviceID string) (io.ReadCloser, error) {
	name := "ha-" + serviceID

	// Find the pod for this deployment
	pods, err := a.clientset.CoreV1().Pods(a.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=" + name,
	})
	if err != nil || len(pods.Items) == 0 {
		return nil, fmt.Errorf("kubernetes: no pods found for logs")
	}

	podName := pods.Items[0].Name
	req := a.clientset.CoreV1().Pods(a.namespace).GetLogs(podName, &corev1.PodLogOptions{
		Follow: true,
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("kubernetes: failed to stream logs: %w", err)
	}

	return stream, nil
}
