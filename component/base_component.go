// Package component provides UI component implementations for GOUI.
package component

import (
	"image"
	"log/slog"
	"math"
	"sync"

	"github.com/millken/goui/event"
	"github.com/millken/goui/window"
	"github.com/millken/thorvg"
	"github.com/millken/yoga"
)

type BaseComponent struct {
	root     Component
	id       string
	mu       sync.RWMutex
	yoga     *yoga.Node
	children []Component
	//events
	eventManager     *event.EventManager
	pendingListeners map[event.Type][]event.EventListener

	parent  Component
	visible bool // 可见性标记
	dirty   bool // 脏标记
	hovered bool // 悬停状态
	win     window.Window
	// layout 脏标记：仅当影响布局的属性变化时置为 true
	layoutDirty bool
	// 样式（内联）
	style Style
	// 计算样式（含继承），惰性计算
	computedStyle Style
	styleDirty    bool
}

func NewBaseComponent(id ...string) *BaseComponent {
	base := &BaseComponent{
		id: func() string {
			if len(id) > 0 {
				return id[0]
			} else {
				return GenerateID()
			}
		}(),
		yoga:             yoga.NewNode(),
		children:         []Component{},
		parent:           nil,
		visible:          true, // 默认可见
		dirty:            true, // 初始状态为需要渲染
		pendingListeners: make(map[event.Type][]event.EventListener),
		styleDirty:       true,
	}
	return base
}

func (b *BaseComponent) SetRoot(c Component) {
	b.mu.Lock()
	b.root = c
	// 复制 children 与 pending 快照
	children := append([]Component(nil), b.children...)
	pending := make(map[event.Type][]event.EventListener, len(b.pendingListeners))
	for et, ls := range b.pendingListeners {
		if len(ls) == 0 {
			continue
		}
		cp := make([]event.EventListener, len(ls))
		copy(cp, ls)
		pending[et] = cp
	}
	// 清空缓存，避免重复
	b.pendingListeners = make(map[event.Type][]event.EventListener)
	b.mu.Unlock()

	// 若能从 root 上拿到 EventManager，则立即注册缓存的监听器
	if rc, ok := c.(RootComponent); ok && rc != nil {
		if em := rc.GetEventManager(); em != nil {
			for et, ls := range pending {
				for _, l := range ls {
					em.RegisterListener(b.ID(), et, l)
				}
			}
			// 清空已注册成功的本地快照
			pending = nil
		}
	}
	// 如果仍无法注册（例如 EM 未就绪），将快照放回缓存
	if len(pending) > 0 {
		b.mu.Lock()
		for et, ls := range pending {
			b.pendingListeners[et] = append(b.pendingListeners[et], ls...)
		}
		b.mu.Unlock()
	}

	// 将根组件的引用向下传播，避免已添加的子组件 Root 仍为 nil
	for _, child := range children {
		child.SetRoot(c)
	}

	// 根变化可能影响样式继承上下文
	b.mu.Lock()
	b.styleDirty = true
	b.mu.Unlock()
}

func (b *BaseComponent) Root() Component {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.root
}

func (b *BaseComponent) SetEventManager(em *event.EventManager) {
	// 设置自身 EM，并尽量刷新缓存，再向下传播
	b.mu.Lock()
	b.eventManager = em
	// 快照 children 与 pending，避免持锁调用外部
	children := append([]Component(nil), b.children...)
	pending := make(map[event.Type][]event.EventListener, len(b.pendingListeners))
	for et, ls := range b.pendingListeners {
		if len(ls) == 0 {
			continue
		}
		cp := make([]event.EventListener, len(ls))
		copy(cp, ls)
		pending[et] = cp
	}
	// 清空缓存，避免重复
	b.pendingListeners = make(map[event.Type][]event.EventListener)
	b.mu.Unlock()

	if em != nil {
		for et, ls := range pending {
			for _, l := range ls {
				em.RegisterListener(b.ID(), et, l)
			}
		}
	} else if len(pending) > 0 {
		// EM 仍为空，把快照放回缓存（极少见）
		b.mu.Lock()
		for et, ls := range pending {
			b.pendingListeners[et] = append(b.pendingListeners[et], ls...)
		}
		b.mu.Unlock()
	}

	// 将 EM 向下传播，促使子组件也能刷新缓存
	for _, child := range children {
		child.SetEventManager(em)
	}
}

