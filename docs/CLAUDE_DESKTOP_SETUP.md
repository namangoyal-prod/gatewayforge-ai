# Claude Desktop Integration Setup

This guide shows you how to use GatewayForge AI with **Claude Desktop** instead of API keys. You'll authenticate through your Claude account and use your existing Claude subscription.

---

## 🎯 Benefits of Claude Desktop Integration

✅ **No API Key Required** - Use your existing Claude login
✅ **Uses Your Subscription** - Claude Pro or Enterprise
✅ **Familiar Interface** - Same Claude you already use
✅ **MCP Protocol** - Modern, secure communication
✅ **Cost Effective** - No separate API billing

---

## 📋 Prerequisites

1. **Claude Desktop App** (Free or Pro account)
2. **MCP Enabled** in Claude Desktop settings
3. **GatewayForge AI** platform installed

---

## 🚀 Setup Steps

### Step 1: Install Claude Desktop

If you don't have Claude Desktop installed:

**Download:**
- Visit: https://claude.ai/download
- Download for macOS
- Install the app
- Sign in with your Anthropic account

**Or via Homebrew:**
```bash
# Not yet available via Homebrew - use website download
```

### Step 2: Enable MCP in Claude Desktop

1. **Open Claude Desktop**
2. **Go to Settings** (⌘ + ,)
3. **Enable Developer Mode**:
   - Click on "Advanced"
   - Toggle "Enable Developer Mode"
4. **Enable MCP Server**:
   - Go to "MCP" tab
   - Toggle "Enable MCP Server"
   - Note the port (default: 52828)
5. **Configure GatewayForge Skills**:
   - Click "Add Skills"
   - Navigate to your GatewayForge installation
   - Select the `skills/` directory
   - Claude will load all 5 skills:
     - BRD Harmonizer
     - PRD Generator
     - Coding Agent
     - Test Agent
     - Deploy Agent

### Step 3: Configure GatewayForge

Update your `.env` file:

```bash
# Edit .env
nano .env

# Add/Update these lines:
CLAUDE_MCP_ENDPOINT=http://localhost:52828
CLAUDE_AUTH_MODE=desktop
CLAUDE_API_KEY=not_required_with_mcp

# Remove or comment out:
# CLAUDE_API_KEY=your_api_key_here
```

### Step 4: Verify Connection

Test the MCP connection:

```bash
# From gatewayforge-ai directory
cd backend/services

# Run health check
go run ../cmd/mcp-health-check/main.go
```

**Expected Output:**
```
✅ Claude Desktop detected
✅ MCP Server accessible at http://localhost:52828
✅ Skills loaded: 5/5
   - brd-harmonizer
   - prd-generator
   - coding-agent
   - test-agent
   - deploy-agent
✅ Authentication: Active (logged in as: your.email@domain.com)

🎉 Claude Desktop MCP is ready!
```

### Step 5: Start GatewayForge

Now start the platform as usual:

```bash
# Option A: With Docker
make docker-up

# Option B: Individual services
make dev-backend
make dev-frontend
```

---

## 🔧 Configuring Claude Desktop Skills

### Manual Skill Installation

If auto-detection doesn't work, manually add skills:

1. **Locate Skills Directory**:
   ```bash
   pwd
   # Should show: /Users/naman.goyal/Documents/vault/gatewayforge-ai

   ls skills/
   # Shows: brd-harmonizer  coding-agent  deploy-agent  prd-generator  test-agent
   ```

2. **Add to Claude Desktop**:
   - Open Claude Desktop
   - Settings → MCP → Skills
   - Click "Add Skill Directory"
   - Select `/path/to/gatewayforge-ai/skills`
   - Click "Load Skills"

3. **Verify Skills Loaded**:
   - You should see all 5 skills in Claude Desktop
   - Each skill will have a description and parameters

### Skill Configuration File

Claude Desktop uses a skills configuration file. Create it:

```bash
# Create MCP skills config for Claude Desktop
cat > ~/.claude/skills.json << 'EOF'
{
  "skills": [
    {
      "name": "brd-harmonizer",
      "path": "/Users/naman.goyal/Documents/vault/gatewayforge-ai/skills/brd-harmonizer",
      "enabled": true,
      "description": "Validates BRD documents against Razorpay standards"
    },
    {
      "name": "prd-generator",
      "path": "/Users/naman.goyal/Documents/vault/gatewayforge-ai/skills/prd-generator",
      "enabled": true,
      "description": "Generates production-grade PRDs from BRDs"
    },
    {
      "name": "coding-agent",
      "path": "/Users/naman.goyal/Documents/vault/gatewayforge-ai/skills/coding-agent",
      "enabled": true,
      "description": "Generates production-ready code across repositories"
    },
    {
      "name": "test-agent",
      "path": "/Users/naman.goyal/Documents/vault/gatewayforge-ai/skills/test-agent",
      "enabled": true,
      "description": "Generates comprehensive test suites"
    },
    {
      "name": "deploy-agent",
      "path": "/Users/naman.goyal/Documents/vault/gatewayforge-ai/skills/deploy-agent",
      "enabled": true,
      "description": "Automates dev stack deployment"
    }
  ]
}
EOF
```

---

## 🧪 Testing the Integration

### Test 1: Health Check

```bash
# Test MCP connectivity
curl http://localhost:52828/health

# Expected: {"status": "ok", "mcp_version": "1.0"}
```

