# PromMySQL

PromMySQL is a **simple Golang application** designed for continuous data push into MySQL while exporting key metrics to Prometheus. This project is particularly useful for scenarios where seamless data migration with minimal downtime is crucial.

## Features

- **Continuous Data Push:** The application maintains an infinite loop, pushing data into MySQL at regular intervals.

- **Prometheus Metrics Export:** Key metrics related to MySQL connections, throughput, success count, and failure count are exported to Prometheus.

## Getting Started

### Prerequisites

- Golang installed on your machine.
- Access to a MySQL database.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/tasszz2k/PromMySQL.git
   ```

2. Navigate to the project directory:

   ```bash
   cd PromMySQL
   ```

3. Install dependencies:

   ```bash
   go get -u github.com/go-sql-driver/mysql
   go get -u github.com/prometheus/client_golang/prometheus
   ```

4. Modify the `main.go` file with your MySQL connection details and adjust the data push logic if needed.

5. Run the application:

   ```bash
   go run main.go
   ```

### Docker

1. **Build the Docker image:**

   Open a terminal in the directory where your Dockerfile is located and run:

   ```bash
   docker build -t prommysql .
   ```

2. **Run the Docker container:**

   ```bash
   docker run -p 8080:8080 -e MYSQL_USERNAME=<your_username> -e MYSQL_PASSWORD=<your_password> -e MYSQL_HOST=<your_host> -e MYSQL_PORT=<your_port> -e MYSQL_DB_NAME=<your_db_name> -e SLEEP_INTERVAL_MILLISECOND=<your_sleep_interval> prommysql
   ```

   Replace `<your_username>`, `<your_password>`, `<your_host>`, `<your_port>`,
   `<your_db_name>`, and `<your_sleep_interval>` with your actual MySQL configuration.

   This command maps port 8080 from the container to port 8080 on your host machine
   and sets the necessary environment variables.

This Dockerfile assumes that your Golang application uses environment variables for MySQL configuration,
as indicated by your Golang code. Make sure to customize the environment variables accordingly.

### Kubernetes

1. **Applying the Manifests:**

Apply the deployment and service:

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

2. **Checking the Status:**

Check the status of your deployment:

```bash
kubectl get deployments
kubectl get pods
```

Check the status of your service:

```bash
kubectl get services
```

3. **Accessing Your Application:**

If you've used `LoadBalancer` type, wait until the external IP is assigned:

```bash
kubectl get services -w
```

Once you have the external IP, you can access your application using that IP.

If you've used `NodePort` or `ClusterIP`, you'll need to access the application through the respective Node IP or Cluster IP and the assigned port.

Remember to replace placeholder values with your actual configuration

### Usage

The application continuously pushes data into MySQL and exposes Prometheus metrics at `http://localhost:8080/metrics`.

### Metrics

- `mysql_connections_active`: Number of active MySQL connections.
- `mysql_throughput`: Number of data entries pushed to MySQL.
- `mysql_success_count`: Number of successful data pushes to MySQL.
- `mysql_fail_count`: Number of failed data pushes to MySQL.

## Contributing

If you'd like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

- Thanks to the [Prometheus](https://prometheus.io/) and [Golang](https://golang.org/) communities for their excellent tools and libraries.

Feel free to customize this template further based on your project's unique aspects. If you have any specific sections or details you'd like to include, let me know!