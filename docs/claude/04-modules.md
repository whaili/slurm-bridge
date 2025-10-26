# 04 - 模块依赖与数据流

## 模块依赖关系

### 1. 核心模块依赖图

```mermaid
graph TD
    %% 入口组件
    A[cmd/scheduler/main.go] --> B[slurmbridge.Plugin]
    C[cmd/admission/main.go] --> D[admission.PodAdmission]
    E[cmd/controllers/main.go] --> F[node.Controller]
    E --> G[pod.Controller]

    %% 内部模块
    B --> H[config.Config]
    D --> H
    F --> H
    G --> H

    %% Slurm 集成模块
    B --> I[slurmcontrol.SlurmControlInterface]
    F --> I
    G --> I

    %% 工具模块
    B --> J[utils/slurmjobir.SlurmJobIR]
    I --> J
    B --> K[utils/placeholderinfo.PlaceholderInfo]

    %% 常量模块
    B --> L[wellknown.Labels & Annotations]
    D --> L

    %% 外部依赖
    H --> M[k8s.io/apimachinery]
    B --> N[k8s.io/kube-scheduler]
    D --> O[sigs.k8s.io/controller-runtime]
    F --> O
    G --> O
    I --> P[github.com/SlinkyProject/slurm-client]
```

### 2. 组件间依赖分析

#### 2.1 cmd/ 组件独立性
```mermaid
graph LR
    subgraph "入口组件"
        A[cmd/scheduler] -->|注册插件| B[slurmbridge]
        C[cmd/admission] -->|Webhook| D[admission.PodAdmission]
        E[cmd/controllers] -->|控制器| F[node.Controller]
        E -->|控制器| G[pod.Controller]
    end

    subgraph "共享依赖"
        H[config.Config]
        I[slurmcontrol.Interface]
        J[wellknown]
    end

    B --> H
    D --> H
    F --> H
    G --> H
    B --> I
    F --> I
    G --> I
    B --> J
    D --> J
```

### 3. 外部依赖层次

#### 3.1 Kubernetes 生态依赖
```mermaid
graph TD
    subgraph "Kubernetes 核心"
        A[k8s.io/api] -->|API 定义| B[Pod/Node/JobSet/LWS]
        C[k8s.io/apimachinery] -->|对象模型| D[metav1/v1]
        E[k8s.io/client-go] -->|客户端| F[Kubernetes API]
        G[sigs.k8s.io/controller-runtime] -->|控制器框架| H[Manager/Webhook]
    end

    subgraph "调度器生态"
        I[k8s.io/kube-scheduler] -->|调度器核心| J[插件框架]
        K[sigs.k8s.io/scheduler-plugins] -->|调度器插件| L[框架扩展]
        M[sigs.k8s.io/jobset] -->|JobSet API| N[作业管理]
        O[sigs.k8s.io/lws] -->|LeaderWorkerSet| P[工作负载管理]
    end

    subgraph "Slurm 生态"
        Q[github.com/SlinkyProject/slurm-client] -->|Slurm 客户端| R[REST API]
    end
```

## 重要数据结构

### 1. 配置管理结构

#### Config 结构体
```go
// 文件: internal/config/config.go
type Config struct {
    SchedulerName            string                // 调度器名称，用于识别 Slurm Bridge 调度器
    SlurmRestApi             string                // Slurm REST API 服务地址
    ManagedNamespaces        []string              // 需要管理的 Kubernetes 命名空间列表
    ManagedNamespaceSelector *metav1.LabelSelector // 命名空间选择器，用于动态选择命名空间
    MCSLabel                 string                // 多类别安全标签
    Partition                string                // Slurm 分区名称
}
```

### 2. Slurm 作业中间表示

