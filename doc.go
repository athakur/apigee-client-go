// Package apigee provides a Go client library for the Google Apigee Management API.
//
// The library supports CRUD operations for:
//   - API Products
//   - App Groups
//   - App Group Apps
//   - App Credentials (Keys)
//   - Key Value Maps (KVMs) at organization and environment levels
//   - Target Servers (environment-level backend endpoints)
//
// # Authentication
//
// By default, the client uses Google Application Default Credentials (ADC).
// You can configure credentials in several ways:
//
//   - Set GOOGLE_APPLICATION_CREDENTIALS environment variable
//   - Run on Google Cloud with an attached service account
//   - Use gcloud auth application-default login
//
// The client requires the https://www.googleapis.com/auth/cloud-platform OAuth scope.
//
// # Creating a Client
//
// Create a client with default authentication:
//
//	ctx := context.Background()
//	client, err := apigee.NewClient(ctx, "my-organization")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Create a client with custom options:
//
//	client, err := apigee.NewClient(ctx, "my-organization",
//	    apigee.WithHTTPClient(customHTTPClient),
//	    apigee.WithBaseURL("https://custom-endpoint.example.com/v1"),
//	    apigee.WithTokenSource(customTokenSource),
//	)
//
// # Working with API Products
//
// API Products define the API resources and quotas available to applications:
//
//	// Create a product
//	product, err := client.APIProducts.Create(ctx, &apigee.APIProduct{
//	    Name:         "my-product",
//	    DisplayName:  "My Product",
//	    ApprovalType: "auto",
//	    Environments: []string{"test", "prod"},
//	})
//
//	// List products with pagination
//	resp, err := client.APIProducts.List(ctx, &apigee.ListOptions{PageSize: 100})
//
//	// Get, update, delete
//	product, err := client.APIProducts.Get(ctx, "my-product")
//	product, err = client.APIProducts.Update(ctx, "my-product", product)
//	err = client.APIProducts.Delete(ctx, "my-product")
//
// # Working with App Groups
//
// App Groups organize collections of apps:
//
//	// Create an app group
//	appGroup, err := client.AppGroups.Create(ctx, &apigee.AppGroup{
//	    Name:        "my-group",
//	    DisplayName: "My App Group",
//	})
//
//	// List, get, update, delete
//	groups, err := client.AppGroups.List(ctx, nil)
//	group, err := client.AppGroups.Get(ctx, "my-group")
//	group, err = client.AppGroups.Update(ctx, "my-group", group)
//	err = client.AppGroups.Delete(ctx, "my-group")
//
// # Working with Apps
//
// Apps represent client applications within an App Group:
//
//	// Create an app
//	app, err := client.AppGroupApps.Create(ctx, "my-group", &apigee.AppGroupApp{
//	    Name: "my-app",
//	    Attributes: []apigee.Attribute{
//	        {Name: "team", Value: "engineering"},
//	    },
//	})
//
//	// List, get, update, delete
//	apps, err := client.AppGroupApps.List(ctx, "my-group", nil)
//	app, err := client.AppGroupApps.Get(ctx, "my-group", "my-app")
//	app, err = client.AppGroupApps.Update(ctx, "my-group", "my-app", app)
//	err = client.AppGroupApps.Delete(ctx, "my-group", "my-app")
//
// # Working with App Credentials
//
// Credentials (keys) provide authentication for apps to access APIs:
//
//	// Generate a key (Apigee creates key/secret)
//	app, err := client.AppGroupAppKeys.Generate(ctx, "my-group", "my-app",
//	    &apigee.AppGroupAppKeyGenerateRequest{
//	        APIProducts:  []string{"my-product"},
//	        KeyExpiresIn: -1, // never expires
//	    },
//	)
//	// New key is in app.Credentials
//
//	// Create a key with custom values
//	key, err := client.AppGroupAppKeys.Create(ctx, "my-group", "my-app",
//	    &apigee.AppGroupAppKeyCreateRequest{
//	        ConsumerKey:    "custom-key",
//	        ConsumerSecret: "custom-secret",
//	    },
//	)
//
//	// Get, update, delete
//	key, err := client.AppGroupAppKeys.Get(ctx, "my-group", "my-app", "consumer-key")
//	key, err = client.AppGroupAppKeys.Update(ctx, "my-group", "my-app", "consumer-key", key)
//	err = client.AppGroupAppKeys.Delete(ctx, "my-group", "my-app", "consumer-key")
//
// # Working with Key Value Maps (KVMs)
//
// KVMs store configuration data at organization or environment level:
//
//	// Organization-level KVM
//	kvm, err := client.KeyValueMaps.Create(ctx, &apigee.KeyValueMap{
//	    Name:      "my-config",
//	    Encrypted: true,
//	})
//
//	// Add entry to org KVM
//	entry, err := client.KeyValueMapEntries.Create(ctx, "my-config", &apigee.KeyValueEntry{
//	    Name:  "api-key",
//	    Value: "secret123",
//	})
//
//	// Get, update, delete entries
//	entry, err := client.KeyValueMapEntries.Get(ctx, "my-config", "api-key")
//	entry, err = client.KeyValueMapEntries.Update(ctx, "my-config", "api-key", entry)
//	err = client.KeyValueMapEntries.Delete(ctx, "my-config", "api-key")
//
//	// Environment-level KVM
//	kvm, err := client.EnvKeyValueMaps.Create(ctx, "prod", &apigee.KeyValueMap{
//	    Name:      "env-config",
//	    Encrypted: true,
//	})
//
//	// Add entry to env KVM
//	entry, err := client.EnvKeyValueMapEntries.Create(ctx, "prod", "env-config", &apigee.KeyValueEntry{
//	    Name:  "db-host",
//	    Value: "localhost:5432",
//	})
//
//	// List, get, update, delete env KVM entries
//	entries, err := client.EnvKeyValueMapEntries.List(ctx, "prod", "env-config")
//	entry, err := client.EnvKeyValueMapEntries.Get(ctx, "prod", "env-config", "db-host")
//	entry, err = client.EnvKeyValueMapEntries.Update(ctx, "prod", "env-config", "db-host", entry)
//	err = client.EnvKeyValueMapEntries.Delete(ctx, "prod", "env-config", "db-host")
//
// # Working with Target Servers
//
// Target Servers define backend endpoints for load balancing and failover:
//
//	// Create a target server
//	ts, err := client.TargetServers.Create(ctx, "prod", &apigee.TargetServer{
//	    Name:      "backend-api",
//	    Host:      "api.example.com",
//	    Port:      443,
//	    Protocol:  "HTTP",
//	    IsEnabled: true,
//	    SSLInfo: &apigee.SSLInfo{
//	        Enabled: true,
//	    },
//	})
//
//	// List target servers (returns names only)
//	list, err := client.TargetServers.List(ctx, "prod")
//	for _, name := range list.TargetServerNames {
//	    fmt.Println(name)
//	}
//
//	// Get, update, delete
//	ts, err := client.TargetServers.Get(ctx, "prod", "backend-api")
//	ts.Port = 8443
//	ts, err = client.TargetServers.Update(ctx, "prod", "backend-api", ts)
//	err = client.TargetServers.Delete(ctx, "prod", "backend-api")
//
// # Error Handling
//
// The library provides helper functions to check for common error types:
//
//	app, err := client.AppGroupApps.Get(ctx, "my-group", "my-app")
//	if apigee.IsNotFound(err) {
//	    // Handle 404 Not Found
//	}
//	if apigee.IsConflict(err) {
//	    // Handle 409 Conflict
//	}
//	if apigee.IsForbidden(err) {
//	    // Handle 403 Forbidden
//	}
//	if apigee.IsUnauthorized(err) {
//	    // Handle 401 Unauthorized
//	}
//
// Access full error details:
//
//	if apiErr, ok := err.(*apigee.Error); ok {
//	    fmt.Printf("Status: %d, Code: %d, Message: %s\n",
//	        apiErr.StatusCode, apiErr.Code, apiErr.Message)
//	}
//
// # Pagination
//
// List operations support pagination:
//
//	opts := &apigee.ListOptions{PageSize: 100}
//	for {
//	    resp, err := client.APIProducts.List(ctx, opts)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    // Process resp.APIProducts
//	    if resp.NextPageToken == "" {
//	        break
//	    }
//	    opts.PageToken = resp.NextPageToken
//	}
package apigee
