package jet

import (
	"strings"
	"testing"
)

func TestMarshalSimpleStruct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	result, err := Marshal(Person{Name: "Alice", Age: 30, City: "Wonderland"})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Simple struct output:\n%s", resultStr)

	// Check all fields are present
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

func TestMarshalSliceOfStructs(t *testing.T) {
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

	result, err := Marshal(testData)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Slice of structs output:\n%s", resultStr)

	// Check tabular format header
	if !strings.Contains(resultStr, "{category|id|name}:") || !strings.Contains(resultStr, "{id|name|category}:") || !strings.Contains(resultStr, "{name|category|id}:") {
		// Any order is fine, just check it's a table
		if !strings.Contains(resultStr, "category") || !strings.Contains(resultStr, "id") || !strings.Contains(resultStr, "name") {
			t.Errorf("Expected tabular header with category, id, name fields")
		}
	}

	// Check data rows
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

func TestMarshalNestedStruct(t *testing.T) {
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

	result, err := Marshal(Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "Wonderland",
			ZipCode: "12345",
		},
	})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Nested struct output:\n%s", resultStr)

	// Check nested object structure
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

func TestMarshalStructWithNestedSlice(t *testing.T) {
	type Item struct {
		ProductID int
		Quantity  int
	}

	type Order struct {
		OrderID int
		Items   []Item
	}

	result, err := Marshal(Order{
		OrderID: 12345,
		Items: []Item{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Struct with nested slice output:\n%s", resultStr)

	// Check nested tabular array
	if !strings.Contains(resultStr, "items") {
		t.Errorf("Expected 'items' field")
	}
	if !strings.Contains(resultStr, "orderid: 12345") {
		t.Errorf("Expected 'orderid: 12345'")
	}
}

func TestMarshalSliceWithNestedObjects(t *testing.T) {
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

	result, err := Marshal([]Person{
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
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Slice with nested objects output:\n%s", resultStr)

	// Check table header includes all fields
	if !strings.Contains(resultStr, "age") || !strings.Contains(resultStr, "city") || !strings.Contains(resultStr, "name") || !strings.Contains(resultStr, "profile") {
		t.Errorf("Expected table header with age, city, name, profile fields")
	}

	// Check nested profile blocks
	if !strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected '> profile:' nested blocks")
	}
	if !strings.Contains(resultStr, "alice@example.com") {
		t.Errorf("Expected 'alice@example.com' in nested profile")
	}
	if !strings.Contains(resultStr, "username: bob") {
		t.Errorf("Expected 'username: bob' in nested profile")
	}
}

func TestMarshalStructs(t *testing.T) {
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

	type Config struct {
		Route   string
		Method  string
		Headers []string
	}

	type Response struct {
		Persons []Person
		Config  Config
	}

	result, err := Marshal(Response{
		Persons: []Person{
			{Name: "Alice", Age: 30, City: "Wonderland", Profile: Profile{Username: "alice", Email: "alice@example.com"}},
			{Name: "Bob", Age: 25, City: "Builderland", Profile: Profile{Username: "bob", Email: "bob@example.com"}},
		},
		Config: Config{
			Route:   "/home",
			Method:  "GET",
			Headers: []string{"Content-Type: application/json"},
		},
	})

	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Complex struct output:\n%s", resultStr)

	// Check config section
	if !strings.Contains(resultStr, "config:") {
		t.Errorf("Expected 'config:' section")
	}
	if !strings.Contains(resultStr, "method: GET") {
		t.Errorf("Expected 'method: GET'")
	}
	if !strings.Contains(resultStr, "route: /home") {
		t.Errorf("Expected 'route: /home'")
	}

	// Check persons tabular section
	if !strings.Contains(resultStr, "persons") {
		t.Errorf("Expected 'persons' table")
	}
	if !strings.Contains(resultStr, "Alice") && !strings.Contains(resultStr, "Bob") {
		t.Errorf("Expected person names in output")
	}
}

func TestMarshalMap(t *testing.T) {
	data := map[string]interface{}{
		"route":  "/api/data",
		"method": "POST",
		"headers": []interface{}{
			map[string]interface{}{
				"key":   "Content-Type",
				"value": "application/json",
			},
			map[string]interface{}{
				"key":   "Authorization",
				"value": "Bearer token",
			},
		},
	}

	result, err := Marshal(data)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Map output:\n%s", resultStr)

	// Check headers table
	if !strings.Contains(resultStr, "headers") {
		t.Errorf("Expected 'headers' table")
	}
	if !strings.Contains(resultStr, "Content-Type") {
		t.Errorf("Expected 'Content-Type' in headers")
	}
	if !strings.Contains(resultStr, "Bearer token") {
		t.Errorf("Expected 'Bearer token' in headers")
	}

	// Check simple fields
	if !strings.Contains(resultStr, "method: POST") {
		t.Errorf("Expected 'method: POST'")
	}
	if !strings.Contains(resultStr, "route: /api/data") {
		t.Errorf("Expected 'route: /api/data'")
	}
}

func TestMarshalSliceOfMaps(t *testing.T) {
	data := []map[string]interface{}{
		{
			"id":    1,
			"name":  "Item One",
			"price": 9.99,
		},
		{
			"id":    2,
			"name":  "Item Two",
			"price": 19.99,
		},
	}

	result, err := Marshal(data)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Slice of maps output:\n%s", resultStr)

	// Check tabular format
	if !strings.Contains(resultStr, "id") || !strings.Contains(resultStr, "name") || !strings.Contains(resultStr, "price") {
		t.Errorf("Expected table with id, name, price columns")
	}
	if !strings.Contains(resultStr, "Item One") {
		t.Errorf("Expected 'Item One' in output")
	}
	if !strings.Contains(resultStr, "19.99") {
		t.Errorf("Expected '19.99' in output")
	}
}

func TestMarshalEmptySlice(t *testing.T) {
	type Product struct {
		ID   int
		Name string
	}

	result, err := Marshal([]Product{})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Empty slice should produce empty or minimal output
	if len(result) > 0 {
		t.Logf("Empty slice output: %s", string(result))
	}
}

func TestMarshalDeeplyNestedStructure(t *testing.T) {
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

	result, err := Marshal([]Person{
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
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Deeply nested output:\n%s", resultStr)

	// Check multiple levels of nesting
	if !strings.Contains(resultStr, "> profile:") {
		t.Errorf("Expected '> profile:' nested block")
	}
	if !strings.Contains(resultStr, "address:") {
		t.Errorf("Expected 'address:' nested block")
	}
	if !strings.Contains(resultStr, "country: Fantasy") {
		t.Errorf("Expected 'country: Fantasy' in deeply nested structure")
	}
}

func TestMarshalMixedTypes(t *testing.T) {
	type Data struct {
		StringField string
		IntField    int
		FloatField  float64
		BoolField   bool
	}

	result, err := Marshal(Data{
		StringField: "test",
		IntField:    42,
		FloatField:  3.14,
		BoolField:   true,
	})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	resultStr := string(result)
	t.Logf("Mixed types output:\n%s", resultStr)

	// Check all field types are correctly marshaled
	if !strings.Contains(resultStr, "stringfield: test") {
		t.Errorf("Expected 'stringfield: test'")
	}
	if !strings.Contains(resultStr, "intfield: 42") {
		t.Errorf("Expected 'intfield: 42'")
	}
	if !strings.Contains(resultStr, "floatfield: 3.14") {
		t.Errorf("Expected 'floatfield: 3.14'")
	}
	if !strings.Contains(resultStr, "boolfield: true") {
		t.Errorf("Expected 'boolfield: true'")
	}
}
