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

1. We add a **Label** to the *dependency* that's pointing **back to the referer**
   - That label is a list of all referers (by their type, name and namespace).

Why?

- So that dependency changes trigger reconciliation

What about cleanup ?

- Using a **finalizer** on the referer, we can make sure to clean that Label of the reference back.

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