func (b *BaseComponent) GetEventManager() *event.EventManager {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.eventManager
}

// SetWindow 设置窗口引用，并向下传播
func (b *BaseComponent) SetWindow(w window.Window) {
	b.mu.Lock()
	b.win = w
	children := append([]Component(nil), b.children...)
	b.mu.Unlock()
	for _, c := range children {
		if rc, ok := c.(RootComponent); ok {
			rc.SetWindow(w)
		} else if bc, ok := c.(*BaseComponent); ok {
			bc.SetWindow(w)
		}
	}
}

// GetWindow 获取窗口引用
func (b *BaseComponent) GetWindow() window.Window {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.win
}

func (b *BaseComponent) SetID(id string) {
	b.id = id
}

func (b *BaseComponent) ID() string {
	return b.id
}

func (b *BaseComponent) IsVisible() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.visible
}

func (b *BaseComponent) SetVisible(visible bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.visible != visible {
		b.visible = visible
		b.MarkDirty(false) // 可见性改变时标记为脏
	}
}

func (b *BaseComponent) GetYogaNode() *yoga.Node {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.yoga
}

// YogaUpdate 在受控的临界区内对 Yoga 节点进行样式更新，并集中标记渲染脏与布局脏。
// 注意：传入的更新函数必须只对 node 做无阻塞操作，避免回调中再调用组件方法造成潜在死锁。
func (b *BaseComponent) YogaUpdate(update func(node *yoga.Node)) {
	b.mu.Lock()
	if b.yoga != nil && update != nil {
		update(b.yoga)
		b.dirty = true
		b.layoutDirty = true
	}
	b.mu.Unlock()
}

// -------- Style（内联样式）管理与应用 --------

// SetStyle 覆盖设置组件样式并应用到 Yoga
func (b *BaseComponent) SetStyle(s Style) {
	b.mu.Lock()
	b.style = s
	b.styleDirty = true
	b.mu.Unlock()
	b.applyStyleToYoga()
}

// UpdateStyle 在闭包中修改样式并应用到 Yoga
func (b *BaseComponent) UpdateStyle(mut func(s *Style)) {
	if mut == nil {
		return
	}
	b.mu.Lock()
	mut(&b.style)
	b.styleDirty = true
	b.mu.Unlock()
	b.applyStyleToYoga()
}

// Style 返回样式的副本（避免外部在无锁状态下修改）
func (b *BaseComponent) Style() Style {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.style
}

// ComputedStyle 返回计算后的样式（包含可继承属性）。
// 当前继承的属性：Color, Font(字体族), FontSize, Cursor。
func (b *BaseComponent) ComputedStyle() Style {
	// 快速路径：未标记样式脏，返回快照
	b.mu.RLock()
	dirty := b.styleDirty
	cached := b.computedStyle
	inline := b.style
	parent := b.parent
	b.mu.RUnlock()
	if !dirty {
		return cached
	}

	// 取父的 computedStyle（避免持锁递归）
	var parentCS Style
	if parent != nil {
		if p, ok := parent.(interface{ ComputedStyle() Style }); ok {
			parentCS = p.ComputedStyle()
		}
	}

	// 从自身内联开始，补足可继承属性
	out := inline
	if out.Color == nil {
		out.Color = parentCS.Color
	}
	if out.Font == nil {
		out.Font = parentCS.Font
	}
	if out.FontSize == nil {
		out.FontSize = parentCS.FontSize
	}
	if out.Cursor == nil {
		out.Cursor = parentCS.Cursor
	}

	b.mu.Lock()
	b.computedStyle = out
	b.styleDirty = false
	b.mu.Unlock()
	return out
}

