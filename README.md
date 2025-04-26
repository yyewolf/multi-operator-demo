# When One Operator Isn’t Enough: Building Composable Kubernetes Controllers

This repository contains the code for the talk "When One Operator Isn’t Enough: Building Composable Kubernetes Controllers".

## Testing

In order to run the material presented, you need to have the following tools installed:
- [Docker](https://www.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kustomize](https://kustomize.io/)
- [kind](https://kind.sigs.k8s.io/)
- [helm](https://helm.sh/)

A Makefile is provided to help you with the setup. Running `make deploy` does the following:
- Create a kind cluster named "multi-operator-demo"
- Install the prerequisites for the demo
  - Gateway API CRDs
  - Envoy Gateway
  - Maintenance Service (for the Maintenance operator)
  - Cert Manager
  - Workload image (for the App operator)
- Install the operators
  - Maintenance operator
  - App operator
  - Route operator

This step can take a few minutes because it needs to build the images and push them to the kind cluster.

## Operators

Each operator contains a README file explaining how each of them works. The operators are:
- [App Operator](./app)
- [Maintenance Operator](./maintenance)
- [Route Operator](./route)

## Library

Each operator uses the [library](./library) to help with the development of the operators. It is where most of the logic is implemented. Using generics, it is possible to create a single function that can be used by all the operators to, for example, reconcile any resource type.

## Presentation

The presentation slides are available in the [presentation](./presentation) folder. The slides are in Markdown format.

To help with the presentation, the [instrumentation](./instrumentation) folder contains a python script that starts a web server containing a control panel, it is used to switch between different states of the demo. The script is not part of the demo and is not required to run it. It is only used to help with the presentation.
