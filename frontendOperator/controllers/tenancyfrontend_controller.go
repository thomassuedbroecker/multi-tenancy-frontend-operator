/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	// Add to read error from Kubernetes
	"k8s.io/apimachinery/pkg/api/errors"

	// Add to read deployments from Kubernetes
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	// Add to define the own deployment 'yaml' configuration
	"github.com/thomassuedbroecker/multi-tenancy-frontend-operator/api/v1alpha1"
	multitenancyv1alpha1 "github.com/thomassuedbroecker/multi-tenancy-frontend-operator/api/v1alpha1"
	"github.com/thomassuedbroecker/multi-tenancy-frontend-operator/helpers"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"fmt" // Basic functionalities

	"k8s.io/apimachinery/pkg/util/intstr" // Because of the cluster service target port definition

	// Get deployments
	"k8s.io/client-go/kubernetes"
)

// TenancyFrontendReconciler reconciles a TenancyFrontend object
type TenancyFrontendReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var customLogger bool = true

//+kubebuilder:rbac:groups=multitenancy.example.net,resources=tenancyfrontends,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multitenancy.example.net,resources=tenancyfrontends/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multitenancy.example.net,resources=tenancyfrontends/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TenancyFrontend object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *TenancyFrontendReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Verify if a CRD of TenancyFrontend exists")
	tenancyfrontend := &multitenancyv1alpha1.TenancyFrontend{}

	// Get objects inside the Kubernetes namespace
	namespace := tenancyfrontend.Namespace
	kind := "TenancyFrontend"
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	// "Verify if a CR of TenancyFrontend exists"
	err := r.Get(ctx, req.NamespacedName, tenancyfrontend)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("TenancyFrontend resource not found. Ignoring since object must be deleted")

			helpers.CustomLogs(err.Error(), ctx, customLogger)
			// delete secret "appid.client-id-frontend"
			targetSecretName := "appid.client-id-frontend"
			verifysecret, errsec := VerifySecretTenancyFrontend(clientset, ctx, namespace, targetSecretName)
			if errsec != nil {
				return ctrl.Result{}, errsec
			}

			helpers.CustomLogs("Try to delete secret", ctx, customLogger)
			if verifysecret != nil {
				errsec = r.Delete(context.TODO(), verifysecret, &client.DeleteOptions{})
				if errsec != nil {
					return ctrl.Result{}, errsec
				}
			}

			// delete secret "appid.discovery-endpoint"
			targetSecretName = "appid.discovery-endpoint"
			verifysecret, errsec = VerifySecretTenancyFrontend(clientset, ctx, namespace, targetSecretName)
			if errsec != nil {
				return ctrl.Result{}, errsec
			}

			helpers.CustomLogs("Try to delete secret", ctx, customLogger)
			if verifysecret != nil {

				errsec = r.Delete(context.TODO(), verifysecret, &client.DeleteOptions{})
				if errsec != nil {
					return ctrl.Result{}, errsec
				}
			}

			// delete services
			verifyservice, errsec := VerifyServiceTenancyFrontend(clientset, ctx, namespace)
			if errsec != nil {
				return ctrl.Result{}, errsec
			}

			helpers.CustomLogs("Try to delete service", ctx, customLogger)
			if verifyservice != nil {

				errsec = r.Delete(context.TODO(), verifyservice, &client.DeleteOptions{})
				if errsec != nil {
					return ctrl.Result{}, errsec
				}
			}

			// delete services
			verifyservice, errsec = VerifyServiceTenancyFrontend(clientset, ctx, namespace)
			if errsec != nil {
				return ctrl.Result{}, errsec
			}

			helpers.CustomLogs("Try to delete service", ctx, customLogger)
			if verifyservice != nil {

				errsec = r.Delete(context.TODO(), verifyservice, &client.DeleteOptions{})
				if errsec != nil {
					return ctrl.Result{}, errsec
				}
			}

			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get TenancyFrontend")
		return ctrl.Result{}, err
	}

	// "Verify if TenancyFrontend deployment exists"
	deployment_exists, err := VerifyDeploymentExists(clientset, ctx, namespace, kind)

	if err != nil {
		return ctrl.Result{}, nil
	}

	if !deployment_exists {
		// Check if the deployment already exists, if not create a new one
		logger.Info("Verify if the deployment already exists, if not create a new one")

		found := &appsv1.Deployment{}
		err = r.Get(ctx, types.NamespacedName{Name: tenancyfrontend.Name, Namespace: tenancyfrontend.Namespace}, found)

		if err != nil && errors.IsNotFound(err) {

			// Define a new deployment
			dep := r.deploymentForTenancyFronted(tenancyfrontend, ctx)
			logger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			err = r.Create(ctx, dep)
			if err != nil {
				logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}

			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil

		} else if err != nil {

			logger.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err

		}
	}

	if deployment_exists {
		//*****************************************
		// Define service NodePort
		servPort := &corev1.Service{}
		helpers.CustomLogs("Define service NodePort", ctx, customLogger)

		//*****************************************
		// Create service NodePort
		helpers.CustomLogs("Create service NodePort", ctx, customLogger)

		targetServPort, err := defineServiceNodePort(tenancyfrontend.Name, tenancyfrontend.Namespace)

		// Error creating replicating the secret - requeue the request.
		if err != nil {
			return ctrl.Result{}, err
		}

		err = r.Get(context.TODO(), types.NamespacedName{Name: targetServPort.Name, Namespace: targetServPort.Namespace}, servPort)
		if err != nil && errors.IsNotFound(err) {
			logger.Info(fmt.Sprintf("Target service port %s doesn't exist, creating it", targetServPort.Name))
			err = r.Create(context.TODO(), targetServPort)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {
			logger.Info(fmt.Sprintf("Target service port %s exists, updating it now", targetServPort))
			err = r.Update(context.TODO(), targetServPort)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		//*****************************************
		// Define cluster
		servClust := &corev1.Service{}

		//*****************************************
		// Create service cluster
		helpers.CustomLogs("Create service Cluster IP", ctx, customLogger)

		targetServClust, err := defineServiceClust(tenancyfrontend.Name, tenancyfrontend.Namespace)

		// Error creating replicating the service cluster - requeue the request.
		if err != nil {
			return ctrl.Result{}, err
		}

		err = r.Get(context.TODO(), types.NamespacedName{Name: targetServClust.Name, Namespace: targetServClust.Namespace}, servClust)

		if err != nil && errors.IsNotFound(err) {
			logger.Info(fmt.Sprintf("Target service cluster %s doesn't exist, creating it", targetServClust.Name))
			err = r.Create(context.TODO(), targetServClust)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {
			logger.Info(fmt.Sprintf("Target service cluster %s exists, updating it now", targetServClust))
			err = r.Update(context.TODO(), targetServClust)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		//*****************************************
		// Define secret
		helpers.CustomLogs("Define secret", ctx, customLogger)
		secret := &corev1.Secret{}

		createUpdateSecrect()
		deleteSecrect()

		//*****************************************
		// Create secret appid.client-id-frontend
		helpers.CustomLogs("Create secret appid.client-id-frontend", ctx, customLogger)

		targetSecretName := "appid.client-id-frontend"
		clientId := "b12a05c3-8164-45d9-a1b8-af1dedf8ccc3"

		verifysecret, err := VerifySecretDeployment(clientset, ctx, namespace, tenancyfrontend.Name, targetSecretName)

		if err != nil {
			return ctrl.Result{}, err
		}

		if !verifysecret {
			targetSecret, err := defineSecret(targetSecretName, tenancyfrontend.Namespace, "VUE_APPID_CLIENT_ID", clientId, tenancyfrontend.Name)
			// Error creating replicating the secret - requeue the request.
			if err != nil {
				return ctrl.Result{}, err
			}

			err = r.Get(context.TODO(), types.NamespacedName{Name: targetSecret.Name, Namespace: targetSecret.Namespace}, secret)
			secretErr := verifySecrectStatus(ctx, r, targetSecretName, targetSecret, err)
			if secretErr != nil && errors.IsNotFound(secretErr) {
				return ctrl.Result{}, secretErr
			}
		}

		//*****************************************
		// Create secret appid.discovery-endpoint
		targetSecretName = "appid.discovery-endpoint"
		discoveryEndpoint := "https://eu-de.appid.cloud.ibm.com/oauth/v4/3793e3f8-ed31-42c9-9294-bc415fc58ab7/.well-known/openid-configuration"

		verifysecret, err = VerifySecretDeployment(clientset, ctx, namespace, tenancyfrontend.Name, targetSecretName)

		if err != nil {
			return ctrl.Result{}, err
		}

		if !verifysecret {

			targetSecret, err := defineSecret(targetSecretName, tenancyfrontend.Namespace, "VUE_APPID_DISCOVERYENDPOINT", discoveryEndpoint, tenancyfrontend.Name)
			// Error creating replicating the secret - requeue the request.
			if err != nil {
				return ctrl.Result{}, err
			}

			err = r.Get(context.TODO(), types.NamespacedName{Name: targetSecret.Name, Namespace: targetSecret.Namespace}, secret)
			secretErr := verifySecrectStatus(ctx, r, targetSecretName, targetSecret, err)
			if secretErr != nil && errors.IsNotFound(secretErr) {
				return ctrl.Result{}, secretErr
			}
		}
	}

	logger.Info("Just return nil")
	return ctrl.Result{}, nil
}

// deploymentForTenancyFronted returns a tenancyfrontend Deployment object
func (r *TenancyFrontendReconciler) deploymentForTenancyFronted(frontend *v1alpha1.TenancyFrontend, ctx context.Context) *appsv1.Deployment {
	logger := log.FromContext(ctx)
	ls := labelsForTenancyFrontend(frontend.Name, frontend.Name)
	replicas := frontend.Spec.Size

	// Just reflect the command in the deployment.yaml
	// for the ReadinessProbe and LivenessProbe
	// command: ["sh", "-c", "curl -s http://localhost:8080"]
	mycommand := make([]string, 3)
	mycommand[0] = "/bin/sh"
	mycommand[1] = "-c"
	mycommand[2] = "curl -s http://localhost:8080"

	// Using the context to log information
	logger.Info("Logging: Creating a new Deployment", "Replicas", replicas)
	message := "Logging: (Name: " + frontend.Name + ") \n"
	logger.Info(message)
	message = "Logging: (Namespace: " + frontend.Namespace + ") \n"
	logger.Info(message)

	for key, value := range ls {
		message = "Logging: (Key: [" + key + "] Value: [" + value + "]) \n"
		logger.Info(message)
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontend.Name,
			Namespace: frontend.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "quay.io/tsuedbroecker/service-frontend:latest",
						Name:  "service-frontend",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "nginx-port",
						}},
						Env: []corev1.EnvVar{{
							Name: "VUE_APPID_DISCOVERYENDPOINT",
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "appid.discovery-endpoint",
									},
									Key: "VUE_APPID_DISCOVERYENDPOINT",
								},
							}},
							{Name: "VUE_APPID_CLIENT_ID",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "appid.client-id-frontend",
										},
										Key: "VUE_APPID_CLIENT_ID",
									},
								}},
							{Name: "VUE_APP_API_URL_CATEGORIES",
								Value: "VUE_APP_API_URL_CATEGORIES_VALUE",
							},
							{Name: "VUE_APP_API_URL_PRODUCTS",
								Value: "VUE_APP_API_URL_PRODUCTS_VALUE",
							},
							{Name: "VUE_APP_API_URL_ORDERS",
								Value: "VUE_APP_API_URL_ORDERS_VALUE",
							},
							{Name: "VUE_APP_CATEGORY_NAME",
								Value: "VUE_APP_CATEGORY_NAME_VALUE",
							},
							{Name: "VUE_APP_HEADLINE",
								Value: frontend.Spec.DisplayName,
							},
							{Name: "VUE_APP_ROOT",
								Value: "/",
							}}, // End of Env listed values and Env definition
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{Command: mycommand},
							},
							InitialDelaySeconds: 20,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{Command: mycommand},
							},
							InitialDelaySeconds: 20,
						},
					}}, // Container
				}, // PodSec
			}, // PodTemplateSpec
		}, // Spec
	} // Deployment

	// Set TenancyFrontend instance as the owner and controller
	ctrl.SetControllerReference(frontend, dep, r.Scheme)
	return dep
}

