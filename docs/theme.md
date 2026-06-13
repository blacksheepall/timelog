# 主题系统扩展指南

## 系统架构

当前主题系统基于 **CSS 自定义属性** + **Tailwind CSS 插件** 构建，支持任意数量的主题切换。核心设计原则：

- **语义化 token**：组件只关心「这是什么元素」，不关心「它是什么颜色」
- **集中管理**：所有颜色值在 `style.css` 中定义，组件零硬编码
- **运行时切换**：通过 `document.documentElement.classList` 切换主题，无需刷新页面

## 文件结构

```
web/src/
├── style.css              # 主题变量定义（所有主题在此）
├── tailwind.config.js     # Token → Tailwind 工具类映射
├── App.vue                # 主题切换逻辑 + 切换按钮
└── index.html             # FOUC 防护（优先加载脚本）
```

## Token 层级

### 背景层级（Background Layers）

| Token | 语义 | Light 值 | Dark 值 | 使用场景 |
|-------|------|----------|---------|----------|
| `--color-bg-base` | 页面底层 | `#f3f4f6` | `#0f172a` | `body` 背景 |
| `--color-bg-surface` | 卡片/面板 | `#ffffff` | `#1e293b` | 卡片、弹窗、表单项 |
| `--color-bg-elevated` | 悬停/高亮 | `#f9fafb` | `#334155` | 悬停态、列表 zebra |
| `--color-bg-overlay` | 浮层 | `#ffffff` | `#1e293b` | 下拉菜单、tooltip |

### 文本层级（Text Layers）

| Token | 语义 | Light 值 | Dark 值 | 使用场景 |
|-------|------|----------|---------|----------|
| `--color-text-primary` | 主文本 | `#111827` | `#f1f5f9` | 标题、正文 |
| `--color-text-secondary` | 次要文本 | `#6b7280` | `#94a3b8` | 描述、时间戳 |
| `--color-text-tertiary` | 辅助文本 | `#9ca3af` | `#64748b` | 占位符、禁用态 |
| `--color-text-muted` | 弱化文本 | `#374151` | `#cbd5e1` | 按钮文字、label |

### 边框层级（Border Layers）

| Token | 语义 | Light 值 | Dark 值 | 使用场景 |
|-------|------|----------|---------|----------|
| `--color-border-default` | 默认边框 | `#e5e7eb` | `#334155` | 卡片边框、分割线 |
| `--color-border-subtle` | 弱边框 | `#f3f4f6` | `#1e293b` | 内部分隔、细线 |

### 品牌/功能色（Brand & Functional）

| Token | 语义 | Light 值 | Dark 值 | 使用场景 |
|-------|------|----------|---------|----------|
| `--color-brand` | 品牌主色 | `#2563eb` | `#3b82f6` | 主按钮、链接 |
| `--color-brand-hover` | 品牌悬停 | `#1d4ed8` | `#60a5fa` | 按钮 hover |
| `--color-success` | 成功 | `#16a34a` | `#4ade80` | 成功提示 |
| `--color-danger` | 危险 | `#dc2626` | `#f87171` | 删除、错误 |
| `--color-success-bg` | 成功背景 | `#f0fdf4` | `#052e16` | 成功通知背景 |
| `--color-danger-bg` | 危险背景 | `#fef2f2` | `#450a0a` | 错误通知背景 |
| `--color-brand-bg` | 品牌背景 | `#eff6ff` | `#172554` | 选中态、导航激活 |

## 新增主题的步骤

以新增 **「高对比度（High Contrast）」** 主题为例：

### 步骤 1：在 `style.css` 中定义主题变量

