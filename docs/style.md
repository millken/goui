# Style 使用说明

本文档介绍组件 Style 的使用方式、可用属性、单位与 Yoga 的映射关系，以及实践示例与建议。

## 核心概念

- 每个组件包含一个内联样式 `Style`（字段在 `BaseComponent` 内）。
- 使用 `SetStyle(Style)` 覆盖样式，或使用 `UpdateStyle(func(*Style))` 增量修改样式。
- 影响布局的样式（尺寸、边距、定位、Flex 等）会自动标记布局脏（layoutDirty），在下一次布局计算时生效；纯视觉样式（颜色、圆角等）会标记重绘（dirty）而不触发布局。
- Style 会映射到 Yoga 的对应属性；对于百分比/auto，使用 Yoga 的 Set*Percent/Set*Auto API。

## 维度（Dimension）与单位（Unit）

- 维度类型：`type Dimension struct { Value float32; Unit Unit }`
- 单位枚举：
  - `UnitPoint`：物理像素/点（px）
  - `UnitPercent`：百分比（0..100）
  - `UnitAuto`：auto（由 Yoga/布局决定）
  - `UnitUndefined`：未定义（不设置该属性）
- 辅助构造：
  - `Px(v float32)`：像素
  - `Pc(v float32)`：百分比（0..100）
  - `AutoDim()`：auto
  - `UndefinedDim()`：未定义

## 可用样式字段（节选）

布局尺寸：
- `Width` / `Height`：支持 px/percent/auto
- `MinWidth` / `MinHeight`：支持 px/percent
- `MaxWidth` / `MaxHeight`：支持 px/percent

Flex 布局：
- `Flex`, `FlexGrow`, `FlexShrink`
- `FlexBasis`（px/percent/auto）
- `FlexDirection`, `JustifyContent`, `AlignItems`, `AlignSelf`, `AlignContent`, `FlexWrap`
- `Gap`（容器统一行列间距，当前 px）

定位与偏移：
- `Position`（`yoga.PositionType`）
- `Left` / `Top` / `Right` / `Bottom`（px/percent）

边距与内边距：
- `Margin` / `Padding` 使用 `Edges`
  - 支持简写：`All`、`Horizontal`、`Vertical` 与单边 `Left/Top/Right/Bottom`
  - Margin 支持 px/percent/auto；Padding 支持 px/percent（无 auto 概念）

视觉属性：
- `BackgroundColor`（字符串，如 `#RRGGBB` / `#RRGGBBAA` / `rgb(...)` / `rgba(...)`）
- `BorderRadius`（统一圆角，px；后续可扩展四角）
- `Border`（仅宽度，颜色见 `BorderColor`，暂未绘制完整边框）
- `Color`（文本颜色）
- 其他：`Opacity`、`Display`、`Overflow`、`AspectRatio` 等（按需支持）

## 映射到 Yoga 的规则（摘要）

- Width/Height：
  - px -> `SetWidth/SetHeight`
  - percent -> `SetWidthPercent/SetHeightPercent`
  - auto -> `SetWidthAuto/SetHeightAuto`
- Min/Max 宽高：
  - px -> `SetMinWidth/SetMinHeight/SetMaxWidth/SetMaxHeight`
  - percent -> 对应的 `*Percent` 变体
- Flex/FlexBasis：
  - `Flex`, `FlexGrow`, `FlexShrink` 直接映射
  - FlexBasis：px/percent/auto -> `SetFlexBasis/SetFlexBasisPercent/SetFlexBasisAuto`
- Padding：px/percent -> `SetPadding/SetPaddingPercent`
- Margin：px/percent/auto -> `SetMargin/SetMarginPercent/SetMarginAuto`
- Position 偏移：px/percent -> `SetPosition/SetPositionPercent`
- Gap：px -> `SetGap(yoga.GutterAll, v)`

## 使用方式

- 覆盖设置：
```go
c.SetStyle(component.Style{
    Width:  component.Pc(50),     // 宽度 50%
    Height: component.AutoDim(),  // 高度 auto
})
```

- 增量更新：
```go
c.UpdateStyle(func(s *component.Style) {
    s.Margin = &component.Edges{Left: component.AutoDim(), Right: component.AutoDim()} // 水平居中
    s.Padding = &component.Edges{All: component.Px(8)}
})
```

- 配合视觉属性：
```go
c.UpdateStyle(func(s *component.Style) {
    hex := "#3366FF"
    s.BackgroundColor = &hex
    s.BorderRadius = component.Px(8)
})
```

## 实战示例

1) 子项宽度 50%，父固定宽度：
```go
parent.SetStyle(component.Style{ Width: component.Px(400) })
child.SetStyle(component.Style{ Width: component.Pc(50) })
```

2) 使用 margin:auto 实现水平居中：
```go
child.UpdateStyle(func(s *component.Style) {
    s.Margin = &component.Edges{ Left: component.AutoDim(), Right: component.AutoDim() }
})
```

3) Padding 简写与单边：
```go
c.UpdateStyle(func(s *component.Style) {
    s.Padding = &component.Edges{ All: component.Px(12) }         // 四边 12px
    // 或者：s.Padding = &component.Edges{ Horizontal: component.Px(16), Vertical: component.Px(8) }
    // 或者：s.Padding = &component.Edges{ Left: component.Px(8), Right: component.Px(8) }
})
```

4) Button 视觉通过 Style：
```go
btn.UpdateStyle(func(s *component.Style) {
    bg := "#2196F3"   // 背景
    fg := "#FFFFFFFF" // 文本
    s.BackgroundColor = &bg
    s.Color = &fg
    s.BorderRadius = component.Px(6)
})
```

## 最佳实践与注意事项

- 尽量通过 `UpdateStyle` 汇总修改，避免多次触发布局标记。
- 百分比单位相对于父级可用空间计算；请确保父级具有确定的尺寸（固定或由布局决定）。
- 使用 `Auto`（如 `Width: AutoDim()`）时，由 Yoga 根据内容或上下文决定；与 `Min/Max` 组合时注意边界效果。
- `Margin: auto` 常用于居中（左右均为 auto）。
- 视觉属性变化会触发重绘，但不会额外强制布局；布局相关变化会自动标记 layoutDirty。
- 圆角目前为统一半径；如需四角不同，可在后续扩展。

## 相关类型位置

- `Style` 与 `Dimension/Edges` 定义：`component/style.go`
- 应用 Style → Yoga 的入口：`BaseComponent.applyStyleToYoga()`（`component/base_component.go`）
- 视觉使用示例：`component/button.go`、`component/box.go`

如需更多示例或特性支持（如四角圆角、阴影、边框绘制等），欢迎提需求，我们可以继续扩展。