//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package groupresourcequotaenforcer

import (
	"context"

	// IBMDEV
	"reflect"

	operatorv1alpha1 "github.com/IBM/ibm-mt-grqenforcer-operator/pkg/apis/operator/v1alpha1"
	"github.com/IBM/ibm-mt-grqenforcer-operator/version"

	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_groupresourcequotaenforcer")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

/*
IBMDEV
This type and global variable provide the names for all the resources created/managed by this operator.
As such, it provides an overview of the required pieces for the enforcement component.
Note however that the certSecret isn't currently used (secret name is provided by CR param instead).
*/
type grqeNameSuffix struct {
	//certSecret       string	/*currently unused, see commented out code blocks*/
	grqeDeployment     string
	grqeService        string
	grqeWebhook        string
	bridgeDeployment   string
	bridgeService      string
	serviceAccount     string
	clusterRole        string
	clusterRoleBinding string
}

var suffix = grqeNameSuffix{
	//certSecret:       "-grqe-crt",	/*currently unused, see commented out code blocks*/
	grqeDeployment:     "-grqe-dep",
	grqeService:        "-grqe-svc",
	grqeWebhook:        ".mt-grqe.ibm",
	bridgeDeployment:   "-grqb-dep",
	bridgeService:      "-grqb-svc",
	serviceAccount:     "-grqe-svcacct",
	clusterRole:        "-grqe-crole",
	clusterRoleBinding: "-grqe-crbinding",
}

// Add creates a new GroupResourceQuotaEnforcer Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileGroupResourceQuotaEnforcer{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("groupresourcequotaenforcer-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource GroupResourceQuotaEnforcer
	err = c.Watch(&source.Kind{Type: &operatorv1alpha1.GroupResourceQuotaEnforcer{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Deployments and requeue the owner GroupResourceQuotaEnforcer
	// IBMDEV: Done
	secondaryResourceTypes := []runtime.Object{
		&appsv1.Deployment{},
		&corev1.Service{},
		/* &corev1.Secret{}, */
		&corev1.ServiceAccount{},
		&rbacv1.ClusterRole{},
		&rbacv1.ClusterRoleBinding{},
		&admissionv1beta1.MutatingWebhookConfiguration{},
	}
	for _, restype := range secondaryResourceTypes {
		log.Info("Watching", "restype", restype)
		//err = c.Watch(&kind, &handler.EnqueueRequestForOwner{
		err = c.Watch(&source.Kind{Type: restype}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &operatorv1alpha1.GroupResourceQuotaEnforcer{},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// blank assignment to verify that ReconcileGroupResourceQuotaEnforcer implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileGroupResourceQuotaEnforcer{}

// ReconcileGroupResourceQuotaEnforcer reconciles a GroupResourceQuotaEnforcer object
type ReconcileGroupResourceQuotaEnforcer struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a GroupResourceQuotaEnforcer object and makes changes based on the state read
// and what is in the GroupResourceQuotaEnforcer.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileGroupResourceQuotaEnforcer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// IBMDEV As operator/CR are cluster scoped, request.Namespace is expected to be empty string
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	// IBMDEV

	// Fetch the CR instance
	cr := &operatorv1alpha1.GroupResourceQuotaEnforcer{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cr)
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

	var recResult reconcile.Result
	var recErr error

	// TODO:
	// Reconcile - create grqe-ca-issuer (self signed)
	// Reconcile - create grqe-ca-cert (isCA) from grqe-ca-issuer
	// Reconcile - create grqe-issuer referencing grqe-ca
	// Reconcile - create gqre-cert from grqe-issuer
	// Update deployments to use gqre-cert-secret created by certmanager
	// Update CRD to no longer take certSecret spec param

	// Reconcile the expected deployment
	recResult, recErr = r.deploymentForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Reconcile the expected grq enforcer service
	recResult, recErr = r.grqEnforcerServiceForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Reconcile the expected bridge deployment
	recResult, recErr = r.bridgeDeploymentForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Reconcile the expected bridge service
	recResult, recErr = r.bridgeServiceForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Reconcile the expected ServiceAccount
	recResult, recErr = r.serviceAccountForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Reconcile the expected Role
	recResult, recErr = r.roleForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Reconcile the expected RoleBinding
	recResult, recErr = r.roleBindingForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// Reconcile the expected MutatingWebhookConfig
	recResult, recErr = r.webhookConfigForCR(cr)
	if recErr != nil || recResult.Requeue {
		return recResult, recErr
	}

	// If necessary, update the CR status with the pod names
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(cr.Spec.InstanceNamespace),
		client.MatchingLabels(labelsForDeployment(cr.Name)),
	}
	if err = r.client.List(context.TODO(), podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "GroupResourceQuotaEnforcer.Namespace", cr.Spec.InstanceNamespace, "GroupResourceQuotaEnforcer.Name", cr.Name)
		return reconcile.Result{}, err
	}
	podNames := []string{}
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	// Get bridge pods too
	listOpts = []client.ListOption{
		client.InNamespace(cr.Spec.InstanceNamespace),
		client.MatchingLabels(labelsForBridgeDeployment(cr.Name)),
	}
	if err = r.client.List(context.TODO(), podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list bridge pods", "Namespace", cr.Spec.InstanceNamespace, "Name", cr.Name)
		return reconcile.Result{}, err
	}
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}

	if !reflect.DeepEqual(podNames, cr.Status.Nodes) {
		reqLogger.Info("Updating CR status", "Name", cr.Name)
		cr.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), cr)
		if err != nil {
			reqLogger.Error(err, "Failed to update status")
			return reconcile.Result{}, err
		}
	}

	reqLogger.Info("Reconciliation successful!", "Name", cr.Name)
	return reconcile.Result{}, nil
}

