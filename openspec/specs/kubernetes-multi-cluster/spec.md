# kubernetes-multi-cluster Specification

## Purpose
TBD - created by archiving change add-kubernetes-multi-cluster. Update Purpose after archive.
## Requirements
### Requirement: Kubernetes integration storage

The system SHALL store cluster credentials as encrypted JSON on `integrations` rows with canonical kind `kubernetes`, containing a `kubeconfig` string (YAML).

#### Scenario: Config present

- **WHEN** an operator saves a kubernetes integration with non-empty `kubeconfig`
- **THEN** the configuration is encrypted at rest and not returned on public integration APIs

### Requirement: Cluster resource APIs

The system SHALL expose authenticated APIs under `/api/v1/k8s/clusters/:id` to list namespaces, normalized workloads (Deployment, StatefulSet, DaemonSet), pods with status metadata, pod logs, events, and nodes for the integration id.

#### Scenario: Cluster unreachable

- **WHEN** the Kubernetes API is unreachable or credentials are invalid
- **THEN** the API returns a non-2xx response with an error message and does not panic

### Requirement: RBAC

The system SHALL require permission `kubernetes:*` for all `/api/v1/k8s` routes and expose menu entries only to authorized users.

#### Scenario: Denied user

- **WHEN** a user lacks `kubernetes:*`
- **THEN** kubernetes UI routes and APIs are denied consistent with other integrations

