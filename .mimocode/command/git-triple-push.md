---
name: git-triple-push
description: 推送到 Gitea + GitHub + Gitee 三个远程仓库
---

# Git 三远程推送

将当前分支推送到三个远程仓库：Gitea、GitHub、Gitee。

## 使用方式

```
git-triple-push [branch]
```

不指定分支时默认推送当前分支。

## 执行步骤

1. 确认当前分支和提交状态
2. 推送到 origin (Gitea)
3. 推送到 github
4. 推送到 gitee

## 远程仓库配置

| 名称 | 地址 |
|------|------|
| origin (Gitea) | http://localhost:3001/hlw/{repo}.git |
| github | https://github.com/hlw422/{repo}.git |
| gitee | https://gitee.com/hlw422/{repo}.git |

## 命令

```bash
git push origin {branch}
git push github {branch}
git push gitee {branch}
```

## 新项目初始化

首次设置远程：
```bash
git remote add origin http://localhost:3001/hlw/{repo}.git
git remote add github https://github.com/hlw422/{repo}.git
git remote add gitee https://gitee.com/hlw422/{repo}.git
```

在 Gitea 创建仓库（需要 API token）：
```bash
curl -X POST "http://localhost:3001/api/v1/user/repos" \
  -H "Content-Type: application/json" \
  -d '{"name": "{repo}", "description": "...", "private": false}'
```

GitHub/Gitee 需要手动在网页创建仓库，或使用 gh/gitee CLI。
