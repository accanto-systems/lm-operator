package alm

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	comv1alpha1 "github.com/orgs/accanto-systems/lm-operator/pkg/apis/com/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileALM) installConductor(cr *comv1alpha1.ALM, service serviceDeploymentInfo, reqLogger logr.Logger) (reconcile.Result, error) {
	cmName := fmt.Sprintf("%s-%s-cm", cr.Name, service.serviceName)

	// Check if this Statefulset already exists
	// statefulsetName := fmt.Sprintf("%s-%s", cr.Name, service.serviceName)
	statefulsetName := service.serviceName
	found, err := r.statefulsetExists(cr, service)
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
					data["spring_profiles_active"] = "security,vault," + cr.Spec.SpringProfilesActive
				} else {
					data["spring_profiles_active"] = "nosecurity,vault," + cr.Spec.SpringProfilesActive
				}
			} else {
				if cr.Spec.Secure {
					data["spring_profiles_active"] = "security,vault"
				} else {
					data["spring_profiles_active"] = "nosecurity,vault"
				}
			}

			// data["spring_cloud_vault_host"] = {{ .Values.app.config.configServer.vault.host }}
			// data["spring_cloud_config_server_vault_host"] = {{ .Values.app.config.configServer.vault.host }}

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

		reqLogger.Info(fmt.Sprintf("Creating a new %s Statefulset", service.serviceName), "Namespace", cr.Namespace, "Name", statefulsetName)

		var volumes []corev1.Volume
		volumes = append(volumes,
			corev1.Volume{
				Name: "vault-cert",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "vault-cert",
					},
				},
			})

		var volumeMounts []corev1.VolumeMount
		volumeMounts = append(volumeMounts,
			corev1.VolumeMount{
				Name:      "vault-cert",
				MountPath: "/var/lm/vault/certs",
			})

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
				},
				corev1.VolumeMount{
					Name:      "lm-keystore",
					MountPath: "/var/lm/keystore",
				})
		}

		statefulset := buildStatefulset(cr.Namespace, statefulsetName, cr, service, volumeMounts, volumes,
			[]corev1.EnvVar{
				corev1.EnvVar{
					Name:  "eureka_instance_hostname",
					Value: "${HOSTNAME}.conductor",
				},
				corev1.EnvVar{
					Name:  "numReplicas",
					Value: strconv.Itoa(int(service.numReplicas)),
				},
				corev1.EnvVar{
					Name:  "secure",
					Value: strconv.FormatBool(cr.Spec.Secure),
				},
				corev1.EnvVar{
					Name: "SPRING_CLOUD_CONFIG_SERVER_VAULT_TOKEN",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "vault-token",
							},
							Key: "lmToken",
						},
					},
				},
				corev1.EnvVar{
					Name: "SPRING_CLOUD_VAULT_TOKEN",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "vault-token",
							},
							Key: "lmToken",
						},
					},
				},
			})

		if err := controllerutil.SetControllerReference(cr, statefulset, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), statefulset)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Failed to create a new %s Statefulset", service.serviceName), "Namespace", cr.Namespace, "Name", statefulsetName, "Error", err)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Statefulset", service.serviceName), "Namespace", cr.Namespace, "Name", statefulsetName)

		// Created successfully - don't requeue and create service
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return r.service(cr, service, reqLogger)
}