// MarkStyleDirty 标记自身/子树需要重算计算样式（继承影响）。
func (b *BaseComponent) MarkStyleDirty(recursive bool) {
	b.mu.Lock()
	b.styleDirty = true
	children := append([]Component(nil), b.children...)
	b.mu.Unlock()
	if recursive {
		for _, child := range children {
			if bc, ok := child.(*BaseComponent); ok {
				bc.MarkStyleDirty(true)
			}
		}
	}
}

// 将样式字段映射到 Yoga 节点；当前仅处理 px/Auto 的常用属性，后续可扩展百分比等
func (b *BaseComponent) applyStyleToYoga() {
	s := b.Style()
	b.YogaUpdate(func(n *yoga.Node) {
		// 尺寸（支持 px/percent/auto）
		if s.Width != nil {
			switch s.Width.Unit {
			case UnitPoint:
				n.SetWidth(s.Width.Value)
			case UnitPercent:
				n.SetWidthPercent(s.Width.Value)
			case UnitAuto:
				n.SetWidthAuto()
			}
		}
		if s.Height != nil {
			switch s.Height.Unit {
			case UnitPoint:
				n.SetHeight(s.Height.Value)
			case UnitPercent:
				n.SetHeightPercent(s.Height.Value)
			case UnitAuto:
				n.SetHeightAuto()
			}
		}
		if s.MinWidth != nil {
			switch s.MinWidth.Unit {
			case UnitPoint:
				n.SetMinWidth(s.MinWidth.Value)
			case UnitPercent:
				n.SetMinWidthPercent(s.MinWidth.Value)
			}
		}
		if s.MinHeight != nil {
			switch s.MinHeight.Unit {
			case UnitPoint:
				n.SetMinHeight(s.MinHeight.Value)
			case UnitPercent:
				n.SetMinHeightPercent(s.MinHeight.Value)
			}
		}
		if s.MaxWidth != nil {
			switch s.MaxWidth.Unit {
			case UnitPoint:
				n.SetMaxWidth(s.MaxWidth.Value)
			case UnitPercent:
				n.SetMaxWidthPercent(s.MaxWidth.Value)
			}
		}
		if s.MaxHeight != nil {
			switch s.MaxHeight.Unit {
			case UnitPoint:
				n.SetMaxHeight(s.MaxHeight.Value)
			case UnitPercent:
				n.SetMaxHeightPercent(s.MaxHeight.Value)
			}
		}

		// Flex 属性
		if s.Flex != nil {
			n.SetFlex(*s.Flex)
		}
		if s.FlexGrow != nil {
			n.SetFlexGrow(*s.FlexGrow)
		}
		if s.FlexShrink != nil {
			n.SetFlexShrink(*s.FlexShrink)
		}
		if s.FlexBasis != nil {
			switch s.FlexBasis.Unit {
			case UnitPoint:
				n.SetFlexBasis(s.FlexBasis.Value)
			case UnitPercent:
				n.SetFlexBasisPercent(s.FlexBasis.Value)
			case UnitAuto:
				n.SetFlexBasisAuto()
			}
		}
		if s.AlignSelf != nil {
			n.SetAlignSelf(*s.AlignSelf)
		}
		if s.FlexDirection != nil {
			n.SetFlexDirection(*s.FlexDirection)
		}
		if s.JustifyContent != nil {
			n.SetJustifyContent(*s.JustifyContent)
		}
		if s.AlignItems != nil {
			n.SetAlignItems(*s.AlignItems)
		}
		if s.AlignContent != nil {
			n.SetAlignContent(*s.AlignContent)
		}
		if s.FlexWrap != nil {
			n.SetFlexWrap(*s.FlexWrap)
		}

		// 内边距（支持 px/percent；padding 无 auto 概念）
		if s.Padding != nil {
			for _, edge := range []yoga.Edge{yoga.EdgeLeft, yoga.EdgeTop, yoga.EdgeRight, yoga.EdgeBottom} {
				if d := edgeDim(s.Padding, edge); d != nil {
					switch d.Unit {
					case UnitPoint:
						n.SetPadding(edge, d.Value)
					case UnitPercent:
						n.SetPaddingPercent(edge, d.Value)
					}
				}
			}
		}
		// 边距（支持 px/percent/auto）
		if s.Margin != nil {
			for _, edge := range []yoga.Edge{yoga.EdgeLeft, yoga.EdgeTop, yoga.EdgeRight, yoga.EdgeBottom} {
				if d := edgeDim(s.Margin, edge); d != nil {
					switch d.Unit {
					case UnitPoint:
						n.SetMargin(edge, d.Value)
					case UnitPercent:
						n.SetMarginPercent(edge, d.Value)
					case UnitAuto:
						n.SetMarginAuto(edge)
					}
				}
			}
		}

		// 位置与偏移
		if s.Position != nil {
			n.SetPositionType(*s.Position)
		}
		if s.Left != nil {
			switch s.Left.Unit {
			case UnitPoint:
				n.SetPosition(yoga.EdgeLeft, s.Left.Value)
			case UnitPercent:
				n.SetPositionPercent(yoga.EdgeLeft, s.Left.Value)
			}
		}
		if s.Top != nil {
			switch s.Top.Unit {
			case UnitPoint:
				n.SetPosition(yoga.EdgeTop, s.Top.Value)
			case UnitPercent:
				n.SetPositionPercent(yoga.EdgeTop, s.Top.Value)
			}
		}
		if s.Right != nil {
			switch s.Right.Unit {
			case UnitPoint:
				n.SetPosition(yoga.EdgeRight, s.Right.Value)
			case UnitPercent:
				n.SetPositionPercent(yoga.EdgeRight, s.Right.Value)
			}
		}
		if s.Bottom != nil {
			switch s.Bottom.Unit {
			case UnitPoint:
				n.SetPosition(yoga.EdgeBottom, s.Bottom.Value)
			case UnitPercent:
				n.SetPositionPercent(yoga.EdgeBottom, s.Bottom.Value)
			}
		}

		// Display / Overflow / Gap（仅 px 的统一 gap）
		if s.Display != nil {
			n.SetDisplay(*s.Display)
		}
		if s.Overflow != nil {
			n.SetOverflow(*s.Overflow)
		}
		if v, ok := dimPx(s.Gap); ok {
			n.SetGap(yoga.GutterAll, v)
		}
	})
	// 即便没有 Yoga 属性变化，纯视觉样式（如背景色、圆角）也需要重绘
	b.MarkDirty(false)
}