```css
@layer base {
  :root { /* 默认亮色，已存在 */ }
  .dark { /* 暗色，已存在 */ }

  /* 新增主题：高对比度 */
  .high-contrast {
    /* 背景 */
    --color-bg-base: #000000;
    --color-bg-surface: #000000;
    --color-bg-elevated: #1a1a1a;
    --color-bg-overlay: #000000;

    /* 文本 —— 高对比度要求：前景与背景对比度 ≥ 7:1 */
    --color-text-primary: #ffffff;
    --color-text-secondary: #e5e5e5;
    --color-text-tertiary: #a0a0a0;
    --color-text-muted: #d0d0d0;

    /* 边框 */
    --color-border-default: #ffffff;
    --color-border-subtle: #666666;

    /* 品牌 —— 高对比度下用鲜艳色 */
    --color-brand: #00ff00;      /* 荧光绿 */
    --color-brand-hover: #00cc00;
    --color-success: #00ff00;
    --color-danger: #ff0000;

    /* 状态背景 */
    --color-success-bg: #001a00;
    --color-success-border: #00ff00;
    --color-danger-bg: #1a0000;
    --color-danger-border: #ff0000;
    --color-brand-bg: #001a00;

    /* 阴影 */
    --shadow-sm: 0 1px 2px 0 rgb(255 255 255 / 0.1);
    --shadow-md: 0 4px 6px -1px rgb(255 255 255 / 0.2);
  }
}
```

### 步骤 2：注册主题到类型系统

在 `web/src/stores/settings.ts` 中扩展类型：

```typescript
export interface SettingsState {
  // ... 现有字段
  theme: 'light' | 'dark' | 'auto' | 'high-contrast'  // 新增主题名
}
```

在 `web/src/composables/useSettings.ts` 中同步：

```typescript
const theme = computed({
  get: () => settings.value.theme,
  set: value => settingsStore.updateSetting('theme', value),
})
```

### 步骤 3：添加切换逻辑到 `App.vue`

找到 `themeCycle` 函数，把新主题加入循环：

```typescript
const themeCycle = () => {
  const cycle = ['light', 'dark', 'auto', 'high-contrast'] as const  // 新增
  const current = theme.value
  const next = cycle[(cycle.indexOf(current) + 1) % cycle.length]
  updateSetting('theme', next)
  applyTheme(next)  // 复用已有的 applyTheme
}

const themeIcon = computed(() => {
  if (theme.value === 'dark') return MoonIcon
  if (theme.value === 'light') return SunIcon
  if (theme.value === 'high-contrast') return EyeIcon  // 新增图标
  return ComputerDesktopIcon
})

const themeLabel = computed(() => {
  if (theme.value === 'dark') return '暗色'
  if (theme.value === 'light') return '亮色'
  if (theme.value === 'high-contrast') return '高对比'  // 新增
  return '自动'
})
```

### 步骤 4：更新 `index.html` 的 FOUC 防护

在 `<head>` 的内联脚本中，确保新主题在页面加载时立即生效：

```html
<script>
  (function() {
    const theme = localStorage.getItem('timelog-settings');
    if (theme) {
      try {
        const parsed = JSON.parse(theme);
        const t = parsed.theme || 'auto';
        const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        
        // 新主题直接应用，不需要判断系统
        if (t === 'high-contrast') {
          document.documentElement.classList.add('high-contrast');
        } else if (t === 'dark' || (t === 'auto' && systemDark)) {
          document.documentElement.classList.add('dark');
        }
      } catch (e) {}
    }
  })();
</script>
```

### 步骤 5：验证构建

```bash
cd web && npx vite build
```

---

## 自定义主题色（品牌色）

如果你想让主题不仅仅是深浅变化，而是改变**品牌色**（比如把蓝色改成橙色），有两种方式：

### 方式 A：覆盖品牌变量（推荐）

在 `style.css` 中，只覆盖品牌相关的变量，保留背景和文本结构：

```css
.orange-theme {
  /* 继承暗色/亮色的背景/文本，只改品牌色 */
  --color-brand: #f97316;      /* orange-500 */
  --color-brand-hover: #ea580c; /* orange-600 */
  --color-brand-bg: #fff7ed;    /* orange-50 */
  --color-success: #22c55e;
  --color-danger: #ef4444;
}
```

> ⚠️ 注意：如果主题是继承亮色/暗色的，需要把 `.orange-theme` 和 `light`/`dark` 一起加在 `html` 上：`class="dark orange-theme"`。这要求设计时考虑**叠加组合**。

