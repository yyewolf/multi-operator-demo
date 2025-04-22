<!-- Slide 7: Child Resources vs. Dependencies -->
## Managing Child Resources vs Dependencies

In a typical operator:

- **Child resources** (like Deployments, Services, PVCs...)  
  are created by the controller and marked with an **OwnerReference**.

🌀 This means:
- Changes or deletions trigger a **reconciliation**
- Kubernetes **automatically garbage-collects** them when the parent is deleted

---

<!-- Slide 8: Handling Dependencies -->
## Handling Dependencies (the Other Way Around)

But what about when our resource depends on another one?

To track and react to **external dependencies**:

1. We add an **OwnerReference** to the *dependency* — but pointing **back to us**
2. We set a **finalizer** on our custom resource

Why?

- So that dependency changes trigger reconciliation
- And so we can **clean up owner references** before we’re deleted

---

<!-- Slide 9: Finalizer Flow -->
## Finalizer Flow

🧼 Finalizer ensures a clean delete lifecycle:

1. User deletes the CR  
2. Kubernetes sets `deletionTimestamp` but keeps the resource alive  
3. Our controller:
   - Removes all `ownerReferences` from dependent resources
   - Then removes the **finalizer**  
4. Now Kubernetes deletes the CR

✅ No orphaned ownerReferences on dependent object that should not be deleted  
✅ No premature garbage collection  