//IBMDEV
func meteringAnnotations() map[string]string {
	return map[string]string{
		"productName":    version.Name,
		"productID":      version.Name,
		"productVersion": version.Version,
	}
}

//IBMDEV
func labelsForDeployment(crName string) map[string]string {
	return map[string]string{"app": "groupresourcequotaenforcer", "groupresourcequotaenforcer_cr": crName}
}

//IBMDEV
func labelsForBridgeDeployment(crName string) map[string]string {
	return map[string]string{"app": "iam-bridge", "groupresourcequotaenforcer_cr": crName}
}

// IBMDEV deploymentForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) deploymentForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	ls := labelsForDeployment(cr.Name)

	int32_1 := int32(1)
	int32_420 := int32(420)

	var imagePullSecrets []corev1.LocalObjectReference
	if cr.Spec.ImagePullSecret != "" {
		imagePullSecrets = append(
			imagePullSecrets,
			corev1.LocalObjectReference{
				Name: cr.Spec.ImagePullSecret,
			},
		)
	}

	expectedRes := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + suffix.grqeDeployment,
			Namespace: cr.Spec.InstanceNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &int32_1,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      ls,
					Annotations: meteringAnnotations(),
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cr.Name + suffix.serviceAccount,
					Containers: []corev1.Container{{
						Image:           cr.Spec.ImageRegistry + "/ibm-mt-groupresourcequota:" + version.Version,
						Name:            "ibm-mt-grq-enforcer",
						ImagePullPolicy: "IfNotPresent",
						Args: []string{
							"-port=7443",
							"-tlsCertFile=/etc/webhook/certs/cert.pem",
							"-tlsKeyFile=/etc/webhook/certs/key.pem",
							"-alsologtostderr",
							"-v=4",
							"2>&1",
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 7443,
							Protocol:      "TCP",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "webhook-certs",
							MountPath: "/etc/webhook/certs",
							ReadOnly:  true,
						}},
						//LivenessProbe: TBD,
						//ReadinessProbe: TBD,
					}},
					ImagePullSecrets: imagePullSecrets,
					Volumes: []corev1.Volume{{
						Name: "webhook-certs",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								//SecretName: cr.Name + suffix.certSecret,
								SecretName:  cr.Spec.CertSecret,
								DefaultMode: &int32_420,
							},
						},
					}},
				},
			},
		},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If deployment does not exist, create it and requeue
	foundDeployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name, Namespace: cr.Spec.InstanceNamespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", expectedRes.Namespace, "Deployment.Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", expectedRes.Namespace, "Deployment.Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	} else if !reflect.DeepEqual(foundDeployment.Spec.Template.Spec.Volumes, expectedRes.Spec.Template.Spec.Volumes) ||
		len(foundDeployment.Spec.Template.Spec.Containers) != len(expectedRes.Spec.Template.Spec.Containers) ||
		!reflect.DeepEqual(foundDeployment.Spec.Template.Spec.Containers[0].Name, expectedRes.Spec.Template.Spec.Containers[0].Name) ||
		!reflect.DeepEqual(foundDeployment.Spec.Template.Spec.Containers[0].Image, expectedRes.Spec.Template.Spec.Containers[0].Image) ||
		!reflect.DeepEqual(foundDeployment.Spec.Template.Spec.Containers[0].Args, expectedRes.Spec.Template.Spec.Containers[0].Args) ||
		!reflect.DeepEqual(foundDeployment.Spec.Template.Spec.Containers[0].VolumeMounts, expectedRes.Spec.Template.Spec.Containers[0].VolumeMounts) {
		// Spec is incorrect, update it and requeue
		reqLogger.Info("Found deployment spec is incorrect", "Found", foundDeployment.Spec.Template.Spec, "Expected", expectedRes.Spec.Template.Spec)
		foundDeployment.Spec.Template.Spec.Volumes = expectedRes.Spec.Template.Spec.Volumes
		foundDeployment.Spec.Template.Spec.Containers = expectedRes.Spec.Template.Spec.Containers
		err = r.client.Update(context.TODO(), foundDeployment)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Namespace", cr.Spec.InstanceNamespace, "Name", foundDeployment.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// No reconcile was necessary
	return reconcile.Result{}, nil
}

