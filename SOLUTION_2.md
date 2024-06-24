# technical test solution

## Introduction

This Golang-based solution implements a RESTful API for a blog app. It provides the ability to create, update, delete, and retrieve blog posts, with each post consisting of a title, content, and an author. The solution implements CRUD (Create, Read, Update, Delete) operations for blog posts, and utilizes an in-memory data store for holding the blog posts data.

## Prerequisite: Build

Before running the application, you need to build it. Use the following command to build the application:

```bash
make build
```

This command compiles the application and prepares it for execution. Ensure you have a suitable Go environment set up, and all the necessary dependencies are installed.

## Hexagonal Architecture

The solution is based on the Hexagonal Architecture (also known as Ports and Adapters) pattern. This architecture separates the core logic of the application from the outside concerns like UI, database, and external interfaces. This separation allows independent evolution of the application and its infrastructure code.

In this architecture, the application is at the center (inside the hexagon) and all interactions with external systems or services happen through ports. The adapters at the ports convert the external interactions into a form that the application can handle.

This architecture brings advantages such as:

- **Maintainability**: By separating concerns, the system becomes easier to maintain. Changes in one part of the system (e.g., the database layer) do not affect other parts of the system (e.g., the business logic).
- **Readability**: The code is more understandable, as each part of the system has a specific role.
- **Testability**: It's easier to write unit tests for the application, as each part of the system can be tested independently.

However, it also has some disadvantages:

- **Verbosity**: The separation of concerns may result in more code to write. For example, you might need to write interfaces and implementations for each port.

## API-First Approach

The solution uses an API-first approach, meaning that the API is designed and documented before any coding begins. This provides a clear understanding of the system's requirements and functionality from the outset.

The server implementation is generated from an OpenAPI specification, which is a language-agnostic interface to the RESTful API. This allows the API to be understood and used by individuals who may not be familiar with Golang.

## DTO Mapping

Data Transfer Object (DTO) mapping is used to map between layers in the application. This is done using MapStruct, a code generator that simplifies the implementation of mappings between Java bean types. This leads to a cleaner code base and removes the need for manual mapping.

## Command Handling

Command handling is done using Cobra, a library for creating powerful modern CLI applications. It is used in combination with Viper, a library for Go that handles configuration.

## Configuration

Configuration in this solution is handled by Viper. It provides a unified way to manage configuration files, environment variables, command-line flags, and more.

## Authentication

The solution includes a basic authentication mechanism using an in-memory store. The authentication settings can be adjusted in the `app.yaml` configuration file.

## Usage

After building the application, you can start the server with the following command:

```bash
./app technical --initial-posts-path blog_data.json
```

This command starts the server and loads the initial blog posts from the `blog_data.json` file.

## Roadmap

The following features are planned for future development:

- **Logging**: Implementation of a robust logging system to aid in the debugging process and to provide insights about the application's behavior and performance.

- **Metrics Collection & Grafana Dashboard**: Incorporation of a metrics collection system to gather and analyze data about the application's performance. This data will be visualized using a Grafana dashboard, providing a user-friendly interface for monitoring the application.

- **Tracing**: Addition of a tracing mechanism to track the flow of requests through the application. This will help in identifying bottlenecks and improving the overall performance of the application.