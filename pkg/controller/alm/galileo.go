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

type galileoImportConfig struct {
	ReplicatorFactor    int32
	NumStandbyReplicas  int32
	NumShards           int32
	NumReplicas         int32
	ProcessingGuarantee string
}

// processingGuarantee at_least_once
// numShards 1
// numReplicas 0
// numStandbyReplicas 1
// replicationFactor 1
func galileoConfig(processingGuarantee string, numShards int32, numReplicas int32, numStandbyReplicas int32, replicatorFactor int32) (string, error) {
	galileoConfigImport := galileoImportConfig{
		ReplicatorFactor:    replicatorFactor,
		NumStandbyReplicas:  numStandbyReplicas,
		NumShards:           numShards,
		NumReplicas:         numReplicas,
		ProcessingGuarantee: processingGuarantee,
	}

	t, err := template.New("galileoConfig").Parse("alm:\n" +
		"  galileo:\n" +
		"    ldu:\n" +
		"      streams:\n" +
		"        processing.guarantee: {{.ProcessingGuarantee}}\n" +
		"        replication.factor: {{.ReplicatorFactor}}\n" +
		"        num.standby.replicas: {{.NumStandbyReplicas}}\n" +
		"    janus:\n" +
		"      cluster.max-partitions: 4\n" +
		"      storage.cql.replication-factor: {{.ReplicatorFactor}}\n" +
		"      index:\n" +
		"        search:\n" +
		"          elasticsearch.create.ext.index.number_of_replicas: {{.NumStandbyReplicas}}\n" +
		"          elasticsearch.create.ext.index.number_of_shards: {{.NumShards}}\n")

	if err != nil {
		return "", err
	}
	var galileoTpl bytes.Buffer
	if err := t.Execute(&galileoTpl, galileoConfigImport); err != nil {
		return "", err
	}

	return galileoTpl.String(), nil
}

func (r *ReconcileALM) installGalileo(cr *comv1alpha1.ALM, service serviceDeploymentInfo, reqLogger logr.Logger) (reconcile.Result, error) {
	// dockerImage := fmt.Sprintf("%s/%s:%s", cr.Spec.DockerRepo, serviceDeploymentInfo.imageName, serviceDeploymentInfo.imageVersion)
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

		reqLogger.Info(fmt.Sprintf("Creating a new %s Statefulset", service.serviceName), "Namespace", cr.Namespace, "Name", statefulsetName)

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
				})

			volumeMounts = append(volumeMounts,
				corev1.VolumeMount{
					Name:      "lm-certs",
					MountPath: "/var/lm/certs",
				})
		}

		statefulset := buildStatefulset(cr.Namespace, statefulsetName, cr, service, volumeMounts, volumes, []corev1.EnvVar{})

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
