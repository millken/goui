# goui v2 接口规范

本目录定义一套最小且清晰的 GUI 接口：Window、Renderer、Canvas、Component。
目标：职责分离、后端可替换、组件可组合、易于测试。

## 协作关系简述

- Window：应用外壳，持有 Renderer 与根 Component，驱动事件/帧循环。
- Renderer：渲染管线，负责 BeginFrame/EndFrame，提供 Canvas。
- Canvas：绘制抽象，提供最小 2D 原语与状态栈。
- Component：组件抽象，具备 Measure/Layout/Draw，容器负责递归子树。

## 一帧的调用顺序（示例）

1. Window 检测到需要重绘，调用 Renderer.BeginFrame() 拿到 Canvas
2. Window 触发布局：root.Measure(...) → root.Layout(...)
3. Window 触发绘制：从根组件开始依次调用 Draw(Canvas)
4. Window 通知 Renderer.EndFrame() 提交结果

```
Window
  ├─ Renderer.BeginFrame() → Canvas
  ├─ root.Measure(maxW,maxH)
  ├─ root.Layout(0,0,w,h)
  ├─ root.Draw(Canvas)   // 容器负责对子组件递归
  └─ Renderer.EndFrame()
```

## 接口文件

- window.go：窗口最小契约（尺寸、生命周期、挂载、请求重绘）
- renderer.go：渲染器（Begin/EndFrame、Resize、Cleanup）
- canvas.go：绘制抽象（状态、变换、基本图元、文本）
- component.go：组件（测量、布局、绘制、子节点）

## 扩展建议

- 输入事件：另行定义 Event 接口与分发策略，不耦合渲染与布局
- 高 DPI：由 Window 提供 DIP 与物理像素换算，Renderer 内部处理像素对齐
- 样式系统：可在 Component 组合层引入 Style 对象，而不污染 Canvas API
