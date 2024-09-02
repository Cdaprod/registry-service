# REGS - Registry Service

The REGS service manages the registration and discovery of services and plugins within the Core Data Archetype (CDA) system. It provides dynamic service management, enabling easy registration, deregistration, and discovery of services.

---

# Registry Service

Welcome to the **Registry Service**--the beating heart of the **Core Data Archetype (CDA)**! This innovative system design is all about creating a flexible, scalable, and easy-to-maintain framework for managing complex data operations. With the Registry Service, you get a powerful registry and dynamic plugin system, allowing seamless integration with tools like Docker and GitHub, making your CDA ecosystem more powerful and cohesive.

## Overview

The **Core Data Archetype (CDA)** is your go-to framework for simplifying and organizing core operations and data management. At its core is the **Registry Service**, a central hub that keeps everything in sync. It ensures all your services within the CDA ecosystem can dynamically register, manage, and query data, keeping everything consistent and reliable. This is the foundation for building robust, scalable applications that grow with you.

## Features

- **Dynamic Registration and Discovery**: Easily extend your application by dynamically loading and registering plugins.
- **Centralized State Management**: One-stop-shop for managing all registered entities, ensuring smooth operations and consistency.
- **Extensible Plugin System**: Plug and play! Add new functionalities without touching the core code, thanks to our flexible plugin system.
- **Seamless Integration with Tools**: Effortlessly integrates with Docker, GitHub, and more, allowing you to manage repositories, containers, images, networks, and beyond with ease.

## Components

### 1. **Registry Service**

Provides a centralized registry for managing entities like Docker containers, images, GitHub repositories, etc. It includes:

- **Registerable Items**: Abstract representation of any operational entity that can be managed within CDA.
- **Persistent State Management**: Maintains the state of registered items across service restarts.
- **Metadata Storage**: Flexible metadata storage to accommodate different operational details such as URLs, ports, network settings, etc.

### 2. **Builtin and External Plugins**

Plugins are dynamically loaded components that extend the functionality of the core registry service:

- **Builtin Plugins**: Core components like Docker and GitHub integrations, essential for basic operations.
- **External Plugins**: Custom plugins that can be developed and integrated to add new capabilities or enhance existing ones.

### 3. **Operational Logic Modules**

Modules that perform specific tasks such as managing Docker containers, pulling images, and handling GitHub repositories. Each module ensures that any operation performed is consistent with the state maintained in the registry.

## Getting Started

### Prerequisites

- Go 1.16 or later
- Docker installed and running

### Installation

Clone the repository:

```bash
git clone https://github.com/Cdaprod/registry-service.git
cd registry-service
```

### Usage

1. **Initialize the Registry Service:**

   ```bash
   go run cmd/server/main.go
   ```

2. **Integrate with `repocate-service` or any other service:**

   - Use the `registry-service` to register and manage your operational entities.
   - Extend functionality by developing custom plugins and dynamically loading them into the registry.

### Example

```go
package main

import (
    "github.com/Cdaprod/registry-service/internal/registry"
    "github.com/cdaprod/repocate/internal/container"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    reg := registry.NewCentralRegistry()

    // Initialize and register a Docker container
    if err := container.InitializeAndRegisterContainer("repocate-default", "cdaprod/repocate-dev:1.0.0-arm64", reg); err != nil {
        logger.Fatal("Failed to initialize and register container", zap.Error(err))
    }

    // Additional logic...
}
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests for any improvements or new features.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

By utilizing the **Registry Service** within the **Core Data Archetype (CDA)** framework, you can achieve a high level of modularity, scalability, and maintainability in your data-centric applications.
```

### Key Sections in README.md

1. **Overview**: Introduces the Registry Service as a core part of CDA.
2. **Features**: Highlights the key capabilities and benefits.
3. **Components**: Describes the primary components (Registry Service, Plugins, Operational Modules).
4. **Getting Started**: Provides installation and basic usage instructions.
5. **Example**: Demonstrates how to use the registry service in practice.
6. **Contributing and License**: Standard sections for community involvement and licensing information.

This README serves as a comprehensive guide to understanding, using, and extending the `registry-service` within the context of your CDA framework.