// labelsForTenancyFrontend returns the labels for selecting the resources
// belonging to the given tenancyfrontend CR name.
func labelsForTenancyFrontend(name_app string, name_cr string) map[string]string {
	return map[string]string{"app": name_app, "tenancyfrontend_cr": name_cr}
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenancyFrontendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multitenancyv1alpha1.TenancyFrontend{}).
		Complete(r)
}

// ********************************************************
// additional functions

// Create Secret definition
func defineSecret(name string, namespace string, key string, value string, deploymentname string) (*corev1.Secret, error) {
	secretdata := make(map[string]string)
	secretdata[key] = value

	// Define map for the labels
	mlabel := make(map[string]string)
	key = "deployment"
	value = deploymentname
	mlabel[key] = value
	key = "tenancyfrontend"
	value = "yes"
	mlabel[key] = value

	return &corev1.Secret{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: mlabel},
		Immutable:  new(bool),
		Data:       map[string][]byte{},
		StringData: secretdata,
		Type:       "Opaque",
	}, nil
}

// Create Service NodePort definition
func defineServiceNodePort(name string, namespace string) (*corev1.Service, error) {
	// Define map for the selector
	mselector := make(map[string]string)
	key := "app"
	value := name
	mselector[key] = value

	// Define map for the labels
	mlabel := make(map[string]string)
	key = "app"
	value = "service-frontend"
	mlabel[key] = value

	var port int32 = 8080

	return &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: mlabel},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{{
				Port: port,
				Name: "http",
			}},
			Selector: mselector,
		},
	}, nil
}

