# Route Operator

The Route operator is responsible for abstracting the routing logic. It is used to register HTTPRoute based on targets, a spec might look like this:

```yaml
apiVersion: route.multi.ch/v1
kind: Route
metadata:
  name: route-sample
spec:
  hostnames:
    - test.yewolf.fr
  targetRefs:
    - kind: App
      name: app-sample-1
      pathPrefix: /app-1
    - kind: App
      name: app-sample-2
      pathPrefix: /app-2
```

## Design

Route is meant to not import any other operator, it should not know about the types of its possible targets. The only requirement for a target is to implement the `routeContract` in its status. This contract is used to generate the HTTPRoute.

The entity responsible for creating the Route is also not expected to know about the target's version. The Route operator, through a webhook, will default them to the preferred version of the cluster.
