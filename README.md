# Apigee Client Go

A Go client library for the [Google Apigee Management API](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest).

## Installation

```bash
go get github.com/athakur/apigee-client-go
```

## Authentication

By default, the client uses [Google Application Default Credentials (ADC)](https://cloud.google.com/docs/authentication/application-default-credentials). Ensure you have one of the following configured:

- `GOOGLE_APPLICATION_CREDENTIALS` environment variable pointing to a service account key file
- Running on Google Cloud with a service account attached
- Authenticated via `gcloud auth application-default login`

The client requires the `https://www.googleapis.com/auth/cloud-platform` OAuth scope.

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    apigee "github.com/athakur/apigee-client-go"
)

func main() {
    ctx := context.Background()

    // Create a client for your Apigee organization
    client, err := apigee.NewClient(ctx, "my-org")
    if err != nil {
        log.Fatal(err)
    }

    // List all API products
    products, err := client.APIProducts.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, p := range products.APIProducts {
        fmt.Println(p.Name)
    }
}
```

## Client Options

The client can be customized using functional options:

```go
// Use a custom HTTP client
client, err := apigee.NewClient(ctx, "my-org",
    apigee.WithHTTPClient(myHTTPClient),
)

// Use a custom base URL (e.g., for testing)
client, err := apigee.NewClient(ctx, "my-org",
    apigee.WithBaseURL("https://custom-apigee-endpoint.example.com/v1"),
)

// Use a custom OAuth2 token source
client, err := apigee.NewClient(ctx, "my-org",
    apigee.WithTokenSource(myTokenSource),
)
```

## API Reference

### API Products

API Products define the API resources and quotas available to client applications.

```go
// Create an API product
product, err := client.APIProducts.Create(ctx, &apigee.APIProduct{
    Name:         "my-product",
    DisplayName:  "My Product",
    Description:  "A sample API product",
    ApprovalType: "auto",  // or "manual"
    Environments: []string{"test", "prod"},
    Proxies:      []string{"my-proxy"},
    Quota:        "1000",
    QuotaInterval: "1",
    QuotaTimeUnit: "hour",
    Attributes: []apigee.Attribute{
        {Name: "access", Value: "public"},
    },
})

// Get an API product
product, err := client.APIProducts.Get(ctx, "my-product")

// Update an API product
product.Description = "Updated description"
product, err = client.APIProducts.Update(ctx, "my-product", product)

// List API products with pagination
products, err := client.APIProducts.List(ctx, &apigee.ListOptions{
    PageSize:  100,
    PageToken: "",  // Use NextPageToken from previous response
})

// Delete an API product
err = client.APIProducts.Delete(ctx, "my-product")
```

### App Groups

App Groups allow you to organize and manage collections of apps.

```go
// Create an app group
appGroup, err := client.AppGroups.Create(ctx, &apigee.AppGroup{
    Name:        "my-group",
    DisplayName: "My App Group",
    ChannelID:   "web",
    ChannelURI:  "https://example.com",
    Attributes: []apigee.Attribute{
        {Name: "team", Value: "engineering"},
    },
})

// Get an app group
appGroup, err := client.AppGroups.Get(ctx, "my-group")

// Update an app group
appGroup.DisplayName = "Updated Name"
appGroup, err = client.AppGroups.Update(ctx, "my-group", appGroup)

// List app groups
appGroups, err := client.AppGroups.List(ctx, &apigee.ListOptions{
    PageSize: 50,
})

// Delete an app group
err = client.AppGroups.Delete(ctx, "my-group")
```

### App Group Apps

Apps within an App Group represent client applications that consume your APIs.

```go
// Create an app
app, err := client.AppGroupApps.Create(ctx, "my-group", &apigee.AppGroupApp{
    Name:       "my-app",
    CallbackURL: "https://example.com/callback",
    Attributes: []apigee.Attribute{
        {Name: "environment", Value: "production"},
    },
})

// Get an app
app, err := client.AppGroupApps.Get(ctx, "my-group", "my-app")

// Update an app
app.CallbackURL = "https://example.com/new-callback"
app, err = client.AppGroupApps.Update(ctx, "my-group", "my-app", app)

// List apps in an app group
apps, err := client.AppGroupApps.List(ctx, "my-group", nil)

// Delete an app
err = client.AppGroupApps.Delete(ctx, "my-group", "my-app")
```

### App Group App Keys (Credentials)

Keys provide authentication credentials for apps to access APIs.

```go
// Generate a new key (Apigee generates the key/secret)
// This preserves existing app attributes and settings
app, err := client.AppGroupAppKeys.Generate(ctx, "my-group", "my-app",
    &apigee.AppGroupAppKeyGenerateRequest{
        APIProducts:  []string{"my-product"},
        KeyExpiresIn: -1,  // -1 = never expires, or milliseconds
    },
)
// The new key is in app.Credentials

