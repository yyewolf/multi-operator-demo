## Kubernetes

<div class="flex-center pl-45 pb-5 mt-[-60px]">
  <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/3/39/Kubernetes_logo_without_workmark.svg/1280px-Kubernetes_logo_without_workmark.svg.png" class="h-20 rounded" />
</div>

- You describe your system with **YAML specs** (Deployments, Services, etc.)
- All components interact with the **API server** as the single source of truth
- Desired state is stored in **etcd** (distributed key–value store)

- Controllers continuously compare **desired state** (YAML) vs **actual state** (Pods, Nodes) and reconcile differences
- Handles scheduling, scaling, rolling updates, networking, storage, health checks... all automatically

---

## Kubernetes Cluster

<div class="absolute-center" style="position: absolute; top: 50%; left:50%; transform: translate(-50%, -50%);">
```mermaid {theme: 'neutral', scale: 0.8}
graph TB
  subgraph CP[☸️ Control Plane]
    APIS[📡 kube-apiserver]
    ETCD[(🗄️ etcd)]
  end

  subgraph N1[🖥️ Worker Node #1]
    Redis[🟥 Redis]
    DB[(🗄️ PostgreSQL)]
  end

  subgraph N2[🖥️ Worker Node #2]
    App[📦 Application]
  end

  CP --> N1
  CP --> N2
```
</div>

---

## Kubernetes Cluster (close up)

<div class="absolute-center" style="position: absolute; top: 50%; left:50%; transform: translate(-50%, -50%);">
```mermaid {theme: 'neutral', scale: 0.7}
graph TB
  subgraph CP[☸️ Control Plane]
    APIS[📡 kube-apiserver]
    ETCD[(🗄️ etcd)]
    CM[⚙️ kube-controller-manager]
    SCH[📊 kube-scheduler]
  end

  subgraph N1[🖥️ Worker Node #1]
    Kubelet1[🧩 kubelet]
    Pods1[📦 Pods]
  end

  subgraph N2[🖥️ Worker Node #2]
    Kubelet2[🧩 kubelet]
    Pods2[📦 Pods]
  end

  APIS --> ETCD
  APIS --> CM
  APIS --> SCH

  APIS --> Kubelet1
  APIS --> Kubelet2

  Kubelet1 --> Pods1
  Kubelet2 --> Pods2
```
</div>

---

## "Simple" Deployment

<div class="absolute-center" style="position: absolute; top: 50%; left:50%; transform: translate(-50%, -50%);">
```mermaid {theme: 'neutral', scale: 0.7}
  graph TD
    A[📦 Deployment] -->|Manages| B[Pods]
    B -->|Runs on| C((🖥️ K8S Node))
    
    B -->|Uses| D[💾 PersistentVolumeClaim]
    D -->|Binds to| E[🧩 PersistentVolume]
    E -->|Stores on| F[(Storage Backend<br>e.g., NFS, AWS EBS)]

    B -->|Exposed via| G[🌐 Service]
    G -->|Load Balances to| B

    H[🔐 NetworkPolicy] -->|Controls| B
    H -->|Egress & Ingress| G

    style A fill:#42b883,stroke:#389d70,color:white
    style G fill:#2c8ce0,stroke:#226db3,color:white
    style D fill:#d67736,stroke:#b8622c,color:white
    style H fill:#a250c5,stroke:#8442a6,color:white
    style B fill:#333,stroke:#fff,stroke-width:2px,color:white
```
</div>

---

## Deployment

<div class="absolute-center" style="position: absolute; top: 50%; left:50%; transform: translate(-50%, -50%);">
```mermaid {theme: 'neutral', scale: 0.7}
  graph TD
      A[📦 Deployment] -->|Creates & Updates| B[📑 ReplicaSet]
      B -->|Maintains| C[🔹 Pod 1]
      B -->|Maintains| D[🔹 Pod 2]
      B -->|Maintains| E[🔹 Pod 3]

      subgraph Pod_Template[🧩 Pod Template]
          C1[🐳 Containers]
      end

      A -->|Defines| Pod_Template

      %% Self-healing
      C -.If Pod fails, ReplicaSet replaces it.-> B

      %% Scaling
      A -->|Scales replicas| B

      %% Rolling Updates
      A -->|Manages| F[♻️ Rolling Updates]

      style A fill:#42b883,stroke:#389d70,color:white
```
</div>

---

## Deployment

```yaml {1-2|4|6-8|9|14-19|17|19|all}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 2 # tells deployment to run 2 pods matching the template
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```

---

## Core Resources

<div class="p-5"></div>

### 🔹 Workloads
- **Pod** 🐳 — smallest deployable unit, runs containers  
- **ReplicaSet** 📑 — ensures a specified number of Pods  
- **Deployment** 📦 — declarative updates for Pods/ReplicaSets  
- **StatefulSet** 📚 — manages stateful apps, stable IDs & storage  
- **DaemonSet** ⚙️ — runs a Pod on every Node  
- **Job** ⏳ — run Pods until completion  
- **CronJob** ⏰ — run Pods on a schedule  

### 🌐 Networking
- **Service** 🌐 — stable networking endpoint for Pods  
- **NetworkPolicy** 🔐 — control traffic rules  

---

## Core Resources

<div class="p-5"></div>

### 💾 Storage
- **PersistentVolume (PV)** 🧩 — cluster-wide storage resource  
- **PersistentVolumeClaim (PVC)** 💾 — request for storage  
- **ConfigMap** 📝 — store config as key-value pairs  
- **Secret** 🔑 — store sensitive data  

### 🖥️ Cluster
- **Node** 🖥️ — a worker machine  
- **Namespace** 📂 — logical partition of cluster resources  
- **ServiceAccount** 👤 — identity  
- **Role / ClusterRole** & **RoleBinding / ClusterRoleBinding** 🔒 — RBAC access control  