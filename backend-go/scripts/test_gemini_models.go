package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	// Get API key from environment or argument
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" && len(os.Args) > 1 {
		apiKey = os.Args[1]
	}

	if apiKey == "" {
		fmt.Println("Usage:")
		fmt.Println("  go run test_gemini_models.go YOUR_API_KEY")
		fmt.Println("Or set GEMINI_API_KEY environment variable")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  go run test_gemini_models.go AIzaSy...")
		os.Exit(1)
	}

	ctx := context.Background()

	// Create client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("❌ Failed to create client: %v", err)
	}
	defer client.Close()

	fmt.Println("============================================================")
	fmt.Println("GEMINI AVAILABLE MODELS")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Models available for your API key:")
	fmt.Println("------------------------------------------------------------")

	// List models
	iter := client.ListModels(ctx)
	count := 0

	for {
		model, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("❌ Error listing models: %v", err)
		}

		// Check if model supports generateContent
		supportsGenerate := false
		for _, method := range model.SupportedGenerationMethods {
			if method == "generateContent" {
				supportsGenerate = true
				break
			}
		}

		if supportsGenerate {
			count++
			fmt.Printf("✓ %s\n", model.Name)
			fmt.Printf("  Display name: %s\n", model.DisplayName)
			fmt.Printf("  Description: %s\n", model.Description)
			fmt.Println()
		}
	}

	fmt.Println("============================================================")
	fmt.Printf("Total: %d models support generateContent\n", count)
	fmt.Println("============================================================")

	if count > 0 {
		fmt.Println("\nRECOMMENDATIONS:")
		fmt.Println("- For fast, free: Use models with 'flash' in the name")
		fmt.Println("- For best quality: Use models with 'pro' in the name")
		fmt.Println("\nUpdate analysis/gemini.go line 50 with one of the model names above.")
	} else {
		fmt.Println("\n⚠️  No models found!")
		fmt.Println("\nPossible issues:")
		fmt.Println("- Invalid API key")
		fmt.Println("- Network connection problem")
		fmt.Println("- API quota exceeded")
	}
}
