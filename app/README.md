# App Operator

The App operator is responsible for deploying an end-user application and its configuration. The design here is minimal and can look like this :

```yaml
apiVersion: app.multi.ch/v1
kind: App
metadata:
  name: app-sample
spec:
  port: 8080
  command: |
    /bin/bash -c "python3 -m http.server 8080"
```

## Design

The process of the user is handled by supervisor, which is a process manager for Linux. It is responsible for starting the application and monitoring its health. The operator uses the `supervisor` image to run the application. 

For the purpose of this demo, the `workload` image already contains two example applications, this is not really necessary as you could give access to the container to the end-user for him to build his app.

The [Dockerfile](./docker/Dockerfile) here runs as root, which is not a good practice either.

Supervisor also offers an http server to remotely control the process. This *could* be used to check health, but in this demo we are only using the `start`, `stop` and `restart` commands. To do this, we add a APIService to Kubernetes in order to expose custom endpoints to the operator. The APIService can be called with the following command:

```bash
$ kubectl get --raw="/apis/agent.app.multi.ch/v1/namespaces/<NS>/apps/<NAME>/<ACTION>"
```

This demonstration does not take into account the security standpoint of our implementation and ignores other problems such as :
- How to change the technology (python) of the application
- How to handle updating the runtime
- How to make sure the end-user cannot break anything
- How to handle the lifecycle of the application
- Resource Requests/Limits
- Probably much more...

## Contract

App is meant to be *routeable* and exposed to the world via HTTP. As such, it implements Route's contract which looks like this : 

```yaml
status:
  routeContract:
    serviceRef:
      name: app-sample
      port: 80
```