// Create a key with custom key/secret values
key, err := client.AppGroupAppKeys.Create(ctx, "my-group", "my-app",
    &apigee.AppGroupAppKeyCreateRequest{
        ConsumerKey:    "my-custom-key",
        ConsumerSecret: "my-custom-secret",
    },
)

// Get a key
key, err := client.AppGroupAppKeys.Get(ctx, "my-group", "my-app", "consumer-key")

// Update a key (e.g., to change status or API products)
key.Status = "revoked"
key, err = client.AppGroupAppKeys.Update(ctx, "my-group", "my-app", "consumer-key", key)

// Update API product approval status for a key
// action can be "approve" or "revoke"
err = client.AppGroupAppKeys.UpdateAPIProductStatus(ctx, "my-group", "my-app",
    "consumer-key", "my-product", "approve")

// Delete a key
err = client.AppGroupAppKeys.Delete(ctx, "my-group", "my-app", "consumer-key")
```

### Key Value Maps (KVMs)

Key Value Maps store configuration data at organization or environment level.

#### Organization-Level KVMs

```go
// Create an organization-level KVM
kvm, err := client.KeyValueMaps.Create(ctx, &apigee.KeyValueMap{
    Name:      "my-config",
    Encrypted: true,
})

// Get a KVM
kvm, err := client.KeyValueMaps.Get(ctx, "my-config")

// List all KVMs (returns names only)
kvms, err := client.KeyValueMaps.List(ctx)
for _, name := range kvms.KeyValueMapNames {
    fmt.Println(name)
}

// Delete a KVM
err = client.KeyValueMaps.Delete(ctx, "my-config")
```

#### Organization-Level KVM Entries

```go
// Add an entry to a KVM
entry, err := client.KeyValueMapEntries.Create(ctx, "my-config", &apigee.KeyValueEntry{
    Name:  "api-key",
    Value: "secret123",
})

// Get an entry
entry, err := client.KeyValueMapEntries.Get(ctx, "my-config", "api-key")

// Update an entry
entry, err = client.KeyValueMapEntries.Update(ctx, "my-config", "api-key", &apigee.KeyValueEntry{
    Name:  "api-key",
    Value: "new-secret456",
})

// List all entries in a KVM
entries, err := client.KeyValueMapEntries.List(ctx, "my-config")

// Delete an entry
err = client.KeyValueMapEntries.Delete(ctx, "my-config", "api-key")
```

#### Environment-Level KVMs

```go
// Create an environment-level KVM
kvm, err := client.EnvKeyValueMaps.Create(ctx, "prod", &apigee.KeyValueMap{
    Name:      "env-config",
    Encrypted: true,
})

// Get an environment KVM
kvm, err := client.EnvKeyValueMaps.Get(ctx, "prod", "env-config")

// List all KVMs in an environment
kvms, err := client.EnvKeyValueMaps.List(ctx, "prod")

// Delete an environment KVM
err = client.EnvKeyValueMaps.Delete(ctx, "prod", "env-config")
```

#### Environment-Level KVM Entries

```go
// Add an entry to an environment KVM
entry, err := client.EnvKeyValueMapEntries.Create(ctx, "prod", "env-config", &apigee.KeyValueEntry{
    Name:  "db-host",
    Value: "localhost:5432",
})

// Get an entry
entry, err := client.EnvKeyValueMapEntries.Get(ctx, "prod", "env-config", "db-host")

// Update an entry
entry, err = client.EnvKeyValueMapEntries.Update(ctx, "prod", "env-config", "db-host", &apigee.KeyValueEntry{
    Name:  "db-host",
    Value: "db.example.com:5432",
})

// List all entries in an environment KVM
entries, err := client.EnvKeyValueMapEntries.List(ctx, "prod", "env-config")

// Delete an entry
err = client.EnvKeyValueMapEntries.Delete(ctx, "prod", "env-config", "db-host")
```

### Target Servers

Target Servers define backend server endpoints for load balancing and failover. They are environment-scoped resources.

```go
// Create a target server
ts, err := client.TargetServers.Create(ctx, "prod", &apigee.TargetServer{
    Name:        "backend-api",
    Host:        "api.example.com",
    Port:        443,
    Protocol:    "HTTP",
    IsEnabled:   true,
    Description: "Primary backend API server",
    SSLInfo: &apigee.SSLInfo{
        Enabled: true,
    },
})