// 工具：解析 Dimension（仅 px 与 auto 支持；percent 留待后续）
func dimPx(d *Dimension) (float32, bool) {
	if d == nil {
		return 0, false
	}
	switch d.Unit {
	case UnitPoint:
		return d.Value, true
	case UnitAuto:
		// 使用 Yoga 的 Auto（特定 API，如 SetWidthAuto），此处跳过，由上层专门处理
		return 0, false
	default:
		return 0, false
	}
}

// 工具：将 Edges 解析为具体的 LTRB（仅 px）
func resolveEdgesPx(e *Edges) (l, t, r, btm float32) {
	if e == nil {
		return
	}
	var hv, vv float32
	if v, ok := dimPx(e.All); ok {
		l, t, r, btm = v, v, v, v
	}
	if v, ok := dimPx(e.Horizontal); ok {
		hv = v
		l, r = v, v
	}
	if v, ok := dimPx(e.Vertical); ok {
		vv = v
		t, btm = v, v
	}
	if v, ok := dimPx(e.Left); ok {
		l = v
	}
	if v, ok := dimPx(e.Top); ok {
		t = v
	}
	if v, ok := dimPx(e.Right); ok {
		r = v
	}
	if v, ok := dimPx(e.Bottom); ok {
		btm = v
	}
	// 若仅设置了水平/垂直而未单独覆盖，已在上方处理
	_ = hv
	_ = vv
	return
}

