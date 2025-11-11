package jet

import (
	"strings"
	"testing"
)

func TestMarshalNormalizedSimpleStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	result, err := MarshalNormalized(Person{Name: "Alice", Age: 30, City: "Wonderland"})
	if err != nil {
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized simple struct output:\n%s", resultStr)

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

func TestMarshalNormalizedSliceOfStructs(t *testing.T) {
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

	result, err := MarshalNormalized(testData)
	if err != nil {
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized slice of structs output:\n%s", resultStr)

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
}

func TestMarshalNormalizedNestedStruct(t *testing.T) {
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

	result, err := MarshalNormalized(Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "Wonderland",
			ZipCode: "12345",
		},
	})
	if err != nil {
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized nested struct output:\n%s", resultStr)

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

func TestMarshalNormalizedSliceWithNestedObjects(t *testing.T) {
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

	result, err := MarshalNormalized([]Person{
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
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized slice with nested objects output:\n%s", resultStr)

	// Check normalized schema in header
	if !strings.Contains(resultStr, "profile{email|username}") && !strings.Contains(resultStr, "profile{username|email}") {
		t.Errorf("Expected normalized schema 'profile{email|username}' or similar in header")
	}

	// Check nested profile blocks still use >
	if !strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected '> profile:' nested blocks")
	}

	// Check that nested values are pipe-delimited, not key:value format
	lines := strings.Split(resultStr, "\n")
	foundPipeDelimited := false
	for i, line := range lines {
		if strings.Contains(line, "> profile:") && i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			// Should be pipe-delimited values, not "email: ..."
			if strings.Contains(nextLine, "|") && !strings.Contains(nextLine, "email:") && !strings.Contains(nextLine, "username:") {
				foundPipeDelimited = true
				break
			}
		}
	}

	if !foundPipeDelimited {
		t.Errorf("Expected pipe-delimited values in nested profile block")
	}
}

func TestMarshalNormalizedDeeplyNestedStructure(t *testing.T) {
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

	result, err := MarshalNormalized([]Person{
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
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized deeply nested output:\n%s", resultStr)

	// Check multiple levels of nesting
	if !strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected '> profile:' nested block")
	}
	if !strings.Contains(resultStr, "address:") {
		t.Errorf("Expected 'address:' nested block")
	}

	// The deepest level should still be key:value since Address is nested within Profile
	if !strings.Contains(resultStr, "country: Fantasy") {
		t.Errorf("Expected 'country: Fantasy' in deeply nested structure")
	}
}

func TestMarshalNormalizedMixedTypes(t *testing.T) {
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

	result, err := MarshalNormalized([]Person{
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
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized mixed types output:\n%s", resultStr)

	// Check that different types are pipe-delimited in nested block
	if !strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected '> profile:' nested block")
	}

	// Should have pipe-delimited values with different types
	lines := strings.Split(resultStr, "\n")
	for i, line := range lines {
		if strings.Contains(line, "> profile:") && i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			if strings.Contains(nextLine, "|") {
				t.Logf("Found pipe-delimited profile values: %s", nextLine)
				// Should contain the values but not the keys
				if strings.Contains(nextLine, "95.5") && strings.Contains(nextLine, "alice") {
					return // Success
				}
			}
		}
	}
	t.Errorf("Expected pipe-delimited mixed type values in nested profile block")
}

func TestMarshalNormalizedComplexStructure(t *testing.T) {
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

	result, err := MarshalNormalized([]Order{
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
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized complex structure output:\n%s", resultStr)

	// Check for nested address with pipe-delimited values
	if !strings.Contains(resultStr, "address{") {
		t.Errorf("Expected normalized address schema in header")
	}

	// Check for nested items table
	if !strings.Contains(resultStr, "items") {
		t.Errorf("Expected 'items' nested table")
	}
}

func TestMarshalNormalized(t *testing.T) {
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

	result, err := MarshalNormalized(data)
	if err != nil {
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Normalized output:\n%s", resultStr)

	// Check that the header contains the normalized profile schema with pipe delimiter
	if !strings.Contains(resultStr, "profile{email|username}") && !strings.Contains(resultStr, "profile{username|email}") {
		t.Errorf("Expected normalized schema 'profile{email|username}' in header, got:\n%s", resultStr)
	}

	// Check that the nested block still uses > profile: but with pipe-delimited values
	if !strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected '> profile:' nesting in normalized format, got:\n%s", resultStr)
	}

	// Check that the nested values are pipe-delimited (not key: value format)
	lines := strings.Split(resultStr, "\n")
	foundPipeDelimitedNested := false
	for i, line := range lines {
		if strings.Contains(line, "> profile:") && i+1 < len(lines) {
			nextLine := lines[i+1]
			// The next line should contain pipe-delimited values, not "email: alice@example.com"
			if strings.Contains(nextLine, "|") && !strings.Contains(nextLine, "email:") {
				foundPipeDelimitedNested = true
				break
			}
		}
	}

	if !foundPipeDelimitedNested {
		t.Errorf("Expected pipe-delimited values in nested block, got:\n%s", resultStr)
	}
}

func TestMarshalNormalizedVsFlattened(t *testing.T) {
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

	// Normalized format
	normalizedResult, err := MarshalNormalized(data)
	if err != nil {
		t.Fatalf("MarshalNormalized failed: %v", err)
	}

	// Flattened format
	flattenedResult, err := MarshalFlattened(data)
	if err != nil {
		t.Fatalf("MarshalFlattened failed: %v", err)
	}

	t.Logf("\n=== Normal Format ===\n%s", string(normalResult))
	t.Logf("\n=== Normalized Format ===\n%s", string(normalizedResult))
	t.Logf("\n=== Flattened Format ===\n%s", string(flattenedResult))

	// Normalized should be between normal and flattened in size
	t.Logf("\nSizes: Normal=%d, Normalized=%d, Flattened=%d bytes",
		len(normalResult), len(normalizedResult), len(flattenedResult))

	// Flattened should be smallest
	if len(flattenedResult) >= len(normalizedResult) {
		t.Logf("Note: Flattened (%d) should typically be smaller than Normalized (%d)",
			len(flattenedResult), len(normalizedResult))
	}

	// Normalized should be smaller than normal (fewer lines due to pipe-delimited nested values)
	if len(normalizedResult) >= len(normalResult) {
		t.Logf("Note: Normalized (%d) should typically be smaller than Normal (%d)",
			len(normalizedResult), len(normalResult))
	}
}
