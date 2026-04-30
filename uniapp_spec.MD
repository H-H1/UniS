**uni-app 项目开发规范文档（Frontend Spec）**，使用 **Markdown 格式**，可直接保存为 `.md` 文件。

这份规范适用于：
- Vue 2 / Vue 3 版本的 uni-app
- 跨端开发（App、小程序、H5）

---

```markdown
# uni-app 项目开发规范文档（Frontend Spec）

## 1. 项目概述

- **项目名称**：____________（请填写）
- **技术栈**：uni-app + Vue 3 / Vue 2 + Vite / webpack + Pinia / Vuex + uView / uni-ui
- **跨端平台**：App（iOS/Android） / 微信小程序 / H5（根据实际勾选）
- **核心目标**：多端复用、高性能、可维护

---

## 2. 项目结构

```
├── pages/                  # 页面目录（主包）
│   └── index/
│       └── index.vue
├── pagesSub/               # 分包目录（可选）
├── components/             # 通用组件
│   ├── base/               # 基础组件（按钮、输入框等）
│   └── business/           # 业务组件
├── api/                    # 接口请求封装
│   ├── modules/            # 按模块划分的接口
│   └── request.js          # 请求拦截器
├── store/                  # 状态管理（Pinia / Vuex）
├── utils/                  # 工具函数
│   ├── index.js
│   └── storage.js          # 本地存储封装
├── static/                 # 静态资源（图片、字体等）
├── uni_modules/            # uni-app 插件市场模块
├── styles/                 # 全局样式 / 主题变量
│   └── variables.scss
├── manifest.json           # 应用配置
├── pages.json              # 页面路由与窗口配置
├── uni.scss                # 全局 SCSS 变量
├── App.vue                 # 应用入口
├── main.js                 # 入口文件
└── vite.config.js          # 构建配置（Vite）
```

---

## 3. 编码规范

### 3.1 Vue 组件规范

- 组件命名：
  - 页面：`kebab-case`（文件 `user-info.vue`）
  - 通用组件：`BaseXxx` 或 `TheXxx`
- 组件结构顺序（推荐）：
  ```vue
  <template>...</template>
  <script setup>  // 或 export default {}
  <style lang="scss" scoped>
  ```
- Props 必须声明类型和默认值
- 事件命名使用 `kebab-case`（`@update-value`）

### 3.2 样式规范

- 使用 `scoped` 避免样式污染
- 优先使用 `rpx`（响应式像素）适配多端
- 全局变量定义在 `uni.scss` 或 `styles/variables.scss`
- 禁止使用 `!important`（除非覆盖第三方样式）

### 3.3 JS/TS 规范

- 使用 `ES6+` 语法（`const/let`、箭头函数、解构、模板字符串）
- 推荐使用 TypeScript（类型安全）
- 禁止 `var`，统一使用 `const` / `let`
- 异步优先使用 `async/await`，避免回调地狱
- 每个函数应有明确注释（复杂逻辑）

---

## 4. API 请求规范

- 统一在 `api/request.js` 中封装请求拦截器：
  - 自动添加 token
  - 统一 loading 管理
  - 统一错误处理（如 401 跳转登录）
- 接口按模块放在 `api/modules/` 下：
  ```js
  // api/modules/user.js
  import request from '@/api/request'
  export const getUserInfo = (id) => request.get(`/user/${id}`)
  ```
- 所有请求方法返回 Promise

---

## 5. 路由与页面跳转

- 使用 uni-app 原生路由：
  - `uni.navigateTo`（保留当前页面）
  - `uni.redirectTo`（关闭当前页面）
  - `uni.switchTab`（tab 页）
  - `uni.reLaunch`（关闭所有页面）
- 路由路径使用绝对路径（以 `/` 开头）
- 页面参数传递：
  - 简单参数通过 query 传递
  - 复杂数据使用全局 store 或 `eventChannel`

---

## 6. 状态管理（Pinia / Vuex）

- 推荐 **Pinia**（Vue 3 + 更好的 TS 支持）
- 按模块拆分 store：
  ```
  store/
    modules/
      user.js
      cart.js
    index.js
  ```
- 禁止在组件中直接修改 state，必须通过 actions

---

## 7. 多端兼容规范

- 条件编译：使用 `#ifdef` / `#ifndef` 处理平台差异
  ```vue
  <!-- #ifdef MP-WEIXIN -->
  <button open-type="getPhoneNumber">微信手机号</button>
  <!-- #endif -->
  ```
- 尽量使用 uni-app 跨端 API（如 `uni.setStorage` 代替 `localStorage`）
- 样式差异通过 **不同端样式文件** 或 **条件编译样式** 解决

---

## 8. 性能优化规范

- 分包加载：
  - 主包 ≤ 2MB（小程序限制）
  - 非首屏页面放入分包
- 图片优化：
  - 使用 WebP 格式（App、H5 支持）
  - 小图标使用 iconfont / SVG
  - 大图使用 CDN + 懒加载
- 列表渲染：
  - `v-for` 必须提供 `key`
  - 长列表使用虚拟滚动（如 `recycle-view`）
- 避免在 `onLoad` / `onShow` 中执行耗时同步操作

---

## 9. 安全规范

- 敏感信息（token、密钥）不硬编码，使用环境变量或加密存储
- 禁止将用户密码明文存储（本地或传输）
- 小程序 webview 需配置业务域名白名单
- 网络请求必须使用 HTTPS（生产环境）

---

## 10. 测试规范

- 单元测试：推荐 `vitest`（Vite 项目）
- 端到端测试（可选）：`miniprogram-automator`（微信小程序）
- 真机调试：至少覆盖 iOS + Android 各一款主流机型

---

## 11. 构建与发布

- 环境区分：
  - 开发环境：`npm run dev:mp-weixin` / `dev:h5`
  - 生产环境：`npm run build:mp-weixin` / `build:h5`
- 环境变量：
  - 使用 `.env.dev` / `.env.prod` 文件（Vite 模式）
- 发布前检查：
  - 取消 console.log（可使用 `vite-plugin-remove-console`）
  - 检查分包大小
  - 小程序上传前执行 `eslint` + `prettier`

---

## 12. 文档要求

- `README.md`：项目介绍 / 启动方式 / 环境配置
- `docs/`：
  - 接口文档（可选对接 Apifox / YApi）
  - 组件库使用说明
  - 常见问题（FAQ）

---

## 13. 版本管理

- 语义化版本：`v主版本.次版本.修订号`
- Git 提交规范（Conventional Commits）：
  - `feat:` 新功能
  - `fix:` 修复
  - `style:` 样式调整
  - `docs:` 文档
  - `perf:` 性能

---

## 14. 工具链推荐

| 工具           | 用途                        |
|----------------|-----------------------------|
| HBuilderX      | uni-app 官方 IDE             |
| VSCode         | 替代 IDE（安装 uni-app 插件）|
| eslint + prettier | 代码规范检查与格式化        |
| husky + lint-staged | Git 提交前检查           |
| Sass / Less    | CSS 预处理器                |

---

## 15. 附则

- 本规范由项目组共同维护
- 如遇特殊场景与规范冲突，需在代码注释或文档中说明原因
- 规范每季度 review 一次，持续改进
```

---

你可以直接将上面的完整 Markdown 内容复制到 `uniapp_spec.md` 文件中。  
如需针对 Vue2/Vue3、特定组件库（如 uView）或 App 端进行个性化调整，请告诉我。