package alm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	resty "github.com/go-resty/resty/v2"
	comv1alpha1 "github.com/orgs/accanto-systems/lm-operator/pkg/apis/com/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ishtarImportConfig struct {
	CassandraReplicatorFactor int32
}

type Auth struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
	Scope        string
}

type LMSecurityCtrl struct {
	restClient  *resty.Client
	lmBase      string
	username    string
	password    string
	loginResult *Auth
	loginTime   time.Time
}

type Ishtar struct {
	restClient     *resty.Client
	LMSecurityCtrl *LMSecurityCtrl
}

type CreateAssemblyBody struct {
	AssemblyName   string            `json:"assemblyName"`
	DescriptorName string            `json:"descriptorName"`
	IntendedState  string            `json:"intendedState"`
	Properties     map[string]string `json:"properties"`
}

func NewIshtar(client *resty.Client) *Ishtar {
	// client := resty.New()
	// client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// client.SetTimeout(2 * time.Minute)

	lmSecurityCtrl := LMSecurityCtrl{
		restClient: client,
		lmBase:     "https://nimrod:8290",
		username:   "jack",
		password:   "jack",
	}
	return &Ishtar{
		restClient:     client,
		LMSecurityCtrl: &lmSecurityCtrl,
	}
}

type login struct {
	Username string
	Password string
}

func (c *LMSecurityCtrl) login(username string, password string) (*Auth, error) {
	url := fmt.Sprintf("%s/api/login", c.lmBase)
	data := login{
		Username: username,
		Password: password,
	}

	t, err := template.New("login").Parse(`{"username":"{{.Username}}", "password":"{{.Password}}"}`)
	if err != nil {
		return nil, err
	}

	var postTpl bytes.Buffer
	if err := t.Execute(&postTpl, data); err != nil {
		return nil, err
	}

	body := postTpl.String()
	log.Info(fmt.Sprintf("Login %s", body))

	resp, err := c.restClient.R().
		EnableTrace().
		SetResult(&Auth{}).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusOK {
		loginResult := (*resp.Result().(*Auth))
		return &loginResult, nil
	}

	return nil, fmt.Errorf("%s", resp)
}

func (c *LMSecurityCtrl) getAccessToken() (string, error) {
	if c.needNewToken() {
		log.Info("Requesting new access token")
		result, err := c.login(c.username, c.password)
		if err != nil {
			return "", err
		}

		c.loginResult = result
		c.loginTime = time.Now()
	}

	return c.loginResult.AccessToken, nil
}

func (c *LMSecurityCtrl) needNewToken() bool {
	if c.loginResult == nil {
		log.Info("No current access token, must request one")
		return true
	}

	log.Info("Checking if access token has expired")
	expirationSeconds := c.loginResult.ExpiresIn
	loggedInTime := time.Now().Sub(c.loginTime).Seconds()
	log.Info(fmt.Sprintf("Logged in for %f seconds, token had an expiration time of %d seconds", loggedInTime, expirationSeconds))
	if int(loggedInTime) >= int(expirationSeconds) {
		log.Info("Token expired, must request a new one")
		return true
	}
	// If the token expires within 1 second, wait and get a new one
	if int32(loggedInTime) >= (expirationSeconds - 1) {
		log.Info("Expires in less than 1 second, waiting before requesting a new Token")
		time.Sleep(2 * time.Second)
		return true
	}

	return false
}

