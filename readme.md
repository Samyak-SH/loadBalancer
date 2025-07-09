
# Go Load Balancer

A lightweight and efficient load balancer written in Go, supporting **Round Robin**, **Sticky Session**, and **IP Hashing** algorithms. This project includes a test Node.js backend server for demonstration purposes.

---

## 🚀 Getting Started

You can either run the project **manually** or using **Docker Compose**.

---

### 🛠️ Option 1: Manual Setup

#### 1. Clone the Repository

```bash
git clone https://github.com/Samyak-SH/loadBalancer.git
cd loadBalancer
```

#### 2. Download Go Dependencies

Use the `go.mod` file to install all required Go packages:

```bash
go mod download
```

#### 3. Setup Environment Variables (Required for Sticky Session Algorithm)

Create a `.env` file in the root directory of the project and add the following:

```env
SECRET_KEY=your_secret_key_here
```

This key is used to sign cookies for sticky session management.

#### 4. Start the Test Node Server (optional)

Navigate to the `nodeserver` directory and install dependencies:

```bash
cd nodeserver
npm install
npm start
```

This spins up a simple Node.js server to simulate backend nodes for the load balancer.

---

### ⚙️ Option 2: Docker Compose (Recommended)

#### 1. Clone the Repository

```bash
git clone https://github.com/Samyak-SH/loadBalancer.git
cd loadBalancer
```

#### 2. Update Configuration via Environment Variables

You **do not need to modify `config.json` manually**. Instead, configure the load balancer using environment variables by editing the `docker-compose.yml` file or a `.env` file.

Available environment variables:

```env
PORT=8080
SERVERS=http://localhost:3000,http://localhost:3001,http://localhost:3002,http://localhost:3003
ALGORITHM=1
VIRTUAL_NODE_COUNT=5
HEALTH_CHECK_INTERVAL=5
SECRET_KEY=your_secret_key_here
```

Each variable maps to the same field in `config.json`.

#### 3. Start the Load Balancer

```bash
docker compose up --build
```

This builds and starts the container, automatically generating `config.json` from the environment variables and launching the load balancer.

---

## ⚙️ Configuration (Manual Mode Only)

If not using Docker Compose, you can edit the `config.json` manually:

```json
{
  "PORT": 8080,
  "Servers": [
    "http://localhost:3000",
    "http://localhost:3001",
    "http://localhost:3002",
    "http://localhost:3003"
  ],
  "Algorithm": 1,
  "VirtualNodeCount": 5,
  "HealthCheckInterval": 5,
  "SECRET_KEY": "your_secret_key_here"
}
```

### Configuration Options:

- **PORT**: The port on which the load balancer will run.
- **Servers**: List of backend server URLs.
- **Algorithm**:
  - `1` → Round Robin
  - `2` → Sticky Session
  - `3` → IP Hashing
- **VirtualNodeCount**: *(Only for IP Hashing)* Number of virtual nodes per real server.
- **HealthCheckInterval**: Interval (in seconds) for backend health checks.
- **SECRET_KEY**: Used to sign cookies for sticky sessions.

---

## 🔄 Load Balancing Algorithms

### 🔁 Round Robin (Algorithm = 1)

Distributes incoming requests evenly by cycling through the backend servers.

### 📌 Sticky Sessions (Algorithm = 2)

Uses cookies signed with the `SECRET_KEY` to consistently route requests from the same client to the same backend.

### 🌐 IP Hashing (Algorithm = 3)

Routes based on client IP using **consistent hashing**. Ensures:
- Stable request routing
- Minimal disruption during server changes

---

## 📎 License

MIT
