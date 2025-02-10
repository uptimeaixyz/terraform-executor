# API Reference

## Table of contents
- [Executor Service](#executor-service)
    - [CreateProject](#createproject)
    - [DeleteProject](#deleteproject)
    - [AddProviders](#addproviders)
    - [AppendCode](#appendcode)
    - [Plan](#plan)
    - [Apply](#apply)
    - [Destroy](#destroy)
    - [GetStateList](#getstatelist)
    - [ClearCode](#clearcode)
    - [ClearProviders](#clearproviders)
    - [AddSecretEnv](#addsecretenv)
    - [AddSecretVar](#addsecretvar)
    - [ClearSecretVars](#clearsecretvars)
    - [ClearSecretEnv](#clearsecretenv)
    - [GetMainTf](#getmaintf)

## Executor Service

The `Executor` service provides methods to manage Terraform operations such as planning, applying, destroying infrastructure, and managing projects.

### CreateProject

Creates a new project.

**Request:** `CreateProjectRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `CreateProjectResponse`
- `bool success`: Whether the project creation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Create a new project
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/CreateProject
```

### DeleteProject

Deletes a project.

**Request:** `DeleteProjectRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `DeleteProjectResponse`
- `bool success`: Whether the project deletion was successful
- `string error`: Error message, if any

**Example:**
```bash
# Delete a project
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/DeleteProject
```

### AddProviders

Adds providers to the Terraform configuration.

**Request:** `AddProvidersRequest`
- `string user_id`: User identifier
- `string project`: Name of the project
- `repeated Provider providers`: List of providers to add
    - `string name`: Name of the provider
    - `string source`: Source of the provider
    - `string version`: Version of the provider

**Response:** `AddProvidersResponse`
- `bool success`: Whether the provider addition was successful
- `string error`: Error message, if any

**Example:**
```bash
# Add providers to the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a",
    "providers": [
        {
            "name": "aws",
            "source": "hashicorp/aws",
            "version": "3.0.0"
        },
        {
            "name": "digitalocean",
            "source": "digitalocean/digitalocean",
            "version": "~> 2.0"
        }
    ]
}' localhost:50051 executor.Executor/AddProviders
```

### AppendCode
Appends code to the configuration file.

**Request:** `AppendCodeRequest`
- `string user_id`: User identifier
- `string project`: Name of the project
- `string code`: The Terraform configuration code snippet

**Response:** `AppendCodeResponse`
- `bool success`: Whether the code append operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Append code to the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a",
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-terraform-bucket\"\n  tags = {\n    Environment = \"Dev\"\n  }\n}"
}' localhost:50051 executor.Executor/AppendCode
```

### Plan

Generates a Terraform plan and returns the result.

**Request:** `PlanRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `PlanResponse`
- `bool success`: Whether the plan generation was successful
- `string plan_output`: The output of `terraform plan`
- `string error`: Error message, if any

**Example:**
```bash
# Run a Terraform plan
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/Plan
```

### Apply

Applies the Terraform plan and returns the execution result.

**Request:** `ApplyRequest`
- `string user_id`: User identifier
- `string project`: Name of the project
- `string plan_file`: Path to the saved Terraform plan file (optional)

**Response:** `ApplyResponse`
- `bool success`: Whether the apply operation was successful
- `string apply_output`: The output of `terraform apply`
- `string error`: Error message, if any

**Example:**
```bash
# Apply the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/Apply
```

### Destroy

Destroys the Terraform-managed infrastructure and returns the result.

**Request:** `DestroyRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `DestroyResponse`
- `bool success`: Whether the destroy operation was successful
- `string destroy_output`: The output of `terraform destroy`
- `string error`: Error message, if any

**Example:**
```bash
# Destroy the Terraform-managed infrastructure
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/Destroy
```

### GetStateList

Retrieves the Terraform state list.

**Request:** `GetStateListRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `GetStateListResponse`
- `bool success`: Whether the state list retrieval was successful
- `string state_list_output`: The Terraform state file content
- `string error`: Error message, if any

**Example:**
```bash
# Get the Terraform state list
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/GetStateList
```

### ClearCode

Clears the main.tf file.

**Request:** `ClearCodeRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `ClearCodeResponse`
- `bool success`: Whether the clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/ClearCode
```

### ClearProviders

Clears providers from the Terraform configuration.

**Request:** `ClearProvidersRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `ClearProvidersResponse`
- `bool success`: Whether the provider clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear providers from the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/ClearProviders
```

### AddSecretEnv

Adds secret environment variables to the Terraform configuration.

**Request:** `AddSecretEnvRequest`
- `string user_id`: User identifier
- `string project`: Name of the project
- `repeated SecretEnv secrets`: List of secret environment variables to add
    - `string key`: Key of the secret environment variable
    - `string value`: Value of the secret environment variable

**Response:** `AddSecretEnvResponse`
- `bool success`: Whether the secret environment variable addition was successful
- `string error`: Error message, if any

**Example:**
```bash
# Add secret environment variables to the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a",
    "secrets": [
        {
            "key": "AWS_ACCESS_KEY_ID",
            "value": "your-access-key-id"
        },
        {
            "key": "AWS_SECRET_ACCESS_KEY",
            "value": "your-secret-access-key"
        }
    ]
}' localhost:50051 executor.Executor/AddSecretEnv
```

### AddSecretVar

Adds secret variables to the Terraform configuration.

**Request:** `AddSecretVarRequest`
- `string user_id`: User identifier
- `string project`: Name of the project
- `repeated SecretVar secrets`: List of secret variables to add
    - `string key`: Key of the secret variable
    - `string value`: Value of the secret variable

**Response:** `AddSecretVarResponse`
- `bool success`: Whether the secret variable addition was successful
- `string error`: Error message, if any

**Example:**
```bash
# Add secret variables to the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a",
    "secrets": [
        {
            "key": "db_password",
            "value": "your-db-password"
        },
        {
            "key": "api_key",
            "value": "your-api-key"
        }
    ]
}' localhost:50051 executor.Executor/AddSecretVar
```

### ClearSecretVars

Clears secret variables from the Terraform configuration.

**Request:** `ClearSecretVarsRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `ClearSecretVarsResponse`
- `bool success`: Whether the secret variable clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear secret variables from the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/ClearSecretVars
```

### ClearSecretEnv

Clears secret environment variables from the Terraform configuration.

**Request:** `ClearSecretEnvRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `ClearSecretEnvResponse`
- `bool success`: Whether the secret environment variable clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear secret environment variables from the Terraform configuration
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/ClearSecretEnv
```

### ClearWorkspace

Clears the Terraform workspace.

**Request:** `ClearWorkspaceRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `ClearWorkspaceResponse`
- `bool success`: Whether the workspace clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear the Terraform workspace
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/ClearWorkspace
```

### GetMainTf

Gets the content of the main.tf file.

**Request:** `GetMainTfRequest`
- `string user_id`: User identifier
- `string project`: Name of the project

**Response:** `GetMainTfResponse`
- `bool success`: Whether the operation was successful
- `string content`: Content of the main.tf file
- `string error`: Error message, if any

**Example:**
```bash
# Get the content of the main.tf file
grpcurl -plaintext -d '{
    "user_id": "user123",
    "project": "project-a"
}' localhost:50051 executor.Executor/GetMainTf
```
