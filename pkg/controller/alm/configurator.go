package alm

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"

	comv1alpha1 "github.com/orgs/accanto-systems/lm-operator/pkg/apis/com/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type janus struct {
	CassandraHostname string
	ESHostname        string
}

type SecurityConfig struct {
	Enabled                            string
	SecurityKeyStorePassword           string
	SecurityNimrodClientSecret         string
	SecurityNimrodAccessTokenValidity  string
	SecurityNimrodRefreshTokenValidity string
	SecurityDokiClientSecret           string
	SecurityDokiAccessTokenValidity    string
	SecurityDokiRoles                  string
	SecurityLdapEnabled                string
	SecurityLdapConfigPassword         string
	SecurityLdapManagerPassword        string
	SecurityLdapDomain                 string
	SecurityUIHost                     string
	SecurityAPIHost                    string
	SecurityBrentHost                  string
	CassandraUsername                  string
	CassandraPassword                  string
	SecurityClientCredentialsConfig    string
}

type configuratorConfig struct {
	KafkaConfig      string
	TopicsConfig     string
	JanusgraphConfig string
	SecurityConfig   SecurityConfig
}

func buildConfiguratorCM(name string, cr *comv1alpha1.ALM, configuratorDeploymentInfo configuratorDeploymentInfo) (*corev1.ConfigMap, error) {
	janus := janus{
		ESHostname:        "foundation-elasticsearch-client:9200",
		CassandraHostname: "foundation-cassandra",
	}

	t, err := template.New("janus").Parse("alm:\n" +
		"  janus:\n" +
		"    storage:\n" +
		"      hostname: \"{{.CassandraHostname}}\"\n" +
		"    cluster:\n" +
		"      max-partitions: 4\n" +
		"    index:\n" +
		"      search:\n" +
		"        hostname: \"{{.ESHostname}}\"\n" +
		"        elasticsearch.create.ext.index.number_of_shards: 1\n" +
		"        elasticsearch.create.ext.index.number_of_replicas: 0\n")
	if err != nil {
		return nil, err
	}

	var janusTpl bytes.Buffer
	if err := t.Execute(&janusTpl, janus); err != nil {
		return nil, err
	}

	type kafkaConfig struct {
		KafkaSource  string
		KafkaVersion string
		ZookeeperURL string
	}

	kafka := kafkaConfig{
		KafkaSource:  "https://archive.apache.org/dist/kafka/2.0.0/kafka_2.11-2.0.0.tgz",
		KafkaVersion: "kafka_2.11-2.0.0",
		ZookeeperURL: "foundation-zookeeper:2181",
	}

	t, err = template.New("kafka").Parse("kafka_source: \"{{.KafkaSource}}\"\n" +
		"kafka_version: \"{{.KafkaVersion}}\"\n" +
		"zookeeper_url: \"{{.ZookeeperURL}}\"\n")

	if err != nil {
		return nil, err
	}

	var kafkaTpl bytes.Buffer
	if err := t.Execute(&kafkaTpl, kafka); err != nil {
		return nil, err
	}

	type topicsConfig struct {
		NumPartitions     int32
		ReplicationFactor int32
	}

	topics := topicsConfig{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	t, err = template.New("topics").Parse("topics:\n" +
		"  alm__health:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"cleanup.policy=compact\"\n" +
		"  alm__metric:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"cleanup.policy=compact\"\n" +
		"  alm__metric-integrity:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"cleanup.policy=compact\"\n" +
		"  alm__policy:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"cleanup.policy=compact\"\n" +
		"  alm__policyAction:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__policyStatusHeal:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__policyStatusScale:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: 1{{.NumPartitions}}\n" +
		"  alm__clock:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__descriptorChange:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"cleanup.policy=compact\"\n" +
		"  alm__processStateChange:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__processTasksStateChange:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__stateChange:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__serviceStateTransition:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__taskUpdate:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__load:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__integrity:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__integrityMissing:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  info:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__clockticks5:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"retention.ms=3600000\"\n" +
		"  alm__clockticks10:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"retention.ms=3600000\"\n" +
		"  alm__clockticks15:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"retention.ms=3600000\"\n" +
		"  alm__clockticks30:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"retention.ms=3600000\"\n" +
		"  alm__clockticks60:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"retention.ms=3600000\"\n" +
		"  alm__clockticksother:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__tick:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"retention.ms=3600000\"\n" +
		"  alm__clocktickTypes:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"cleanup.policy=compact\"\n" +
		"  alm__processRestart:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  alm__stateChange__ldu:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"    config: \"retention.ms=31536000000\"\n" +
		"  lm_vim_infrastructure_task_events:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n" +
		"  lm_vnfc_lifecycle_execution_events:\n" +
		"    replication_factor: {{.ReplicationFactor}}\n" +
		"    partitions: {{.NumPartitions}}\n")

	if err != nil {
		return nil, err
	}

	var topicsTpl bytes.Buffer
	if err := t.Execute(&topicsTpl, topics); err != nil {
		return nil, err
	}

	type clientCredentialsConfig struct {
		LMClientID         string
		LMClientSecret     string
		LMGrantTypes       string
		LMRoles            string
		NimrodClientID     string
		NimrodClientSecret string
		NimrodGrantTypes   string
		DokiClientID       string
		DokiClientSecret   string
		DokiGrantTypes     string
		DokiRoles          string
	}

	clientCredentials := clientCredentialsConfig{
		LMClientID:         "LmClient",
		LMClientSecret:     "pass123",
		LMGrantTypes:       "client_credentials",
		LMRoles:            "SLMAdmin",
		NimrodClientID:     "NimrodClient",
		NimrodClientSecret: "pass123",
		NimrodGrantTypes:   "password,refresh_token",
		DokiClientID:       "DokiClient",
		DokiClientSecret:   "pass123",
		DokiGrantTypes:     "client_credentials",
		DokiRoles:          "BehaviourScenarioExecute",
	}

	t, err = template.New("clientCredentials").Parse("- clientId: {{.LMClientID}}\n" +
		"  clientSecret: {{.LMClientSecret}}\n" +
		"  grantTypes: {{.LMGrantTypes}}\n" +
		"  roles: {{.LMRoles}}\n" +
		"- clientId: {{.NimrodClientID}}\n" +
		"  clientSecret: {{.NimrodClientSecret}}\n" +
		"  grantTypes: {{.NimrodGrantTypes}}\n" +
		"- clientId: {{.DokiClientID}}\n" +
		"  clientSecret: {{.DokiClientSecret}}\n" +
		"  grantTypes: {{.DokiGrantTypes}}\n" +
		"  roles: {{.DokiRoles}}\n")
	if err != nil {
		return nil, err
	}

	var clientCredentialsTpl bytes.Buffer
	if err := t.Execute(&clientCredentialsTpl, clientCredentials); err != nil {
		return nil, err
	}

	config := configuratorConfig{
		KafkaConfig:      kafkaTpl.String(),
		TopicsConfig:     topicsTpl.String(),
		JanusgraphConfig: janusTpl.String(),
		SecurityConfig: SecurityConfig{
			Enabled:                            strconv.FormatBool(cr.Spec.Secure),
			SecurityKeyStorePassword:           "keypass",
			SecurityNimrodClientSecret:         "pass123",
			SecurityNimrodAccessTokenValidity:  "1200",  // 20 minutes
			SecurityNimrodRefreshTokenValidity: "30600", // 8.5 hours
			SecurityDokiClientSecret:           "pass123",
			SecurityDokiAccessTokenValidity:    "1200", // 20 minutes
			SecurityDokiRoles:                  "BehaviourScenarioExecute",
			SecurityLdapEnabled:                "true",
			SecurityLdapConfigPassword:         "config",
			SecurityLdapManagerPassword:        "lmadmin",
			SecurityLdapDomain:                 "lm.com",
			SecurityUIHost:                     "ui.lm",
			SecurityAPIHost:                    "app.lm",
			SecurityBrentHost:                  "brent.lm",
			CassandraUsername:                  "",
			CassandraPassword:                  "",
			SecurityClientCredentialsConfig:    clientCredentialsTpl.String(),
		},
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cr.Namespace,
			Name:      name,
		},
		Data: map[string]string{
			"kafka_config.yaml":                  config.KafkaConfig,
			"topics.yaml":                        config.TopicsConfig,
			"janusgraph_config.yaml":             config.JanusgraphConfig,
			"securityEnabled":                    config.SecurityConfig.Enabled,
			"securityClientCredentials":          config.SecurityConfig.SecurityClientCredentialsConfig,
			"securityKeyStorePassword":           config.SecurityConfig.SecurityKeyStorePassword,
			"securityNimrodClientSecret":         config.SecurityConfig.SecurityNimrodClientSecret,
			"securityNimrodAccessTokenValidity":  config.SecurityConfig.SecurityNimrodAccessTokenValidity,
			"securityNimrodRefreshTokenValidity": config.SecurityConfig.SecurityNimrodRefreshTokenValidity,
			"securityDokiClientSecret":           config.SecurityConfig.SecurityDokiClientSecret,
			"securityDokiAccessTokenValidity":    config.SecurityConfig.SecurityDokiAccessTokenValidity,
			"securityDokiRoles":                  config.SecurityConfig.SecurityDokiRoles,
			"securityLdapEnabled":                config.SecurityConfig.SecurityLdapEnabled,
			"securityLdapConfigPassword":         config.SecurityConfig.SecurityLdapConfigPassword,
			"securityLdapManagerPassword":        config.SecurityConfig.SecurityLdapManagerPassword,
			"securityLdapDomain":                 config.SecurityConfig.SecurityLdapDomain,
			"securityUiHost":                     config.SecurityConfig.SecurityUIHost,
			"securityApiHost":                    config.SecurityConfig.SecurityAPIHost,
			"securityBrentHost":                  config.SecurityConfig.SecurityBrentHost,
			"cassandraUsername":                  config.SecurityConfig.CassandraUsername,
			"cassandraPassword":                  config.SecurityConfig.CassandraPassword,
		},
	}, nil
}