// IBMDEV bridgeDeploymentForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) bridgeDeploymentForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	ls := labelsForBridgeDeployment(cr.Name)

	int32_1 := int32(1)
	int32_420 := int32(420)

	var imagePullSecrets []corev1.LocalObjectReference
	if cr.Spec.ImagePullSecret != "" {
		imagePullSecrets = append(
			imagePullSecrets,
			corev1.LocalObjectReference{
				Name: cr.Spec.ImagePullSecret,
			},
		)
	}

	expectedRes := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + suffix.bridgeDeployment,
			Namespace: cr.Spec.InstanceNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &int32_1,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cr.Name + suffix.serviceAccount,
					Containers: []corev1.Container{{
						Image:           cr.Spec.ImageRegistry + "/ibm-mt-iam-bridge:" + version.Version,
						Name:            "ibm-mt-iam-bridge",
						ImagePullPolicy: "IfNotPresent",
						Args: []string{
							"-port=7443",
							"-tlsCertFile=/etc/webhook/certs/cert.pem",
							"-tlsKeyFile=/etc/webhook/certs/key.pem",
							"-alsologtostderr",
							"-v=4",
							"2>&1",
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 7443,
							Protocol:      "TCP",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "webhook-certs",
							MountPath: "/etc/webhook/certs",
							ReadOnly:  true,
						}},
						//LivenessProbe: TBD,
						//ReadinessProbe: TBD,
					}},
					ImagePullSecrets: imagePullSecrets,
					Volumes: []corev1.Volume{{
						Name: "webhook-certs",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								//SecretName: cr.Name + suffix.certSecret,
								SecretName:  cr.Spec.CertSecret,
								DefaultMode: &int32_420,
							},
						},
					}},
				},
			},
		},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If deployment does not exist, create it and requeue
	foundBridgeDep := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name, Namespace: cr.Spec.InstanceNamespace}, foundBridgeDep)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", expectedRes.Namespace, "Deployment.Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", expectedRes.Namespace, "Deployment.Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	} else if !reflect.DeepEqual(foundBridgeDep.Spec.Template.Spec.Volumes, expectedRes.Spec.Template.Spec.Volumes) ||
		len(foundBridgeDep.Spec.Template.Spec.Containers) != len(expectedRes.Spec.Template.Spec.Containers) ||
		!reflect.DeepEqual(foundBridgeDep.Spec.Template.Spec.Containers[0].Name, expectedRes.Spec.Template.Spec.Containers[0].Name) ||
		!reflect.DeepEqual(foundBridgeDep.Spec.Template.Spec.Containers[0].Image, expectedRes.Spec.Template.Spec.Containers[0].Image) ||
		!reflect.DeepEqual(foundBridgeDep.Spec.Template.Spec.Containers[0].Args, expectedRes.Spec.Template.Spec.Containers[0].Args) ||
		!reflect.DeepEqual(foundBridgeDep.Spec.Template.Spec.Containers[0].Ports, expectedRes.Spec.Template.Spec.Containers[0].Ports) ||
		!reflect.DeepEqual(foundBridgeDep.Spec.Template.Spec.Containers[0].VolumeMounts, expectedRes.Spec.Template.Spec.Containers[0].VolumeMounts) {
		// Spec is incorrect, update it and requeue
		reqLogger.Info("Found deployment spec is incorrect", "Found", foundBridgeDep.Spec.Template.Spec, "Expected", expectedRes.Spec.Template.Spec)
		foundBridgeDep.Spec.Template.Spec.Volumes = expectedRes.Spec.Template.Spec.Volumes
		foundBridgeDep.Spec.Template.Spec.Containers = expectedRes.Spec.Template.Spec.Containers
		err = r.client.Update(context.TODO(), foundBridgeDep)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Namespace", cr.Spec.InstanceNamespace, "Name", foundBridgeDep.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// No reconcile was necessary
	return reconcile.Result{}, nil
}