func ishtarConfig(cassandraReplicatorFactor int32) (string, error) {
	ishtarConfigImport := ishtarImportConfig{
		CassandraReplicatorFactor: cassandraReplicatorFactor,
	}

	t, err := template.New("ishtarConfig").Parse("alm:\n" +
		"  ishtar:\n" +
		"    cassandra:\n" +
		"      keyspaceManager:\n" +
		"        replicationFactor: {{.CassandraReplicatorFactor}}\n" +
		"  roles:\n" +
		"    SLMAdmin:\n" +
		"      ldapGroups:\n" +
		"        - SLMAdmin\n" +
		"      privileges:\n" +
		"        NsinstsMgt: read,write,execute\n" +
		"        VnfInstsMgt: read,write,execute\n" +
		"        nsDesMgt: read,write,execute\n" +
		"        VnfDesMgt: read,write,execute\n" +
		"        DeployLocMgt: read,write,execute\n" +
		"        VduMgt: read,write,execute\n" +
		"        IntentReqslMgt: read,execute\n" +
		"        IntentReqsOps: read,execute\n" +
		"        SlmAdmin: read,write,execute\n" +
		"        MaintModeOride: read,execute\n" +
		"        VduDesMgt: read,write,execute\n" +
		"        VduGrpMgt: read,write,execute\n" +
		"        VduInstsMgt: read,write,execute\n" +
		"        BehvrScenExec: read,write,execute\n" +
		"        BehvrScenDes: read,write\n" +
		"        RmDrvr: read,write\n" +
		"        ResourcePkg: write\n" +
		"    Portal:\n" +
		"      ldapGroups:\n" +
		"        - Portal\n" +
		"      privileges:\n" +
		"        NsinstsMgt: read,write,execute\n" +
		"        VduDesMgt: read\n" +
		"        VduGrpMgt: read\n" +
		"        VduInstsMgt: read\n" +
		"        VnfInstsMgt: read\n" +
		"        nsDesMgt: read\n" +
		"        DeployLocMgt: read\n" +
		"        IntentReqslMgt: read,execute\n" +
		"        BehvrScenExec: read,write,execute\n" +
		"        BehvrScenDes: read,write\n" +
		"        ResourcePkg: write\n" +
		"    ReadOnly:\n" +
		"      ldapGroups:\n" +
		"        - ReadOnly\n" +
		"      privileges:\n" +
		"        NsinstsMgt: read\n" +
		"        VduDesMgt: read\n" +
		"        VduGrpMgt: read\n" +
		"        VduInstsMgt: read\n" +
		"        VnfInstsMgt: read\n" +
		"        nsDesMgt: read\n" +
		"        VnfDesMgt: read\n" +
		"        DeployLocMgt: read\n" +
		"        VduMgt: read\n" +
		"        BehvrScenExec: read\n" +
		"        BehvrScenDes: read\n" +
		"    RootSecAdmin:\n" +
		"      ldapGroups:\n" +
		"        - RootSecAdmin\n" +
		"      privileges:\n" +
		"        SecAdmin: read,write,execute\n" +
		"    BehaviourScenarioExecute:\n" +
		"      privileges:\n" +
		"        NsinstsMgt: read,write\n" +
		"        IntentReqslMgt: execute\n" +
		"        IntentReqsOps: execute\n")

	if err != nil {
		return "", err
	}
	var ishtarTpl bytes.Buffer
	if err := t.Execute(&ishtarTpl, ishtarConfigImport); err != nil {
		return "", err
	}

	return ishtarTpl.String(), nil
}

func (r *ReconcileALM) installIshtar(cr *comv1alpha1.ALM, service serviceDeploymentInfo, reqLogger logr.Logger) (reconcile.Result, error) {
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
				},
				corev1.Volume{
					Name: "lm-client-credentials",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "lm-client-credentials",
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
				},
				corev1.VolumeMount{
					Name:      "lm-client-credentials",
					MountPath: "/var/lm/bootstrap",
				})
		}

		deployment := buildDeployment(cr.Namespace, deploymentName, cr, service, volumeMounts, volumes)

		if err := controllerutil.SetControllerReference(cr, deployment, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Failed to create a new %s Deployment", service.serviceName), "Namespace", cr.Namespace, "Name", deploymentName, "Error", err)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Deployment", service.serviceName), "Namespace", cr.Namespace, "Name", deploymentName)

		// Created successfully - don't requeue

		ingressName := "ishtar-ingress"
		ingress := buildIngress(cr.Spec.Secure, cr.Namespace, ingressName, "app.lm", 8280, "ishtar", "ishtar-tls")

		if err := controllerutil.SetControllerReference(cr, ingress, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), ingress)
		if err != nil {
			reqLogger.Info(fmt.Sprintf("Failed to create a new %s Ingress", ingressName), "Namespace", cr.Namespace, "Name", ingressName, "Error", err)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Ingress", ingressName), "Namespace", cr.Namespace, "Name", ingressName)

		// Created successfully - don't requeue and create service
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return r.service(cr, service, reqLogger)
}

