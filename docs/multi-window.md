架构目标与约束
桌面：支持多窗口、每窗独立事件与渲染管线、共享全局资源（字体/主题等），避免线程/GL上下文踩坑。
移动：仅单窗口（或同一上下文下的“视图/页面”切换），API 体验尽量与桌面一致（用“会话”概念屏蔽差异）。
兼容性：保留现有单窗口 NewApp(window, root) 路径，新增“多会话 Application”而不破坏既有代码。
性能：空闲窗口不绘制；批量拉取事件；避免不必要的上下文切换；VSync/TargetFPS 按窗可配置（默认开启）。
核心抽象
Application（新）
职责：持有多个窗口会话，统一主循环与生命周期，管理全局资源（ThorVG Init/Term、FontManager 等）。
关键方法：
NewSession(opts) -> Session
Run() / Quit()
GlobalStyle/Theme 设置（可选）
Session（每窗一份）
成员：Window、Renderer、RootComponent、EventManager、配置（VSync/TargetFPS/Stats）
生命周期：Open/Close、ResetScene、帧渲染（Begin→Layout→Paint→End→Swap）
脏检查：有事件/布局脏/尺寸变化才渲染
WindowFactory（平台创建器，可选）
桌面：每次返回新 Window
移动：拒绝第二个（或返回同一 Window 的“视图”代理），由 Application 做上限约束
建议保留现有 App 作为单窗包装：NewApp(window, root) 实际内部构建 Application + 1 个 Session，以兼容历史。

事件与渲染时序（多窗）
事件侧
主循环每帧做一次全局 Poll（GLFW 是全局函数），然后对每个 Session 调用 w.FetchEvents(0) 读取该窗队列，逐条 Post 到其 EventManager。
鼠标移动合并策略保持在 backend（已实现）；每窗独立队列，互不干扰。
渲染侧（单线程顺序渲染，先简后繁）
对每个 Session：
若尺寸变更：ResetScene(w,h) 并标记根组件 Dirty
是否需要布局：检查 layoutDirty/尺寸变更 → root.Layout → ClearLayoutDirtyRecursive
将该窗 GL 上下文置为当前（见下节）→ renderer.BeginFrame → root.Paint → renderer.EndFrame → window.SwapBuffers
优化：空闲窗只渲染背景（或完全跳过）以减少开销；统计每窗渲染/跳帧率。
可选进阶：后续可演进为“每窗独立渲染 goroutine + OS 线程”，但首版建议单线程依次渲染，避免多上下文并发复杂度。
GL 上下文与 Renderer 绑定
需求：在调用 BeginFrame/ResetScene 之前，确保“该 Session 对应的 GL 上下文”已经 current。
建议最小增量接口（后续再实现）：
window.Window 增加 MakeContextCurrent()（桌面必需；移动可 no-op）
GLFW 后端实现：w.native.MakeContextCurrent()
渲染主循环中每次切换窗时先调用 MakeContextCurrent，再走 BeginFrame/ResetScene；GLRenderer 依赖 C.getCurrentGLContext 不需变更。
VSync：当前实现已通过 SetVSync(…) 配置；需在“对应窗口上下文 current”时调用才能生效（GLFW 语义）。
资源与共享
ThorVG Init/Term：仅 Application 级别调用一次；Session.Cleanup 不 Term 全局，只销毁各自 Canvas。
字体/样式：FontManager/Theme 维持全局实例，Session 共享；避免重复加载字体。
事件管理：每 Session 自有 EventManager，避免不同窗互相吞事件；如需广播事件，Application 可提供简单总线，按需分发至会话。
移动端映射（单窗口）
策略：Application 限制 Session 上限为 1；NewSession 在已存在时返回错误（或返回同一个 Session）。
窗口级行为在移动端退化为“视图/页面栈（Navigator/ViewStack）”：
Application 提供 PushView/PopView，UI 等效于新开/关闭窗口，但仍在同一 Window/Context 内切换 Root（或子树），并做转场动画。
API 一致性：业务层依然“创建会话/打开”，移动端由 Application 将其映射为视图切换，避免平台分支蔓延。
用户体验建议（桌面）
窗口管理
新窗口默认记忆上一次位置与尺寸；避免层叠严格重叠；按平台规范处理焦点与前置。
关闭行为：最后一个窗口关闭即退出（Windows/Linux）；macOS 可选保留应用存活（配置项）。
模态/对话框
在 Application 提供 ShowModal(parent, content)；期间禁用父窗交互，Esc/回车快捷键，点击背景关闭可配。
统一快捷键
每窗公共快捷键（如 Ctrl+W 关闭、Ctrl+Tab 窗间切换）由 Application 注册并转发至当前活动窗口。
性能与电源
非活动窗降低帧率或完全停止渲染（仅尺寸变更/显式刷新时渲染）；统计信息按窗输出。
最小实现路线（建议两步走）
步骤一：为多窗做接口/循环准备（低风险改动）
增：window.Window.MakeContextCurrent()
改：在现有 App.Run 前置一次 glfw 全局 Poll（或保持现状），并为未来多窗循环预留 Session 列表与遍历骨架（保持单窗运行）。
保：现有 NewApp 完全兼容。
步骤二：引入 Application/Session（功能启用）
新：application.go（Application/Session 定义与主循环）
移：将现有 App 的循环/布局/渲染逻辑抽成 per-session 方法（最小复制/复用）
包装：NewApp(window, root) 内部创建 Application+1 Session 并返回一个薄包装以调用 Application.Run。