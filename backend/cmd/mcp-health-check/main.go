package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/razorpay/gatewayforge-ai/services"
)

func main() {
	fmt.Println("🔍 GatewayForge AI - Claude Desktop MCP Health Check")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println()

	// Create MCP client
	client := services.NewClaudeMCPClient()

	fmt.Printf("📍 MCP Endpoint: %s\n", client.MCPEndpoint)
	fmt.Println()

	// Test 1: Check if Claude Desktop is running
	fmt.Print("1️⃣  Checking if Claude Desktop is running... ")
	if client.IsClaudeDesktopRunning() {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
		fmt.Println()
		fmt.Println("⚠️  Claude Desktop is not running or MCP is not enabled.")
		fmt.Println()
		fmt.Println("Please:")
		fmt.Println("  1. Start Claude Desktop app")
		fmt.Println("  2. Go to Settings → Advanced → Enable Developer Mode")
		fmt.Println("  3. Go to Settings → MCP → Enable MCP Server")
		fmt.Println()
		fmt.Println("Then run this health check again.")
		os.Exit(1)
	}

	// Test 2: Test MCP connectivity
	fmt.Print("2️⃣  Testing MCP connectivity... ")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testPrompt := "Hello, this is a test from GatewayForge AI. Please respond with 'OK'."
	response, err := client.SendPrompt(ctx, testPrompt, "You are a helpful assistant.")
	if err != nil {
		fmt.Println("❌ FAIL")
		fmt.Printf("   Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ PASS")
	fmt.Printf("   Response: %s\n", truncate(response, 100))
	fmt.Println()

	// Test 3: Check if skills are accessible
	fmt.Println("3️⃣  Checking GatewayForge skills...")

	skills := []string{
		"brd-harmonizer",
		"prd-generator",
		"coding-agent",
		"test-agent",
		"deploy-agent",
	}

	skillsFound := 0
	for _, skill := range skills {
		fmt.Printf("   ├─ %s... ", skill)

		// Try to invoke skill with minimal input
		skillInput := map[string]interface{}{
			"test": true,
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := client.InvokeSkill(ctx2, skill, skillInput)
		cancel2()

		if err != nil {
			// Skill might not be loaded or requires specific input
			// This is expected for initial setup
			fmt.Println("⚠️  Not configured (see setup guide)")
		} else {
			fmt.Println("✅ Available")
			skillsFound++
		}
	}
	fmt.Println()

	if skillsFound == 0 {
		fmt.Println("ℹ️  No skills configured yet. This is normal for first-time setup.")
		fmt.Println()
		fmt.Println("To configure skills:")
		fmt.Println("  1. Open Claude Desktop")
		fmt.Println("  2. Go to Settings → MCP → Skills")
		fmt.Println("  3. Add skills directory: /path/to/gatewayforge-ai/skills")
		fmt.Println()
		fmt.Println("See docs/CLAUDE_DESKTOP_SETUP.md for detailed instructions.")
	} else {
		fmt.Printf("✅ %d/%d skills configured\n", skillsFound, len(skills))
	}

	fmt.Println()
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Println()

	if skillsFound == len(skills) {
		fmt.Println("🎉 All systems ready! You can now use GatewayForge AI with Claude Desktop.")
	} else if skillsFound > 0 {
		fmt.Println("⚠️  Some skills need configuration. See docs/CLAUDE_DESKTOP_SETUP.md")
	} else {
		fmt.Println("ℹ️  Claude Desktop is running, but skills need configuration.")
		fmt.Println("   Follow: docs/CLAUDE_DESKTOP_SETUP.md")
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  • Configure skills (if needed): See docs/CLAUDE_DESKTOP_SETUP.md")
	fmt.Println("  • Start GatewayForge: make docker-up")
	fmt.Println("  • Access frontend: http://localhost:5173")
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