#### SlurmJobIR 结构体
```go
// 文件: internal/utils/slurmjobir/slurmjobir.go
type SlurmJobIR struct {
    RootPOM metav1.PartialObjectMetadata  // 根对象元数据（JobSet/PodGroup/LWS/Job/Pod）
    Pods    corev1.PodList                 // 包含的所有 Pod 列表
    JobInfo SlurmJobIRJobInfo              // 作业配置信息
}

type SlurmJobIRJobInfo struct {
    Account      *string  // 账户
    CpuPerTask   *int32   // 每任务 CPU
    Constraints  *string  // 约束条件
    Gres         *string  // GPU 资源
    GroupId      *string  // 组 ID
    JobName      *string  // 作业名称
    Licenses     *string  // 许可证
    MemPerNode   *int64   // 每节点内存（MB）
    MinNodes     *int32   // 最小节点数
    MaxNodes     *int32   // 最大节点数
    Partition    *string  // Slurm 分区
    QOS          *string  // 服务质量
    Reservation  *string  // 预留
    TasksPerNode *int32   // 每节点任务数
    TimeLimit    *int32   // 时间限制（分钟）
    UserId       *string  // 用户 ID
    Wckey        *string  // 工作密钥
}
```

### 3. Slurm 控制接口

#### SlurmControlInterface（调度器用）
```go
// 文件: internal/scheduler/plugins/slurmbridge/slurmcontrol/slurmcontrol.go
type SlurmControlInterface interface {
    DeleteJob(ctx context.Context, pod *corev1.Pod) error                      // 删除 Slurm 作业
    GetJobsForPods(ctx context.Context) (*map[string]PlaceholderJob, error)    // 获取所有 Pod 对应的作业
    GetJob(ctx context.Context, pod *corev1.Pod) (*PlaceholderJob, error)      // 获取单个作业
    SubmitJob(ctx context.Context, pod *corev1.Pod, slurmJobIR *SlurmJobIR) (int32, error) // 提交作业
    UpdateJob(ctx context.Context, pod *corev1.Pod, slurmJobIR *SlurmJobIR) (int32, error) // 更新作业
}

type PlaceholderJob struct {
    JobId int32  // Slurm 作业 ID
    Nodes string // 分配的节点列表
}
```

#### SlurmControlInterface（控制器用）
```go
// 文件: internal/controller/node/slurmcontrol/slurmcontrol.go
type SlurmControlInterface interface {
    GetNodeNames(ctx context.Context) ([]string, error)                        // 获取所有节点名称
    MakeNodeDrain(ctx context.Context, node *corev1.Node, reason string) error    // 设置节点为 DRAIN 状态
    MakeNodeUndrain(ctx context.Context, node *corev1.Node, reason string) error  // 设置节点为 UNDRAIN 状态
    IsNodeDrain(ctx context.Context, node *corev1.Node) (bool, error)           // 检查节点是否为 DRAIN 状态
}
```

### 4. 控制器结构体

#### NodeReconciler 结构体
```go
// 文件: internal/controller/node/node_controller.go
type NodeReconciler struct {
    client.Client                         // Kubernetes 客户端
    Scheme        *runtime.Scheme          // Kubernetes scheme
    SchedulerName string                   // 调度器名称
    SlurmClient   slurmclient.Client      // Slurm 客户端
    EventCh       chan event.GenericEvent  // 事件通道
    slurmControl  slurmcontrol.SlurmControlInterface // Slurm 控制接口
    eventRecorder record.EventRecorderLogger // 事件记录器
}
```

#### PodReconciler 结构体
```go
// 文件: internal/controller/pod/pod_controller.go
type PodReconciler struct {
    client.Client                         // Kubernetes 客户端
    Scheme        *runtime.Scheme          // Kubernetes scheme
    SchedulerName string                   // 调度器名称
    SlurmClient   slurmclient.Client      // Slurm 客户端
    EventCh       chan event.GenericEvent  // 事件通道
    slurmControl  slurmcontrol.SlurmControlInterface // Slurm 控制接口
    eventRecorder record.EventRecorderLogger // 事件记录器
}
```

### 5. 插件结构体

#### SlurmBridge 结构体
```go
// 文件: internal/scheduler/plugins/slurmbridge/slurmbridge.go
type SlurmBridge struct {
    client.Client                         // Kubernetes 客户端
    schedulerName string                   // 调度器名称
    slurmControl  slurmcontrol.SlurmControlInterface // Slurm 控制接口
    handle        framework.Handle          // 调度器句柄
}
```

### 6. Webhook 结构体

#### PodAdmission 结构体
```go
// 文件: internal/admission/admission.go
type PodAdmission struct {
    client.Client                         // Kubernetes 客户端
    SchedulerName            string        // 调度器名称
    ManagedNamespaces        []string      // 管理的命名空间列表
    ManagedNamespaceSelector *metav1.LabelSelector // 命名空间选择器
}
```

