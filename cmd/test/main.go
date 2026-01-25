package main

import (
	"context"
	"fmt"
	"log"
	"time"

	apigee "github.com/athakur/apigee-client-go"
)

func main() {
	ctx := context.Background()
	org := "techkur-dev"

	fmt.Printf("Creating Apigee client for organization: %s\n", org)
	client, err := apigee.NewClient(ctx, org)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use unique names based on timestamp
	suffix := fmt.Sprintf("%d", time.Now().Unix())
	productName := "test-product-" + suffix
	appGroupName := "test-appgroup-" + suffix
	appName := "test-app-" + suffix

	// Track what we've created for cleanup
	var productCreated, appGroupCreated, appCreated bool

	// Cleanup function
	defer func() {
		fmt.Println("\n--- Cleanup ---")
		if appCreated {
			fmt.Printf("Deleting app: %s\n", appName)
			if err := client.AppGroupApps.Delete(ctx, appGroupName, appName); err != nil {
				fmt.Printf("  Warning: failed to delete app: %v\n", err)
			} else {
				fmt.Println("  Deleted successfully")
			}
		}
		if appGroupCreated {
			fmt.Printf("Deleting app group: %s\n", appGroupName)
			if err := client.AppGroups.Delete(ctx, appGroupName); err != nil {
				fmt.Printf("  Warning: failed to delete app group: %v\n", err)
			} else {
				fmt.Println("  Deleted successfully")
			}
		}
		if productCreated {
			fmt.Printf("Deleting API product: %s\n", productName)
			if err := client.APIProducts.Delete(ctx, productName); err != nil {
				fmt.Printf("  Warning: failed to delete product: %v\n", err)
			} else {
				fmt.Println("  Deleted successfully")
			}
		}
	}()

	// 1. Create API Product (needed for app keys)
	fmt.Println("\n--- Create API Product ---")
	product, err := client.APIProducts.Create(ctx, &apigee.APIProduct{
		Name:         productName,
		DisplayName:  "Test Product " + suffix,
		Description:  "Test product for Go client testing",
		ApprovalType: "auto",
		Attributes: []apigee.Attribute{
			{Name: "test-attr", Value: "test-value"},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create API product: %v", err)
	}
	productCreated = true
	fmt.Printf("Created API product: %s\n", product.Name)

	// 2. Create AppGroup
	fmt.Println("\n--- Create AppGroup ---")
	appGroup, err := client.AppGroups.Create(ctx, &apigee.AppGroup{
		Name:        appGroupName,
		DisplayName: "Test App Group " + suffix,
		Attributes: []apigee.Attribute{
			{Name: "env", Value: "test"},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create app group: %v", err)
	}
	appGroupCreated = true
	fmt.Printf("Created app group: %s (ID: %s)\n", appGroup.Name, appGroup.AppGroupID)

	// 3. Get AppGroup
	fmt.Println("\n--- Get AppGroup ---")
	appGroup, err = client.AppGroups.Get(ctx, appGroupName)
	if err != nil {
		log.Fatalf("Failed to get app group: %v", err)
	}
	fmt.Printf("Got app group: %s, Status: %s\n", appGroup.Name, appGroup.Status)

	// 4. List AppGroups
	fmt.Println("\n--- List AppGroups ---")
	appGroups, err := client.AppGroups.List(ctx, &apigee.ListOptions{PageSize: 5})
	if err != nil {
		log.Fatalf("Failed to list app groups: %v", err)
	}
	fmt.Printf("Found %d app groups\n", len(appGroups.AppGroups))
	for _, ag := range appGroups.AppGroups {
		fmt.Printf("  - %s\n", ag.Name)
	}

	// 5. Update AppGroup
	fmt.Println("\n--- Update AppGroup ---")
	appGroup.DisplayName = "Updated Test App Group " + suffix
	appGroup.Attributes = append(appGroup.Attributes, apigee.Attribute{Name: "updated", Value: "true"})
	appGroup, err = client.AppGroups.Update(ctx, appGroupName, appGroup)
	if err != nil {
		log.Fatalf("Failed to update app group: %v", err)
	}
	fmt.Printf("Updated app group display name: %s\n", appGroup.DisplayName)

	// 6. Create AppGroupApp
	fmt.Println("\n--- Create AppGroupApp ---")
	app, err := client.AppGroupApps.Create(ctx, appGroupName, &apigee.AppGroupApp{
		Name: appName,
		Attributes: []apigee.Attribute{
			{Name: "app-attr", Value: "app-value"},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	appCreated = true
	fmt.Printf("Created app: %s (ID: %s)\n", app.Name, app.AppID)
	fmt.Printf("Initial credentials count: %d\n", len(app.Credentials))

	// 7. Get AppGroupApp
	fmt.Println("\n--- Get AppGroupApp ---")
	app, err = client.AppGroupApps.Get(ctx, appGroupName, appName)
	if err != nil {
		log.Fatalf("Failed to get app: %v", err)
	}
	fmt.Printf("Got app: %s, Status: %s\n", app.Name, app.Status)

	// 8. List AppGroupApps
	fmt.Println("\n--- List AppGroupApps ---")
	apps, err := client.AppGroupApps.List(ctx, appGroupName, nil)
	if err != nil {
		log.Fatalf("Failed to list apps: %v", err)
	}
	fmt.Printf("Found %d apps in app group\n", len(apps.AppGroupApps))
	for _, a := range apps.AppGroupApps {
		fmt.Printf("  - %s\n", a.Name)
	}

	// 9. Generate a key for the app (Apigee generates key/secret)
	fmt.Println("\n--- Generate AppGroupAppKey ---")
	app, err = client.AppGroupAppKeys.Generate(ctx, appGroupName, appName, &apigee.AppGroupAppKeyGenerateRequest{
		APIProducts:  []string{productName},
		KeyExpiresIn: -1, // never expires
	})
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}
	fmt.Printf("Generated key, app now has %d credential(s)\n", len(app.Credentials))
	if len(app.Credentials) > 0 {
		cred := app.Credentials[len(app.Credentials)-1]
		fmt.Printf("  Consumer Key: %s...\n", cred.ConsumerKey[:16])
		fmt.Printf("  Status: %s\n", cred.Status)
		fmt.Printf("  API Products: %d\n", len(cred.APIProducts))
	}

	// 10. Get the generated key
	if len(app.Credentials) > 0 {
		consumerKey := app.Credentials[len(app.Credentials)-1].ConsumerKey
		fmt.Println("\n--- Get AppGroupAppKey ---")
		key, err := client.AppGroupAppKeys.Get(ctx, appGroupName, appName, consumerKey)
		if err != nil {
			log.Fatalf("Failed to get key: %v", err)
		}
		fmt.Printf("Got key: %s...\n", key.ConsumerKey[:16])
		fmt.Printf("  Status: %s\n", key.Status)
		fmt.Printf("  Issued At: %s\n", key.IssuedAt)
	}

	// 11. Create a key with custom key/secret
	fmt.Println("\n--- Create AppGroupAppKey with custom values ---")
	customKey, err := client.AppGroupAppKeys.Create(ctx, appGroupName, appName, &apigee.AppGroupAppKeyCreateRequest{
		ConsumerKey:    "custom-key-" + suffix,
		ConsumerSecret: "custom-secret-" + suffix,
	})
	if err != nil {
		log.Fatalf("Failed to create custom key: %v", err)
	}
	fmt.Printf("Created custom key: %s\n", customKey.ConsumerKey)

	// 12. Delete the custom key
	fmt.Println("\n--- Delete custom AppGroupAppKey ---")
	err = client.AppGroupAppKeys.Delete(ctx, appGroupName, appName, customKey.ConsumerKey)
	if err != nil {
		log.Fatalf("Failed to delete custom key: %v", err)
	}
	fmt.Println("Deleted custom key successfully")

	fmt.Println("\n--- All tests passed! ---")
}