// Get a target server
ts, err := client.TargetServers.Get(ctx, "prod", "backend-api")

// Update a target server
ts.Port = 8443
ts, err = client.TargetServers.Update(ctx, "prod", "backend-api", ts)

// List target servers (returns names only)
list, err := client.TargetServers.List(ctx, "prod")
for _, name := range list.TargetServerNames {
    fmt.Println(name)
}

// Delete a target server
err = client.TargetServers.Delete(ctx, "prod", "backend-api")
```

#### Target Server with mTLS

```go
// Create a target server with mutual TLS
ts, err := client.TargetServers.Create(ctx, "prod", &apigee.TargetServer{
    Name:      "secure-backend",
    Host:      "secure.example.com",
    Port:      443,
    Protocol:  "HTTP",
    IsEnabled: true,
    SSLInfo: &apigee.SSLInfo{
        Enabled:           true,
        ClientAuthEnabled: true,
        KeyStore:          "my-keystore",
        KeyAlias:          "my-key",
        TrustStore:        "my-truststore",
        Protocols:         []string{"TLSv1.2", "TLSv1.3"},
        CommonName: &apigee.CommonName{
            Value:         "*.example.com",
            WildcardMatch: true,
        },
    },
})
```

## Error Handling

The library provides helper functions to check for common HTTP error types:

```go
app, err := client.AppGroupApps.Get(ctx, "my-group", "my-app")
if err != nil {
    switch {
    case apigee.IsNotFound(err):
        // 404 - Resource doesn't exist
        fmt.Println("App not found")
    case apigee.IsConflict(err):
        // 409 - Resource already exists or conflict
        fmt.Println("Conflict occurred")
    case apigee.IsForbidden(err):
        // 403 - Permission denied
        fmt.Println("Access forbidden")
    case apigee.IsUnauthorized(err):
        // 401 - Authentication failed
        fmt.Println("Unauthorized")
    default:
        // Other error
        fmt.Printf("Error: %v\n", err)
    }
}
```

You can also access the full error details:

```go
if apiErr, ok := err.(*apigee.Error); ok {
    fmt.Printf("Status Code: %d\n", apiErr.StatusCode)
    fmt.Printf("Error Code: %d\n", apiErr.Code)
    fmt.Printf("Message: %s\n", apiErr.Message)
    fmt.Printf("Status: %s\n", apiErr.Status)
}
```

## Pagination

List operations support pagination through `ListOptions`:

```go
var allProducts []apigee.APIProduct

opts := &apigee.ListOptions{
    PageSize: 100,
}

for {
    resp, err := client.APIProducts.List(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }

    allProducts = append(allProducts, resp.APIProducts...)

    if resp.NextPageToken == "" {
        break
    }
    opts.PageToken = resp.NextPageToken
}
```

## Types Reference

### Core Types

| Type | Description |
|------|-------------|
| `Client` | The main API client |
| `Attribute` | Name-value pair for custom attributes |
| `ListOptions` | Pagination options (PageSize, PageToken, Filter) |
| `Error` | API error with StatusCode, Code, Message, Status |

### Resource Types

| Type | Description |
|------|-------------|
| `APIProduct` | API product definition |
| `AppGroup` | App group definition |
| `AppGroupApp` | App within an app group |
| `AppGroupAppKey` | Credential (key/secret) for an app |
| `APIProductRef` | Reference to an API product with approval status |
| `KeyValueMap` | Key Value Map (KVM) definition |
| `KeyValueEntry` | Entry (key-value pair) in a KVM |
| `TargetServer` | Target server definition for backend endpoints |
| `SSLInfo` | TLS/SSL configuration for target servers |
| `CommonName` | Certificate common name configuration |

## API Endpoints

| Resource | Endpoint |
|----------|----------|
| API Products | `organizations/{org}/apiproducts` |
| App Groups | `organizations/{org}/appgroups` |
| App Group Apps | `organizations/{org}/appgroups/{group}/apps` |
| App Keys | `organizations/{org}/appgroups/{group}/apps/{app}/keys` |
| Org KVMs | `organizations/{org}/keyvaluemaps` |
| Org KVM Entries | `organizations/{org}/keyvaluemaps/{kvm}/entries` |
| Env KVMs | `organizations/{org}/environments/{env}/keyvaluemaps` |
| Env KVM Entries | `organizations/{org}/environments/{env}/keyvaluemaps/{kvm}/entries` |
| Target Servers | `organizations/{org}/environments/{env}/targetservers` |

## License

MIT License
