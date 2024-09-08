# Simple Bank

此專案是一個後端 API 服務，使用 Go 語言開發，並實現了簡單的交易功能。專案架構基於 gin 作為主要的 Web 框架，結合 PostgreSQL 資料庫進行資料儲存，並透過 JWT 和 Paseto 實現用戶認證與授權。專案已實現 CI/CD，包括自動化測試和部署至 AWS。

## 目錄

- [技術架構](#技術架構)
- [安裝與設定](#安裝與設定)
- [使用指令](#使用指令)
- [CI/CD 配置](#cicd-配置)


## 技術架構

- **Go**：專案使用 Go 語言作為核心開發語言。
- **Gin**：使用 Gin 作為 Web 框架，提供高效能的 API 服務。
- **PostgreSQL**：作為主要的關聯式資料庫。
- **sqlc**：自動生成 Go 程式碼，直接對應 SQL 查詢語句，提升開發效率並減少潛在錯誤。
- **JWT / Paseto**：用於實現用戶認證與授權機制，增強系統安全性。
- **Testify 和 Mock**：使用 testify 進行單元測試和斷言，並使用 mock 創建模擬物件，以便測試中替代真實依賴，確保程式碼穩定性和可靠性。
- **Docker**：容器化專案，便於部署與維護。
- **Kubernetes (EKS)**：透過 Amazon EKS 部署應用程式，實現高可用性與彈性擴展。
- **CI/CD**：透過 GitHub Actions 完成自動化測試、建置與部署流程，確保程式碼品質。


## 安裝與設定

1. **複製專案**
    
    ```bash
    git clone https://github.com/LucyYeung/simplebank.git
    ```
    
2. **安裝依賴套件**
    
    使用 Go Modules 管理依賴：
    
    ```bash
    go mod tidy
    ```
3. **修改 `app.env` 檔案**
    將下列變數填入環境變數中。 
    ```bash
    DB_DRIVER=postgres
    DB_SOURCE=postgresql://root:password@localhost:5432/simple_bank?sslmode=disable
    SERVER_ADDRESS=0.0.0.0:8080
    TOKEN_SYMMETRIC_KEY=your_symmetric_key
    ACCESS_TOKEN_DURATION=15m
    ```
4. **設定資料庫**
    
    安裝並啟動 PostgreSQL，然後創建資料庫，並將資料庫連接資訊填入環境變數。

    ```bash
    docker compose up
    ```


## 使用指令

- **啟動專案**
    
    ```bash
    make server
    ```
    
- **執行測試**
    
    ```bash
    make test
    ```
    
- **資料庫遷移**
    
    ```bash
    make migrateup
    ```
    

## CI/CD 配置

本專案使用 GitHub Actions 進行 CI/CD，包含測試和自動部署兩個工作流程。

### 測試工作流程

`.github/workflows/test.yml` 會在每次推送或合併請求時觸發，並執行以下步驟：

1. 設定 Go 環境並安裝依賴。
2. 啟動 PostgreSQL 服務並運行資料庫遷移。
3. 執行單元測試。

### 部署工作流程

部署工作流程定義在 `.github/workflows/deployment.yml` 檔案中，當程式碼推送至 `master` 分支時，自動觸發部署至 AWS 環境。此流程會進行映像檔的建置、推送至 **Amazon ECR**，並透過 kubectl 部署到 **Amazon EKS**。


1. **Checkout Repo**
    - **步驟名稱**: Checkout repo
    - **動作**: 使用 `actions/checkout@v3` 拉取最新的程式碼。
    - **目的**: 確保工作流程運行時，能夠取得最新的專案程式碼，以便後續步驟進行建置和部署。
2. **安裝 kubectl**
    - **步驟名稱**: Install kubectl
    - **動作**: 使用 `azure/setup-kubectl@v3` 安裝指定版本的 `kubectl` 工具。
    - **目的**: `kubectl` 是 Kubernetes 的命令行工具，透過它可以管理 Kubernetes 叢集上的應用程式，執行部署、調整擴展以及查看執行狀態等。
3. **配置 AWS 認證**
    - **步驟名稱**: Configure AWS credentials
    - **動作**: 使用 `aws-actions/configure-aws-credentials@v4` 配置 AWS 認證，並設定執行身份角色與 AWS 區域。
    - **參數**:
        - `role-to-assume`: 指定需要扮演的 IAM 角色，例如 `simplebank-github-actions-role`。
        - `aws-region`: 設定部署的 AWS 區域，如 `ap-northeast-1`。
    - **目的**: 為後續步驟提供 AWS 的訪問權限，讓 GitHub Actions 能夠與 AWS 服務進行互動。
4. **登入 Amazon ECR**
    - **步驟名稱**: Login to Amazon ECR
    - **動作**: 使用 `aws-actions/amazon-ecr-login@v2` 登入 AWS ECR（Elastic Container Registry）。
    - **目的**: 允許接下來的步驟能夠推送 Docker 映像檔至 Amazon ECR，ECR 是 AWS 的 Docker 映像檔儲存庫。
5. **讀取機密並保存到環境檔案**
    - **步驟名稱**: Load secrets and save to app.env
    - **動作**: 使用 AWS Secrets Manager 讀取儲存的機密資訊，並將其保存為 `app.env` 檔案。
    - **目的**: 讀取並保存應用程式執行所需的環境變數，例如資料庫密碼、API 金鑰等敏感資料，以確保部署的應用程式有正確的配置。
6. **建置、標籤並推送 Docker 映像檔至 Amazon ECR**
    - **步驟名稱**: Build, tag, and push docker image to Amazon ECR
    - **動作**: 執行 Docker 命令建置映像檔，並加上最新的 Git SHA 作為標籤，最後推送至 ECR。
    - **環境變數**:
        - `REGISTRY`: ECR 註冊表地址，由前一個步驟的輸出獲取。
        - `REPOSITORY`: 專案的映像檔儲存庫名稱，例如 `simplebank`。
        - `IMAGE_TAG`: 映像檔標籤，使用 Git SHA 以便於版本管理。
    - **目的**: 將應用程式的最新版本封裝為 Docker 映像檔，並推送至 ECR，供後續部署使用。
7. **更新 kube config**
    - **步驟名稱**: Update kube config
    - **動作**: 使用 AWS CLI 更新本地的 kubeconfig 設定，連接到指定的 EKS 叢集。
    - **命令**:
        
        ```bash
        aws eks update-kubeconfig --name simple-bank-eks --region ap-northeast-1
        ```
        
    - **目的**: 設定本地的 kubeconfig，以便後續能夠使用 kubectl 與指定的 EKS 叢集互動。
8. **部署映像檔至 Amazon EKS**
    - **步驟名稱**: Deploy image to Amazon EKS
    - **動作**: 使用 `kubectl apply` 命令，將映像檔部署至 EKS 叢集，並應用多個 Kubernetes 資源配置檔案：
        - `eks/aws-auth.yaml`: 用於設定 IAM 角色與 Kubernetes 角色之間的映射關係。
        - `eks/deployment.yaml`: 用於定義 Kubernetes 部署資源，負責將應用程式容器部署到 Amazon EKS 叢集。
        - `eks/service.yaml`: 定義 Kubernetes 服務，用於暴露應用程式。
        - `eks/issuer.yaml`: 配置認證發行者。
        - `eks/ingress-nginx.yaml`: 配置 NGINX Ingress 控制器。
        - `eks/ingress-http.yaml`: 配置 HTTP Ingress 路由。
    - **目的**: 將最新的應用程式映像檔部署至 EKS 叢集，並應用所需的網路、服務及認證配置，確保應用程式正確運行並能夠對外提供服務。