func (i *Ishtar) CreateAssembly(reqLogger logr.Logger, assembly CreateAssemblyBody) (string, error) {
	accessToken, err := i.LMSecurityCtrl.getAccessToken()
	if err != nil {
		reqLogger.Error(err, "Unable to get access token")
		return "", err
	}

	bytes, err := json.Marshal(assembly)
	if err != nil {
		reqLogger.Error(err, "Unable to create assembly template")
		return "", err
	}
	assemblyJSON := string(bytes)

	// t, err := template.New("createAssembly").Parse(`{"assemblyName":"{{.AssemblyName}}", "descriptorName":"{{.DescriptorName}}", "intendedState": "{{.IntendedState}}", "properties": { For {{range $k,$v := .Properties}} "{{$k}}": "{{$v}}" {{end}}} }`)
	// if err != nil {
	// 	reqLogger.Error(err, "Unable to create assembly template")
	// 	return "", err
	// }
	// var postTpl bytes.Buffer
	// if err := t.Execute(&postTpl, assembly); err != nil {
	// 	reqLogger.Error(err, "Unable to create assembly template")
	// 	return "", err
	// }
	// assemblyJSON := postTpl.String()

	reqLogger.Info(fmt.Sprintf("Create assembly %s", assemblyJSON))
	reqLogger.Info(fmt.Sprintf("Access token %s", accessToken))

	resp, err := i.restClient.R().
		EnableTrace().
		SetBody(assemblyJSON).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Post("https://ishtar:8280/api/intent/createAssembly")
	if err != nil {
		reqLogger.Error(err, "Unable to create assembly")
		return "", err
	}

	reqLogger.Info(fmt.Sprintf("Create assembly status %d", resp.StatusCode()))

	if resp.StatusCode() != http.StatusCreated {
		return "", fmt.Errorf("Create assembly failed %s %s", resp.Body(), string(resp.StatusCode()))
	}

	location := resp.Header().Get(http.CanonicalHeaderKey("Location"))
	ss := strings.Split(location, "/")
	return ss[len(ss)-1], nil
}

type HealthStatus struct {
	Status string
}

func (i *Ishtar) Health(reqLogger logr.Logger) (bool, error) {
	accessToken, err := i.LMSecurityCtrl.getAccessToken()
	if err != nil {
		reqLogger.Error(err, "Unable to get access token")
		return false, err
	}

	resp, err := i.restClient.R().
		EnableTrace().
		SetResult(&HealthStatus{}).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Get("https://ishtar:8280/management/health")
	if err != nil {
		reqLogger.Error(err, "Unable to get Ishtar health")
		return false, err
	}

	healthStatus := (*resp.Result().(*HealthStatus))

	// statusBody, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	reqLogger.Error(err, "Unable to get Ishtar health")
	// 	return false, err
	// }

	reqLogger.Info(fmt.Sprintf("Ishtar health status %s", healthStatus.Status))

	return healthStatus.Status == "UP", nil
}

func (i *Ishtar) GetAssemblyStatus(reqLogger logr.Logger, processID string) (string, error) {
	accessToken, err := i.LMSecurityCtrl.getAccessToken()
	if err != nil {
		reqLogger.Error(err, "Unable to get access token")
		return "", err
	}

	resp, err := i.restClient.R().
		EnableTrace().
		SetResult(map[string]interface{}{}).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Get(fmt.Sprintf("https://ishtar:8280/api/processes/%s", processID))
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("Get process failed %s %s", resp.Body(), string(resp.StatusCode()))
	}

	result := (*resp.Result().(*map[string]interface{}))
	return result["status"].(string), nil
}
