# 05 - 架构文档

## 系统整体架构综述

Slurm Bridge 是一个**双向桥接系统**，实现了 Kubernetes 与 Slurm 两大工作负载管理系统的深度集成。该系统采用**微服务架构**，由三个核心组件构成，通过统一的配置管理和事件驱动机制实现状态同步。

### 架构核心特性

#### 1. 双向桥接机制
- **Kubernetes → Slurm**: Pod 提交时创建 Slurm 占位符作业
- **Slurm → Kubernetes**: 作业状态变化反向同步到 Pod 状态
- **双向状态同步**: 确保两个系统状态的一致性

#### 2. 占位符作业机制
- 通过 Placeholder Job 实现资源预留
- 支持复杂工作负载（JobSet、PodGroup、LWS等）
- 实现了 Kubernetes 调度器与 Slurm 调度器的无缝集成

#### 3. 分层服务架构
- **应用层**: Admission Webhook、Controller、Scheduler
- **逻辑层**: Slurm Bridge 插件、状态管理、同步逻辑
- **接口层**: Kubernetes API、Slurm REST API
- **存储层**: Kubernetes 元数据、Slurm 作业信息

#### 4. 事件驱动架构
- 基于 Kubernetes Informer 的事件监听
- Controller 模式的状态协调机制
- 异步事件处理确保高性能

## 顶层目录表

| 目录 | 作用 | 关键文件 |
|------|------|----------|
| `cmd/` | 应用程序入口点 | `cmd/scheduler/main.go`<br>`cmd/admission/main.go`<br>`cmd/controllers/main.go` |
| `internal/` | 内部应用代码 | |
| &nbsp;&nbsp;`scheduler/plugins/slurmbridge/` | 核心 Slurm Bridge 插件逻辑 | `slurmbridge.go`<br>`slurmcontrol/slurmcontrol.go` |
| &nbsp;&nbsp;`admission/` | Webhook 验证和变更逻辑 | `admission.go` |
| &nbsp;&nbsp;`controller/` | Kubernetes 控制器实现 | `node/node_controller.go`<br>`pod/pod_controller.go` |
| &nbsp;&nbsp;`config/` | 配置管理 | `config.go`<br>`defaults.go` |
| &nbsp;&nbsp;`utils/` | 工具函数 | `slurmjobir/slurmjobir.go`<br>`placeholderinfo/`<br>`durationstore/` |
| &nbsp;&nbsp;`wellknown/` | 常量定义 | `annotations.go`<br>`labels.go`<br>`finalizers.go` |
| `config/` | Kubernetes 清单和 RBAC 配置 | `crd/bases/`<br>`rbac/` |
| `helm/` | Helm 图表部署包 | `slurm-bridge/Chart.yaml`<br>`slurm-bridge/values.yaml` |
| `docs/` | 项目文档 | `architecture.md`<br>`quickstart.md`<br>`workload.md` |
| `hack/` | 构建和开发脚本 | `build-images.sh`<br>`deploy-local.sh` |

## 启动流程图

```mermaid
graph TB
    subgraph "启动阶段"
        A[程序启动] --> B[参数解析]
        B --> C[配置加载]
    end

    subgraph "组件初始化"
        C --> D1{组件类型}

        %% Scheduler 启动
        D1 -->|Scheduler| E1[Kubernetes 调度器命令]
        E1 --> F1[注册 SlurmBridge 插件]
        F1 --> G1[启动调度器服务]

        %% Admission 启动
        D1 -->|Admission| E2[读取配置文件]
        E2 --> F2[创建 Controller Runtime Manager]
        F2 --> G2[设置 TLS 配置]
        G2 --> H2[注册 Pod Admission Webhook]
        H2 --> I2[添加健康检查]
        I2 --> J2[启动 Webhook Manager]

        %% Controllers 启动
        D1 -->|Controllers| E3[读取配置文件]
        E3 --> F3[创建 Controller Runtime Manager]
        F3 --> G3[创建 Slurm 客户端]
        G3 --> H3[启动 Slurm 客户端]
        H3 --> I3[设置 Node 控制器]
        I3 --> J3[设置 Pod 控制器]
        J3 --> K3[添加健康检查]
        K3 --> L3[启动 Controller Manager]
    end

    subgraph "运行阶段"
        G1 --> M[调度循环]
        J2 --> N[Webhook 服务]
        L3 --> O[控制器循环]
    end

    subgraph "服务就绪"
        M --> P[服务运行中]
        N --> P
        O --> P
    end

    subgraph "���能服务"
        P --> Q[监听 Pod 事件]
        P --> R[处理 Slurm 交互]
        P --> S[状态同步]
        P --> T[健康检查]
        P --> U[指标收集]
    end

    %% 安全配置
    G2 --> V[HTTP/2 安全配置]
    G3 --> V
    V --> W[默认禁用 HTTP/2]
    W --> X[防止安全漏洞]
```

