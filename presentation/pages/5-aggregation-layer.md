<!-- Slide 10: Beyond Reconciliation -->
## Beyond Reconciliation: Custom APIs

Sometimes, reconciling specs isn't enough.  
We needed to expose **actionable business logic** via APIs.

So we used the **Kubernetes Aggregation Layer**  
to serve **custom subresources** directly from our Operators.

This allowed us to build what we call... **Agents**.

---

<!-- Slide 11: Agents Overview -->
## What Are Agents?

🧠 Agents are lightweight HTTP APIs embedded in our Operators.

They expose **operator-specific actions** via the Kubernetes API Server:

- Exposed as custom API services (`/apis/agent.webapp.infomaniak.com/v1/namespaces/default/webapps/myapp/run`)
- Implemented as **Subresources** handled by the Operator's HTTP server

Think of it like:  
💡 CRD = Spec definition  
🛠 Agent = Live control interface

---

<!-- Slide 12: Our Two Agents -->
## Our Two Agents

We currently deploy two agents:

### 🗄 Storage Agent 

- Returns status about the actual storage usage  


### 🖥 WebApp Agent 

- Control the lifecycle of the user's app in real time

---

<!-- Slide 13: Why Use the Aggregation Layer? -->
## Why the Aggregation Layer?

✅ Seamless integration into Kubernetes API  
✅ Standard auth via Kubernetes RBAC  
✅ Reusable tools: `kubectl`, audit logs, etc.  
✅ No need to reconciliate in order to get real-time status.  