// edgeDim 根据 CSS 优先级从 Edges 中解析某个具体边的 Dimension
// 优先级：Specific(Left/Top/Right/Bottom) > Axis(Horizontal/Vertical) > All
func edgeDim(e *Edges, edge yoga.Edge) *Dimension {
	if e == nil {
		return nil
	}
	var base *Dimension
	// All 作为默认值
	if e.All != nil {
		base = e.All
	}
	switch edge {
	case yoga.EdgeLeft, yoga.EdgeRight:
		if e.Horizontal != nil {
			base = e.Horizontal
		}
	case yoga.EdgeTop, yoga.EdgeBottom:
		if e.Vertical != nil {
			base = e.Vertical
		}
	}
	// 最终由具体边覆盖
	switch edge {
	case yoga.EdgeLeft:
		if e.Left != nil {
			base = e.Left
		}
	case yoga.EdgeTop:
		if e.Top != nil {
			base = e.Top
		}
	case yoga.EdgeRight:
		if e.Right != nil {
			base = e.Right
		}
	case yoga.EdgeBottom:
		if e.Bottom != nil {
			base = e.Bottom
		}
	}
	return base
}

// SetWidth 设置组件的宽度（便捷方法，避免直接操作 yoga.Node）
func (b *BaseComponent) SetWidth(w float32) {
	b.YogaUpdate(func(n *yoga.Node) { n.SetWidth(w) })
}

// SetHeight 设置组件的高度（便捷方法，避免直接操作 yoga.Node）
func (b *BaseComponent) SetHeight(h float32) {
	b.YogaUpdate(func(n *yoga.Node) { n.SetHeight(h) })
}

// SetMinWidth 设置组件的最小宽度（便捷方法，避免直接操作 yoga.Node）
func (b *BaseComponent) SetMinWidth(w float32) {
	b.YogaUpdate(func(n *yoga.Node) { n.SetMinWidth(w) })
}

// SetMinHeight 设置组件的最小高度（便捷方法，避免直接操作 yoga.Node）
func (b *BaseComponent) SetMinHeight(h float32) {
	b.YogaUpdate(func(n *yoga.Node) { n.SetMinHeight(h) })
}

// SetFlexGrow 设置组件的伸展比例（便捷方法，避免直接操作 yoga.Node）
func (b *BaseComponent) SetFlexGrow(g float32) {
	b.YogaUpdate(func(n *yoga.Node) { n.SetFlexGrow(g) })
}

// SetAlignSelf 设置交叉轴上的自对齐（便捷方法，避免直接操作 yoga.Node）
func (b *BaseComponent) SetAlignSelf(a Align) {
	b.YogaUpdate(func(n *yoga.Node) { n.SetAlignSelf(a.toYoga()) })
}

func (b *BaseComponent) Layout(width, height float32, direction yoga.Direction) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.yoga.CalculateLayout(width, height, direction)
}

func (b *BaseComponent) ComputedLayout() yoga.Layout {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.yoga.GetComputedLayout()
}

func (b *BaseComponent) AddChild(c Component) {
	// 在加锁前抓取子组件 Yoga 节点，避免在持有父锁时再去锁子组件
	var childNode *yoga.Node
	type yogaNodeProvider interface{ GetYogaNode() *yoga.Node }
	if yp, ok := any(c).(yogaNodeProvider); ok {
		childNode = yp.GetYogaNode()
	}

	// 并发安全地修改 children，并将 Yoga 子节点挂接到父节点
	b.mu.Lock()
	b.children = append(b.children, c)
	if b.yoga != nil && childNode != nil {
		idx := b.yoga.GetChildCount()
		b.yoga.InsertChild(childNode, idx)
	}
	// 标记自身为脏
	b.dirty = true
	b.layoutDirty = true
	b.mu.Unlock()

	// 在未持有父锁的情况下设置父子关系，避免潜在的锁顺序反转
	c.SetParent(b)
	if bc, ok := c.(*BaseComponent); ok {
		bc.MarkStyleDirty(true)
	}

	// Root 传递：若当前 Root 已设置则传递；否则先用自身作为临时 Root，后续 SetRoot 会向下传播
	if r := b.Root(); r != nil {
		c.SetRoot(r)
	} else {
		c.SetRoot(b)
	}

	// 若父组件已关联事件管理器，则将其同步给新子组件，刷新其待注册监听器
	if em := b.GetEventManager(); em != nil {
		c.SetEventManager(em)
	}
}

func (b *BaseComponent) Children() []Component {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.children
}