func buildLmConfigImportCm(namespace string, name string) (*corev1.ConfigMap, error) {
	watchtowerCfg, watchtowerCfgErr := watchtowerConfig(1)
	if watchtowerCfgErr != nil {
		return nil, watchtowerCfgErr
	}

	galileoCfg, galileoCfgErr := galileoConfig("at_least_once", 1, 0, 1, 1)
	if galileoCfgErr != nil {
		return nil, galileoCfgErr
	}

	talledegaCfg, talledegaCfgErr := talledegaConfig(1, 0, 1)
	if talledegaCfgErr != nil {
		return nil, talledegaCfgErr
	}

	brentCfg, brentCfgErr := brentConfig(1, 0, 1)
	if brentCfgErr != nil {
		return nil, brentCfgErr
	}

	apolloCfg, apolloCfgErr := apolloConfig(1, 0, 1)
	if apolloCfgErr != nil {
		return nil, apolloCfgErr
	}

	nimrodCfg, nimrodCfgErr := nimrodConfig(1)
	if nimrodCfgErr != nil {
		return nil, nimrodCfgErr
	}

	ishtarCfg, ishtarCfgErr := ishtarConfig(1)
	if ishtarCfgErr != nil {
		return nil, ishtarCfgErr
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Data: map[string]string{
			"watchtower.yaml": watchtowerCfg,
			"galileo.yaml":    galileoCfg,
			"talledega.yaml":  talledegaCfg,
			"brent.yaml":      brentCfg,
			"apollo.yaml":     apolloCfg,
			"nimrod.yaml":     nimrodCfg,
			"ishtar.yaml":     ishtarCfg,
		},
	}, nil
}

