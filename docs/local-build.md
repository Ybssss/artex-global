# 本地构建与部署

本文件说明如何从本仓库源码构建镜像，并让 `docker compose` 使用它。

## 为什么需要这一步

上游 `docker-compose.yml` 固定引用预构建镜像：

```yaml
image: ghcr.io/jiwoochris/artex-global:${ARTEX_TAG:-latest}
```

而 `Dockerfile` 只做 `COPY dist/${TARGETARCH}/artex`——二进制由 CI 单独编译后放进
`dist/`，这个目录又被 `.gitignore` 排除。结果是在本仓库改了 Go 代码之后：

```bash
docker compose up -d          # 重建容器
```

容器会用**镜像里的旧二进制**覆盖 `/app/artex`，你手工 `docker cp` 进去的修复随之丢失。
这不是 bug，是上游「运行期镜像 / 构建产物分离」设计的正常行为——但对本地开发很致命。

## 解决方案

`Dockerfile.build` 在保留上游运行期层（工具链、Playwright、`start.sh` 全部照搬）的前提下，
前面加了两个构建阶段：

| 阶段 | 作用 |
|---|---|
| `web` | `npm ci` + `NEXT_EXPORT=1 next build`，产出静态站点 |
| `build` | `go build -tags embedui`，把前端内嵌进单二进制 |
| 最终 | 与上游 `Dockerfile` 完全一致的运行期镜像 |

于是 `docker build` 一步到位，不依赖预置的 `dist/`。

`docker-compose.local.yml` 是一个覆盖层，把 `image` 指向本地构建的标签：

```yaml
services:
  artex:
    image: artex-global:local
    pull_policy: never
```

## 构建

```bash
docker build -f Dockerfile.build -t artex-global:local .
```

首次构建约 5-10 分钟（主要是 Playwright 与 Chromium 下载）。后续改了 Go 代码，
只有 `build` 阶段之后会重跑，通常 1-2 分钟。

## 启用

两种方式，任选其一。

### 方式一：每次显式指定

```bash
docker compose -f docker-compose.yml -f docker-compose.local.yml up -d
```

### 方式二：设为默认（推荐）

在 `.env` 里加一行：

```bash
# Linux / macOS
COMPOSE_FILE=docker-compose.yml:docker-compose.local.yml
```

```bash
# Windows（分隔符是分号）
COMPOSE_FILE=docker-compose.yml;docker-compose.local.yml
```

之后平常的 `docker compose up -d` 就会自动带上覆盖层。

> `.env` 已被 `.gitignore` 排除，这一行不会进版本库。

## 验证修复已固化

重建容器后确认二进制来自镜像而非旧版：

```bash
docker compose up -d --force-recreate artex
docker exec artex-global-artex-1 ls -la /app/artex
```

期望看到的大小应与 `Dockerfile.build` 编出来的产物一致（`-ldflags "-s -w"` 剥离符号后
约 40 MB）。如果看到明显不同的旧大小，说明覆盖层没有生效——检查 `COMPOSE_FILE`
的分隔符，或用 `docker compose config | grep image` 确认解析结果。

## 更新流程

改了代码之后：

```bash
docker build -f Dockerfile.build -t artex-global:local .
docker compose up -d --force-recreate artex
```

不需要再 `docker cp`。

## 注意

- `docker-compose.local.yml` 这个名字没被 `.gitignore` 的 `*.local` 规则命中，
  所以能正常提交。`Dockerfile` 最初想叫 `Dockerfile.local`，但会被 `*.local` 忽略，
  因此改名为 `Dockerfile.build`。
- 镜像约 2.3 GB，主要是 Playwright + Chromium。若不需要浏览器能力，可以删掉
  运行期阶段的 Playwright 安装层，镜像会降到约 300 MB。
- `.dockerignore` 已排除 `data/`（1.2 GB）、`web/node_modules`（735 MB）等，
  构建上下文很小，不会把业务数据打进镜像。
