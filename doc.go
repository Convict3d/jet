// Package jet provides a compact, human-readable data serialization format
// optimized for LLM token efficiency.
//
// Jet is a lightweight alternative to JSON/YAML that uses tabular representation
// to dramatically reduce token count and byte size while maintaining readability.
// It's perfect for LLM applications where token efficiency is critical.
//
// # Features
//
//   - Token Efficient: Up to 36% fewer tokens compared to JSON for nested data
//   - Human Readable: Clear, tabular format that's easy to understand
//   - Two Format Modes: Normal and Normalized to suit different use cases
//   - Type Safe: Full Go struct support with reflection-based marshaling
//   - Zero Dependencies: Pure Go implementation
//
// # Basic Usage
//
// Marshal Go structs to Jet format:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	    City string
//	}
//
//	data := []Person{
//	    {Name: "Alice", Age: 30, City: "Wonderland"},
//	    {Name: "Bob", Age: 25, City: "Builderland"},
//	}
//
//	// Normal format - readable with nested blocks
//	result, err := jet.Marshal(data)
//
//	// Normalized format - pipe-delimited nested values
//	normalized, err := jet.MarshalNormalized(data)
//
// # Format Modes
//
// Normal Format - Best for readability:
//
//	persons{age|city|name|profile}:
//	  30|Wonderland|Alice
//	    > profile:
//	    email: alice@example.com
//	    username: alice
//
// Normalized Format - Best for balance between readability and compression:
//
//	persons{age|city|name|profile{email|username}}:
//	  30|Wonderland|Alice
//	    > profile:
//	    alice@example.com|alice
//
// # Struct Tags
//
// Use jet tags to customize field names:
//
//	type User struct {
//	    Name     string `jet:"username"`
//	    Email    string `jet:"email"`
//	    Internal string `jet:"-"`  // Skip this field
//	}
//
// # Use Cases
//
//   - LLM Tool Responses: Minimize token usage in function calling
//   - Multi-Agent Systems: Efficient data exchange between agents
//   - API Responses: Compact data transfer
//   - Configuration Files: Human-readable with minimal footprint
//   - Data Pipelines: Efficient serialization for data processing
package jet
