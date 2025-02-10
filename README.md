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
│   └── test/          # Tests
├── internal/
│   └── executor/      # Core executor logic
│   └── k8s/           # Kubernetes client
│   └── awsclient/     # Aws client
├── docs/              # API documentation
├── pkg/
│   └── utils/         # Utility functions
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
or
```
docker build -t terraform-executor .
```

## API endpoints
For detailed API documentation, please refer to the [API Documentation](docs/API.md).

## Configuration

### Required AWS resources and IAM policies
- S3 bucket for storing Terraform configurations and state files
- Boundary policy for iam roles "s3-boundary" with content:
    ```json
    {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "LimitedListBucketAccess",
        "Effect": "Allow",
        "Action": "s3:ListBucket",
        "Resource": "arn:aws:s3:::<bucket-name>",
        "Condition": {
          "StringLike": {
            "s3:prefix": "${aws:PrincipalTag/UserId}/*"
          }
        }
      },
      {
        "Sid": "LimitedS3ObjectActions",
        "Effect": "Allow",
        "Action": [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ],
        "Resource": "arn:aws:s3:::<bucket-name>/${aws:PrincipalTag/UserId}/*"
      }
    ]
  }
    ```
- IAM user and Policy for executor service with these permissions:
    ```json
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "AllowAssumeTargetRole",
                "Effect": "Allow",
                "Action": "sts:AssumeRole",
                "Resource": "arn:aws:iam::066889832768:role/app/uptimeai/*"
            },
            {
                "Sid": "AllowRoleCreationWithBoundary",
                "Effect": "Allow",
                "Action": "iam:CreateRole",
                "Resource": "arn:aws:iam::066889832768:role/app/*",
                "Condition": {
                    "StringEquals": {
                        "iam:PermissionsBoundary": "arn:aws:iam::066889832768:policy/s3-boundary"
                    }
                }
            },
            {
                "Sid": "AllowAttachingRolePolicy",
                "Effect": "Allow",
                "Action": "iam:AttachRolePolicy",
                "Resource": "arn:aws:iam::066889832768:role/app/*",
                "Condition": {
                    "StringEquals": {
                        "iam:PolicyArn": "arn:aws:iam::066889832768:policy/s3-boundary"
                    }
                }
            },
            {
                "Sid": "AllowTaggingForBoundaryUsage",
                "Effect": "Allow",
                "Action": [
                    "iam:TagRole",
                    "iam:UntagRole",
                    "iam:ListRoleTags",
                    "iam:ListRoles"
                ],
                "Resource": "arn:aws:iam::066889832768:role/app/*"
            },
            {
                "Sid": "AllowGetRole",
                "Effect": "Allow",
                "Action": "iam:GetRole",
                "Resource": "arn:aws:iam::066889832768:*"
            },
            {
                "Sid": "AllowFullS3AccessToSpecificBucket",
                "Effect": "Allow",
                "Action": "s3:*",
                "Resource": [
                    "arn:aws:s3:::uptimeai-test-bucket",
                    "arn:aws:s3:::uptimeai-test-bucket/*"
                ]
            }
        ]
    }
    ```


The server configuration can be modified using the following environment variables:
```bash
# AWS credentials for IAM management and S3
AWS_ACCESS_KEY_ID=access_key_id
AWS_SECRET_ACCESS_KEY=secret_access_key
AWS_REGION=eu-west-3

# S3 bucket
BUCKET_NAME=uptimeai-test-bucket

# Kubernetes configuration (if running outside the cluster)
KUBECONFIG=/path/to/kubeconfig
```

## Test
```bash
go run cmd/test/main.go
```

## Run
1. Start the gRPC server
    ```
    ./terraform-executor
    ```
    or
    ```
    docker run -p 50051:50051 terraform-executor
    ```
2. Send gRPC requests with grpcurl or any gRPC client

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

