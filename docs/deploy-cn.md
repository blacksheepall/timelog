# 部署说明

## 自动部署

推送 Git Tag 触发 GitHub Actions 自动构建和部署：

```bash
git tag v1.0.0
git push origin v1.0.0
```

## GitHub Secrets 配置

| Secret           | 说明                      |
| ---------------- | ------------------------- |
| `DEPLOY_HOST`    | 服务器地址                |
| `DEPLOY_USER`    | SSH 用户名                |
| `DEPLOY_SSH_KEY` | SSH 私钥                  |
| `DEPLOY_PATH`    | artifact 上传目录         |
| `DEPLOY_SCRIPT`  | 服务器上的部署脚本路径    |
| `DEPLOY_PORT`    | SSH 端口（默认 22） |

## 服务器准备

1. 创建部署目录和配置文件：

```bash
mkdir -p /var/www/timelog
cp config-example.yml /var/www/timelog/config.yml
# 编辑配置文件
```

应用本身以纯 HTTP 运行。如需 HTTPS（例如启用 passkey），请在其前方部署 Nginx、Traefik 等反向代理来终止 TLS，并以 HTTP 转发到应用端口。

2. 创建 systemd 服务 `/etc/systemd/system/timelog.service`：

```ini
[Unit]
Description=Timelog Service
After=network.target

[Service]
Type=simple
WorkingDirectory=/var/www/timelog
ExecStart=/var/www/timelog/main
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

3. 重载并启用服务：

```bash
systemctl daemon-reload
systemctl enable timelog.service
```