## 核心调用链时序图

```mermaid
sequenceDiagram
    participant User as 用户
    participant K8sAPI as Kubernetes API
    participant AW as Admission Webhook
    participant SB as Slurm Bridge Scheduler
    participant Slurm as Slurm REST API
    participant CC as Controllers
    participant NodeC as Node Controller
    participant PodC as Pod Controller
    participant K8sStore as Kubernetes Storage

    User->>K8sAPI: 创建 Pod 请求
    K8sAPI->>AW: 调用 MutatingWebhook
    AW->>AW: 检查命名空间和标签
    AW->>K8sAPI: 返回修改后的 Pod
    K8sAPI->>K8sStore: 存储 Pod (添加调度器标签)

    K8sAPI->>SB: Pod 进入调度队列
    SB->>SB: PreFilter 阶段
    SB->>SB: 验证 Pod 配置
    SB->>Slurm: 检查占位符作业

    alt 作业不存在
        SB->>Slurm: 提交占位符作业
        Slurm->>SB: 返回 JobId
        SB->>K8sStore: 更新 Pod 标签
        SB->>SB: 返回 Pending 状态
    else 作业存在
        SB->>Slurm: 获取分配的节点
        Slurm->>SB: 返回节点列表
        SB->>K8sStore: 注解节点信息
        SB->>SB: 返回 Filter 结果
        SB->>K8sAPI: 绑定指令
    end

    K8sAPI->>K8sStore: 绑定 Pod 到节点
    K8sStore->>PodC: Pod 变更事件
    PodC->>PodC: Sync() 协调
    PodC->>Slurm: 检查作业状态
    PodC->>K8sStore: 更新 Pod 状态

    K8sStore->>NodeC: Node 变更事件
    NodeC->>NodeC: Sync() 协调
    NodeC->>Slurm: 检查节点状态
    NodeC->>Slurm: 设置节点状态

    loop 状态同步循环
        Slurm->>PodC: 作业状态变更
        Slurm->>NodeC: 节点状态变更
        PodC->>K8sStore: 更新相关 Pod
        NodeC->>K8sStore: 更新相关 Node
    end

    %% 错误处理分支
    alt Pod 终止
        PodC->>Slurm: 删除对应作业
        Slurm->>PodC: 确认删除
    end

    alt 节点故障
        NodeC->>Slurm: 节点排水
        Slurm->>NodeC: 确认排水
        NodeC->>K8sStore: 更新节点状态
    end
```

## 模块依赖关系图

