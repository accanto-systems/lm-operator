package alm

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/go-logr/logr"
	comv1alpha1 "github.com/orgs/accanto-systems/lm-operator/pkg/apis/com/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type nimrodImportConfig struct {
	CassandraReplicatorFactor int32
}

func nimrodConfig(cassandraReplicatorFactor int32) (string, error) {
	nimrodConfigImport := nimrodImportConfig{
		CassandraReplicatorFactor: cassandraReplicatorFactor,
	}

	t, err := template.New("nimrodConfig").Parse("alm:\n" +
		"  nimrod:\n" +
		"    cassandra:\n" +
		"      keyspaceManager:\n" +
		"        replicationFactor: {{.CassandraReplicatorFactor}}\n")

	if err != nil {
		return "", err
	}
	var nimrodTpl bytes.Buffer
	if err := t.Execute(&nimrodTpl, nimrodConfigImport); err != nil {
		return "", err
	}

	return nimrodTpl.String(), nil
}

func (r *ReconcileALM) installNimrod(cr *comv1alpha1.ALM, service nimrodServiceDeploymentInfo, reqLogger logr.Logger) (reconcile.Result, error) {
	cmName := fmt.Sprintf("%s-%s-cm", cr.Name, service.serviceName)

	// Check if this Deployment already exists
	// deploymentName := fmt.Sprintf("%s-%s", cr.Name, service.serviceName)
	deploymentName := service.serviceName
	found, err := r.deploymentExists(cr, service.serviceDeploymentInfo)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !found {
		// Check if this CM already exists
		foundCm := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: cmName, Namespace: cr.Namespace}, foundCm)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Creating a new %s ConfigMap", service.serviceName), "Namespace", cr.Namespace, "Name", cmName)

			data := make(map[string]string)
			data["eureka_instance_ipAddress"] = service.imageName
			data["spring_profiles_include"] = "prod,kubernetes"
			data["spring_cloud_config_failFast"] = "true"
			data["LOG_FOLDER"] = "/var/lm/logs"
			data["spring_cloud_config_label"] = cr.Spec.SpringCloudConfigLabel
			data["JVM_OPTIONS"] = fmt.Sprintf("-Xmx%s", service.heap)
			if cr.Spec.SpringProfilesActive != "" {
				if cr.Spec.Secure {
					data["spring_profiles_active"] = "security," + cr.Spec.SpringProfilesActive
				} else {
					data["spring_profiles_active"] = "nosecurity," + cr.Spec.SpringProfilesActive
				}
			} else {
				if cr.Spec.Secure {
					data["spring_profiles_active"] = "security"
				} else {
					data["spring_profiles_active"] = "nosecurity"
				}
			}

			cm := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: cr.Namespace,
					Name:      cmName,
				},
				Data: data,
			}

			if err := controllerutil.SetControllerReference(cr, cm, r.scheme); err != nil {
				return reconcile.Result{}, err
			}

			err = r.client.Create(context.TODO(), cm)
			if err != nil {
				reqLogger.Info(fmt.Sprintf("Failed to create a new %s ConfigMap", service.serviceName), "Namespace", cr.Namespace, "Name", cmName, "Error", err)
				return reconcile.Result{}, err
			}

			reqLogger.Info(fmt.Sprintf("Created a new %s ConfigMap", service.serviceName), "Namespace", cr.Namespace, "Name", cmName)
			// Conductor CM created successfully - drop through and don't requeue
		} else if err != nil {
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Creating a new %s Deployment", service.serviceName), "Namespace", cr.Namespace, "Name", deploymentName)

		var volumeMounts []corev1.VolumeMount
		var volumes []corev1.Volume

		if cr.Spec.Secure {
			volumes = append(volumes,
				corev1.Volume{
					Name: "lm-certs",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "lm-certs",
						},
					},
				},
				corev1.Volume{
					Name: "lm-keystore",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "lm-keystore",
						},
					},
				})

			volumeMounts = append(volumeMounts,
				corev1.VolumeMount{
					Name:      "lm-certs",
					MountPath: "/var/lm/certs",
				}, corev1.VolumeMount{
					Name:      "lm-keystore",
					MountPath: "/var/lm/keystore",
				})
		}

		if service.themesConfigMap != "" {
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      "themes",
				MountPath: "/var/lm/themes",
			}, corev1.VolumeMount{
				Name:      "themesbinary",
				MountPath: "/var/lm/themesbinary",
			})

			volumes = append(volumes,
				corev1.Volume{
					Name: "themes",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				corev1.Volume{
					Name: "themesbinary",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: service.themesConfigMap},
						},
					},
				})
		}

		if service.localesConfigMap != "" {
			volumeMounts = append(volumeMounts,
				corev1.VolumeMount{
					Name:      "locales",
					MountPath: "/var/lm/locales",
				},
				corev1.VolumeMount{
					Name:      "localesbinary",
					MountPath: "/var/lm/localesbinary",
				})

			volumes = append(volumes,
				corev1.Volume{
					Name: "locales",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				corev1.Volume{
					Name: "localesbinary",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{Name: service.localesConfigMap},
						},
					},
				})
		}

		deployment := buildDeployment(cr.Namespace, deploymentName, cr, service.serviceDeploymentInfo,
			volumeMounts, volumes)

		if err := controllerutil.SetControllerReference(cr, deployment, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Failed to create a new %s Deployment", service.serviceName), "Namespace", cr.Namespace, "Name", deploymentName, "Error", err)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Deployment", service.serviceName), "Namespace", cr.Namespace, "Name", deploymentName)

		ingressName := "nimrod-ingress"
		ingress := buildIngress(cr.Spec.Secure, cr.Namespace, ingressName, "ui.lm", 8290, "nimrod", "nimrod-tls")

		if err := controllerutil.SetControllerReference(cr, ingress, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), ingress)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Failed to create a new %s Ingress", service.serviceName), "Namespace", cr.Namespace, "Name", ingressName, "Error", err)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Ingress", ingressName), "Namespace", cr.Namespace, "Name", ingressName)

		// Created successfully - don't requeue and create service
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return r.service(cr, service.serviceDeploymentInfo, reqLogger)
}
