# Terraform Executor
A gRPC service that provides remote execution capabilities for Terraform operations. This service enables secure and controlled execution of Terraform commands through a well-defined API interface.

## Features
- Remote execution of Terraform commands (plan, apply, destroy)
- Workspace isolation for multiple configurations
- gRPC API interface with protobuf message definitions
- Support for:
    - terraform plan - Generate execution plans
    - terraform apply - Apply configurations
    - terraform destroy - Clean up resources
- Error handling and output capture
- Configurable workspaces

## Project Structure
```
terraform-executor/
├── api/
│   ├── handlers/      # gRPC service handlers
│   └── proto/         # Protocol buffer definitions
├── cmd/
│   └── server/        # Server implementation
├── internal/
│   └── executor/      # Core executor logic
├── pkg/
│   └── utils/         # Utility functions
└── default_workspace/ # Default Terraform configuration
```

## Prerequisites
- Go 1.22 or later
- Protocol Buffers compiler
- Terraform CLI
- gRPC tools (optional, for testing)

## Installation
```
git clone <repository-url>
cd terraform-executor
go mod download
```

## Development
1. Define the protobuf message definitions in `api/proto/executor.proto`
2. Generate the gRPC service code using the `protoc` compiler
    ```
    protoc --go_out=. --go-grpc_out=. api/proto/executor.proto
    ```
3. Implement the service handlers in `api/handlers`
4. Implement the core executor logic in `internal/executor`
5. Build the server implementation in `cmd/server`

## Build

```
go build -o terraform-executor cmd/server/main.go
```

## API endpoints
The service exposes these gRPC endpoints:
- Plan: Generate an execution plan
    ```
    rpc Plan(PlanRequest) returns (PlanResponse) {}
    ```
- Apply: Apply the configuration
    ```
    rpc Apply(ApplyRequest) returns (ApplyResponse) {}
    ```
- Destroy: Destroy the resources
    ```
    rpc Destroy(DestroyRequest) returns (DestroyResponse) {}
    ```
- GetStateList: Get the current state of the resources
    ```
    rpc GetStateList(GetStateListRequest) returns (GetStateListResponse) {}
    ```
- Clear: clear terraform files
    ```
    rpc Clear(ClearRequest) returns (ClearResponse) {}
    ```

TO DO:
- GetOutput: Get the output of the last operation
    ```
    rpc GetOutput(GetOutputRequest) returns (GetOutputResponse) {}
    ```
- GetLogs: Get the logs of the last operation
    ```
    rpc GetLogs(GetLogsRequest) returns (GetLogsResponse) {}
    ```

## Usage
1. Start the gRPC server
    ```
    ./terraform-executor
    ```
2. Send gRPC requests with grpcurl or any gRPC client

## Example usage with grpcurl

Generate a plan:

    grpcurl -plaintext -d '{"workspace": "default_workspace", "code": "resource \"example\" {...}"}' \
        localhost:50051 executor.Executor/Plan

Apply the configuration:

    grpcurl -plaintext -d '{"workspace": "default_workspace", "code": "resource \"example\" {...}"}' \
        localhost:50051 executor.Executor/Apply

Destroy the resources:

    grpcurl -plaintext -d '{"workspace": "default_workspace", "code": "resource \"example\" {...}"}' \
        localhost:50051 executor.Executor/Destroy

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

