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
	kvmName := "test-kvm-" + suffix

	targetServerName := "test-ts-" + suffix
	envName := "eval" // Default environment for testing

	// Track what we've created for cleanup
	var productCreated, appGroupCreated, appCreated, kvmCreated, targetServerCreated bool

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
		if kvmCreated {
			fmt.Printf("Deleting KVM: %s\n", kvmName)
			if err := client.KeyValueMaps.Delete(ctx, kvmName); err != nil {
				fmt.Printf("  Warning: failed to delete KVM: %v\n", err)
			} else {
				fmt.Println("  Deleted successfully")
			}
		}
		if targetServerCreated {
			fmt.Printf("Deleting Target Server: %s\n", targetServerName)
			if err := client.TargetServers.Delete(ctx, envName, targetServerName); err != nil {
				fmt.Printf("  Warning: failed to delete target server: %v\n", err)
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

	// ==================== KVM Tests ====================

	// 13. Create Organization-level KVM
	fmt.Println("\n--- Create Organization KVM ---")
	kvm, err := client.KeyValueMaps.Create(ctx, &apigee.KeyValueMap{
		Name:      kvmName,
		Encrypted: true,
	})
	if err != nil {
		log.Fatalf("Failed to create KVM: %v", err)
	}
	kvmCreated = true
	fmt.Printf("Created KVM: %s (Encrypted: %v)\n", kvm.Name, kvm.Encrypted)

	// 14. Get Organization-level KVM
	fmt.Println("\n--- Get Organization KVM ---")
	kvm, err = client.KeyValueMaps.Get(ctx, kvmName)
	if err != nil {
		log.Fatalf("Failed to get KVM: %v", err)
	}
	fmt.Printf("Got KVM: %s\n", kvm.Name)

	// 15. List Organization-level KVMs
	fmt.Println("\n--- List Organization KVMs ---")
	kvms, err := client.KeyValueMaps.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list KVMs: %v", err)
	}
	fmt.Printf("Found %d KVMs\n", len(kvms.KeyValueMapNames))
	for _, name := range kvms.KeyValueMapNames {
		fmt.Printf("  - %s\n", name)
	}

	// 16. Create KVM Entry
	fmt.Println("\n--- Create KVM Entry ---")
	entry, err := client.KeyValueMapEntries.Create(ctx, kvmName, &apigee.KeyValueEntry{
		Name:  "test-key",
		Value: "test-value-123",
	})
	if err != nil {
		log.Fatalf("Failed to create KVM entry: %v", err)
	}
	fmt.Printf("Created entry: %s = %s\n", entry.Name, entry.Value)

	// 17. Get KVM Entry
	fmt.Println("\n--- Get KVM Entry ---")
	entry, err = client.KeyValueMapEntries.Get(ctx, kvmName, "test-key")
	if err != nil {
		log.Fatalf("Failed to get KVM entry: %v", err)
	}
	fmt.Printf("Got entry: %s = %s\n", entry.Name, entry.Value)

	// 18. Update KVM Entry
	fmt.Println("\n--- Update KVM Entry ---")
	entry, err = client.KeyValueMapEntries.Update(ctx, kvmName, "test-key", &apigee.KeyValueEntry{
		Name:  "test-key",
		Value: "updated-value-456",
	})
	if err != nil {
		log.Fatalf("Failed to update KVM entry: %v", err)
	}
	fmt.Printf("Updated entry: %s = %s\n", entry.Name, entry.Value)

	// 19. List KVM Entries
	fmt.Println("\n--- List KVM Entries ---")
	kvmEntries, err := client.KeyValueMapEntries.List(ctx, kvmName)
	if err != nil {
		log.Fatalf("Failed to list KVM entries: %v", err)
	}
	fmt.Printf("Found %d entries\n", len(kvmEntries.KeyValueEntries))
	for _, e := range kvmEntries.KeyValueEntries {
		fmt.Printf("  - %s = %s\n", e.Name, e.Value)
	}

	// 20. Delete KVM Entry
	fmt.Println("\n--- Delete KVM Entry ---")
	err = client.KeyValueMapEntries.Delete(ctx, kvmName, "test-key")
	if err != nil {
		log.Fatalf("Failed to delete KVM entry: %v", err)
	}
	fmt.Println("Deleted KVM entry successfully")

	// ==================== Target Server Tests ====================

	// 21. Create Target Server
	fmt.Println("\n--- Create Target Server ---")
	ts, err := client.TargetServers.Create(ctx, envName, &apigee.TargetServer{
		Name:        targetServerName,
		Host:        "api.example.com",
		Port:        443,
		Protocol:    "HTTP",
		IsEnabled:   true,
		Description: "Test target server for Go client testing",
		SSLInfo: &apigee.SSLInfo{
			Enabled: true,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create target server: %v", err)
	}
	targetServerCreated = true
	fmt.Printf("Created target server: %s (Host: %s, Port: %d)\n", ts.Name, ts.Host, ts.Port)

	// 22. Get Target Server
	fmt.Println("\n--- Get Target Server ---")
	ts, err = client.TargetServers.Get(ctx, envName, targetServerName)
	if err != nil {
		log.Fatalf("Failed to get target server: %v", err)
	}
	fmt.Printf("Got target server: %s, Host: %s, Port: %d, Enabled: %v\n", ts.Name, ts.Host, ts.Port, ts.IsEnabled)

	// 23. List Target Servers
	fmt.Println("\n--- List Target Servers ---")
	tsList, err := client.TargetServers.List(ctx, envName)
	if err != nil {
		log.Fatalf("Failed to list target servers: %v", err)
	}
	fmt.Printf("Found %d target servers\n", len(tsList.TargetServerNames))
	for _, name := range tsList.TargetServerNames {
		fmt.Printf("  - %s\n", name)
	}

	// 24. Update Target Server
	fmt.Println("\n--- Update Target Server ---")
	ts.Port = 8443
	ts.Description = "Updated test target server"
	ts, err = client.TargetServers.Update(ctx, envName, targetServerName, ts)
	if err != nil {
		log.Fatalf("Failed to update target server: %v", err)
	}
	fmt.Printf("Updated target server: %s, Port: %d, Description: %s\n", ts.Name, ts.Port, ts.Description)

	// 25. Delete Target Server
	fmt.Println("\n--- Delete Target Server ---")
	err = client.TargetServers.Delete(ctx, envName, targetServerName)
	if err != nil {
		log.Fatalf("Failed to delete target server: %v", err)
	}
	targetServerCreated = false
	fmt.Println("Deleted target server successfully")

	fmt.Println("\n--- All tests passed! ---")
}
