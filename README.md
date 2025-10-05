# Access Control Web Service

This project is a Go-based web service for managing and controlling access tickets with a simple modern frontend and Telegram notifications.

## Overview

The service provides two main features:

- **Control Page:** Allows entering a barrier/gate number and sending an open request. The request triggers a Telegram message notification with timestamp, gate number, and unique access ticket ID.

- **Admin Page:** Protected by basic authentication, this page lists active access tickets with their expiration times and provides controls to issue new tickets, prolong existing ones, or revoke them. Each ticket has a unique URL to access the control page without authentication.

## Features
- Access tickets expire after 6 hours by default, with prolong and revoke capabilities.
- Telegram bot integration for real-time notifications.
- Configuration via YAML file.
- Embedded HTML templates for easy deployment and performance.
- Special expired ticket page informing users if the ticket is no longer valid.

## Configuration

Example `config.yaml`:
```yaml
telegram_bot_token: "123456789:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
telegram_chat_id: 1234567890
admin_user: "admin"
admin_password: "password123"
server_addr: ":8080"
public_address: "http://mypage.com"
```

## Running the Service

1. Place your configuration in `config.yaml`.

2. Build and run the service:
```
go build -o access-service ./cmd/server
./access-service -config config.yaml
```

3. Access the admin panel at `http://localhost:8080/admin` and log in with the configured admin credentials.

4. Use the admin panel to create access tickets. Each ticket URL allows direct access to the control page.