// Create Service ClusterIP definition
func defineServiceClust(name string, namespace string) (*corev1.Service, error) {
	// Define map for the selector
	mselector := make(map[string]string)
	key := "app"
	value := name
	mselector[key] = value

	// Define map for the labels
	mlabel := make(map[string]string)
	key = "app"
	value = "service-frontend"
	mlabel[key] = value

	var port int32 = 80
	var targetPort int32 = 8080
	var clustserv = name + "clusterip"

	return &corev1.Service{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "Service"},
		ObjectMeta: metav1.ObjectMeta{Name: clustserv, Namespace: namespace, Labels: mlabel},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{{
				Port:       port,
				TargetPort: intstr.IntOrString{IntVal: targetPort},
			}},
			Selector: mselector,
		},
	}, nil
}

// Do all the tests for the status
func verifySecrectStatus(ctx context.Context, r *TenancyFrontendReconciler, targetSecretName string, targetSecret *v1.Secret, err error) error {
	logger := log.FromContext(ctx)

	if err != nil && errors.IsNotFound(err) {
		logger.Info(fmt.Sprintf("Target secret %s doesn't exist, creating it", targetSecretName))
		err = r.Create(context.TODO(), targetSecret)
		if err != nil {
			return err
		}
	} else {
		logger.Info(fmt.Sprintf("Target secret %s exists, updating it now", targetSecretName))
		err = r.Update(context.TODO(), targetSecret)
		if err != nil {
			return err
		}
	}

	return err
}