// IBMDEV grqEnforcerServiceForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) grqEnforcerServiceForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	ls := labelsForDeployment(cr.Name)

	expectedRes := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + suffix.grqeService,
			Namespace: cr.Spec.InstanceNamespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{{
				Port:       443,
				TargetPort: intstr.FromInt(7443),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If service does not exist, create it and requeue
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name, Namespace: cr.Spec.InstanceNamespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)

		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	} else if !reflect.DeepEqual(foundService.Spec.Ports, expectedRes.Spec.Ports) ||
		!reflect.DeepEqual(foundService.Spec.Selector, expectedRes.Spec.Selector) {
		// Spec is incorrect, update it and requeue
		reqLogger.Info("Found service spec is incorrect", "Found", foundService.Spec, "Expected", expectedRes.Spec)
		foundService.Spec.Ports = expectedRes.Spec.Ports
		foundService.Spec.Selector = expectedRes.Spec.Selector
		err = r.client.Update(context.TODO(), foundService)
		if err != nil {
			reqLogger.Error(err, "Failed to update Service", "Namespace", cr.Spec.InstanceNamespace, "Name", foundService.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// No reconcile was necessary
	return reconcile.Result{}, nil
}

// IBMDEV bridgeServiceForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) bridgeServiceForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	ls := labelsForBridgeDeployment(cr.Name)

	expectedRes := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + suffix.bridgeService,
			Namespace: cr.Spec.InstanceNamespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{{
				Port:       443,
				TargetPort: intstr.FromInt(7443),
				Protocol:   corev1.ProtocolTCP,
			}},
		},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If service does not exist, create it and requeue
	foundService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name, Namespace: cr.Spec.InstanceNamespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)

		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	} else if !reflect.DeepEqual(foundService.Spec.Ports, expectedRes.Spec.Ports) ||
		!reflect.DeepEqual(foundService.Spec.Selector, expectedRes.Spec.Selector) {
		// Spec is incorrect, update it and requeue
		reqLogger.Info("Found service spec is incorrect", "Found", foundService.Spec, "Expected", expectedRes.Spec)
		foundService.Spec.Ports = expectedRes.Spec.Ports
		foundService.Spec.Selector = expectedRes.Spec.Selector
		err = r.client.Update(context.TODO(), foundService)
		if err != nil {
			reqLogger.Error(err, "Failed to update Service", "Namespace", cr.Spec.InstanceNamespace, "Name", foundService.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// No reconcile was necessary
	return reconcile.Result{}, nil
}

// IBMDEV serviceAccountForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) serviceAccountForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	expectedRes := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + suffix.serviceAccount,
			Namespace: cr.Spec.InstanceNamespace,
		},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If ServiceAccount does not exist, create it and requeue
	foundSvcAcct := &corev1.ServiceAccount{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name, Namespace: cr.Spec.InstanceNamespace}, foundSvcAcct)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ServiceAccount", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new ServiceAccount", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get ServiceAccount")
		return reconcile.Result{}, err
	}
	// No extra validation of the service account required

	// No reconcile was necessary
	return reconcile.Result{}, nil
}

