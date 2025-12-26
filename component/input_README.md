# Input 组件

Input 组件提供了一个文本输入框，支持文本输入、光标移动、删除等基本功能。

## 功能特性

### 文本输入
- 文本输入与显示
- 支持中英文输入
- 占位符文本
- 禁用状态

### 光标控制
- 光标显示与闪烁
- 点击定位光标
- 自动文本滚动（长文本时）

### 文本选择
- 鼠标拖拽选择文本
- Shift+左右箭头扩展选择
- Shift+Home/End快速扩展选择
- 选中文本高亮显示

### 样式
- 自定义样式（颜色、字体、边框等）
- 焦点状态样式

### 键盘快捷键
- `Backspace` - 删除光标前的字符
- `Left/Right Arrow` - 左右移动光标
- `Shift+Left/Right Arrow` - 扩展选择
- `Home` - 光标移到开头
- `End` - 光标移到末尾
- `Shift+Home` - 扩展选择到开头
- `Shift+End` - 扩展选择到结尾
- `Escape` - 失去焦点

## 基本使用

```go
import "github.com/millken/goui/component"

// 创建输入框
input := component.NewInput()

// 设置占位符
input.SetPlaceholder("请输入文本...")

// 获取输入的文本
text := input.Text()

// 设置文本
input.SetText("初始文本")

// 清空文本
input.Clear()

// 设置禁用状态
input.SetDisabled(true)
```


## 自定义样式

```go
// 设置字体大小
input.SetTextSize(16)

// 设置字体家族
input.SetFontFamily("Arial")

// 使用自定义主题
customTheme := component.InputTheme{
    BackgroundColor:  gouicolor.NewColor(255, 255, 255, 255),
    BorderColor:      gouicolor.NewColor(200, 200, 200, 255),
    FocusColor:       gouicolor.NewColor(33, 150, 243, 255),
    TextColor:        gouicolor.NewColor(30, 30, 30, 255),
    PlaceholderColor: gouicolor.NewColor(180, 180, 180, 255),
    DisabledColor:    gouicolor.NewColor(235, 235, 235, 255),
    FontSize:         14.0,
    FontFamily:       "Arial",
    PaddingX:         12,
    PaddingY:         8,
    CornerRadius:     4,
    BorderWidth:      1.5,
}
input.SetTheme(customTheme)
```

## 示例

运行示例程序：

```bash
cd examples/input_demo
go run main.go
```

示例程序展示了以下功能：
- 普通输入框
- 禁用输入框
- 带有初始值的输入框
- 自定义样式的输入框
- 获取和清空输入框内容
- 鼠标拖拽选择文本
- Shift+方向键扩展选择

## 注意事项

- 点击输入框后会获得焦点，显示光标
- 焦点状态下，光标会闪烁
- 失去焦点后，光标会消失
- 禁用状态下，无法输入文本
