package alm

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	comv1alpha1 "github.com/orgs/accanto-systems/lm-operator/pkg/apis/com/v1alpha1"
	v1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileALM) service(cr *comv1alpha1.ALM, serviceDeploymentInfo serviceDeploymentInfo, reqLogger logr.Logger) (reconcile.Result, error) {
	reqLogger.Info(fmt.Sprintf("Creating service %s on port %d", serviceDeploymentInfo.serviceName, serviceDeploymentInfo.port), "Namespace", cr.Namespace, "Name", serviceDeploymentInfo.serviceName)

	var port corev1.ServicePort
	if serviceDeploymentInfo.nodePort <= 0 {
		port = corev1.ServicePort{
			Name:       "http",
			Protocol:   corev1.ProtocolTCP,
			Port:       serviceDeploymentInfo.port,
			TargetPort: intstr.FromInt(serviceDeploymentInfo.targetPort),
		}
	} else if serviceDeploymentInfo.nodePort > 0 {
		port = corev1.ServicePort{
			Name:       "http",
			Protocol:   corev1.ProtocolTCP,
			Port:       serviceDeploymentInfo.port,
			TargetPort: intstr.FromInt(serviceDeploymentInfo.targetPort),
			NodePort:   serviceDeploymentInfo.nodePort,
		}
	} else {
		// nothing to do
	}

	// serviceName := fmt.Sprintf("%s-%s-service", cr.Name, serviceDeploymentInfo.serviceName)
	serviceName := serviceDeploymentInfo.serviceName

	// Check if this Service already exists
	foundService := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: serviceName, Namespace: cr.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info(fmt.Sprintf("Creating a new %s Service", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", serviceName)

		service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: cr.Namespace,
				Name:      serviceName,
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{port},
				Selector: map[string]string{
					"app": serviceDeploymentInfo.serviceName,
				},
				Type: "NodePort",
			},
		}

		if err := controllerutil.SetControllerReference(cr, service, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Failed to create a new %s Service", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", serviceName, "Error", err)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Service", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", serviceName)
		// Conductor CM created successfully - drop through and don't requeue
		// return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileALM) statefulsetExists(cr *comv1alpha1.ALM, service serviceDeploymentInfo) (bool, error) {
	namespace := cr.Namespace
	// name := fmt.Sprintf("%s-%s", cr.Name, service.serviceName)
	name := service.serviceName

	found := &v1beta1.StatefulSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (r *ReconcileALM) deploymentExists(cr *comv1alpha1.ALM, service serviceDeploymentInfo) (bool, error) {
	return r.deploymentByNameExists(cr, service.serviceName)
}

func (r *ReconcileALM) deploymentByNameExists(cr *comv1alpha1.ALM, name string) (bool, error) {
	namespace := cr.Namespace

	found := &extv1beta1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func buildDeployment(namespace string, statefulsetName string, cr *comv1alpha1.ALM, service serviceDeploymentInfo,
	volumeMounts []corev1.VolumeMount, volumes []corev1.Volume) *extv1beta1.Deployment {
	dockerImage := fmt.Sprintf("%s/%s:%s", cr.Spec.DockerRepo, service.imageName, service.imageVersion)
	// deploymentName := fmt.Sprintf("%s-%s", cr.Name, service.serviceName)
	deploymentName := service.serviceName

	return &extv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": service.serviceName,
			},
		},
		Spec: extv1beta1.DeploymentSpec{
			Replicas: int32Ptr(service.numReplicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": service.serviceName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%s", cr.Name, service.serviceName),
					Namespace: cr.Namespace,
					Labels: map[string]string{
						"app": service.serviceName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  service.serviceName,
							Image: dockerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: service.port,
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{Name: fmt.Sprintf("%s-%s-cm", cr.Name, service.serviceName)},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "INSTANCE_ID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "SERVER_ID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse(service.cpuRequests),
									"memory": resource.MustParse(service.memoryRequests),
								},
							},
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
}

func buildStatefulset(namespace string, statefulsetName string, cr *comv1alpha1.ALM, service serviceDeploymentInfo,
	volumeMounts []corev1.VolumeMount, volumes []corev1.Volume, additionalEnv []corev1.EnvVar) *v1beta1.StatefulSet {
	dockerImage := fmt.Sprintf("%s/%s:%s", cr.Spec.DockerRepo, service.imageName, service.imageVersion)

	var env []corev1.EnvVar
	env = append(env,
		corev1.EnvVar{
			Name: "INSTANCE_ID",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		corev1.EnvVar{
			Name: "SERVER_ID",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		})
	env = append(env, additionalEnv...)

	return &v1beta1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulsetName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": service.serviceName,
			},
		},
		Spec: v1beta1.StatefulSetSpec{
			Replicas: int32Ptr(service.numReplicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": service.serviceName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%s", cr.Name, service.serviceName),
					Namespace: namespace,
					Labels: map[string]string{
						"app": service.serviceName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  service.serviceName,
							Image: dockerImage,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: service.port,
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{Name: fmt.Sprintf("%s-%s-cm", cr.Name, service.serviceName)},
									},
								},
							},
							Env: env,
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									"cpu":    resource.MustParse(service.cpuRequests),
									"memory": resource.MustParse(service.memoryRequests),
								},
							},
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
}

func buildIngress(secure bool, namespace string, name string, ingressHost string, port int, serviceName string, externalCertSecretName string) *extv1beta1.Ingress {
	annotations := make(map[string]string)
	annotations["ingress.kubernetes.io/rewrite-target"] = "/"
	annotations["nginx.org/websocket-services"] = serviceName
	if secure {
		annotations["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
		annotations["ingress.kubernetes.io/secure-backends"] = "true"
	}
	return &extv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      map[string]string{},
			Annotations: annotations,
		},
		Spec: extv1beta1.IngressSpec{
			Rules: []extv1beta1.IngressRule{
				{
					Host: ingressHost,
					IngressRuleValue: extv1beta1.IngressRuleValue{
						HTTP: &extv1beta1.HTTPIngressRuleValue{
							Paths: []extv1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: extv1beta1.IngressBackend{
										ServiceName: serviceName,
										ServicePort: intstr.FromInt(port),
									},
								},
							},
						},
					},
				},
			},
			TLS: []extv1beta1.IngressTLS{
				{
					Hosts:      []string{ingressHost},
					SecretName: externalCertSecretName,
				},
			},
		},
	}
}