// IBMDEV roleForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) roleForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	expectedRes := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name + suffix.clusterRole,
		},
		Rules: []rbacv1.PolicyRule{{
			Verbs:     []string{"get", "watch", "list", "create", "update"},
			APIGroups: []string{""},
			Resources: []string{"resourcequotas", "namespaces", "groupresourcequotas"},
		}},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If Role does not exist, create it and requeue
	foundRole := &rbacv1.ClusterRole{}
	// Note: clusterroles are cluster-scoped, so this does not search using namespace (unlike other resources above)
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name}, foundRole)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new "+reflect.TypeOf(expectedRes).String(), "Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new "+reflect.TypeOf(expectedRes).String(), "Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get "+reflect.TypeOf(expectedRes).String())
		return reconcile.Result{}, err
	} else if !reflect.DeepEqual(foundRole.Rules, expectedRes.Rules) {
		// Spec is incorrect, update it and requeue
		reqLogger.Info("Found role is incorrect", "Found", foundRole.Rules, "Expected", expectedRes.Rules)
		foundRole.Rules = expectedRes.Rules
		err = r.client.Update(context.TODO(), foundRole)
		if err != nil {
			reqLogger.Error(err, "Failed to update role", "Name", foundRole.Name)
			return reconcile.Result{}, err
		}
		// Updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// No reconcile was necessary
	return reconcile.Result{}, nil
}

// IBMDEV rolebindingForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) roleBindingForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	expectedRes := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name + suffix.clusterRoleBinding,
		},
		Subjects: []rbacv1.Subject{{
			APIGroup:  "",
			Kind:      "ServiceAccount",
			Name:      cr.Name + suffix.serviceAccount,
			Namespace: cr.Spec.InstanceNamespace,
		}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     cr.Name + suffix.clusterRole,
		},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If RoleBinding does not exist, create it and requeue
	foundRoleBinding := &rbacv1.ClusterRoleBinding{}
	// Note: clusterrolebindings are cluster-scoped, so this does not search using namespace (unlike other resources above)
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name}, foundRoleBinding)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new "+reflect.TypeOf(expectedRes).String(), "Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new "+reflect.TypeOf(expectedRes).String(), "Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get "+reflect.TypeOf(expectedRes).String())
		return reconcile.Result{}, err
	} else if !reflect.DeepEqual(foundRoleBinding.Subjects, expectedRes.Subjects) ||
		!reflect.DeepEqual(foundRoleBinding.RoleRef, expectedRes.RoleRef) {
		// Spec is incorrect, update it and requeue
		reqLogger.Info("Found rolebinding is incorrect", "Found", foundRoleBinding, "Expected", expectedRes)
		err = r.client.Delete(context.TODO(), foundRoleBinding)
		if err != nil {
			reqLogger.Error(err, "Failed to delete rolebinding", "Name", foundRoleBinding.Name)
			return reconcile.Result{}, err
		}
		// Deleted - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// No reconcile was necessary
	return reconcile.Result{}, nil
}

