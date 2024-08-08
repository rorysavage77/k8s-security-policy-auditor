# Kubernetes Security Policy Auditor


`k8s-security-policy-auditor` is a Kubernetes controller designed to audit security policies and configurations within a Kubernetes cluster. It helps ensure that your cluster resources, such as ConfigMaps, Secrets, Roles, and RoleBindings, adhere to best security practices.

## Features

- Audits Kubernetes ConfigMaps for sensitive data.
- Audits Secrets to ensure encryption and minimal exposure.
- Checks Roles and RoleBindings for excessive permissions and role misconfigurations.
- Provides insights and logs for detected security issues.

## Getting Started

### Prerequisites

- Kubernetes cluster (v1.20 or later recommended)
- kubectl command-line tool
- Docker for building the container image
- GitHub account for using GitHub Container Registry (GHCR)

### Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/rorysavage77/k8s-security-policy-auditor.git
cd k8s-security-policy-auditor
```

#### 2. Build and Push Docker Image

Ensure you have the Dockerfile set up and use GitHub Actions for CI/CD:

GitHub Actions Workflow: Ensure you have the workflow set up in .github/workflows/go.yml to build and push your image to GHCR.

#### 3. Deploy to Kubernetes

1. Create a Docker Image Pull Secret:

If your image is private, create a pull secret:
```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<your-username> \
  --docker-password=<your-token> \
  --docker-email=<your-email>
```
2. Apply the RBAC Configuration:

Ensure your service account has the necessary permissions by applying the RBAC manifest:
```bash
kubectl apply -f rbac-config.yaml
```

3. Deploy the Auditor:

Deploy the auditor to your Kubernetes cluster:

```bash
kubectl apply -f k8s-security-policy-auditor.yaml
``````

#### Usage
Once deployed, the auditor will start monitoring your cluster's ConfigMaps, Secrets, Roles, and RoleBindings for security issues. It logs findings to the standard output, which you can view using:

```bash
kubectl logs -f deployment/k8s-security-policy-auditor
```

#### Configuration
You can customize the auditor by modifying environment variables and resource configurations within the k8s-security-policy-auditor.yaml manifest.

#### Run Tests:

Ensure you have the necessary Go tools installed:

```bash
go test -v ./...
```


#### Build Locally:

Build the binary:

```bash
go build -o k8s-security-policy-auditor .
``````

#### Example logs

```shell
2024-08-07T21:56:59Z	INFO	Auditing RoleBinding	{"Namespace": "litmus", "Name": "litmus-admin-ops-role-binding"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "kyverno", "Request.Name": "sh.helm.release.v1.kyverno.v1", "Namespace": "kyverno", "Name": "sh.helm.release.v1.kyverno.v1"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "kyverno", "Name": "sh.helm.release.v1.kyverno.v1"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "litmus", "Request.Name": "chaos-litmus-admin-secret", "Namespace": "litmus", "Name": "chaos-litmus-admin-secret"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "litmus", "Name": "chaos-litmus-admin-secret"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "litmus", "Request.Name": "chaos-mongodb", "Namespace": "litmus", "Name": "chaos-mongodb"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "litmus", "Name": "chaos-mongodb"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "default", "Request.Name": "ghcr-secret", "Namespace": "default", "Name": "ghcr-secret"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "default", "Name": "ghcr-secret"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "kyverno", "Request.Name": "kyverno-cleanup-controller.kyverno.svc.kyverno-tls-pair", "Namespace": "kyverno", "Name": "kyverno-cleanup-controller.kyverno.svc.kyverno-tls-pair"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "kyverno", "Name": "kyverno-cleanup-controller.kyverno.svc.kyverno-tls-pair"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "kyverno", "Request.Name": "kyverno-svc.kyverno.svc.kyverno-tls-ca", "Namespace": "kyverno", "Name": "kyverno-svc.kyverno.svc.kyverno-tls-ca"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "kyverno", "Name": "kyverno-svc.kyverno.svc.kyverno-tls-ca"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "kyverno", "Request.Name": "kyverno-svc.kyverno.svc.kyverno-tls-pair", "Namespace": "kyverno", "Name": "kyverno-svc.kyverno.svc.kyverno-tls-pair"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "kyverno", "Name": "kyverno-svc.kyverno.svc.kyverno-tls-pair"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "kyverno", "Request.Name": "kyverno-cleanup-controller.kyverno.svc.kyverno-tls-ca", "Namespace": "kyverno", "Name": "kyverno-cleanup-controller.kyverno.svc.kyverno-tls-ca"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "kyverno", "Name": "kyverno-cleanup-controller.kyverno.svc.kyverno-tls-ca"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "litmus", "Request.Name": "sh.helm.release.v1.chaos.v1", "Namespace": "litmus", "Name": "sh.helm.release.v1.chaos.v1"}
2024-08-07T21:56:59Z	INFO	Auditing Secret	{"Namespace": "litmus", "Name": "sh.helm.release.v1.chaos.v1"}
2024-08-07T21:56:59Z	INFO	Reconciling resource	{"Request.Namespace": "litmus", "Request.Name": "subscriber-secret", "Namespace": "litmus", "Name": "subscriber-secret"}
```


## Contributing
Contributions are welcome! Please open an issue or submit a pull request with your changes. Ensure that your code is well-documented and includes tests where applicable.

## License
This project is licensed under the MIT License. See the LICENSE file for more information.

## Acknowledgments
Special thanks to the Kubernetes community for providing an extensive set of tools and libraries that make projects like this possible.