```mermaid
graph TD
    %% 入口层
    subgraph "Application Entry Points"
        A[cmd/scheduler/main.go]
        B[cmd/admission/main.go]
        C[cmd/controllers/main.go]
    end

    %% 核心业务层
    subgraph "Core Business Logic"
        D[slurmbridge.Plugin]
        E[admission.PodAdmission]
        F[node.Controller]
        G[pod.Controller]
        H[config.Config]
    end

    %% 接口层
    subgraph "Interfaces"
        I[slurmcontrol.SlurmControlInterface]
        J[slurmclient.Client]
    end

    %% 工具层
    subgraph "Utilities"
        K[utils/slurmjobir.SlurmJobIR]
        L[utils/placeholderinfo.PlaceholderInfo]
        M[utils/durationstore.DurationStore]
        N[wellknown.Labels & Annotations]
    end

    %% Kubernetes 生态
    subgraph "Kubernetes Ecosystem"
        O[k8s.io/kube-scheduler]
        P[sigs.k8s.io/controller-runtime]
        Q[k8s.io/api]
        R[k8s.io/apimachinery]
        S[sigs.k8s.io/scheduler-plugins]
        T[sigs.k8s.io/jobset]
        U[sigs.k8s.io/lws]
    end

    %% Slurm 生态
    subgraph "Slurm Ecosystem"
        V[github.com/SlinkyProject/slurm-client]
    end

    %% 依赖关系
    A --> O
    A --> D
    B --> P
    B --> E
    C --> P
    C --> F
    C --> G

    D --> H
    D --> I
    D --> K
    D --> N

    E --> H
    E --> N

    F --> I
    F --> H
    F --> M

    G --> I
    G --> H
    G --> L

    I --> J
    I --> K

    %% 外部依赖
    D --> O
    D --> S
    D --> T
    D --> U

    F --> V
    G --> V

    %% 工具依赖
    K --> R
    I --> R

    %% 样式定义
    classDef entry fill:#e1f5fe,stroke:#0277bd,stroke-width:2px
    classDef business fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef interface fill:#e8f5e8,stroke:#2e7d32,stroke-width:2px
    classDef util fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef k8s fill:#e8eaf6,stroke:#3949ab,stroke-width:2px
    classDef slurm fill:#fce4ec,stroke:#c2185b,stroke-width:2px

    class A,B,C entry
    class D,E,F,G business
    class I,J interface
    class K,L,M,N util
    class O,P,Q,R,S,T,U k8s
    class V slurm
```

## 外部依赖

### 1. 数据库/存储依赖

#### Kubernetes 集群
```mermaid
graph LR
    subgraph "Kubernetes 存储"
        A[etcd] -->|数据存储| B[API Server]
        B --> C[Pod 元数据]
        B --> D[Node 元数据]
        B --> E[ConfigMap]
        B --> F[Secret]
    end
```

#### Slurm 数据存储
```mermaid
graph LR
    subgraph "Slurm 存储"
        A[Slurm DB] -->|作业信息| B[Slurm Controller]
        B --> C[Job 配置]
        B --> D[节点状态]
        B --> E[用户权限]
    end
```

### 2. API 依赖

#### Kubernetes API
```go
// 核心依赖
- k8s.io/api v0.34.1        // Kubernetes API 定义
- k8s.io/apimachinery v0.34.1 // Kubernetes 基础对象
- k8s.io/client-go v0.34.1  // Kubernetes 客户端
- sigs.k8s.io/controller-runtime v0.22.1 // 控制器运行时
```

#### Slurm REST API
```go
// Slurm REST API 客户端
- github.com/SlinkyProject/slurm-client v0.4.1 // Slurm 客户端库
```

### 3. 消息队列/事件系统

#### Kubernetes 事件系统
```mermaid
graph TD
    A[Pod 事件] --> B[Kubernetes Informer]
    B --> C[Controller]
    C --> D[协调逻辑]

    E[Node 事件] --> B
    F[ConfigMap 事件] --> B

    D --> G[状态更新]
    G --> B
```

#### Slurm 事件系统
```mermaid
graph TD
    A[Slurm 作业状态] --> B[Slurm Client]
    B --> C[Pod Controller]
    C --> D[Pod 状态同步]

    E[Slurm 节点状态] --> B
    B --> F[Node Controller]
    F --> G[Node 状态同步]
```

### 4. 第三方 API

#### 监控和日志
```go
// 监控依赖
- k8s.io/component-base/metrics/prometheus/clientgo // Prometheus 指标
- k8s.io/klog/v2 // 日志记录

// 工具依赖
- github.com/onsi/ginkgo/v2 v2.25.3 // 测试框架
- github.com/onsi/gomega v1.38.2  // 断言库
```

#### 网络和工具
```go
// 网络工具
- github.com/puttsk/hostlist v0.1.0 // 主机列表处理
- k8s.io/utils v0.0.0-20250820121507-0af2bda4dd1d // 工具库
```

## 配置项

### 1. 主配置文件