// IBMDEV webhookConfigForCR returns (reconcile.Result, error)
func (r *ReconcileGroupResourceQuotaEnforcer) webhookConfigForCR(cr *operatorv1alpha1.GroupResourceQuotaEnforcer) (reconcile.Result, error) {
	reqLogger := log.WithValues("cr.Name", cr.Name)

	path := "/mutate"
	int32_443 := int32(443)
	service := admissionv1beta1.ServiceReference{
		Namespace: cr.Spec.InstanceNamespace,
		Name:      cr.Name + suffix.grqeService,
		Path:      &path,
		Port:      &int32_443,
	}
	scope := admissionv1beta1.AllScopes

	expectedRes := &admissionv1beta1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name + suffix.grqeWebhook,
		},
		Webhooks: []admissionv1beta1.MutatingWebhook{{
			Name: cr.Name + suffix.grqeWebhook,
			ClientConfig: admissionv1beta1.WebhookClientConfig{
				Service: &service,
			},
			Rules: []admissionv1beta1.RuleWithOperations{{
				Operations: []admissionv1beta1.OperationType{
					admissionv1beta1.Create,
					admissionv1beta1.Update,
					admissionv1beta1.Delete,
				},
				Rule: admissionv1beta1.Rule{
					APIGroups:   []string{"*"},
					APIVersions: []string{"*"},
					Resources:   []string{"resourcequotas"},
					Scope:       &scope,
				},
			}},
		}},
	}
	// Set CR instance as the owner and controller
	err := controllerutil.SetControllerReference(cr, expectedRes, r.scheme)
	if err != nil {
		reqLogger.Error(err, "Failed to define expected resource")
		return reconcile.Result{}, err
	}

	// If MutatingWebhookConfig does not exist, create it and requeue
	foundWebhookConfig := &admissionv1beta1.MutatingWebhookConfiguration{}
	// Note: mutatingwebhookconfigurations are cluster-scoped, so this does not search using namespace (unlike other resources above)
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: expectedRes.Name}, foundWebhookConfig)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new MutatingWebhookConfig", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
		err = r.client.Create(context.TODO(), expectedRes)
		if err != nil && errors.IsAlreadyExists(err) {
			// Already exists from previous reconcile, requeue.
			return reconcile.Result{Requeue: true}, nil
		} else if err != nil {
			reqLogger.Error(err, "Failed to create new MutatingWebhookConfig", "Namespace", cr.Spec.InstanceNamespace, "Name", expectedRes.Name)
			return reconcile.Result{}, err
		}
		// Created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get MutatingWebhookConfig")
		return reconcile.Result{}, err
	} else if len(foundWebhookConfig.Webhooks) != len(expectedRes.Webhooks) ||
		!reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.Service.Name, expectedRes.Webhooks[0].ClientConfig.Service.Name) ||
		!reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.Service.Namespace, expectedRes.Webhooks[0].ClientConfig.Service.Namespace) ||
		!reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.Service.Path, expectedRes.Webhooks[0].ClientConfig.Service.Path) ||
		!reflect.DeepEqual(foundWebhookConfig.Webhooks[0].Rules, expectedRes.Webhooks[0].Rules) {
		// Spec is incorrect, update it and requeue
		reqLogger.Info("Found MutatingWebhookConfig is incorrect", "Found", foundWebhookConfig.Webhooks, "Expected", expectedRes.Webhooks)
		foundWebhookConfig.Webhooks = expectedRes.Webhooks
		err = r.client.Update(context.TODO(), foundWebhookConfig)
		if err != nil {
			reqLogger.Error(err, "Failed to update MutatingWebhookConfig", "Name", foundWebhookConfig.Name)
			return reconcile.Result{}, err
		}
		// Deleted - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// No reconcile was necessary
	return reconcile.Result{}, nil
}
