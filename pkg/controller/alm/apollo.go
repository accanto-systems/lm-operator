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

type apolloImportConfig struct {
	CassandraReplicatorFactor int32
	ESNumShards               int32
	ESNumReplicas             int32
}

func apolloConfig(esNumShards int32, esNumReplicas int32, cassandraReplicatorFactor int32) (string, error) {
	apolloConfigImport := apolloImportConfig{
		CassandraReplicatorFactor: cassandraReplicatorFactor,
		ESNumShards:               esNumShards,
		ESNumReplicas:             esNumReplicas,
	}

	t, err := template.New("apolloConfig").Parse("alm:\n" +
		"  apollo:\n" +
		"    cassandra:\n" +
		"      keyspaceManager:\n" +
		"        replicationFactor: {{.CassandraReplicatorFactor}}\n" +
		"    janus:\n" +
		"      cluster.max-partitions: 4\n" +
		"      storage.cql.replication-factor: {{.CassandraReplicatorFactor}}\n" +
		"      index:\n" +
		"        search:\n" +
		"          elasticsearch.create.ext.index.number_of_replicas: {{.ESNumReplicas}}\n" +
		"          elasticsearch.create.ext.index.number_of_shards: {{.ESNumShards}}\n")
	if err != nil {
		return "", err
	}
	var apolloTpl bytes.Buffer
	if err := t.Execute(&apolloTpl, apolloConfigImport); err != nil {
		return "", err
	}

	return apolloTpl.String(), nil
}

func (r *ReconcileALM) installApollo(cr *comv1alpha1.ALM, service serviceDeploymentInfo, reqLogger logr.Logger) (reconcile.Result, error) {
	cmName := fmt.Sprintf("%s-%s-cm", cr.Name, service.serviceName)

	// Check if this Deployment already exists
	// deploymentName := fmt.Sprintf("%s-%s", cr.Name, service.serviceName)
	deploymentName := service.serviceName
	found, err := r.deploymentExists(cr, service)
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

		var volumes []corev1.Volume
		var volumeMounts []corev1.VolumeMount

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

		deployment := buildDeployment(cr.Namespace, deploymentName, cr, service,
			[]corev1.VolumeMount{
				{
					Name:      "lm-certs",
					MountPath: "/var/lm/certs",
				},
			},
			[]corev1.Volume{
				{
					Name: "lm-certs",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "lm-certs",
						},
					},
				},
			})

		if err := controllerutil.SetControllerReference(cr, deployment, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Failed to create a new %s Deployment", service.serviceName), "Namespace", cr.Namespace, "Name", deploymentName, "Error", err)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Deployment", service.serviceName), "Namespace", cr.Namespace, "Name", deploymentName)

		// Created successfully - don't requeue and create service
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return r.service(cr, service, reqLogger)
}
