<!-- Slide 4: Multi-Operator Design -->
## Multi-Operator Design

Our platform is powered by **multiple Kubernetes Operators**,  
each focused on a **specific responsibility**:

| Operator      | Responsibility                      |
|---------------|--------------------------------------|
| `webapp`      | End-user workloads (apps & services) |
| `vhostroute`  | HTTP traffic routing                 |
| `sshroute`    | SSH access routing                   |
| `certificate` | Certificate storage & management     |
| `maintenance` | Maintenance page and/or fallback     |
| `storage`     | End-user volume declarations         |

---

<!-- Slide 5: Relationships Between Operators -->
## How They Interact

We designed our operators to work **together** but stay **loosely coupled**:

- `vhostroute` can target multiple `webapp` or `maintenance` resources  
- `vhostroute` uses `certificate` for TLS
- `sshroute` targets a `webapp`
- `webapp` implements **contracts** to expose data for `vhostroute` and `sshroute`
- `webapp` references `storage` for persistent volumes

🧩 Each operator **owns its logic**  
🔌 But they connect through **generic contracts**, not direct dependencies.

This allows us to update webapp independently from vhostroute, assuming the contract does not need to change.

---

<!-- Slide 6: Why Contracts? -->
## Why Use Contracts?

To keep operators **decoupled yet interoperable**,  
we use a **contract-based model**:

- One operator exposes the contract its expecting, and the other one can import it.
- Contracts are exposed via a field in the `.status` section
- Other operators can **read** those fields through reflection using `mapstructure`.

This gives us:

✅ Loose coupling  
✅ Extensibility  
✅ Independent lifecycles  

