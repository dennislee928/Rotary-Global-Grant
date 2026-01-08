# CI/CD Workflow 說明

## Docker 映像構建與推送

本專案使用 GitHub Actions 自動構建並推送 Docker 映像到 GitHub Container Registry (ghcr.io)。

### 觸發條件

Workflow 會在以下情況自動執行：

1. **推送到 main/master 分支**：當 `apps/api/` 或 `apps/web/` 目錄下的檔案變更時
2. **Pull Request**：建立或更新 PR 時（僅構建，不推送）
3. **手動觸發**：在 GitHub Actions 頁面手動執行

### 構建的映像

每次執行會構建兩個 Docker 映像：

- `ghcr.io/<owner>/rotary-grant-api`
- `ghcr.io/<owner>/rotary-grant-web`

### 標籤格式

每個映像會使用以下標籤：

1. **日期後綴標籤**：`YYYYMMDD`（例如：`20250108`）
2. **日期 + SHA 標籤**：`YYYYMMDD-<short-sha>`（例如：`20250108-a1b2c3d`）
3. **latest 標籤**：`latest`

### 使用範例

#### 拉取映像

```bash
# 使用日期後綴
docker pull ghcr.io/<owner>/rotary-grant-api:20250108

# 使用 latest
docker pull ghcr.io/<owner>/rotary-grant-api:latest
```

#### 在 docker-compose 中使用

```yaml
services:
  api:
    image: ghcr.io/<owner>/rotary-grant-api:20250108
    # 或使用 latest
    # image: ghcr.io/<owner>/rotary-grant-api:latest
```

### 權限設定

首次使用時，需要確保：

1. GitHub repository 已啟用 GitHub Packages
2. 如果需要公開映像，在 repository settings > Packages 中設定可見性

### 多平台支援

映像會構建以下平台版本：

- `linux/amd64`
- `linux/arm64`

### 快取機制

使用 GitHub Actions cache 來加速後續構建，減少構建時間。

### 查看構建結果

1. 前往 GitHub repository 的 **Actions** 標籤
2. 選擇 **Build and Push Docker Images** workflow
3. 查看構建日誌和輸出的映像標籤
