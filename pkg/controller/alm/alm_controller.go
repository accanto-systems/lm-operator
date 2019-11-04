package alm

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	resty "github.com/go-resty/resty/v2"
	comv1alpha1 "github.com/orgs/accanto-systems/lm-operator/pkg/apis/com/v1alpha1"
	v1beta1 "k8s.io/api/apps/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_alm")

func int32Ptr(i int32) *int32 { return &i }

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ALM Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetTimeout(2 * time.Minute)
	return &ReconcileALM{client: mgr.GetClient(), scheme: mgr.GetScheme(), ishtar: NewIshtar(client)}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("alm-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// kubernetes.NewForConfig(nil)
	// Watch for changes to primary resource ALM
	err = c.Watch(&source.Kind{Type: &comv1alpha1.ALM{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resources and requeue the owner ALM
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &comv1alpha1.ALM{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &comv1alpha1.ALM{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &extv1beta1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &comv1alpha1.ALM{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &comv1alpha1.ALM{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1beta1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &comv1alpha1.ALM{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileALM implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileALM{}

// ReconcileALM reconciles a ALM object
type ReconcileALM struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	ishtar *Ishtar
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ALM object and makes changes based on the state read
// and what is in the ALM.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileALM) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ALM")

	// Fetch the ALM instance
	instance := &comv1alpha1.ALM{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return r.createALM(request, instance, reqLogger)
}

func (r *ReconcileALM) createALM(request reconcile.Request, instance *comv1alpha1.ALM, reqLogger logr.Logger) (reconcile.Result, error) {
	// we assume that the presence of Daytona means that ALM has been deployed
	// TODO use service name
	found, err := r.deploymentByNameExists(instance, "daytona")
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Error detecting whether Daytona exists for release %s", instance.Name), "Namespace", instance.Namespace)
		return reconcile.Result{}, err
	}

	if !found {
		// no Daytona, try to find configurator
		reqLogger.Info(fmt.Sprintf("Daytona does not exist for release %s", instance.Name), "Namespace", instance.Namespace)

		lmConfiguratorName := fmt.Sprintf("%s-lm-configurator", instance.Name)
		found, err := r.getLMConfigurator(instance.Namespace, lmConfiguratorName)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to get %s Job", lmConfiguratorName), "Namespace", instance.Namespace, "Name", lmConfiguratorName)
			return reconcile.Result{}, err
		}

		// TODO cache this somewhere
		deploymentInfo, err := createDeploymentInfo(instance, reqLogger)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to get release information"))
			return reconcile.Result{}, err
		}

		s, _ := json.MarshalIndent(deploymentInfo, "", "\t")
		reqLogger.Info(fmt.Sprintf("Creating ALM with deployment info %s", string(s)))

		if found == nil {
			// no configurator

			var result reconcile.Result
			if deploymentInfo.configurator.run {
				reqLogger.Info(fmt.Sprintf("Configurator does not exist for release %s, creating it", instance.Name), "Namespace", instance.Namespace)
				result, err = r.createLMConfigurator(request, deploymentInfo, instance, deploymentInfo.configurator, reqLogger)
			} else {
				reqLogger.Info(fmt.Sprintf("Creating LM microservices for release %s", instance.Name), "Namespace", instance.Namespace)
				result, err = r.createMicroservices(deploymentInfo, request, instance, reqLogger)
			}

			// if result.Requeue {
			// 	return reconcile.Result{}, err
			// }

			return result, err
		}

		if found.Status.Succeeded > 0 {
			r.addSecretReference(instance.Namespace, "lm-certs", instance, reqLogger)
			r.addSecretReference(instance.Namespace, "lm-client-credentials", instance, reqLogger)
			r.addSecretReference(instance.Namespace, "lm-keystore", instance, reqLogger)
			r.addSecretReference(instance.Namespace, "nimrod-tls", instance, reqLogger)
			r.addSecretReference(instance.Namespace, "brent-tls", instance, reqLogger)
			r.addSecretReference(instance.Namespace, "ishtar-tls", instance, reqLogger)

			// lm-configurator has completed
			reqLogger.Info("LM-configurator complete %s, creating LM microservices", "Namespace", instance.Namespace, "Name", lmConfiguratorName)
			result, err := r.createMicroservices(deploymentInfo, request, instance, reqLogger)
			// TODO handle err
			if result.Requeue {
				return reconcile.Result{}, err
			}
		} else {
			// re-queue because the lm-configurator Job is not complete
			s, _ := json.MarshalIndent(found.Status, "", "\t")
			reqLogger.Info(fmt.Sprintf("LM-configurator not complete %s, re-queuing", string(s)), "Namespace", instance.Namespace, "Name", lmConfiguratorName)
			return reconcile.Result{Requeue: true}, nil
		}
	}

	// nothing to do, don't re-queue
	reqLogger.Info(fmt.Sprintf("Daytona exists for release %s", instance.Name), "Namespace", instance.Namespace)

	status, err := r.ishtar.Health(reqLogger)
	if err != nil {
		return reconcile.Result{Requeue: true}, err
	}

	instance.Status.IshtarHealthy = status
	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Failed to update ALM status.")
		return reconcile.Result{}, err
	}
	log.Info("ALM status updated")

	return reconcile.Result{}, nil
}

func (r *ReconcileALM) getLMConfigurator(namespace string, lmConfiguratorName string) (*batchv1.Job, error) {
	found := &batchv1.Job{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: lmConfiguratorName, Namespace: namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return found, nil
	}
}

func (r *ReconcileALM) createLMConfigurator(request reconcile.Request, deploymentInfo deploymentInfo, cr *comv1alpha1.ALM, serviceDeploymentInfo configuratorDeploymentInfo, reqLogger logr.Logger) (reconcile.Result, error) {

	lmConfiguratorName := fmt.Sprintf("%s-lm-configurator", cr.Name)
	found, err := r.getLMConfigurator(cr.Namespace, lmConfiguratorName)
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Failed to get %s Job", lmConfiguratorName), "Namespace", cr.Namespace, "Name", lmConfiguratorName)
		return reconcile.Result{}, err
	}

	if found == nil {
		reqLogger.Info(fmt.Sprintf("LM Configurator Job %s not found", lmConfiguratorName))

		// LM Configurator CM
		lmConfiguratorCMName := fmt.Sprintf("%s-%s-cm", cr.Name, deploymentInfo.configurator.serviceName)
		foundCm := &corev1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: lmConfiguratorCMName, Namespace: cr.Namespace}, foundCm)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Creating a new %s ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfiguratorCMName)

			cm, err := buildConfiguratorCM(lmConfiguratorCMName, cr, deploymentInfo.configurator)
			if err != nil {
				reqLogger.Error(err, fmt.Sprintf("Failed to create %s ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfiguratorCMName)
				return reconcile.Result{}, err
			}
			if err := controllerutil.SetControllerReference(cr, cm, r.scheme); err != nil {
				reqLogger.Error(err, fmt.Sprintf("Failed to set owner for %s ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfiguratorCMName)
				return reconcile.Result{}, err
			}

			err = r.client.Create(context.TODO(), cm)
			if err != nil {
				reqLogger.Error(err, fmt.Sprintf("Failed to create %s ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfiguratorCMName)
				return reconcile.Result{}, err
			}

			reqLogger.Info(fmt.Sprintf("Created a new %s ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cm.Namespace, "Name", cm.Name)
			// LM Configurator CM created successfully - drop through and don't requeue
		} else if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to get %s ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfiguratorCMName)
			return reconcile.Result{}, err
		}

		// LM Configurator Config Import CM
		lmConfigImportCmName := fmt.Sprintf("%s-%s-cm", cr.Name, "lm-config-import")
		foundCm = &corev1.ConfigMap{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: lmConfigImportCmName, Namespace: cr.Namespace}, foundCm)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Creating a new %s LM Config Import ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfigImportCmName)

			lmConfigImportCm, err := buildLmConfigImportCm(cr.Namespace, lmConfigImportCmName)
			if err != nil {
				reqLogger.Error(err, fmt.Sprintf("Failed to create a new %s LM Config Import ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfigImportCmName)
				return reconcile.Result{}, err
			}

			if err := controllerutil.SetControllerReference(cr, lmConfigImportCm, r.scheme); err != nil {
				reqLogger.Error(err, fmt.Sprintf("Failed to set parent of %s LM Config Import ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfigImportCmName)
				return reconcile.Result{}, err
			}

			err = r.client.Create(context.TODO(), lmConfigImportCm)
			if err != nil {
				reqLogger.Error(err, fmt.Sprintf("Failed to create a new %s LM Config Import ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfigImportCmName)
				return reconcile.Result{}, err
			}

			reqLogger.Info(fmt.Sprintf("Created a new %s LM Config Import ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", lmConfigImportCm.Namespace, "Name", lmConfigImportCmName)
			// LM Configurator Config Import CM created successfully - drop through and don't requeue
		} else if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to get %s LM Config Import ConfigMap", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfigImportCmName)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Creating a new %s Job", serviceDeploymentInfo.serviceName), "Namespace", cr.Namespace, "Name", lmConfiguratorName)

		job := buildJob(cr, deploymentInfo.configurator, cr.Spec.DockerRepo, cr.Namespace, lmConfiguratorName, lmConfiguratorCMName, lmConfigImportCmName)

		if err := controllerutil.SetControllerReference(cr, job, r.scheme); err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to set parent for new %s Job", "lm-configurator"), "Namespace", cr.Namespace, "Name", lmConfiguratorName)
			return reconcile.Result{}, err
		}

		err = r.client.Create(context.TODO(), job)
		if err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to create a new %s Job", "lm-configurator"), "Namespace", cr.Namespace, "Name", lmConfiguratorName)
			return reconcile.Result{}, err
		}

		reqLogger.Info(fmt.Sprintf("Created a new %s Job", "lm-configurator"), "Namespace", cr.Namespace, "Name", lmConfiguratorName)

		// Created successfully - don't requeue
		return reconcile.Result{}, nil
	}

	if found.Status.Succeeded > 0 {
		r.addSecretReference(cr.Namespace, "lm-certs", cr, reqLogger)
		r.addSecretReference(cr.Namespace, "lm-client-credentials", cr, reqLogger)
		r.addSecretReference(cr.Namespace, "lm-keystore", cr, reqLogger)
		r.addSecretReference(cr.Namespace, "nimrod-tls", cr, reqLogger)
		r.addSecretReference(cr.Namespace, "brent-tls", cr, reqLogger)
		r.addSecretReference(cr.Namespace, "ishtar-tls", cr, reqLogger)

		// lm-configurator has completed

		// create lm configurator config map indicating it has run successfully

		reqLogger.Info("LM-configurator complete %s, creating LM microservices", "Namespace", cr.Namespace, "Name", lmConfiguratorName)
		result, err := r.createMicroservices(deploymentInfo, request, cr, reqLogger)
		if result.Requeue {
			return reconcile.Result{}, err
		}
	} else {
		// re-queue because the lm-configurator Job is not complete
		s, _ := json.MarshalIndent(found.Status, "", "\t")
		reqLogger.Info(fmt.Sprintf("LM-configurator not complete %s, re-queuing", string(s)), "Namespace", cr.Namespace, "Name", lmConfiguratorName)
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileALM) addSecretReference(namespace string, secretName string, cr *comv1alpha1.ALM, reqLogger logr.Logger) {
	foundSecret := &corev1.Secret{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, foundSecret)
	if err == nil && foundSecret != nil {
		if err := controllerutil.SetControllerReference(cr, foundSecret, r.scheme); err != nil {
			reqLogger.Error(err, fmt.Sprintf("Failed to set parent for secret %s", secretName), "Namespace", namespace, "Name", secretName)
		}
	}
}

func (r *ReconcileALM) createMicroservices(deploymentInfo deploymentInfo, request reconcile.Request, instance *comv1alpha1.ALM, reqLogger logr.Logger) (reconcile.Result, error) {
	conductorResult, conductorErr := r.installConductor(instance, deploymentInfo.conductor, reqLogger)
	if conductorResult.Requeue {
		return conductorResult, conductorErr
	}

	apolloResult, apolloErr := r.installApollo(instance, deploymentInfo.apollo, reqLogger)
	if apolloResult.Requeue {
		return apolloResult, apolloErr
	}

	galileoResult, galileoErr := r.installGalileo(instance, deploymentInfo.galileo, reqLogger)
	if galileoResult.Requeue {
		return galileoResult, galileoErr
	}

	talledegaResult, talledegaErr := r.installTalledega(instance, deploymentInfo.talledega, reqLogger)
	if talledegaResult.Requeue {
		return talledegaResult, talledegaErr
	}

	daytonaResult, daytonaErr := r.installDaytona(instance, deploymentInfo.daytona, reqLogger)
	if daytonaResult.Requeue {
		return daytonaResult, daytonaErr
	}

	relayResult, relayErr := r.installRelay(instance, deploymentInfo.relay, reqLogger)
	if relayResult.Requeue {
		return relayResult, relayErr
	}

	watchtowerResult, watchtowerErr := r.installWatchtower(instance, deploymentInfo.watchtower, reqLogger)
	if watchtowerResult.Requeue {
		return watchtowerResult, watchtowerErr
	}

	dokiResult, dokiErr := r.installDoki(instance, deploymentInfo.doki, reqLogger)
	if dokiResult.Requeue {
		return dokiResult, dokiErr
	}

	nimrodResult, nimrodErr := r.installNimrod(instance, deploymentInfo.nimrod, reqLogger)
	if nimrodResult.Requeue {
		return nimrodResult, nimrodErr
	}

	ishtarResult, ishtarErr := r.installIshtar(instance, deploymentInfo.ishtar, reqLogger)
	if ishtarResult.Requeue {
		return ishtarResult, ishtarErr
	}

	brentResult, brentErr := r.installBrent(instance, deploymentInfo.brent, reqLogger)
	if brentResult.Requeue {
		return brentResult, brentErr
	}

	return reconcile.Result{}, nil
}
