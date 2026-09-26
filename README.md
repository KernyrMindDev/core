## KernyrMind

一个Pixso白板的开源平替

---

### 技术栈

前端: Vue+TypeScript  
后端: Go+Gin  
API文档生成: Swagger([swaggo](https://github.com/swaggo/swag))

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
