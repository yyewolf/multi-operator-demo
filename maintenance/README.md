# Maintenance

The Maintenance operator is responsible for redirecting traffic to a maintenance page when the end-user chooses to. A spec might look like this:

```yaml
apiVersion: maintenance.multi.ch/v1
kind: Maintenance
metadata:
  name: maintenance-sample
spec:
  replaces:
    apiVersion: app.multi.ch/v1
    kind: App
    name: app-sample
```

## Design

The Maintenance operator only creates a [Backend](https://gateway.envoyproxy.io/docs/api/extension_types/#backend) object that redirects traffic to a global maintenance service.

The `replaces` field is not used explicitely, it is merely informative to allow one to apply the following logic:
- Initially:
  - Create an App
  - Create a Route
- After a while, to put the App in maintenance:
  - Create a Maintenance for the App
  - Update the Route to point to the Maintenance service
- After a while, to put the App back in production:
  - Get the Maintenance (to get the `replaces` field)
  - In the Route, replace the Maintenance service with the App service
  - Delete the Maintenance

Further designing could add the following features (not implemented in this demo):
- Custom maintenance page per user
- Update the Route automatically when the maintenance page is enabled
- Many more things...

## Contract

Maintenance is meant to be *routeable* and exposed to the world via HTTP. As such, it implements Route's contract which looks like this : 

```yaml
status:
  routeContract:
    backendRef:
    name: maintenance-sample
    port: 80
```
