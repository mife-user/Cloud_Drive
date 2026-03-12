# Cloud_Drive
Lanshan寒假考核
# ___功能___
## 用户系统
- [x] 注册
- [x] 登录
- [x] 头像
- [ ] 个人简介
- [x] 用户鉴权
- [x] 信息修改
## 文件系统
- [x] 文件上传（非切片）
- [x] 修改文件权限
- [x] 文件删除（软）
- [x] 文件删除（硬）
- [x] 文件分享
- [x] 访问分享内容
- [x] 文件预览
- [ ] 文件搜索
- [x] 收藏文件
- [x] 取消收藏
- [x] 查看所有收藏文件
- [ ] 文件标签
- [ ] 敏感内容检查
- [x] 回收站
- [x] 收藏
## 会员
- [x] 限速
- [x] 空间额度
## 其他
### 日志
- [x] zap
### 配置
- [x] viper
# ___技术___
- [x] MYSQL
- [x] Redis
  - [x] 缓存击穿(singleflight)
  - [x] 缓存雪崩(随机过期时间)
  - [x] 缓存穿透(空对象缓存)
- [x] Docker

---

## ___环境要求___

- Go 1.21+
- MySQL 8.0
- Redis 7.x
- Node.js 18+ (前端开发)
- Docker & Docker Compose (生产部署)

---

## ___开发环境___

### 后端启动

**方式一：设置环境变量**

```powershell
# Windows PowerShell
$env:CLOUDPAN_ENV="dev"
$env:CLOUDPAN_MYSQL_DSN="root:123456@tcp(127.0.0.1:3306)/cloudpan?charset=utf8mb4&parseTime=True&loc=Local"
$env:CLOUDPAN_REDIS_HOST="127.0.0.1"
$env:CLOUDPAN_JWT_SECRET="your-secret-key"

go run ./cmd/main/main.go
```

**方式二：修改配置文件**

编辑 `configs/config.dev.yaml`，配置本地 MySQL 和 Redis 连接信息，然后：

```powershell
$env:CLOUDPAN_ENV="dev"
go run ./cmd/main/main.go
```

后端服务运行在 `http://localhost:8080`

### 前端启动

```powershell
cd web
npm install
npm run dev
```

前端开发服务运行在 `http://localhost:5173`

---

## ___生产环境 (Docker 部署)___

### 一键启动

```powershell
# 1. 生成交叉编译的 Linux 可执行文件
$env:CGO_ENABLED=0
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o cloud_drive ./cmd/main/main.go

# 2. 构建并启动所有服务（前端 + 后端 + MySQL + Redis）
docker-compose up -d --build

# 3. 查看服务状态
docker-compose ps
```

### 服务端口

| 服务 | 端口 | 访问地址 |
|------|------|----------|
| 前端 | 80 | `http://localhost` |
| 后端 API | 8080 | `http://localhost:8080` |
| MySQL | 3306 | - |
| Redis | 6379 | - |

### 停止服务

```powershell
docker-compose down      # 仅停止服务
docker-compose down -v   # 停止并删除数据卷
```

---

## ___环境变量说明___

| 变量名 | 说明 | 示例 |
|--------|------|------|
| `CLOUDPAN_ENV` | 环境标识 | `dev` / `prod` |
| `CLOUDPAN_MYSQL_DSN` | MySQL 连接串 | `root:pass@tcp(host:3306)/db?charset=utf8mb4&parseTime=True&loc=Local` |
| `CLOUDPAN_REDIS_HOST` | Redis 地址 | `127.0.0.1` |
| `CLOUDPAN_REDIS_PASSWORD` | Redis 密码 | `your-password` |
| `CLOUDPAN_JWT_SECRET` | JWT 密钥 | `your-secret-key` |

---

## 前端测试链接
[测试链接](https://github.com/mife-user/Cloud_Drive/blob/main/%E5%89%8D%E7%AB%AF%E6%B5%8B%E8%AF%95.md)

