# Prahara Testing Tool 🚀

[![Docker Publish](https://github.com/Muhammad-Ikhwan-Fathulloh/Prahara-Testing-Tool/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/Muhammad-Ikhwan-Fathulloh/Prahara-Testing-Tool/actions/workflows/docker-publish.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Prahara** is a lightweight, high-performance, and secure QA and Load Testing tool. It combines the power of **k6** with a modern **Golang** backend and a **Vue 3** dashboard to provide a seamless testing experience.

🔗 **Repository**: [https://github.com/Muhammad-Ikhwan-Fathulloh/Prahara-Testing-Tool](https://github.com/Muhammad-Ikhwan-Fathulloh/Prahara-Testing-Tool)

---

## 🌟 Key Features

- **⚡ Quick Storm**: Instant dynamic load testing for any URL directly from the dashboard.
- **🌐 URL Registry & Categories**: Centralized management for testing endpoints with categories (Frontend, Backend, API, etc.).
- **🚀 Unified Architecture**: API and Web bundled into a single high-performance Docker image.
- **📈 Native k6 Integration**: Seamlessly execute k6 scripts and visualize metrics in real-time.
- **🕒 Sustainability & History**: All metrics are stored in **InfluxDB**, allowing you to track performance trends over time.
- **🔐 Secure by Design**: RBAC (Role-Based Access Control) with JWT authentication and SQLite persistence.

---

## 🛠️ Technology Stack

- **Backend**: Golang (Gin, GORM with Pure-Go SQLite)
- **Frontend**: Vue 3, Vite, Tailwind CSS, Chart.js
- **Testing Engine**: k6
- **Database**: 
  - **InfluxDB 2.x**: Time-series metrics
  - **SQLite**: User metadata and RBAC
- **Infrastructure**: Docker & GitHub Actions

---

## 🚀 Quick Start

### 1. Prerequisites
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### 2. Installation
Clone the repository:

```bash
git clone https://github.com/Muhammad-Ikhwan-Fathulloh/Prahara-Testing-Tool.git
cd Prahara-Testing-Tool
```

#### Choose Your Deployment Mode:

**A. Unified Mode (Recommended)**
Runs the entire stack (Web + API) in a single container. Best for simple deployments.
```bash
docker-compose up -d
```

**B. Split Mode (Modular)**
Runs Web and API in separate containers. Best for scaling or custom proxying.
Use the split configuration:
```bash
docker-compose -f docker-compose.split.yml up -d
```

### 3. Access the Dashboard
- **Unified Mode**: `http://localhost` (Port 80)
- **Split Mode**: 
  - **Frontend**: `http://localhost` (Port 80)
  - **API**: `http://localhost:3000`
- **InfluxDB Console**: `http://localhost:8086`

**Default Credentials**:
- **Username**: `admin`
- **Password**: `password` (You can register a new user in the app)

---

## 🤖 CI/CD Integration

This project is pre-configured with GitHub Actions to build and push the unified image to Docker Hub.

### Setup Secrets
In your GitHub Repo, go to `Settings > Secrets and variables > Actions` and add:
1. `DOCKERHUB_USERNAME`: Your Docker Hub username.
2. `DOCKERHUB_TOKEN`: Your Docker Hub Personal Access Token.

The workflow will automatically build and push to `${DOCKERHUB_USERNAME}/prahara:latest` on every push to the `main` branch.

---

## 📖 Usage Guide

### Dynamic Load Testing (Quick Storm)
1. Go to the **Dashboard**.
2. Locate the **Quick Storm** widget at the top.
3. Enter your target URL, select the method (GET/POST), and set the VUs/Duration.
4. Hit **Launch** and watch the real-time graphs!

### Endpoint Management
Use the **URL Registry** to keep track of your team's different endpoints. Categorize them for easier filtering during test planning.

### Scripted Testing
Write complex k6 scripts in the **Test Editor** section for advanced scenarios like authentication flows or multi-stage load surges.

---

## 🛡️ License
Distributed under the **MIT License**. See `LICENSE` for more information.

## 🤝 Contribution
Contributions are welcome! Feel free to open issues or submit pull requests.

---
Created with ❤️ by the Prahara Team.
