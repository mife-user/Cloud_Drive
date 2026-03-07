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

## ___Docker 部署___

### 1. 生成交叉编译的 Linux 可执行文件

在 **Windows** 上运行以下命令：

```bash

# 方法 2: 手动构建
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o cloud_drive ./cmd/main/main.go
```

### 2. 构建 Docker 镜像

```bash
docker build -t cloudpan:latest .
```

### 3. 运行服务

```bash
# 使用 docker-compose 启动所有服务（推荐）
docker-compose up -d

# 或单独启动应用
docker run -d -p 8080:8080 --name cloudpan-app cloudpan:latest
```

### 4. 停止服务

```bash
docker-compose down -v
```

## 前端测试链接
[测试链接](https://github.com/mife-user/Cloud_Drive/blob/main/%E5%89%8D%E7%AB%AF%E6%B5%8B%E8%AF%95.md)

