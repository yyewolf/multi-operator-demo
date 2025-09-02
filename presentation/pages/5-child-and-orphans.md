---
layout: two-cols-header
---

# ⚡ Our Challenges

::left::

<div class="!-mt-32">

### 🔄 Resource Tracking
- Keep track of created resources  
- Update when **child objects** change  

<div class="p-4"></div>

### 🔗 Dependency Management
- Detect changes in **dependent objects**  
- Ensure reconciliation reacts correctly  

</div>

::right::

<div class="!-mt-32">

### 🧹 Cleanup & Finalization
- Run **cleanup tasks** before deletion  
- Handle finalizers safely  

<div class="p-4"></div>

### 🚀 Beyond Reconciliation
- Perform **actions outside the reconcile loop**  
- Start / Stop / Refresh configuration ?

</div>

---

<!-- Slide 7: Child Resources vs. Dependencies -->
## Managing Child Resources vs Dependencies

In a typical operator:

- **Child resources** (like Deployments, Services, PVCs...)  
  are created by the controller and marked with an **OwnerReference**.

🌀 This means:
- Changes or deletions trigger a **reconciliation**
- Kubernetes **automatically garbage-collects** them when the parent is deleted

---

Example :

<div class="grid grid-cols-2 gap-8 items-center">

<div>

<div class="p-5"></div>

```yaml {all|11-20}
apiVersion: app.multi.ch/v1
kind: App
metadata:
  name: app-sample
  generation: 2
spec:
  command: |
    /bin/bash -c "flask --app=/data/www/apps/app-2/app.py run --host=0.0.0.0"
  port: 5000
status:
  childResources:
   - apiVersion: apps/v1
      group: apps
      kind: Deployment
      name: app-sample
      namespace: default
      observedGeneration: 2
      status: "True"
      transitionTime: "2025-08-26T13:51:56Z"
      uid: e7252a5b-3f8d-4c9e-a7ec-f8befdd6ef75
```

</div>

<div class="flex flex-col space-y-4">
<div class="p-5"></div>

```yaml {all|7-13}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-sample
  namespace: default
  generation: 1
  ownerReferences:
   - apiVersion: app.multi.ch/v1
      blockOwnerDeletion: true
      controller: true
      kind: App
      name: app-sample
      uid: b27d237f-4e22-479e-8853-94b8d27cf224
  uid: e7252a5b-3f8d-4c9e-a7ec-f8befdd6ef75
```

</div>

</div>

---

<!-- Slide 8: Handling Dependencies -->
## Handling Dependencies (the Other Way Around)

But what about when our resource depends on another one?

To track and react to **external dependencies**:

1. We add a **Label** to the *dependency* that's pointing **back to the referer**
   - That label is a list of all referers (by their type, name and namespace).

Why?

- So that dependency changes trigger reconciliation

What about cleanup ?

- Using a **finalizer** on the referer, we can make sure to clean that Label of the reference back.

---

Example :

<div class="grid grid-cols-2 gap-8 items-center">

<div>

```yaml {all|15-24}
apiVersion: route.multi.ch/v1
kind: Route
metadata:
  generation: 2
  name: route-sample
  namespace: default
spec:
  hostnames:
  - 172.19.0.4
  targetRefs:
  - apiVersion: app.multi.ch/v1
    kind: App
    name: app-sample
    pathPrefix: /
  dependencies:
  - apiVersion: app.multi.ch/v1
    group: app.multi.ch
    kind: App
    name: app-sample
    namespace: default
    observedGeneration: 2
    status: "True"
    transitionTime: "2025-08-26T13:52:13Z"
```

</div>

<div>

```yaml {all|5-16}
apiVersion: app.multi.ch/v1
kind: App
metadata:
  annotations:
    multi.ch/managed-by: |
      [
         {
            "name": "route-sample",
            "namespace": "default",
            "gvk": {
               "Group": "route.multi.ch",
               "Version": "v1",
               "Kind": "Route"
            }
         }
      ]
  name: app-sample
  generation: 2
spec:
  command: |
    /bin/bash -c "flask [...]
  port: 5000
```

</div>

</div>

---

<!-- Slide 9: Finalizer Flow -->
## Finalizer Flow

🧼 Finalizer ensures a clean delete lifecycle:

1. User deletes the CR  
2. Kubernetes sets `deletionTimestamp` but keeps our resource alive  
3. Our controller:
   - Removes the reference stored in the label `managed-by` from dependent resources
   - Then removes the **finalizer**  
4. Now Kubernetes deletes the CR

✅ No orphaned labels on dependent object that should not be deleted  
✅ No reconciliation triggered on deleted parents (that are still being refered to)
