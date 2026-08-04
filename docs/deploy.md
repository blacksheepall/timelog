# Deployment

## Automatic Deployment

Push a Git tag to trigger GitHub Actions build & deploy.

```bash
git tag v1.0.0
git push origin v1.0.0
```

## GitHub Secrets

| Secret           | Description                     |
| ---------------- | ------------------------------- |
| `DEPLOY_HOST`    | Server host                     |
| `DEPLOY_USER`    | SSH username                    |
| `DEPLOY_SSH_KEY` | SSH private key                 |
| `DEPLOY_PATH`    | Upload directory for artifact   |
| `DEPLOY_SCRIPT`  | Path to deploy script on server |
| `DEPLOY_PORT`    | SSH port (default 22) |

## Server Preparation

1. Create deploy directory and config file.

```bash
mkdir -p /var/www/timelog
cp config-example.yml /var/www/timelog/config.yml
# Edit config file
```

The application runs plain HTTP. If you need HTTPS (for example, to enable passkeys), place a reverse proxy such as Nginx or Traefik in front of it to terminate TLS and forward HTTP to the application port.

2. Create systemd service at `/etc/systemd/system/timelog.service`.

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

3. Reload and enable the service.

```bash
systemctl daemon-reload
systemctl enable timelog.service
```