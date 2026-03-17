# Quick Start: GatewayForge AI with Claude Desktop

Use GatewayForge AI with your Claude Desktop login - **no API key required!**

---

## ⚡ Super Quick Setup (5 minutes)

### 1. Install Claude Desktop

Download and install Claude Desktop:
```bash
# Visit and download:
open https://claude.ai/download

# Or if already installed, just open it:
open -a "Claude"
```

### 2. Enable MCP in Claude Desktop

1. Open Claude Desktop
2. Press `⌘ + ,` (Settings)
3. Go to **Advanced** → Toggle **"Enable Developer Mode"**
4. Go to **MCP** → Toggle **"Enable MCP Server"**
5. Note the port (usually 52828)

### 3. Configure GatewayForge

Already done! Your `.env` file is configured to use Claude Desktop:

```bash
cat .env | grep CLAUDE
# Should show:
# CLAUDE_AUTH_MODE=desktop
# CLAUDE_MCP_ENDPOINT=http://localhost:52828
```

### 4. Check Connection

```bash
make mcp-check
```

**Expected output:**
```
🔍 GatewayForge AI - Claude Desktop MCP Health Check
==================================================

📍 MCP Endpoint: http://localhost:52828

1️⃣  Checking if Claude Desktop is running... ✅ PASS
2️⃣  Testing MCP connectivity... ✅ PASS
   Response: OK

🎉 All systems ready! You can now use GatewayForge AI with Claude Desktop.
```

### 5. Start GatewayForge

```bash
# With Docker (recommended)
make docker-up

# Access:
# Frontend: http://localhost:5173
# Backend: http://localhost:8080
```

---

## 🎯 What You Get

✅ **No API Key Needed** - Uses your Claude Desktop login
✅ **Free Usage** - Included in Claude Free/Pro subscription
✅ **Same Claude** - Familiar interface and capabilities
✅ **Secure** - All communication stays local (localhost)
✅ **Easy Setup** - Just enable MCP in settings

---

## 🔧 Troubleshooting

### Problem: "Claude Desktop is not running"

**Solution:**
```bash
# Check if running:
ps aux | grep -i claude

# If not, start it:
open -a "Claude"

# Wait for app to fully start (icon in menu bar)
```

### Problem: "MCP Server not accessible"

**Solution:**
1. Open Claude Desktop
2. Settings (⌘ + ,)
3. Advanced → Enable Developer Mode ✅
4. MCP → Enable MCP Server ✅
5. Check port matches .env (default: 52828)

### Problem: Connection refused

**Solution:**
```bash
# Test if MCP server is listening
curl http://localhost:52828/health

# If error, restart Claude Desktop:
killall Claude
open -a "Claude"

# Wait 10 seconds, then test again
sleep 10
curl http://localhost:52828/health
```

---

## 📊 How It Works

```
You → GatewayForge Frontend → Backend → MCP Client
                                           ↓
                                  Claude Desktop MCP Server
                                           ↓
                                  Claude (with your auth)
                                           ↓
                                   Skills execute (BRD, PRD, Code)
                                           ↓
                                   Results returned
```

**All communication is local (localhost) - secure and private!**

---

## 💡 Tips

1. **Keep Claude Desktop Running**: GatewayForge needs it for AI features
2. **Check Icon**: Look for Claude icon in menu bar (should be active)
3. **Watch Logs**: Claude Desktop shows skill execution in real-time
4. **Test First**: Run `make mcp-check` before uploading real BRDs

---

## 🆚 Claude Desktop vs API Key

| Feature | Claude Desktop | API Key |
|---------|---------------|---------|
| **Cost** | Included in subscription | Pay per use ($15-30/1M tokens) |
| **Setup** | Enable MCP in settings | Sign up, create key |
| **Auth** | Your Claude login | API key |
| **Limits** | Subscription limits | Pay-as-you-go |
| **Best For** | Individual use, testing | Production, high volume |

**Recommendation:** Start with Claude Desktop (easier), upgrade to API key if needed for scale.

---

## ✅ Verification

Before using GatewayForge, verify:

```bash
# 1. Claude Desktop running
ps aux | grep -i claude

# 2. MCP endpoint accessible
curl http://localhost:52828/health

# 3. Health check passes
make mcp-check

# 4. Environment configured
cat .env | grep CLAUDE_AUTH_MODE
# Should show: CLAUDE_AUTH_MODE=desktop
```

---

## 🎉 You're Ready!

Once all checks pass:

1. **Start Platform**: `make docker-up`
2. **Open UI**: http://localhost:5173
3. **Upload BRD**: Click "New Integration"
4. **Watch Magic**: Skills execute via Claude Desktop
5. **View Results**: See validation, PRD, code generation

**No API key needed - using your Claude Desktop authentication!** 🚀

---

## 📚 More Details

For detailed setup instructions, see:
- **Full Guide**: [docs/CLAUDE_DESKTOP_SETUP.md](docs/CLAUDE_DESKTOP_SETUP.md)
- **Main README**: [README.md](README.md)
- **Implementation**: [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)

---

**Questions?**
- Run `make mcp-setup` to open full setup guide
- Run `make mcp-check` to diagnose issues
- Check `docs/CLAUDE_DESKTOP_SETUP.md` for troubleshooting

**Happy building!** 🎨
