## Data Agent

### Description
Data Agent is a lightweight monitoring agent that collects system metrics (CPU, RAM, disks, network) from a host and publishes them to a RabbitMQ queue for further processing and storage in PostgreSQL. It also includes a gRPC service that provides read-only access to the collected metrics.

---

### Features
- Collects CPU, RAM, disk, and network metrics
- Publishes metrics to RabbitMQ reliably
- Acknowledges messages (ack/nack) to ensure data integrity
- Stores metrics in PostgreSQL
- Supports graceful shutdown via system signals
- Thread-safe operation for multiple agents
- gRPC service for read-only access to metrics stored in PostgreSQL

---

### gRPC Service
The Data Agent gRPC service allows clients to query the collected metrics from PostgreSQL.  
Clients can retrieve historical or current metrics for monitoring and analysis purposes.

### Capabilities
- Access metrics by host, time range, or metric type
- Lightweight and fast, suitable for multiple concurrent clients
- Works seamlessly with the Data Agent collector and publisher

---

## Install
1. **Clone the repository**  
   Clone the Data Agent source code from GitHub to your local machine:  
   ```bash
   git clone https://github.com/dbrgmnn/data-agent.git && cd data-agent
   ```  
   This downloads the project files and navigates into the project directory.

2. **Install dependencies**  
   Use Go modules to download and install required dependencies:  
   ```bash
   go mod tidy
   ```  
   This ensures all necessary packages are available for building and running the agent.

3. **Configure environment variables**  
   Copy the example environment file to `.env`:  
   ```bash
   cp .env.example .env
   ```  
   Edit `.env` to set your RabbitMQ, PostgreSQL connection details, and other configuration parameters.

4. **Build the agent binary**  
   Compile the Go source code for Linux AMD64 architecture:  
   ```bash
   GOOS=linux GOARCH=amd64 go build -o bin/agent cmd/agent/main.go
   ```  
   This creates an executable binary in the `bin` directory.

5. **Start services with Docker Compose**  
   Launch RabbitMQ, PostgreSQL, and other dependencies using Docker Compose:  
   ```bash
   docker-compose up --build
   ```  
   This sets up the required infrastructure for the agent to operate.

6. **Deploy the agent binary**  
   Copy the agent binary to the target server, grant it execution permissions, and start it on the server:  
   ```bash
   scp bin/agent user@target-server:/path/to/deploy/
   ssh user@target-server
   chmod +x /path/to/deploy/agent
   /path/to/deploy/agent --url 'amqp://login:password@hostname:5672/' --interval 2
   ```  
   Replace `user`, `target-server`, and `/path/to/deploy/` with appropriate values for your environment.

7. **Run the agent directly**  
   Start the Data Agent, specifying the RabbitMQ URL and collection interval:  
   ```bash
   ./bin/agent --url 'amqp://login:password@hostname:5672/' --interval 2
   ```  
   The agent will begin collecting and publishing metrics every 2 seconds.

8. **Use gRPC calls**  
   Connect to the Data Agent's gRPC service to query collected metrics.  
   Clients can request metrics by host, time range, or metric type for monitoring and analysis.

---

## License
This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.