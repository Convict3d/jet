package jet

import (
	"strings"
	"testing"
)

func TestMarshalFlattenedSimpleStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	result, err := MarshalFlattened(Person{Name: "Alice", Age: 30, City: "Wonderland"})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened simple struct output:\n%s", resultStr)

	// Simple structs should be same as normal format
	if !strings.Contains(resultStr, "age: 30") {
		t.Errorf("Expected 'age: 30' in output")
	}
	if !strings.Contains(resultStr, "city: Wonderland") {
		t.Errorf("Expected 'city: Wonderland' in output")
	}
	if !strings.Contains(resultStr, "name: Alice") {
		t.Errorf("Expected 'name: Alice' in output")
	}
}

func TestMarshalFlattenedSliceOfStructs(t *testing.T) {
	type Product struct {
		ID       int
		Name     string
		Category string
	}

	testData := []Product{
		{ID: 1, Name: "Laptop", Category: "Electronics"},
		{ID: 2, Name: "Mouse", Category: "Electronics"},
		{ID: 3, Name: "Book", Category: "Literature"},
	}

	result, err := MarshalFlattened(testData)
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened slice of structs output:\n%s", resultStr)

	// Check tabular format (should be same as normal for simple structs)
	if !strings.Contains(resultStr, "Laptop") {
		t.Errorf("Expected 'Laptop' in output")
	}
	if !strings.Contains(resultStr, "Electronics") {
		t.Errorf("Expected 'Electronics' in output")
	}
	if !strings.Contains(resultStr, "Literature") {
		t.Errorf("Expected 'Literature' in output")
	}

	// Should NOT have any nested blocks for simple data
	if strings.Contains(resultStr, ">") {
		t.Errorf("Expected no nested blocks for simple structs")
	}
}

func TestMarshalFlattenedNestedStruct(t *testing.T) {
	type Address struct {
		Street  string
		City    string
		ZipCode string
	}

	type Person struct {
		Name    string
		Age     int
		Address Address
	}

	result, err := MarshalFlattened(Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "Wonderland",
			ZipCode: "12345",
		},
	})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened nested struct output:\n%s", resultStr)

	// Check nested object structure (should be same as normal for single struct)
	if !strings.Contains(resultStr, "address:") {
		t.Errorf("Expected 'address:' nested block")
	}
	if !strings.Contains(resultStr, "city: Wonderland") {
		t.Errorf("Expected nested 'city: Wonderland'")
	}
	if !strings.Contains(resultStr, "street: 123 Main St") {
		t.Errorf("Expected nested 'street: 123 Main St'")
	}
}

func TestMarshalFlattenedSliceWithNestedObjects(t *testing.T) {
	type Profile struct {
		Username string
		Email    string
	}

	type Person struct {
		Name    string
		Age     int
		City    string
		Profile Profile
	}

	result, err := MarshalFlattened([]Person{
		{
			Name: "Alice",
			Age:  30,
			City: "Wonderland",
			Profile: Profile{
				Username: "alice",
				Email:    "alice@example.com",
			},
		},
		{
			Name: "Bob",
			Age:  25,
			City: "Builderland",
			Profile: Profile{
				Username: "bob",
				Email:    "bob@example.com",
			},
		},
	})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened slice with nested objects output:\n%s", resultStr)

	// Check flattened schema in header with commas
	if !strings.Contains(resultStr, "profile{email,username}") && !strings.Contains(resultStr, "profile{username,email}") {
		t.Errorf("Expected flattened schema 'profile{email,username}' or similar in header")
	}

	// Check that there are NO nested profile blocks (completely flattened)
	if strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected NO '> profile:' nested blocks in flattened format")
	}

	// Check that values are completely inline
	if !strings.Contains(resultStr, "alice@example.com") {
		t.Errorf("Expected 'alice@example.com' inline in output")
	}
	if !strings.Contains(resultStr, "alice") {
		t.Errorf("Expected 'alice' inline in output")
	}

	// Verify it's all on fewer lines (no nested blocks)
	lines := strings.Split(resultStr, "\n")
	dataLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 && !strings.HasSuffix(trimmed, ":") {
			dataLines++
		}
	}

	// Should have just 2 data rows (one per person)
	if dataLines < 2 {
		t.Errorf("Expected at least 2 data rows in flattened format")
	}
}