// How to handle secrets and configmap for multible deployments?

// 1. What is needed to map a secret to a deployment?
//    -> Using a label with the deploymentname in the secret metadata?
//
// 2. How does a application instance know which secret must be used for it's own deployment?
//
//    -> Include the name of the deployment to the secret name?
//       -> But, how to refect this for the variables for the used application container
//          -> Maybe only one application instance is permitted within one namespace?

// Restrictions:
// - Only one deployment per namespace

func createUpdateSecrect() error {

	// 1. verify if the deployment exists
	// 2. verify if a secrect for the deployment exists
	// 3. if no secret exists prepare a secret for the deployment and create a secret
	// 4. if a secret exists update the secret

	return nil
}

func deleteSecrect() error {

	// 1. get the secret
	// 2. verify does a deployment for that secret exists
	// 3. if no deployment for the secret exists delete the secret

	return nil
}

// Verify if a secret for the TenancyFrontend exists

func VerifySecretTenancyFrontend(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string, secretname string) (*v1.Secret, error) {

	list, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	secret_items := list.Items

	if err != nil {
		helpers.CustomLogs(err.Error(), ctx, customLogger)
		return nil, err
	} else {
		for _, item := range secret_items {
			mlabel := item.Labels
			info := "Label app: [" + mlabel["tenancyfrontend"] + "] Name: [" + item.Name + "]"
			helpers.CustomLogs(info, ctx, customLogger)
			if (secretname == item.Name) && (mlabel["tenancyfrontend"] == "yes") {
				return &item, err
			}
		}
	}

	return nil, err
}

// Verify if a secret for the deployment exists

func VerifySecretDeployment(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string, deploymentname string, secretname string) (bool, error) {

	list, err := clientset.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	secret_items := list.Items

	if err != nil {
		helpers.CustomLogs(err.Error(), ctx, customLogger)
		return false, err
	} else {
		for _, item := range secret_items {
			mlabel := item.Labels
			info := "Label app: [" + mlabel["deployment"] + "] Name: [" + item.Name + "]"
			helpers.CustomLogs(info, ctx, customLogger)
			if (secretname == item.Name) && (mlabel["deployment"] == deploymentname) {
				return true, err
			}
		}
	}

	return false, err
}

// Verify if q deployment exist in namespace

func VerifyDeploymentExists(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string, kind string) (bool, error) {

	list, err := clientset.AppsV1().Deployments(namespace).
		List(ctx, metav1.ListOptions{})
	deployment_items := list.Items

	if err != nil {
		helpers.CustomLogs(err.Error(), ctx, customLogger)
		return false, err
	} else {
		for _, item := range deployment_items {
			mlabel := item.Labels
			info := "Label app: [" + mlabel["app"] + "] Name: [" + item.Name + "] Kind: [" + item.Kind + "]"
			helpers.CustomLogs(info, ctx, customLogger)
			if item.OwnerReferences != nil {
				if kind == item.OwnerReferences[0].Kind {
					return true, err
				}
			}
		}
	}
	return false, nil
}

// Verify if a service for the TenancyFrontend exists

func VerifyServiceTenancyFrontend(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string) (*v1.Service, error) {

	list, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	services_items := list.Items

	if err != nil {
		helpers.CustomLogs(err.Error(), ctx, customLogger)
		return nil, err
	} else {
		for _, item := range services_items {
			mlabel := item.Labels
			info := "Label app: [" + mlabel["app"] + "] Name: [" + item.Name + "]"
			helpers.CustomLogs(info, ctx, customLogger)
			if mlabel["app"] == "service-frontend" {
				return &item, err
			}
		}
	}

	return nil, err
}