### Test 2: Invoke a Skill

```bash
# Test BRD validation skill
curl -X POST http://localhost:52828/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "brd-harmonizer",
      "arguments": {
        "brd_content": "Sample BRD for testing"
      }
    },
    "id": 1
  }'
```

### Test 3: Upload a BRD in GatewayForge

1. Open http://localhost:5173
2. Click "New Integration"
3. Upload a sample BRD
4. Watch Claude Desktop window - you should see skill execution logs
5. Validation results will appear in GatewayForge UI

---

## 🔍 Troubleshooting

### Issue 1: "Claude Desktop is not running"

**Solution:**
```bash
# Check if Claude Desktop is running
ps aux | grep -i claude

# If not running, start it
open -a "Claude"

# Wait for it to fully start (look for icon in menu bar)
```

### Issue 2: "MCP Server not accessible"

**Solution:**
```bash
# Check MCP endpoint
curl http://localhost:52828/health

# If connection refused, check Claude Desktop settings:
# 1. Open Claude Desktop
# 2. Settings → Advanced → Enable Developer Mode
# 3. Settings → MCP → Enable MCP Server
# 4. Note the port number
# 5. Update CLAUDE_MCP_ENDPOINT in .env
```

### Issue 3: "Skills not loaded"

**Solution:**
```bash
# Verify skills directory
ls ~/Documents/vault/gatewayforge-ai/skills/

# Each skill should have SKILL.md file
ls skills/*/SKILL.md

# Restart Claude Desktop
killall Claude
open -a "Claude"

# Re-add skills directory in Settings
```

### Issue 4: "Authentication failed"

**Solution:**
```bash
# Check if logged into Claude Desktop
# Open Claude Desktop → should see chat interface
# If not logged in, click "Sign In" and authenticate

# Verify authentication
curl http://localhost:52828/auth/status
```

### Issue 5: MCP Port Conflict

If port 52828 is in use:

```bash
# Find what's using the port
lsof -i :52828

# Change Claude Desktop MCP port:
# Settings → MCP → Port: 52829 (or any available port)

# Update .env:
CLAUDE_MCP_ENDPOINT=http://localhost:52829
```

---

## 📊 How It Works

```
┌─────────────────────────────────────────────────────────────┐
│                  GatewayForge AI Flow                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  User uploads BRD                                           │
│         ↓                                                   │
│  GatewayForge Backend                                       │
│         ↓                                                   │
│  MCP Client (activities_mcp.go)                             │
│         ↓                                                   │
│  HTTP Request to localhost:52828                            │
│         ↓                                                   │
│  Claude Desktop MCP Server                                  │
│         ↓                                                   │
│  Claude (with your authentication)                          │
│         ↓                                                   │
│  Executes BRD Harmonizer Skill                              │
│         ↓                                                   │
│  Returns validation results                                 │
│         ↓                                                   │
│  GatewayForge displays results                              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### MCP Protocol

GatewayForge uses JSON-RPC 2.0 over HTTP to communicate with Claude Desktop:

```javascript
// Request
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "brd-harmonizer",
    "arguments": { "brd_content": "..." }
  },
  "id": 1
}

// Response
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "quality_score": 85,
    "status": "GREEN",
    "gap_analysis": [...]
  }
}
```

---

## 🔐 Security Considerations

1. **Local Communication Only**: MCP server only listens on localhost by default
2. **Authentication**: Uses your Claude Desktop login session
3. **No API Keys Exposed**: No keys stored in code or environment variables
4. **Audit Trail**: All skill executions logged in Claude Desktop
5. **Permissions**: Skills can only access what Claude Desktop allows

---

## 💰 Cost Comparison

| Method | Cost | Limits |
|--------|------|--------|
| **API Key** | $15-30 per 1M tokens | Pay per use |
| **Claude Desktop (Free)** | $0 | Rate limited |
| **Claude Desktop (Pro - $20/mo)** | $20/month | 5x higher limits |
| **Claude Desktop (Enterprise)** | Custom | Unlimited |

**For GatewayForge:**
- Avg integration uses ~200K tokens
- API method: ~$3-6 per integration
- Desktop method: Included in your subscription

---

## ✅ Verification Checklist

Before using GatewayForge with Claude Desktop, verify:

- [ ] Claude Desktop installed and running
- [ ] Logged into your Claude account
- [ ] Developer Mode enabled
- [ ] MCP Server enabled and running
- [ ] Port 52828 accessible (or custom port configured)
- [ ] Skills directory added and loaded
- [ ] `.env` configured with `CLAUDE_MCP_ENDPOINT`
- [ ] Health check passes
- [ ] Can see skills in Claude Desktop

---

## 🎉 You're Ready!

Once all checks pass, you can:

1. **Start GatewayForge**: `make docker-up`
2. **Open Frontend**: http://localhost:5173
3. **Upload BRDs**: Skills will execute via Claude Desktop
4. **Monitor Execution**: Watch Claude Desktop for real-time skill execution
5. **View Results**: See validation, PRD generation, code generation results

**No API key needed! Uses your Claude Desktop authentication!** 🚀

---

## 📞 Support

- **Claude Desktop Issues**: https://claude.ai/help
- **MCP Protocol**: https://modelcontextprotocol.io
- **GatewayForge Issues**: See README.md

---

**Happy Building!** 🎨
