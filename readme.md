# Go Load Balancer

A lightweight and efficient load balancer written in Go, supporting **Round Robin**, **Sticky Session**, and **IP Hashing** algorithms. This project includes a test Node.js backend server for demonstration purposes.

---

## 🚀 Getting Started

Follow the steps below to set up and run the project on your local machine.

### 1. Clone the Repository

```bash
git clone https://github.com/Samyak-SH/loadBalancer.git
cd loadBalancer
```

### 2. Download Go Dependencies

Use the `go.mod` file to install all required Go packages:

```bash
go mod download
```

### 3. Setup Environment Variables

Create a `.env` file in the root directory of the project and add the following:

```env
SECRET_KEY=your_secret_key_here
```

This key is used to sign cookies for sticky session management.

### 4. Start the Test Node Server (optional)

Navigate to the `nodeserver` directory and install dependencies:

```bash
cd nodeserver
npm install
npm start
```

This spins up a simple Node.js server to simulate backend nodes for the load balancer.

---

## ⚙️ Configuration

Edit the `config.json` file in the root of the project to configure the load balancer behavior:

```json
{
    "PORT" : x,
    "Servers" : [
        "http://localhost:3000",
        "http://localhost:3001",
        "http://localhost:3002",
        "http://localhost:3003"
    ],
    "Algorithm" : x,
    "VirtualNodeCount" : x,
    "HealthCheckInterval" : x
}
```

### Configuration Options:

- **PORT**: The port on which the load balancer will run.
- **Servers**: List of backend server URLs.
- **Algorithm**: The load balancing algorithm to use:
  - `1` → Round Robin
  - `2` → Sticky Session
  - `3` → IP Hashing
- **VirtualNodeCount**: *(Only required for IP Hashing)* Number of virtual nodes per real server in the hash ring. Helps improve load distribution.
- **HealthCheckInterval**: Time interval (in seconds) to perform health checks on backend servers.

---

## 🔄 Load Balancing Algorithms

### 🔁 Round Robin (Algorithm = 1)

Distributes incoming requests evenly by cycling through the list of backend servers.

### 📌 Sticky Sessions (Algorithm = 2)

Uses cookies to persistently route a client’s requests to the same backend server. Cookie values are securely signed using the `SECRET_KEY`.

### 🌐 IP Hashing (Algorithm = 3)

Uses the client’s IP address to determine the backend server using **consistent hashing**. Internally utilizes a **Hash Ring** with `VirtualNodeCount` virtual nodes per real node to:

- Maintain consistent routing for a given IP.
- Minimize disruption when servers are added or removed.

---

## 📎 License

MIT