### 方式 B：完整独立主题（更彻底）

像上面的「高对比度」示例一样，定义完整的变量集合。这种方式完全独立，不依赖其他主题。

---

## 设计原则

### 1. WCAG 对比度

| 使用场景 | 最小对比度 | 推荐对比度 |
|----------|-----------|-----------|
| 正文文本（14px+） | 4.5:1 | 7:1 |
| 大号文本（18px+ / bold 14px+） | 3:1 | 4.5:1 |
| 图标/装饰元素 | 3:1 | - |

工具：
- [WebAIM 对比度检查器](https://webaim.org/resources/contrastchecker/)
- [APCA 对比度](https://www.myndex.com/APCA/)

### 2. 颜色一致性

- 同一主题内，**背景色**必须是连续的亮度阶梯（base → surface → elevated，越来越亮或越来越暗）
- **文本色**也必须是连续的（primary → secondary → tertiary，越来越弱）
- 暗色模式下不要简单反转所有颜色，而是整体降亮，保留品牌色鲜艳度

### 3. 测试清单

新增主题后，至少检查这些页面：

- [ ] Dashboard（Home.vue）— 卡片、图表、进度条
- [ ] TimeLog 列表 — 表格、筛选按钮、表单
- [ ] 任务页 — 状态标签、复选框
- [ ] 分类页 — 树形节点、颜色展示
- [ ] 约束页 — 日历热力图、状态指示
- [ ] 登录/注册 — 渐变背景、卡片、表单
- [ ] 通知 toast — 成功/错误状态
- [ ] 移动端菜单 — 深色背景下的菜单

### 4. 调试技巧

在浏览器 DevTools 中，直接编辑 `html` 的 `class` 来切换主题：

```javascript
// 切换到暗色
document.documentElement.classList.remove('light', 'high-contrast')
document.documentElement.classList.add('dark')

// 切换到高对比度
document.documentElement.classList.remove('light', 'dark')
document.documentElement.classList.add('high-contrast')
```

在 DevTools 的 Elements → Styles 面板中，可以看到所有 CSS 变量的当前值。

---

## 常见问题

**Q: 能不能让用户自定义颜色？**

A: 可以。在 `style.css` 中用 JavaScript 动态修改变量：

```javascript
// 用户选了橙色作为品牌色
document.documentElement.style.setProperty('--color-brand', '#f97316')
document.documentElement.style.setProperty('--color-brand-hover', '#ea580c')
```

这种方式不需要定义新主题，直接覆盖变量即可。但持久化需要存到 `localStorage` 并在 `index.html` 的 FOUC 脚本中恢复。

**Q: 组件里还有 `style="background-color: #xxx"` 怎么办？**

A: 替换为 CSS 变量：

```html
<!-- Before -->
<div :style="{ backgroundColor: category.color }">

<!-- After —— 如果颜色是动态的，保持原样；如果是硬编码，改用变量 -->
<div class="bg-brand">  <!-- 固定品牌色 -->
<div :style="{ backgroundColor: 'var(--color-brand)' }">  <!-- 动态但绑定变量 -->
```

**Q: 新主题只改背景色，不改文本色？**

A: 不推荐。如果新主题的背景色和文本色对比度不够，会导致可读性问题。必须同时调整文本色保证对比度 ≥ 4.5:1。

**Q: 如何在组件内使用 `dark:` 变体？**

A: 既然已经用了 CSS 变量，就不需要 `dark:` 了。Tailwind 的 `dark:` 变体只有在 `darkMode: 'class'` 时才生效，但我们的 token 已经封装了 light/dark 差异。组件直接写 `text-text-primary` 即可，自动适配当前主题。

---

## 参考

- [Tailwind CSS 自定义颜色](https://tailwindcss.com/docs/customizing-colors)
- [CSS 自定义属性](https://developer.mozilla.org/en-US/docs/Web/CSS/--*)
- [WCAG 2.1 对比度标准](https://www.w3.org/WAI/WCAG21/Understanding/contrast-minimum)
- [shadcn/ui 主题设计](https://ui.shadcn.com/themes)