func (b *BaseComponent) Paint(vg *thorvg.Canvas) bool {
	// 快照状态，避免持锁递归调用子组件
	b.mu.Lock()
	children := append([]Component(nil), b.children...)
	visible := b.visible
	isDirty := b.dirty
	if isDirty {
		b.dirty = false
	}
	b.mu.Unlock()

	// 先绘制子组件
	childRendered := false
	for _, child := range children {
		if child.Paint(vg) {
			childRendered = true
		}
	}

	// 如果不可见，跳过自身渲染
	if !visible {
		return childRendered
	}
	// 自身脏则视为已渲染（具体绘制由子类完成）
	if isDirty {
		return true
	}
	return childRendered
}

// MarkDirty 标记组件为脏（需要重绘）
func (b *BaseComponent) MarkDirty(recursive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = true
	if recursive {
		for _, child := range b.children {
			child.MarkDirty(true)
		}
	}
}

// IsDirty 检查组件是否需要重绘
func (b *BaseComponent) IsDirty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dirty
}

// ClearDirty 清除脏标记
func (b *BaseComponent) ClearDirty() {
	b.mu.Lock()
	b.dirty = false
	b.mu.Unlock()
}

// SetParent 设置父组件
func (b *BaseComponent) SetParent(c Component) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.parent = c
	b.styleDirty = true
}

// Parent 获取父组件
func (b *BaseComponent) Parent() Component {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.parent
}

// Bounds 返回组件的边界矩形
func (b *BaseComponent) Bounds() image.Rectangle {
	l := b.ComputedLayout()
	ax, ay := b.absoluteOffset()
	// 清洗可能的 NaN/Inf，避免产生极端边界
	left := sanitizeFloat(ax)
	top := sanitizeFloat(ay)
	width := sanitizeFloat(l.Width)
	height := sanitizeFloat(l.Height)
	right := left + width
	bottom := top + height
	return image.Rectangle{
		Min: image.Point{X: int(left), Y: int(top)},
		Max: image.Point{X: int(right), Y: int(bottom)},
	}
}

// CheckDirtyRecursive 递归检查是否有任何子组件需要重绘（保留向后兼容）
func (b *BaseComponent) CheckDirtyRecursive() bool {
	b.mu.RLock()
	if b.dirty {
		b.mu.RUnlock()
		return true
	}
	// 复制子节点切片，避免持锁递归
	children := append([]Component(nil), b.children...)
	b.mu.RUnlock()
	for _, child := range children {
		if rc, ok := child.(interface{ CheckDirtyRecursive() bool }); ok {
			if rc.CheckDirtyRecursive() {
				return true
			}
		}
	}
	return false
}

// ----- 布局脏标记相关 -----

// MarkLayoutDirty 标记布局为脏；可递归标记子树
func (b *BaseComponent) MarkLayoutDirty(recursive bool) {
	b.mu.Lock()
	b.layoutDirty = true
	children := append([]Component(nil), b.children...)
	b.mu.Unlock()
	if recursive {
		for _, child := range children {
			if bc, ok := child.(*BaseComponent); ok {
				bc.MarkLayoutDirty(true)
			}
		}
	}
}

// IsLayoutDirty 返回是否布局脏
func (b *BaseComponent) IsLayoutDirty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.layoutDirty
}

// ClearLayoutDirty 清除自身布局脏标记
func (b *BaseComponent) ClearLayoutDirty() {
	b.mu.Lock()
	b.layoutDirty = false
	b.mu.Unlock()
}

// CheckLayoutDirtyRecursive 递归检查子树是否存在布局脏
func (b *BaseComponent) CheckLayoutDirtyRecursive() bool {
	b.mu.RLock()
	if b.layoutDirty {
		b.mu.RUnlock()
		return true
	}
	children := append([]Component(nil), b.children...)
	b.mu.RUnlock()
	for _, child := range children {
		if rc, ok := child.(interface{ CheckLayoutDirtyRecursive() bool }); ok {
			if rc.CheckLayoutDirtyRecursive() {
				return true
			}
		}
	}
	return false
}