## 典型请求处理流程

### 1. Pod 提交处理流程

#### 输入：Pod 创建请求
```mermaid
sequenceDiagram
    participant User as 用户
    participant K8sAPI as Kubernetes API
    participant AW as Admission Webhook
    participant SB as Slurm Bridge Scheduler
    participant Slurm as Slurm REST API
    participant K8s as Kubernetes Storage

    User->>K8sAPI: POST /api/v1/namespaces/{namespace}/pods
    K8sAPI->>AW: 调用 MutatingWebhook
    AW->>AW: 检查命名空间是否受管理

    alt 非受管理命名空间
        AW->>K8sAPI: 直接允许通过
    else 受管理命名空间
        AW->>AW: 验证 Pod 配置
        AW->>AW: 设置默认调度器名称
        AW->>K8sAPI: 返回修改后的 Pod
    end

    K8sAPI->>K8s: 存储 Pod (添加 scheduler.slinky.slurm.net/slurm-jobid 标签)
    K8sAPI->>SB: Pod 进入调度队列
    SB->>SB: PreFilter 阶段
    SB->>SB: 验证 Pod 配置
    SB->>Slurm: 检查占位符作业

    alt 作业不存在
        SB->>Slurm: 提交占位符作业
        Slurm->>SB: 返回 JobId
        SB->>K8s: 更新 Pod 标签
        SB->>SB: 返回 Pending 状态
    else 作业存在
        SB->>Slurm: 获取分配的节点
        Slurm->>SB: 返回节点列表
        SB->>K8s: 注解节点信息
        SB->>SB: 返回 Filter 结果
    end

    SB->>K8s: 绑定 Pod 到节点
```

#### 处理层：Handler → Service → Repo
```mermaid
graph TD
    subgraph "处理层"
        A[Admission Webhook Handler] --> B[PodAdmission Service]
        B --> C[Config Repository]

        D[Scheduler Plugin Handler] --> E[SlurmBridge Service]
        E --> F[SlurmControl Service]
        F --> G[Slurm Repository]

        H[Controller Handler] --> I[Node/Pod Service]
        I --> J[SlurmControl Service]
        J --> G
    end

    subgraph "数据层"
        C --> K[Config File / ConfigMap]
        G --> L[Slurm REST API]
        G --> M[Kubernetes API]
    end

    subgraph "返回层"
        B --> N[Validation Response]
        E --> O[Scheduling Result]
        I --> P[Sync Result]
    end
```

### 2. Controller 同步流程

#### 输入：事件监听
```mermaid
sequenceDiagram
    participant Informer as Kubernetes Informer
    participant Controller as Controller
    participant Slurm as Slurm Informer
    participant Service as Service Layer
    participant Storage as Storage

    Informer->>Controller: 资源事件 (Add/Update/Delete)
    Controller->>Controller: Reconcile() 请求
    Controller->>Service: 调用协调逻辑
    Service->>Slurm: 查询 Slurm 状态
    Slurm->>Service: 返回状态信息
    Service->>Storage: 更新 Kubernetes 状态
    Service->>Controller: 返回协调结果
```

## API 接口表格

### 1. Kubernetes API 调用

| 路径 | 方法 | 入参 | 出参 | 中间件 |
|------|------|------|------|--------|
| `/api/v1/namespaces/{namespace}/pods` | POST | Pod 对象 | AdmissionReview | Validation, Mutation |
| `/api/v1/namespaces/{namespace}/pods/{name}` | GET | Pod 名称 | Pod 对象 | Authorization, Metrics |
| `/api/v1/namespaces/{namespace}/nodes/{name}` | GET | Node 名称 | Node 对象 | Authorization, Metrics |
| `/apis/slinky.slurm.net/v1/placeholderjobs` | LIST | ListOptions | PlaceholderJobList | Authorization |
| `/apis/jobset.x-k8s.io/v1alpha2/jobsets` | LIST | ListOptions | JobSetList | Authorization |
| `/apis/lws.x-k8s.io/v1alpha1/leaderworkersets` | LIST | ListOptions | LWSList | Authorization |

### 2. Slurm REST API 调用

