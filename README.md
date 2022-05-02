# Example "Multi Tenancy Frontend Operator"

The operator simply deploys a frontend application which runs on a nginx server. Therefor the operator creates following Kubernetes resources:

- A deployment
- A service
- Some secrets

The source code of the example frontend application is available in the open sourced GitHub project [multi-tenancy-frontend](https://github.com/IBM/multi-tenancy-frontend).

## The project is related to following blog posts:

### LOCAL DEVELOPMENT

Get the operator running on a local machine and using Kubernetes.
I verified my example project (Multi Tenancy Frontend Operator) also for these topics and it worked for me.

* [Install the Operator SDK on macOS](https://suedbroecker.net/2022/02/15/fata0009-failed-to-create-api-unable-to-run-post-scaffold-tasks-of-base-go-kubebuilder-io-v3-exit-status-2/)
* [Develop a simple operator locally](https://suedbroecker.net/2022/02/18/start-to-develop-a-simple-operator-to-deploy-the-frontend-application-of-the-open-source-multi-cloud-asset-to-build-saas%c2%b6/)
* [Debug a GO operator](https://suedbroecker.net/2022/03/01/debug-a-kubernetes-operator-written-in-go/)
* [Run an operator as a deployment](https://suedbroecker.net/2022/03/15/run-an-operator-as-a-deployment/)

### USING THE OPERATOR LIFECYCLE MANAGER (OLM)

I verified my example project (Multi Tenancy Frontend Operator) also for these topics and it worked for me:

* [How to create a bundle?](https://suedbroecker.net/2022/03/16/how-to-create-a-bundle-for-an-operator/)
* [Run the bundle with an Operator Lifecycle Manager (OLM)](https://suedbroecker.net/2022/03/16/run-an-operator-using-a-bundle-with-an-operator-lifecycle-manager-olm/)
* [Deploy an operator without the Operator SDK](https://suedbroecker.net/2022/03/22/deploy-an-operator-without-the-operator-sdk/)
* [Add a new API version to an existing operator](https://suedbroecker.net/2022/03/24/add-a-new-api-version-to-an-existing-operator/)
* [Add a conversion webhook to an operator to convert API versions](https://suedbroecker.net/2022/03/29/add-a-conversion-webhook-to-an-operator-to-convert-api-versions/)

### MONITORING

* [Monitor your custom operator with Prometheus](https://wp.me/paelj4-1iv)
* [Example for an installation and an initial configuration of the Grafana operator](https://wp.me/paelj4-1ld)
* [Access Prometheus queries using the Prometheus HTTP API](https://wp.me/paelj4-1kb)