func TestMarshalFlattenedDeeplyNestedStructure(t *testing.T) {
	type Address struct {
		Street  string
		City    string
		Country string
	}

	type Profile struct {
		Username string
		Address  Address
	}

	type Person struct {
		Name    string
		Profile Profile
	}

	result, err := MarshalFlattened([]Person{
		{
			Name: "Alice",
			Profile: Profile{
				Username: "alice",
				Address: Address{
					Street:  "123 Main St",
					City:    "Wonderland",
					Country: "Fantasy",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened deeply nested output:\n%s", resultStr)

	// Cannot fully flatten deeply nested structures (3+ levels)
	// Should show [nested] placeholder or partial flattening
	if strings.Contains(resultStr, "[nested]") {
		t.Logf("Correctly shows [nested] placeholder for non-flattenable structure")
	}
}

func TestMarshalFlattenedMixedTypes(t *testing.T) {
	type Profile struct {
		Active   bool
		Score    float64
		Username string
	}

	type Person struct {
		Name    string
		Age     int
		Profile Profile
	}

	result, err := MarshalFlattened([]Person{
		{
			Name: "Alice",
			Age:  30,
			Profile: Profile{
				Active:   true,
				Score:    95.5,
				Username: "alice",
			},
		},
	})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened mixed types output:\n%s", resultStr)

	// Check that different types are flattened inline
	if strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected NO nested blocks in flattened format")
	}

	// Should contain the values inline
	if !strings.Contains(resultStr, "alice") {
		t.Errorf("Expected 'alice' inline in output")
	}
	if !strings.Contains(resultStr, "95.5") {
		t.Errorf("Expected '95.5' inline in output")
	}
	if !strings.Contains(resultStr, "true") {
		t.Errorf("Expected 'true' inline in output")
	}
}

func TestMarshalFlattenedComplexStructure(t *testing.T) {
	type OrderItem struct {
		ProductID int
		Quantity  int
		Price     float64
	}

	type ShippingAddress struct {
		Street  string
		City    string
		ZipCode string
	}

	type Order struct {
		OrderID int
		Items   []OrderItem
		Address ShippingAddress
	}

	result, err := MarshalFlattened([]Order{
		{
			OrderID: 1001,
			Items: []OrderItem{
				{ProductID: 100, Quantity: 2, Price: 99.99},
				{ProductID: 101, Quantity: 1, Price: 49.99},
			},
			Address: ShippingAddress{
				Street:  "123 Main St",
				City:    "Wonderland",
				ZipCode: "12345",
			},
		},
	})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened complex structure output:\n%s", resultStr)

	// Check for flattened address schema with commas
	if strings.Contains(resultStr, "address{") {
		t.Logf("Found flattened address schema in header")
	}

	// Items table cannot be flattened (it's a nested table)
	// Should show [table] placeholder or similar
	if strings.Contains(resultStr, "[table]") || strings.Contains(resultStr, "> items") {
		t.Logf("Correctly handles nested table that cannot be fully flattened")
	}
}

func TestMarshalFlattenedEmptySlice(t *testing.T) {
	type Product struct {
		ID   int
		Name string
	}

	result, err := MarshalFlattened([]Product{})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	// Empty slice should produce empty or minimal output
	if len(result) > 0 {
		t.Logf("Empty slice output: %s", string(result))
	}
}

func TestMarshalFlattenedSingleLevel(t *testing.T) {
	type Metadata struct {
		Created string
		Author  string
		Version string
	}

	type Document struct {
		Title    string
		Content  string
		Metadata Metadata
	}

	result, err := MarshalFlattened([]Document{
		{
			Title:   "Report",
			Content: "Lorem ipsum",
			Metadata: Metadata{
				Created: "2025-11-11",
				Author:  "Alice",
				Version: "1.0",
			},
		},
		{
			Title:   "Summary",
			Content: "Brief overview",
			Metadata: Metadata{
				Created: "2025-11-10",
				Author:  "Bob",
				Version: "2.0",
			},
		},
	})
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened single level output:\n%s", resultStr)

	// Check completely flattened (no nested blocks)
	if strings.Contains(resultStr, "> metadata:") {
		t.Errorf("Expected NO nested blocks in flattened format")
	}

	// All metadata values should be inline
	if !strings.Contains(resultStr, "2025-11-11") || !strings.Contains(resultStr, "Alice") || !strings.Contains(resultStr, "1.0") {
		t.Errorf("Expected all metadata values inline")
	}
}

func TestMarshalFlattened(t *testing.T) {
	type Profile struct {
		Email    string
		Username string
	}

	type Person struct {
		Name    string
		Age     int
		City    string
		Profile Profile
	}

	data := []Person{
		{
			Name: "Alice",
			Age:  30,
			City: "Wonderland",
			Profile: Profile{
				Email:    "alice@example.com",
				Username: "alice",
			},
		},
		{
			Name: "Bob",
			Age:  25,
			City: "Builderland",
			Profile: Profile{
				Email:    "bob@example.com",
				Username: "bob",
			},
		},
	}

	result, err := MarshalFlattened(data)
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Flattened output:\n%s", resultStr)

	// Check that the header contains the flattened profile schema
	if !strings.Contains(resultStr, "profile{email,username}") && !strings.Contains(resultStr, "profile{username,email}") {
		t.Errorf("Expected flattened schema 'profile{email,username}' in header, got:\n%s", resultStr)
	}

	// Check that the values are on the same line (no > profile: nesting)
	if strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected flattened format without '> profile:', but found it in:\n%s", resultStr)
	}

	// Check that email and username values are inline
	if !strings.Contains(resultStr, "alice@example.com") {
		t.Errorf("Expected alice@example.com in output")
	}
	if !strings.Contains(resultStr, "alice") {
		t.Errorf("Expected alice username in output")
	}
}

func TestMarshalFlattenedVsNormal(t *testing.T) {
	type Profile struct {
		Email    string
		Username string
	}

	type Person struct {
		Name    string
		Age     int
		City    string
		Profile Profile
	}

	data := []Person{
		{
			Name: "Alice",
			Age:  30,
			City: "Wonderland",
			Profile: Profile{
				Email:    "alice@example.com",
				Username: "alice",
			},
		},
		{
			Name: "Bob",
			Age:  25,
			City: "Builderland",
			Profile: Profile{
				Email:    "bob@example.com",
				Username: "bob",
			},
		},
	}

	// Normal format
	normalResult, err := Marshal(data)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Flattened format
	flatResult, err := MarshalFlattened(data)
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	t.Logf("\n=== Normal Format ===\n%s", string(normalResult))
	t.Logf("\n=== Flattened Format ===\n%s", string(flatResult))

	// Flattened should be shorter
	if len(flatResult) >= len(normalResult) {
		t.Errorf("Expected flattened format to be shorter than normal format")
		t.Errorf("Normal: %d bytes, Flattened: %d bytes", len(normalResult), len(flatResult))
	}

	// Count newlines - flattened should have fewer
	normalLines := strings.Count(string(normalResult), "\n")
	flatLines := strings.Count(string(flatResult), "\n")

	if flatLines >= normalLines {
		t.Errorf("Expected flattened format to have fewer lines than normal format")
		t.Errorf("Normal: %d lines, Flattened: %d lines", normalLines, flatLines)
	}
}

func TestMarshalFlattenedNonFlattenable(t *testing.T) {
	// Test with deeply nested structures that cannot be flattened
	type Address struct {
		Street  string
		City    string
		Country string
	}

	type NestedProfile struct {
		Username string
		Address  Address // This is nested within nested, cannot fully flatten
	}

	type Person struct {
		Name    string
		Profile NestedProfile
	}

	data := []Person{
		{
			Name: "Alice",
			Profile: NestedProfile{
				Username: "alice",
				Address: Address{
					Street:  "123 Main St",
					City:    "Wonderland",
					Country: "Fantasy",
				},
			},
		},
	}

	result, err := MarshalFlattened(data)
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Non-flattenable output:\n%s", resultStr)

	// Should still produce valid output
	if len(resultStr) == 0 {
		t.Errorf("Expected non-empty output")
	}
}
