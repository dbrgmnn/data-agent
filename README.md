# Data Agent

## Description
Data Agent is a lightweight monitoring agent that collects system metrics (CPU, RAM, disks, network) from a host and publishes them to a RabbitMQ queue for further processing and storage in PostgreSQL. It also includes a gRPC service that provides read-only access to the collected metrics.

---

## Features
- Collects CPU, RAM, disk, and network metrics
- Publishes metrics to RabbitMQ reliably
- Acknowledges messages (ack/nack) to ensure data integrity
- Stores metrics in PostgreSQL
- Supports graceful shutdown via system signals
- Thread-safe operation for multiple agents
- gRPC service for read-only access to metrics stored in PostgreSQL

---

## gRPC Service
The Data Agent gRPC service allows clients to query the collected metrics from PostgreSQL.  
Clients can retrieve historical or current metrics for monitoring and analysis purposes.

### Capabilities
- Access metrics by host, time range, or metric type
- Lightweight and fast, suitable for multiple concurrent clients
- Works seamlessly with the Data Agent collector and publisher

---

## License
This project is licensed under the MIT License. See [LICENSE](./LICENSE) for details.