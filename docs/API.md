# API Reference

## Table of contents
- [Executor Service](#executor-service)
    - [CreateContext](#createcontext)
    - [DeleteContext](#deletecontext)
    - [CreateWorkspace](#createworkspace)
    - [DeleteWorkspace](#deleteworkspace)
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
    - [ClearWorkspace](#clearworkspace)
    - [GetMainTf](#getmaintf)

## Executor Service

The `Executor` service provides methods to manage Terraform operations such as planning, applying, destroying infrastructure, and managing contexts and workspaces.

### CreateContext

Creates a new context.

**Request:** `CreateContextRequest`
- `string context`: Name of the context

**Response:** `CreateContextResponse`
- `bool success`: Whether the context creation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Create a new context
grpcurl -plaintext -d '{
    "context": "dev"
}' localhost:50051 executor.Executor/CreateContext
```

### DeleteContext

Deletes a context.

**Request:** `DeleteContextRequest`
- `string context`: Name of the context

**Response:** `DeleteContextResponse`
- `bool success`: Whether the context deletion was successful
- `string error`: Error message, if any

**Example:**
```bash
# Delete a context
grpcurl -plaintext -d '{
    "context": "dev"
}' localhost:50051 executor.Executor/DeleteContext
```

### CreateWorkspace

Creates a workspace within a context.

**Request:** `CreateWorkspaceRequest`
- `string context`: Name of the context
- `string workspace`: Name of the workspace

**Response:** `CreateWorkspaceResponse`
- `bool success`: Whether the workspace creation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Create a new workspace
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a"
}' localhost:50051 executor.Executor/CreateWorkspace
```

### DeleteWorkspace

Deletes a workspace within a context.

**Request:** `DeleteWorkspaceRequest`
- `string context`: Name of the context
- `string workspace`: Name of the workspace

**Response:** `DeleteWorkspaceResponse`
- `bool success`: Whether the workspace deletion was successful
- `string error`: Error message, if any

**Example:**
```bash
# Delete a workspace
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a"
}' localhost:50051 executor.Executor/DeleteWorkspace
```

### AddProviders

Adds providers to the Terraform configuration.

**Request:** `AddProvidersRequest`
- `string context`: Name of the context
- `string workspace`: Name of the workspace
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
    "context": "dev",
    "workspace": "project-a",
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
- `string context`: Name of the context
- `string workspace`: Name of the Terraform workspace
- `string code`: The Terraform configuration code snippet

**Response:** `AppendCodeResponse`
- `bool success`: Whether the code append operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Append code to the Terraform configuration
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a",
    "code": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-terraform-bucket\"\n  tags = {\n    Environment = \"Dev\"\n  }\n}"
}' localhost:50051 executor.Executor/AppendCode
```

### Plan

Generates a Terraform plan and returns the result.

**Request:** `PlanRequest`
- `string context`: Name of the context
- `string workspace`: Name of the Terraform workspace

**Response:** `PlanResponse`
- `bool success`: Whether the plan generation was successful
- `string plan_output`: The output of `terraform plan`
- `string error`: Error message, if any

**Example:**
```bash
# Run a Terraform plan
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a"
}' localhost:50051 executor.Executor/Plan
```

### Apply

Applies the Terraform plan and returns the execution result.

**Request:** `ApplyRequest`
- `string context`: Name of the context
- `string workspace`: Name of the Terraform workspace
- `string plan_file`: Path to the saved Terraform plan file (optional)

**Response:** `ApplyResponse`
- `bool success`: Whether the apply operation was successful
- `string apply_output`: The output of `terraform apply`
- `string error`: Error message, if any

**Example:**
```bash
# Apply the Terraform configuration
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a",
}' localhost:50051 executor.Executor/Apply
```

### Destroy

Destroys the Terraform-managed infrastructure and returns the result.

**Request:** `DestroyRequest`
- `string context`: Name of the context
- `string workspace`: Name of the Terraform workspace

**Response:** `DestroyResponse`
- `bool success`: Whether the destroy operation was successful
- `string destroy_output`: The output of `terraform destroy`
- `string error`: Error message, if any

**Example:**
```bash
# Destroy the Terraform-managed infrastructure
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a",
}' localhost:50051 executor.Executor/Destroy
```

### GetStateList

Retrieves the Terraform state list.

**Request:** `GetStateListRequest`
- `string context`: Name of the context
- `string workspace`: Name of the Terraform workspace

**Response:** `GetStateListResponse`
- `bool success`: Whether the state list retrieval was successful
- `string state_list_output`: The Terraform state file content
- `string error`: Error message, if any

**Example:**
```bash
# Get the Terraform state list
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a",
}' localhost:50051 executor.Executor/GetStateList
```

### ClearCode

Clears the main.tf file.

**Request:** `ClearCodeRequest`
- `string context`: Name of the context
- `string workspace`: Name of the Terraform workspace

**Response:** `ClearCodeResponse`
- `bool success`: Whether the clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear the Terraform configuration
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a",
}' localhost:50051 executor.Executor/ClearCode
```

### ClearProviders

Clears providers from the Terraform configuration.

**Request:** `ClearProvidersRequest`
- `string context`: Name of the context
- `string workspace`: Name of the workspace

**Response:** `ClearProvidersResponse`
- `bool success`: Whether the provider clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear providers from the Terraform configuration
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a"
}' localhost:50051 executor.Executor/ClearProviders
```

### AddSecretEnv

Adds secret environment variables to the Terraform configuration.

**Request:** `AddSecretEnvRequest`
- `string context`: Name of the context
- `string workspace`: Name of the workspace
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
    "context": "dev",
    "workspace": "project-a",
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
- `string context`: Name of the context
- `string workspace`: Name of the workspace
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
    "context": "dev",
    "workspace": "project-a",
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
- `string context`: Name of the context
- `string workspace`: Name of the workspace

**Response:** `ClearSecretVarsResponse`
- `bool success`: Whether the secret variable clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear secret variables from the Terraform configuration
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a"
}' localhost:50051 executor.Executor/ClearSecretVars
```

### ClearWorkspace

Clears the Terraform workspace.

**Request:** `ClearWorkspaceRequest`
- `string context`: Name of the context
- `string workspace`: Name of the workspace

**Response:** `ClearWorkspaceResponse`
- `bool success`: Whether the workspace clear operation was successful
- `string error`: Error message, if any

**Example:**
```bash
# Clear the Terraform workspace
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a"
}' localhost:50051 executor.Executor/ClearWorkspace
```

### GetMainTf

Gets the content of the main.tf file.

**Request:** `GetMainTfRequest`
- `string context`: Name of the context
- `string workspace`: Name of the workspace

**Response:** `GetMainTfResponse`
- `bool success`: Whether the operation was successful
- `string content`: Content of the main.tf file
- `string error`: Error message, if any

**Example:**
```bash
# Get the content of the main.tf file
grpcurl -plaintext -d '{
    "context": "dev",
    "workspace": "project-a"
}' localhost:50051 executor.Executor/GetMainTf
```