#### `/etc/slurm-bridge/config.yaml`
```yaml
# 调度器配置
schedulerName: "slurm-bridge"

# Slurm REST API 配置
slurmRestApi: "http://slurm-controller:6820"

# 管理的命名空间
managedNamespaces:
  - "batch"
  - "hpc"
  - "ai-workloads"

# 命名空间选择器（动态选择命名空间）
managedNamespaceSelector:
  matchLabels:
    slurm.managed: "true"

# 多类别安全标签
mcsLabel: ""

# Slurm 分区
partition: "normal"

# 高级配置
features:
  enableLeaderElection: true
  enableHTTP2: false
  secureMetrics: false
```

### 2. 环境变量配置

#### Slurm 认证
```bash
# Slurm JWT Token
export SLURM_JWT="your-jwt-token-here"

# Slurm REST API 地址
export SLURM_REST_API="http://slurm-controller:6820"
```

#### Kubernetes 认证
```bash
# Kubernetes 集群配置
export KUBECONFIG="/path/to/kubeconfig"

# 命名空间选择
export NAMESPACE_SELECTOR="slurm.managed=true"
```

### 3. 命令行参数配置

#### Scheduler 配置
```bash
# 默认使用 Kubernetes 调度器默认参数
--kubeconfig=/path/to/kubeconfig
--v=2
```

#### Admission Webhook 配置
```bash
# 配置文件
--config=/etc/slurm-bridge/config.yaml

# 服务地址
--metrics-bind-address=:8080
--health-probe-bind-address=:8081

# 领导者选举
--leader-elect=true
--leader-election-id=a1f3cd42.slinky.slurm.net

# 安全配置
--metrics-secure=false
--enable-http2=false
```

#### Controllers 配置
```bash
# 配置文件
--config=/etc/slurm-bridge/config.yaml

# 服务地址
--metrics-bind-address=:8080
--health-probe-bind-address=:8081

# 领导者选举
--leader-elect=true
--leader-election-id=69d5fe47.my.slinky.slurm.net

# 安全配置
--metrics-secure=false
--enable-http2=false
```

### 4. Helm 配置

#### `helm/slurm-bridge/values.yaml`
```yaml
# 部署配置
replicaCount: 3
image:
  repository: slinky.slurm.net/slurm-bridge
  pullPolicy: IfNotPresent
  tag: "1.0.0"

# 资源限制
resources:
  limits:
    cpu: 1
    memory: 1Gi
  requests:
    cpu: 100m
    memory: 128Mi

# 配置
config:
  schedulerName: "slurm-bridge"
  slurmRestApi: "http://slurm-controller:6820"
  managedNamespaces:
    - "batch"
    - "hpc"
    - "ai-workloads"
  managedNamespaceSelector:
    matchLabels:
      slurm.managed: "true"
  mcsLabel: ""
  partition: "normal"

# 服务配置
service:
  type: ClusterIP
  port: 8080

# 安全配置
securityContext:
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000

# RBAC 配置
rbac:
  create: true
  rules: []
```

### 5. RBAC 配置

#### ClusterRole 和 ClusterRoleBinding
```yaml
# 服务账户权限
apiVersion: v1
kind: ServiceAccount
metadata:
  name: slurm-bridge-sa
  namespace: slurm-bridge

# 集群角色
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: slurm-bridge-role
rules:
  - apiGroups: [""]
    resources: ["pods", "nodes", "namespaces"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["slinky.slurm.net"]
    resources: ["placeholderjobs"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["jobset.x-k8s.io"]
    resources: ["jobsets"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
```

## 总结

Slurm Bridge 架构设计体现了以下几个关键原则：

### 1. 分层解耦
- **应用层**: 三个独立组件，各司其职
- **逻辑层**: 通过接口抽象实现松耦合
- **数据层**: 统一的配置管理和状态存储

### 2. 事件驱动
- 基于 Kubernetes Informer 的事件机制
- 异步处理确保高性能
- 双向同步保证状态一致性

### 3. 可扩展性
- 插件化架构，易于扩展新功能
- 配置驱动，支持多种部署模式
- 接口抽象，便于测试和维护

### 4. 高可用性
- 领导者选举机制
- 健康检查和就绪检查
- 自动故障恢复

### 5. 安全性
- TLS 支持
- RBAC 权限控制
- JWT 认证

这种架构设计使得 Slurm Bridge 能够有效协调 Kubernetes 和 Slurm 两个不同的调度系统，实现传统 HPC 工作负载与云原生工作负载的统一管理，为用户提供一致的使用体验。