// ClearLayoutDirtyRecursive 递归清除布局脏标记
func (b *BaseComponent) ClearLayoutDirtyRecursive() {
	b.ClearLayoutDirty()
	b.mu.RLock()
	children := append([]Component(nil), b.children...)
	b.mu.RUnlock()
	for _, child := range children {
		if rc, ok := child.(interface{ ClearLayoutDirtyRecursive() }); ok {
			rc.ClearLayoutDirtyRecursive()
		}
	}
}

// absoluteOffset 通过向上累加 Yoga 计算布局，得到绝对坐标
func (b *BaseComponent) absoluteOffset() (float32, float32) {
	var x, y float32
	var cur Component = b
	for cur != nil {
		l := cur.ComputedLayout()
		x += sanitizeFloat(l.Left)
		y += sanitizeFloat(l.Top)
		cur = cur.Parent()
	}
	return x, y
}

// PointInside 判断窗口坐标 (x,y) 是否在该组件边界内
func (b *BaseComponent) PointInside(x, y float32) bool {
	l := b.ComputedLayout()
	ax, ay := b.absoluteOffset()
	left, top := sanitizeFloat(ax), sanitizeFloat(ay)
	right, bottom := left+sanitizeFloat(l.Width), top+sanitizeFloat(l.Height)
	return x >= left && x <= right && y >= top && y <= bottom
}

// sanitizeFloat 将 NaN/Inf 等非常规数值转换为 0，避免后续整数转换产生极端值
func sanitizeFloat(f float32) float32 {
	df := float64(f)
	if math.IsNaN(df) || math.IsInf(df, 0) {
		return 0
	}
	return f
}

func (b *BaseComponent) RegisterEventListener(eventType event.Type, listener event.EventListener) {
	slog.Debug("Registering event listener", "componentID", b.ID(), "eventType", eventType)

	root := b.Root()
	if root == nil {
		// Root 尚未就绪，暂存监听器，待 SetRoot/EM 可用时注册
		b.mu.Lock()
		b.pendingListeners[eventType] = append(b.pendingListeners[eventType], listener)
		b.mu.Unlock()
		return
	}
	if rootComp, ok := root.(RootComponent); ok {
		em := rootComp.GetEventManager()
		if em != nil {
			em.RegisterListener(b.ID(), eventType, listener)
			return
		}
		// EM 未就绪，暂存
		b.mu.Lock()
		b.pendingListeners[eventType] = append(b.pendingListeners[eventType], listener)
		b.mu.Unlock()
		return
	}
	// root 非 RootComponent，保底暂存
	b.mu.Lock()
	b.pendingListeners[eventType] = append(b.pendingListeners[eventType], listener)
	b.mu.Unlock()
}

func (b *BaseComponent) OnClick(handler func(event event.MouseEvent)) {
	b.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool {
		me, ok := e.(*event.MouseEvent)
		if !ok || !me.Down {
			return false
		}
		// 直接使用自身的 Bounds 做命中测试，避免遍历组件树
		r := b.Bounds()
		if me.X < float32(r.Min.X) || me.X >= float32(r.Max.X) || me.Y < float32(r.Min.Y) || me.Y >= float32(r.Max.Y) {
			return false
		}
		handler(*me)
		return true
	}))
}

func (b *BaseComponent) OnHover(handler func(hover bool)) {
	b.RegisterEventListener(event.EventMouse, event.EventListenerFunc(func(e event.Event) bool {
		me, ok := e.(*event.MouseEvent)
		if !ok {
			return false
		}
		// 仅处理鼠标移动（Count==0），避免按下/释放时影响 Click 事件
		if me.Count != 0 {
			return false
		}
		// 命中测试：直接使用自身 Bounds
		r := b.Bounds()
		inside := me.X >= float32(r.Min.X) && me.X < float32(r.Max.X) && me.Y >= float32(r.Min.Y) && me.Y < float32(r.Max.Y)

		// 仅在状态变化时触发回调
		b.mu.Lock()
		changed := (inside != b.hovered)
		if changed {
			b.hovered = inside
		}
		b.mu.Unlock()
		if changed {
			handler(inside)
		}
		// 悬停状态变化不应吞掉事件，允许后续 Click/Active 等处理
		return false
	}))
}