func buildJob(cr *comv1alpha1.ALM, configuratorDeploymentInfo configuratorDeploymentInfo, dockerRepo string, namespace string, name string,
	lmConfiguratorCMName string, lmConfigImportCmName string) *batchv1.Job {
	dockerImage := fmt.Sprintf("%s/%s:%s", dockerRepo, configuratorDeploymentInfo.imageName, configuratorDeploymentInfo.imageVersion)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "lm-configurator",
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%s", cr.Name, configuratorDeploymentInfo.serviceName),
					Namespace: namespace,
					Labels: map[string]string{
						"app": configuratorDeploymentInfo.serviceName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "lm-configurator",
							Image:           dockerImage,
							ImagePullPolicy: corev1.PullAlways,
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{Name: fmt.Sprintf("%s-%s-cm", cr.Name, configuratorDeploymentInfo.serviceName)},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "VAULT_TOKEN",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "vault-token",
											},
											Key: "lmToken",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "lm-configurator",
									MountPath: "/var/lm-configurator",
								},
								{
									Name:      "lm-config-import",
									MountPath: "/var/config-import",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: "lm-configurator",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: lmConfiguratorCMName},
								},
							},
						},
						{
							Name: "lm-config-import",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: lmConfigImportCmName},
								},
							},
						},
					},
				},
			},
		},
	}
}
