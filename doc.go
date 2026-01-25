// Package apigee provides a Go client library for the Google Apigee Management API.
//
// The library supports CRUD operations for:
//   - API Products
//   - App Groups
//   - App Group Apps
//   - App Credentials (Keys)
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
