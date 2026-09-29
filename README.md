## KernyrMind

一个Pixso白板的开源平替

---
[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL_3.0-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![翻译状态](http://weblate.yearnstudio.cn/widget/kernyrmind/svg-badge.svg)](http://weblate.yearnstudio.cn/engage/kernyrmind/)

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-Framework-008ECF?logo=go&logoColor=white)](https://gin-gonic.com/)

[![Vue.js](https://img.shields.io/badge/Vue.js-3.x-4FC08D?logo=vuedotjs&logoColor=white)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Vite](https://img.shields.io/badge/Vite-Ready-646CFF?logo=vite&logoColor=white)](https://vitejs.dev/)

---

### 启动方式

生成前端资源

```bash
cd web
pnpm build
```

编译并启动后端

```bash
cd ..
go run cmd/server/main.go
```

#### 生成Swagger文档

先安装CLI工具

```bash
go install github.com/swaggo/swag/cmd/swag@latest

```

重新生成文档

```bash
swag init -g cmd/server/main.go
```

项目已嵌入Swagger路由, 直接打开`localhost:8080/swagger/index.html`即可访问(**注意**: 仅在调试模式并且需要配置`port=8080`)
