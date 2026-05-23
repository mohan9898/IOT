# Docker 自动构建与部署指南

本项目已配置 GitHub Actions，实现了 Docker 镜像的自动构建和发布！

## ✅ 已完成的配置

- GitHub Actions 工作流文件：`.github/workflows/docker-publish.yml`
- 自动构建触发器已配置

## 🚀 工作原理

### 何时自动构建？

1. **推送到 `main` 分支** → 自动构建并发布 `latest` 标签
2. **创建版本标签（如 `v1.0.0`）** → 自动构建并发布版本标签
3. **提交 PR 到 `main`** → 自动构建（用于验证，不发布）

### 构建产物

Docker 镜像会发布到：
```
ghcr.io/mohan9898/IOT
```

## 📋 使用步骤

### 第一步：提交并推送工作流文件

我们已经创建了工作流文件，现在需要推送到 GitHub：

```bash
cd /workspace/iot-manager

# 添加新文件
git add .github/workflows/docker-publish.yml
git add DOCKER_BUILD_GUIDE.md

# 提交
git commit -m "🔧 Add GitHub Actions for Docker auto build and publish"

# 推送到 GitHub
git push
```

### 第二步：确保 GitHub Packages 权限正确

推送成功后，可能需要检查仓库设置：

1. 访问仓库页面：https://github.com/mohan9898/IOT
2. 点击 **Settings** → **Actions** → **General**
3. 在 **Workflow permissions** 部分，确保：
   - 选择 **Read and write permissions**
   - 勾选 **Allow GitHub Actions to create and approve pull requests**

### 第三步：触发第一次自动构建

推送完成后，GitHub Actions 会自动开始构建！

查看构建状态：
1. 访问仓库页面：https://github.com/mohan9898/IOT
2. 点击 **Actions** 标签
3. 您会看到 "Build and Publish Docker Image" 工作流正在运行

### 第四步：查看发布的镜像

构建成功后，您可以在以下位置找到镜像：

1. 访问：https://github.com/mohan9898/IOT/pkgs/container/IOT
2. 或者直接访问：https://ghcr.io/mohan9898/IOT

## 📦 使用构建好的镜像

### 拉取镜像

```bash
# 拉取最新版本
docker pull ghcr.io/mohan9898/iot:latest

# 或者拉取特定版本
docker pull ghcr.io/mohan9898/iot:main
```

### 运行容器

```bash
# 使用 docker-compose（推荐）
# 修改您的 docker-compose.yml，把 image 指向 ghcr.io/mohan9898/iot:latest

# 或者直接运行
docker run -d \
  -p 6116:6116 \
  -e JWT_SECRET=your-secret-here \
  -e MQTT_USERNAME=your-mqtt-user \
  -e MQTT_PASSWORD=your-mqtt-password \
  ghcr.io/mohan9898/iot:latest
```

## 🏷️ 版本标签说明

工作流会自动生成以下标签：

| 触发方式 | 标签示例 | 说明 |
|---------|---------|------|
| 推送到 main | `main`, `latest` | 最新开发版 |
| PR 提交 | `pr-123` | PR 验证构建 |
| Git tag v1.0.0 | `v1.0.0`, `1.0`, `latest` | 版本发布 |
| 任何提交 | `sha-abc123` | 基于提交哈希 |

### 创建版本标签发布正式版本

```bash
# 创建并推送标签
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

这会自动触发构建并发布带版本号的镜像！

## 🔍 查看构建日志

1. 访问仓库：https://github.com/mohan9898/IOT
2. 点击 **Actions**
3. 点击对应的工作流运行记录
4. 查看构建步骤和日志

## ⚙️ 自定义配置

如果需要修改工作流，编辑：
`.github/workflows/docker-publish.yml`

常见修改：
- 修改触发条件（`on:` 部分）
- 修改构建标签策略
- 添加更多构建步骤

## 📝 完整示例工作流程

```bash
# 1. 修改代码
# ... 进行您的修改 ...

# 2. 提交并推送
git add .
git commit -m "✨ New feature"
git push

# 3. 自动构建开始！
# 访问 Actions 页面查看进度

# 4. 构建成功后，拉取新镜像
docker pull ghcr.io/mohan9898/iot:latest

# 5. 部署新版本
docker-compose up -d
```

## 🔐 权限说明

- `GITHUB_TOKEN` 会自动提供，无需额外配置
- 确保仓库的 Actions 权限设置为 "Read and write"
- GitHub Packages 会自动关联到您的仓库

## 🎉 完成！

现在您的项目已具备完整的 CI/CD 能力！