| 路径 | 方法 | 入参 | 出参 | 中间件 |
|------|------|------|------|--------|
| `/slurm/v00437/job/` | POST | JobSubmitRequest | JobId | Authentication |
| `/slurm/v00437/job/{jobid}` | GET | JobId | JobInfo | Authentication |
| `/slurm/v00437/job/{jobid}` | PUT | JobUpdateRequest | JobId | Authentication |
| `/slurm/v00437/job/{jobid}` | DELETE | JobId | Success | Authentication |
| `/slurm/v00437/node/` | GET | 无 | NodeList | Authentication |
| `/slurm/v00437/node/{nodename}` | PUT | DrainRequest | Success | Authentication |

### 3. Controller Runtime API

| 接口 | 方法 | 入参 | 出参 | 用途 |
|------|------|------|------|------|
| `NewManager()` | 构造函数 | ManagerOptions | Manager | 创建控制器管理器 |
| `Watch()` | 监听器 | Source, Handler | Watch | 监听资源变化 |
| `Reconcile()` | 协调器 | Request, Reconciler | Result | 执行协调逻辑 |
| `SetupWebhookWithManager()` | Webhook设置 | Webhook, Manager | Webhook | 设置 Admission Webhook |
| `AddHealthzCheck()` | 健康检查 | Name, Check | Status | 添加健康检查端点 |
| `AddReadyzCheck()` | 就绪检查 | Name, Check | Status | 添加就绪检查端点 |

### 4. 内部服务 API

| 服务接口 | 方法 | 入参 | 出参 | 用途 |
|----------|------|------|------|------|
| `PodAdmission.Default()` | 默认值设置 | Context, Pod | AdmissionResponse | 设置默认调度器名称 |
| `PodAdmission.ValidateCreate()` | 创建验证 | Context, Pod | AdmissionResponse | 验证 Pod 创建请求 |
| `SlurmBridge.PreFilter()` | 调度预处理 | Context, State, Pod, NodeInfo | PreFilterResult | 预处理调度决策 |
| `SlurmBridge.Filter()` | 调度过滤 | Context, State, Pod, NodeInfo | FilterResult | 过滤可用节点 |
| `NodeReconciler.Sync()` | 节点同步 | Context, Request | Error | 同步节点状态 |
| `PodReconciler.Sync()` | Pod 同步 | Context, Request | Error | 同步 Pod 状态 |

## 数据存储和返回

### 1. 存储层
```mermaid
graph TD
    subgraph "Kubernetes 存储"
        A[Pod] --> B[元数据]
        A --> C[标签和注解]
        A --> D[配置信息]
        E[Node] --> F[节点状态]
        E --> G[节点标签]
    end

    subgraph "Slurm 存储"
        H[Job] --> I[作业状态]
        H --> J[节点分配]
        H --> K[作业配置]
    end

    subgraph "配置存储"
        L[ConfigMap] --> M[配置文件]
        N[Secret] --> O[认证信息]
    end
```

### 2. 返回格式
```go
// Admission Response
type AdmissionResponse struct {
    Allowed bool             // 是否允许
    Patch   []byte           // 补丁操作
    Result  *metav1.Status  // 状态信息
}

// Scheduling Result
type FilterResult struct {
    NodeNames   sets.Set[string]  // 可选节点集合
    Status      *framework.Status // 调度状态
    Reason      string           // 状态原因
}

// Sync Result
type SyncResult struct {
    Status  reconcile.Result // 协调结果
    Error   error           // 错误信息
    Message string          // 消息
}
```

## 总结

Slurm Bridge 采用了清晰的分层架构设计：

1. **模块化设计**：三个独立组件通过接口解耦，便于维护和扩展
2. **数据流清晰**：从 Pod 提交到资源分配的完整链路，每个环节都有明确的职责
3. **接口抽象**：通过 SlurmControlInterface 等接口实现不同组件的统一交互
4. **事件驱动**：基于 Informer 的事件驱动机制，确保状态一致性
5. **配置统一**：通过 Config 结构体统一管理所有组件的配置信息

这种架构设计使得系统能够有效协调 Kubernetes 和 Slurm 两个不同的调度系统，实现传统 HPC 工作负载与云原生工作负载的统一管理。