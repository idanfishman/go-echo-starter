# Go Echo Starter Project

The Go Echo Starter Project is a comprehensive template designed to kickstart your API development using the Echo framework. Crafted with simplicity and efficiency in mind, it serves as the foundation for building robust and scalable applications. With an emphasis on best practices and code quality, this project adheres to the Uber Go Style Guide, ensuring a high standard of code consistency and maintainability.

## Features

- **Structured Project Layout**: Organized into logical directories, facilitating easy navigation and scalability.
- **Pre-configured Logging**: Utilizes a JSON structured logger for clear and concise logs, aiding in better observability and debugging.
- **Middleware Support**: Includes essential middleware for logging, request ID generation, and request timeout handling, ensuring a smooth and secure request handling process.
- **Health Checks**: Implements liveness and readiness probes to monitor the health and performance of your application.
- **Configuration Management**: Streamlines the management of application configurations, making your project adaptable to different environments.

## Project Structure

- `cmd/appname/`: Hosts the main application code along with its configuration, serving as the entry point of the application.
- `pkg/`: Contains a collection of utility packages designed to provide support across the project, ensuring reusability and modularity.
  - `log/`: Offers a sophisticated logging mechanism that writes JSON structured logs to stdout/stderr.
  - `middleware/`: Comprises middleware functionalities for enhanced logging, request ID generation, and effective request timeout management.
  - `probe/`: Includes handlers for conducting health checks through liveness and readiness probes, vital for maintaining application health in kubernetes environments.
  - `validation/`: Provides utilities for rigorous input validation, ensuring the integrity of the data your application processes.

## Contributing

Join us in contributing to and enhancing this project, making it even more beneficial for the developer community.
