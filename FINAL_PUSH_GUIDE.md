# 最终推送指南

我们已经配置了远程仓库地址：`https://github.com/mohan9898/IOT.git`

现在需要进行身份认证才能推送。以下是几种方式：

## 方式一：使用个人访问令牌（PAT）推荐

### 1. 创建 GitHub 个人访问令牌

1. 访问 https://github.com/settings/tokens
2. 点击 **Generate new token** → **Generate new token (classic)**
3. 填写：
   - Note: `IoT Manager Push`
   - Expiration: 选择合适的过期时间
   - Scopes: 勾选 `repo`（完整仓库权限）
4. 点击 **Generate token**
5. **重要！立即复制保存这个 token，只显示一次！**

### 2. 配置 Git 使用令牌

执行以下命令，替换 `YOUR_TOKEN_HERE` 为您刚复制的 token：

```bash
cd /workspace/iot-manager

# 方式 A：直接在 URL 中包含 token（临时）
git remote set-url origin https://YOUR_TOKEN_HERE@github.com/mohan9898/IOT.git

# 然后推送
git push -u origin main
```

或者使用凭证助手永久保存：

```bash
# 方式 B：使用 git credential helper
git config --global credential.helper store

# 然后执行推送，第一次会要求输入用户名和 token
git push -u origin main
# Username: YOUR_GITHUB_USERNAME
# Password: YOUR_TOKEN_HERE
```

## 方式二：使用 SSH（如果您已配置 SSH key）

如果您在 GitHub 上配置了 SSH 密钥，可以使用 SSH 方式：

```bash
cd /workspace/iot-manager

# 切换到 SSH 地址
git remote set-url origin git@github.com:mohan9898/IOT.git

# 推送
git push -u origin main
```

## 方式三：先拉取（如果仓库已有内容）

如果您的 GitHub 仓库 `mohan9898/IOT` 已经有内容（如 README 文件），需要先拉取：

```bash
cd /workspace/iot-manager

# 先拉取远程内容
git pull origin main --allow-unrelated-histories

# 解决可能的合并冲突后
git push -u origin main
```

## 验证推送成功

推送完成后，访问：
https://github.com/mohan9898/IOT

您应该能看到所有文件已经上传！

## 当前状态

- ✅ 本地仓库已初始化
- ✅ 两次提交已完成
- ✅ 远程仓库已配置为 `https://github.com/mohan9898/IOT.git`
- ⏳ 等待身份认证和推送

请按照以上任一方式完成最后的推送！
