package jet

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

// TestTokenComparison runs a comparison test and prints results
func TestTokenComparison(t *testing.T) {
	// Sample data structure
	type Product struct {
		ID       int
		Name     string
		Category string
	}

	testData := []Product{
		{ID: 1, Name: "Laptop", Category: "Electronics"},
		{ID: 2, Name: "Mouse", Category: "Electronics"},
		{ID: 3, Name: "Book", Category: "Literature"},
		{ID: 4, Name: "Desk", Category: "Furniture"},
		{ID: 5, Name: "Chair", Category: "Furniture"},
	}

	comparison, err := CompareTokens(testData)
	if err != nil {
		t.Fatalf("CompareTokens failed: %v", err)
	}

	comparisonNormalized, err := CompareTokensNormalized(testData)
	if err != nil {
		t.Fatalf("CompareTokensNormalized failed: %v", err)
	}
	// Test flattened format
	comparisonFlat, err := CompareTokensFlattened(testData)
	if err != nil {
		t.Fatalf("CompareTokensFlattened failed: %v", err)
	}

	log.Printf("\n=== Token Comparison Results (Flattened) ===")
	log.Printf("JSON Indented:   	%d bytes, %d tokens", comparison.JSONInBytes, comparison.JSONInTokens)
	log.Printf("JSON:            	%d bytes, %d tokens", comparison.JSONBytes, comparison.JSONTokens)
	log.Printf("Jet:   		%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparison.JetBytes, comparison.JetTokens, comparison.TokenSavings, comparison.TokenInSavings)
	log.Printf("Jet Normalized:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonNormalized.JetBytes, comparisonNormalized.JetTokens, comparisonNormalized.TokenSavings, comparisonNormalized.TokenInSavings)
	log.Printf("Jet Flattened:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonFlat.JetBytes, comparisonFlat.JetTokens, comparisonFlat.TokenSavings, comparisonFlat.TokenInSavings)
	log.Printf("Ultimate Savings: %.2f%% bytes, %.2f%% tokens", comparisonFlat.ByteSavings, comparisonFlat.TokenSavings)
}

func TestTokenComparisonLargeData(t *testing.T) {
	// Generate large sample data
	type Item struct {
		ID    int
		Title string
		Price float64
		Tags  []string
	}

	var largeData []Item
	for i := 1; i <= 1000; i++ {
		largeData = append(largeData, Item{
			ID:    i,
			Title: "Item " + fmt.Sprintf("%d", i),
			Price: float64(i) * 1.5,
			Tags:  []string{"tag1", "tag2", "tag3"},
		})
	}

	comparison, err := CompareTokens(largeData)
	if err != nil {
		t.Fatalf("CompareTokens failed: %v", err)
	}

	comparisonNormalized, err := CompareTokensNormalized(largeData)
	if err != nil {
		t.Fatalf("CompareTokensNormalized failed: %v", err)
	}

	// Test flattened format
	comparisonFlat, err := CompareTokensFlattened(largeData)
	if err != nil {
		t.Fatalf("CompareTokensFlattened failed: %v", err)
	}

	log.Printf("\n=== Large Data Token Comparison Results (Flattened) ===")
	log.Printf("JSON Indented:   	%d bytes, %d tokens", comparisonFlat.JSONInBytes, comparisonFlat.JSONInTokens)
	log.Printf("JSON:            	%d bytes, %d tokens", comparisonFlat.JSONBytes, comparisonFlat.JSONTokens)
	log.Printf("Jet:   		%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparison.JetBytes, comparison.JetTokens, comparison.TokenSavings, comparison.TokenInSavings)
	log.Printf("Normalized:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonNormalized.JetBytes, comparisonNormalized.JetTokens, comparisonNormalized.TokenSavings, comparisonNormalized.TokenInSavings)
	log.Printf("Flattened:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonFlat.JetBytes, comparisonFlat.JetTokens, comparisonFlat.TokenSavings, comparisonFlat.TokenInSavings)
	log.Printf("Ultimate Savings: %.2f%% bytes, %.2f%% tokens", comparisonFlat.ByteSavings, comparisonFlat.TokenSavings)
}

func TestTokenComparisonMultiNested(t *testing.T) {
	// Complex nested structure mimicking real-world data
	type Address struct {
		Street  string
		City    string
		ZipCode string
		Country string
	}

	type OrderItem struct {
		ProductID   int
		ProductName string
		Quantity    int
		Price       float64
		Discount    float64
	}

	type Order struct {
		OrderID     int
		OrderDate   string
		Status      string
		TotalAmount float64
		Items       []OrderItem
	}

	type Customer struct {
		ID            int
		Name          string
		Email         string
		Phone         string
		Address       Address
		Orders        []Order
		IsActive      bool
		LoyaltyPoints int
	}

	// Generate nested data
	var customers []Customer
	for i := 1; i <= 100; i++ {
		var orders []Order
		for j := 1; j <= 5; j++ {
			var items []OrderItem
			for k := 1; k <= 3; k++ {
				items = append(items, OrderItem{
					ProductID:   k * 100,
					ProductName: fmt.Sprintf("Product-%d-%d", j, k),
					Quantity:    k,
					Price:       99.99 * float64(k),
					Discount:    0.10 * float64(k),
				})
			}
			orders = append(orders, Order{
				OrderID:     j * 1000,
				OrderDate:   fmt.Sprintf("2025-01-%02d", j),
				Status:      "completed",
				TotalAmount: 299.97,
				Items:       items,
			})
		}

		customers = append(customers, Customer{
			ID:    i,
			Name:  fmt.Sprintf("Customer %d", i),
			Email: fmt.Sprintf("customer%d@example.com", i),
			Phone: fmt.Sprintf("+1-555-0%03d", i),
			Address: Address{
				Street:  fmt.Sprintf("%d Main St", i*10),
				City:    "Metropolis",
				ZipCode: fmt.Sprintf("100%02d", i),
				Country: "USA",
			},
			Orders:        orders,
			IsActive:      i%2 == 0,
			LoyaltyPoints: i * 100,
		})
	}

	comparison, err := CompareTokens(customers)
	if err != nil {
		t.Fatalf("CompareTokens failed: %v", err)
	}

	comparisonNormalized, err := CompareTokensNormalized(customers)
	if err != nil {
		t.Fatalf("CompareTokensNormalized failed: %v", err)
	}
	// Test flattened format
	comparisonFlat, err := CompareTokensFlattened(customers)
	if err != nil {
		t.Fatalf("CompareTokensFlattened failed: %v", err)
	}

	log.Printf("\n=== Multi-Nested Data Token Comparison Results ===")
	log.Printf("JSON Indented:   	%d bytes, %d tokens", comparison.JSONInBytes, comparison.JSONInTokens)
	log.Printf("JSON:            	%d bytes, %d tokens", comparison.JSONBytes, comparison.JSONTokens)
	log.Printf("Jet:   		%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparison.JetBytes, comparison.JetTokens, comparison.TokenSavings, comparison.TokenInSavings)
	log.Printf("Normalized:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonNormalized.JetBytes, comparisonNormalized.JetTokens, comparisonNormalized.TokenSavings, comparisonNormalized.TokenInSavings)
	log.Printf("Flattened:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonFlat.JetBytes, comparisonFlat.JetTokens, comparisonFlat.TokenSavings, comparisonFlat.TokenInSavings)
	log.Printf("Ultimate Savings: %.2f%% bytes, %.2f%% tokens", comparisonFlat.ByteSavings, comparisonFlat.TokenSavings)
}
func TestTokenComparisonHighColumnCount(t *testing.T) {
	type Customer struct {
		ID            int
		Name          string
		Email         string
		Phone         string
		IsActive      bool
		LoyaltyPoints int
		OrderID       int
		OrderDate     string
		Status        string
		TotalAmount   float64
		ProductID     int
		ProductName   string
		Quantity      int
		Price         float64
		Discount      float64

		Street  string
		City    string
		ZipCode string
		Country string
	}

	// Generate nested data
	var customers []Customer
	for i := 1; i <= 100; i++ {

		customers = append(customers, Customer{
			ID:    i,
			Name:  fmt.Sprintf("Customer %d", i),
			Email: fmt.Sprintf("customer%d@example.com", i),
			Phone: fmt.Sprintf("+1-555-0%03d", i),

			Street:      fmt.Sprintf("%d Main St", i*10),
			City:        "Metropolis",
			ZipCode:     fmt.Sprintf("100%02d", i),
			Country:     "USA",
			OrderID:     rand.Intn(1000),
			OrderDate:   fmt.Sprintf("2025-01-%02d", i),
			Status:      "completed",
			TotalAmount: 299.97,

			ProductID:     rand.Intn(1000),
			ProductName:   fmt.Sprintf("Product-%d-%d", i, rand.Intn(100)),
			Quantity:      rand.Intn(10) + 1,
			Price:         99.99 * rand.Float64(),
			Discount:      0.10 * rand.Float64(),
			IsActive:      i%2 == 0,
			LoyaltyPoints: i * 100,
		})
	}

	comparison, err := CompareTokens(customers)
	if err != nil {
		t.Fatalf("CompareTokens failed: %v", err)
	}

	comparisonNormalized, err := CompareTokensNormalized(customers)
	if err != nil {
		t.Fatalf("CompareTokensNormalized failed: %v", err)
	}
	// Test flattened format
	comparisonFlat, err := CompareTokensFlattened(customers)
	if err != nil {
		t.Fatalf("CompareTokensFlattened failed: %v", err)
	}

	log.Printf("\n=== Multi-Nested Data Token Comparison Results ===")
	log.Printf("JSON Indented:   	%d bytes, %d tokens", comparison.JSONInBytes, comparison.JSONInTokens)
	log.Printf("JSON:            	%d bytes, %d tokens", comparison.JSONBytes, comparison.JSONTokens)
	log.Printf("Jet:   		%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparison.JetBytes, comparison.JetTokens, comparison.TokenSavings, comparison.TokenInSavings)
	log.Printf("Normalized:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonNormalized.JetBytes, comparisonNormalized.JetTokens, comparisonNormalized.TokenSavings, comparisonNormalized.TokenInSavings)
	log.Printf("Flattened:   	%d bytes, %d tokens, %.2f%% tokens, %.2f%% ind tokens", comparisonFlat.JetBytes, comparisonFlat.JetTokens, comparisonFlat.TokenSavings, comparisonFlat.TokenInSavings)
	log.Printf("Ultimate Savings: %.2f%% bytes, %.2f%% tokens", comparisonFlat.ByteSavings, comparisonFlat.TokenSavings)
}
