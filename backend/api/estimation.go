package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ── Request ────────────────────────────────────────────────────────────────────

type EstimationRequest struct {
	BankName        string   `json:"bank_name"`
	IntegrationType string   `json:"integration_type"` // new_gateway | new_payment_method | extension
	PaymentMethods  []string `json:"payment_methods"`  // cards | upi | wallet | bnpl

	EncryptionType string `json:"encryption_type"` // none | aes_cbc | rsa | dukpt
	OnboardingType string `json:"onboarding_type"` // file_based | api_based | both

	HasDevices  bool `json:"has_devices"`
	HasMobileApp bool `json:"has_mobile_app"`
	HasNPCICert bool `json:"has_npci_cert"`

	// Gap 1 — Codebase reuse
	// "none"      = greenfield, no existing infra (old behaviour)
	// "partial"   = some patterns exist (e.g. Montran already integrated)
	// "migration" = significant existing infra (e.g. IDFC Omni migration)
	//   Calibration: Likith's actual 97d vs tool's greenfield 240d → 0.40x ratio
	CodebaseReuse string `json:"codebase_reuse"`

	// Gap 2 — Buffer
	// Percentage of net work days to add as adhoc/oncall/leaves buffer (0–40)
	BufferPct float64 `json:"buffer_pct"`

	// Gap 3 — Phase-based estimation
	// Phase 2 reuses Phase 1 patterns at ~55% of the original effort
	NumPhases            int      `json:"num_phases"`             // 1 or 2
	Phase2PaymentMethods []string `json:"phase2_payment_methods"` // if different from Phase 1
}

// ── Data types ─────────────────────────────────────────────────────────────────

type MilestoneEstimate struct {
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	TraditionalDays float64  `json:"traditional_days"`
	AIDays          float64  `json:"ai_days"`
	SavingsDays     float64  `json:"savings_days"`
	SavingsPct      float64  `json:"savings_pct"`
	AISkillsUsed    []string `json:"ai_skills_used"`
	Notes           string   `json:"notes,omitempty"`
	Phase           int      `json:"phase"` // 1 | 2 | 0 (buffer)
}

type EstimationResponse struct {
	BankName           string `json:"bank_name"`
	IntegrationType    string `json:"integration_type"`
	CodebaseReuseLabel string `json:"codebase_reuse_label"`

	Milestones []MilestoneEstimate `json:"milestones"`

	// Phase subtotals (net work, excl. buffer)
	Phase1Traditional float64 `json:"phase1_traditional_days"`
	Phase1AI          float64 `json:"phase1_ai_days"`
	Phase2Traditional float64 `json:"phase2_traditional_days"`
	Phase2AI          float64 `json:"phase2_ai_days"`

	// Gap 2: Buffer breakdown
	BufferTraditional float64 `json:"buffer_traditional_days"`
	BufferAI          float64 `json:"buffer_ai_days"`

	// Grand totals
	TotalTraditional float64 `json:"total_traditional_days"`
	TotalAI          float64 `json:"total_ai_days"`
	TotalSavings     float64 `json:"total_savings_days"`
	TotalSavingsPct  float64 `json:"total_savings_pct"`

	References  []ReferenceCalib `json:"references"`
	Assumptions []string         `json:"assumptions"`
}

type ReferenceCalib struct {
	Name            string  `json:"name"`
	TraditionalDays float64 `json:"traditional_days"`
	AIDays          float64 `json:"ai_days"`
	SavingsPct      float64 `json:"savings_pct"`
	IntegrationType string  `json:"integration_type"`
	Notes           string  `json:"notes"`
}

// ── Handler ────────────────────────────────────────────────────────────────────

func estimateIntegration(c *gin.Context) {
	var req EstimationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Defaults
	if req.IntegrationType == "" {
		req.IntegrationType = "new_gateway"
	}
	if req.EncryptionType == "" {
		req.EncryptionType = "none"
	}
	if req.OnboardingType == "" {
		req.OnboardingType = "api_based"
	}
	if req.CodebaseReuse == "" {
		req.CodebaseReuse = "none"
	}
	if req.NumPhases < 1 {
		req.NumPhases = 1
	}
	if req.NumPhases > 2 {
		req.NumPhases = 2
	}
	if req.BufferPct < 0 {
		req.BufferPct = 0
	}
	if req.BufferPct > 40 {
		req.BufferPct = 40
	}

	// ── Gap 1: resolve codebase reuse scales ──────────────────────────────────
	// reuseScale: deflates *traditional* days (existing infra = less greenfield work)
	// aiScale:    further deflates *AI* days (existing patterns = AI is even more effective)
	reuseScale, aiScale, reuseLabel := resolveReuse(req.CodebaseReuse)

	// ── Phase 1 ───────────────────────────────────────────────────────────────
	phase1 := buildPhase(req, reuseScale, aiScale, 1.0, 1)
	var p1Trad, p1AI float64
	for _, m := range phase1 {
		p1Trad += m.TraditionalDays
		p1AI += m.AIDays
	}

	// ── Gap 3: Phase 2 (pattern reuse = 55% of Phase 1 scope) ────────────────
	var phase2 []MilestoneEstimate
	var p2Trad, p2AI float64
	if req.NumPhases == 2 {
		p2Req := req
		if len(req.Phase2PaymentMethods) > 0 {
			p2Req.PaymentMethods = req.Phase2PaymentMethods
		}
		// Phase 2 never re-does devices, dashboard, mobile app — already in Phase 1
		p2Req.HasDevices = false
		p2Req.HasMobileApp = false
		p2Req.HasNPCICert = false
		// AI is 8% more effective in Phase 2: engineer already knows the codebase deeply
		const p2Scale = 0.55
		const p2AIBoost = 0.92
		phase2 = buildPhase(p2Req, reuseScale, aiScale*p2AIBoost, p2Scale, 2)
		for _, m := range phase2 {
			p2Trad += m.TraditionalDays
			p2AI += m.AIDays
		}
	}

	// ── Gap 2: Buffer ─────────────────────────────────────────────────────────
	// Computed on net Razorpay work (excludes bank UAT which RZP can't control)
	var bufTrad, bufAI float64
	var bufMilestones []MilestoneEstimate
	if req.BufferPct > 0 {
		netWork := p1Trad + p2Trad
		bufTrad = r1(netWork * req.BufferPct / 100)
		// AI reduces buffer by 40%: auto-generated tests catch regressions earlier,
		// parallel code review cuts wait time, devstack automation reduces env debug loops.
		bufAI = r1(bufTrad * 0.60)
		bufSav := r1(bufTrad - bufAI)
		bufPct := 0.0
		if bufTrad > 0 {
			bufPct = r1((bufSav / bufTrad) * 100)
		}
		bufMilestones = []MilestoneEstimate{{
			Name:     fmt.Sprintf("Adhoc / leaves / oncall buffer (%.0f%% of net work)", req.BufferPct),
			Category: "Buffer",
			TraditionalDays: bufTrad,
			AIDays:          bufAI,
			SavingsDays:     bufSav,
			SavingsPct:      bufPct,
			AISkillsUsed:    []string{"slit-generator-v2:test-executor", "code-review:code-reviewer"},
			Notes: "AI reduces rework loops: auto-generated SLITs catch regressions earlier, " +
				"parallel code reviewers cut PR wait time, devstack automation reduces env debug cycles",
			Phase: 0,
		}}
	}

	// ── Totals ────────────────────────────────────────────────────────────────
	all := append(append(phase1, phase2...), bufMilestones...)
	totalTrad := r1(p1Trad + p2Trad + bufTrad)
	totalAI := r1(p1AI + p2AI + bufAI)
	totalSav := r1(totalTrad - totalAI)
	totalSavPct := 0.0
	if totalTrad > 0 {
		totalSavPct = r1((totalSav / totalTrad) * 100)
	}

	c.JSON(http.StatusOK, EstimationResponse{
		BankName:           req.BankName,
		IntegrationType:    req.IntegrationType,
		CodebaseReuseLabel: reuseLabel,
		Milestones:         all,
		Phase1Traditional:  r1(p1Trad),
		Phase1AI:           r1(p1AI),
		Phase2Traditional:  r1(p2Trad),
		Phase2AI:           r1(p2AI),
		BufferTraditional:  bufTrad,
		BufferAI:           bufAI,
		TotalTraditional:   totalTrad,
		TotalAI:            totalAI,
		TotalSavings:       totalSav,
		TotalSavingsPct:    totalSavPct,
		References:         referenceCalibrations(),
		Assumptions:        buildAssumptions(req),
	})
}

// ── Gap 1: Codebase reuse ──────────────────────────────────────────────────────

// resolveReuse returns (reuseScale, aiScale, label).
//
// reuseScale deflates traditional days:
//   none      → 1.00  (full greenfield estimate)
//   partial   → 0.60  (some existing patterns, ~40% of work already done)
//   migration → 0.40  (significant existing infra; calibrated: Likith 97d vs tool 240d = 0.40x)
//
// aiScale deflates AI days on top (lower = AI is more effective):
//   none      → 1.00  (standard AI reduction factors)
//   partial   → 0.94  (AI reads existing patterns, 6% extra benefit)
//   migration → 0.88  (AI navigates mature codebase, 12% extra benefit)
func resolveReuse(reuse string) (reuseScale, aiScale float64, label string) {
	switch reuse {
	case "partial":
		return 0.60, 0.94, "Partial reuse — some patterns exist (e.g. Montran already integrated)"
	case "migration":
		return 0.40, 0.88, "Migration / adaptation — significant existing infra (calibrated: IDFC Omni actual = 97d, tool greenfield = 240d)"
	default:
		return 1.0, 1.0, "Greenfield — no existing infrastructure"
	}
}

// ── Milestone builder ──────────────────────────────────────────────────────────

// r1 rounds to one decimal place.
func r1(f float64) float64 { return math.Round(f*10) / 10 }

// milestone creates one row applying combined scale and aiScale.
// baseTrad:   greenfield traditional days (before any scaling)
// aiFactor:   fraction of trad that AI takes (0.20 = 80% reduction, 0.46 = 54%, etc.)
// scale:      reuseScale * phase2Scale
// aiScale:    from resolveReuse (further reduces AI days)
func milestone(name, category string, baseTrad, aiFactor float64,
	skills []string, notes string, phase int, scale, aiScale float64) MilestoneEstimate {

	trad := r1(baseTrad * scale)
	ai := r1(baseTrad * scale * aiFactor * aiScale)
	// Floor: AI cannot exceed 95% reduction on any single task
	if ai < trad*0.05 {
		ai = r1(trad * 0.05)
	}
	sav := r1(trad - ai)
	pct := 0.0
	if trad > 0 {
		pct = r1((sav / trad) * 100)
	}
	return MilestoneEstimate{
		Name: name, Category: category,
		TraditionalDays: trad, AIDays: ai,
		SavingsDays: sav, SavingsPct: pct,
		AISkillsUsed: skills, Notes: notes, Phase: phase,
	}
}

// fixedMs creates a milestone with no AI reduction (bank UAT, prod release, etc.)
func fixedMs(name, category string, days float64, notes string, phase int) MilestoneEstimate {
	return MilestoneEstimate{
		Name: name, Category: category,
		TraditionalDays: days, AIDays: days,
		SavingsDays: 0, SavingsPct: 0,
		AISkillsUsed: []string{},
		Notes:        notes,
		Phase:        phase,
	}
}

func hasMethod(methods []string, m string) bool {
	for _, v := range methods {
		if v == m {
			return true
		}
	}
	return false
}

// ── Gap 3: Phase builder ───────────────────────────────────────────────────────

// buildPhase generates milestones for one phase.
//
// scale       = reuseScale × phase2Scale  (phase2Scale=1.0 for Phase 1, 0.55 for Phase 2)
// aiScale     = from resolveReuse, optionally boosted for Phase 2
//
// AI reduction factors (aiFactor values) calibrated from:
//   - IDFC Bank tech spec (Feb 2026): 14-step onboarding, AES-256-CBC, H2H middleware
//   - IDFC Omni migration xlsx: Likith Reddy's granular task-level actuals
//   - Axis Bank Fastag H2H: Rough Sample sheet in Gateway_Integration_CheckList.xlsx
//
//   Speccing          0.20 → 80% reduction  (doc-coauthoring generates spec draft in hours)
//   Code gen (new)    0.46 → 54% reduction  (backend-engineer:work)
//   Code gen (reuse)  0.40 → 60% reduction  (existing patterns + AI = faster)
//   Crypto impl       0.35 → 65% reduction  (AES-CBC/DUKPT patterns pre-built)
//   SLIT tests        0.30 → 70% reduction  (slit-generator-v2 auto-generates from KB)
//   PR review         0.25 → 75% reduction  (parallel code-review:*-reviewer skills)
//   Dashboard/UI      0.35 → 65% reduction  (frontend-design + blade-reviewer)
//   Hardware testing  0.60 → 40% reduction  (physical loop limits automation)
//   Fixed steps       1.00 → 0%  reduction  (bank UAT, prod release)
func buildPhase(req EstimationRequest, reuseScale, aiScale, phase2Scale float64, phase int) []MilestoneEstimate {
	var out []MilestoneEstimate
	sc := reuseScale * phase2Scale // combined scale for this phase

	hasCards := hasMethod(req.PaymentMethods, "cards")
	hasUPI := hasMethod(req.PaymentMethods, "upi")

	payCx := 1.0
	if hasCards && hasUPI {
		payCx = 1.35
	} else if hasUPI {
		payCx = 0.85
	}

	intSc := 1.0
	switch req.IntegrationType {
	case "new_payment_method":
		intSc = 0.60
	case "extension":
		intSc = 0.35
	}

	// ── GATEWAY INTEGRATION ───────────────────────────────────────────────────

	specNote := "AI drafts full tech spec from BRD in hours; bank alignment and review still required"
	if phase == 2 {
		specNote = "Phase 2 spec is delta-only — Phase 1 patterns are reused; primarily covers new gateway-specific differences"
	}
	out = append(out, milestone("Tech Spec & Solutioning", "Gateway Integration",
		7*intSc, 0.20,
		[]string{"doc-coauthoring", "swe-repo-builder", "backend-engineer:brainstorming"},
		specNote, phase, sc, aiScale))

	if hasUPI {
		out = append(out, milestone("Integrations UPI changes", "Gateway Integration",
			15*intSc, 0.47,
			[]string{"backend-engineer:workflows:work", "slit-generator-v2:create-slits", "code-review:code-reviewer"},
			"DQR rearch, Montran callbacks, new PGR gateway, offline use case via Mozart/integrations-go",
			phase, sc, aiScale))
	}

	if hasCards {
		out = append(out, milestone("Terminal config, gateway constants, DUKPT flags", "Gateway Integration",
			10*intSc, 0.40,
			[]string{"backend-engineer:workflows:work", "rzp-discover:payments-processing-platform"},
			"Scrooge void config, DUKPT enable, pricing plan defaults, key exchange skip",
			phase, sc, aiScale))

		out = append(out, milestone("Payment Processing — cards (H2H/middleware)", "Gateway Integration",
			35*intSc*payCx, 0.46,
			[]string{
				"backend-engineer:workflows:plan", "backend-engineer:workflows:work",
				"slit-generator-v2:create-slits", "rzp-discover:payments-processing-platform",
			},
			"New gateway constant in PCP, mware Omni translation, H2H integration per TTM spec v1.32",
			phase, sc, aiScale))
	}

	if hasUPI {
		out = append(out, milestone("Payment Processing — UPI (SQR/DQR/offline flows)", "Gateway Integration",
			20*intSc, 0.46,
			[]string{"backend-engineer:workflows:work", "slit-generator-v2:create-slits", "rzp-discover:upi-merchant-acquiring"},
			"SQR changes in API service, gateway reuse for offline, DQR gateway changes in payments-upi",
			phase, sc, aiScale))
	}

	out = append(out, milestone("PR reviews (all services)", "Gateway Integration",
		4, 0.25,
		[]string{"code-review:code-reviewer", "code-review:security-reviewer", "code-review:performance-reviewer"},
		"Parallel AI reviewers run simultaneously across PCP, Middleware, integrations-go, payments-upi",
		phase, sc, aiScale))

	out = append(out, milestone("Dev Testing + devstack", "Gateway Integration",
		10, 0.40,
		[]string{"backend-engineer:workflows:devstack-test", "slit-generator-v2:test-executor"},
		"Devstack pod spin-up automation + AI-generated test harness",
		phase, sc, aiScale))

	out = append(out, milestone("ITF / SLIT Tests", "Gateway Integration",
		7, 0.30,
		[]string{"slit-generator-v2:create-slits", "slit-generator-v2:test-executor", "slit-generator-v2:test-grader"},
		"SLIT generator produces comprehensive integration tests from knowledge base; grader validates quality",
		phase, sc, aiScale))

	out = append(out, milestone("Alerting & observability setup", "Gateway Integration",
		3, 0.40,
		[]string{"backend-engineer:workflows:work", "rzp-discover:devops"},
		"Alert rules from existing templates; AI generates config from service spec",
		phase, sc, aiScale))

	out = append(out, fixedMs("PR merge & prod release", "Gateway Integration", 2,
		"Fixed: requires deployment window and human sign-off", phase))

	// ── ONBOARDING CHANGES ────────────────────────────────────────────────────

	obSc := 1.0
	if req.OnboardingType == "both" {
		obSc = 1.50
	} else if req.OnboardingType == "file_based" {
		obSc = 0.85
	}

	out = append(out, milestone("Batch Service + PGOS onboarding strategy", "Onboarding Changes",
		22*intSc*obSc, 0.45,
		[]string{"backend-engineer:workflows:plan", "backend-engineer:workflows:work", "rzp-discover:payments-processing-platform"},
		"Batch chunking, PGOS IDFCOnboardingStrategy, 14-step execution: GetOrCreateMerchant → CreateDeviceOrder → MarkActivated",
		phase, sc, aiScale))

	encDays := map[string]float64{"aes_cbc": 4, "rsa": 5, "dukpt": 8}[req.EncryptionType]
	if encDays > 0 {
		out = append(out, milestone(
			fmt.Sprintf("Encryption analysis + %s implementation", req.EncryptionType),
			"Onboarding Changes",
			encDays+5, 0.35,
			[]string{"backend-engineer:workflows:work", "slit-generator-v2:create-slits", "code-review:security-reviewer"},
			"AES-256-CBC: decode 64-hex key, dynamic IV, PKCS7 pad, Base64 encode; DUKPT: HSM BDK config; AI generates from bank spec",
			phase, sc, aiScale))

		out = append(out, milestone("Checksum verification", "Onboarding Changes",
			4, 0.38,
			[]string{"backend-engineer:workflows:work", "code-review:security-reviewer"},
			"",
			phase, sc, aiScale))
	}

	out = append(out, milestone("Onboarding / status / deactivation APIs (3 APIs)", "Onboarding Changes",
		5*intSc+6, 0.40,
		[]string{"backend-engineer:workflows:work", "slit-generator-v2:create-slits"},
		"RegMerchantAPI, AddTerminalAPI, UpdateTerminalAPI (integrations-go) + status check + deactivation flow",
		phase, sc, aiScale))

	out = append(out, milestone("Post-payment — batch settlement + data replication", "Onboarding Changes",
		6, 0.42,
		[]string{"backend-engineer:workflows:work", "rzp-discover:financial-infra"},
		"Batch settlement cron config for cards; data replication to Ezetap/Omni for UPI",
		phase, sc, aiScale))

	out = append(out, milestone("Onboarding dev testing + ITF tests", "Onboarding Changes",
		11, 0.35,
		[]string{"backend-engineer:workflows:devstack-test", "slit-generator-v2:create-slits", "slit-generator-v2:test-executor"},
		"E2E onboarding test: file upload → batch processing → PGOS → device order → Ezetap sync",
		phase, sc, aiScale))

	// Bank UAT: fixed in Phase 1, omit in Phase 2 (covered by Phase 1 UAT)
	if phase == 1 {
		out = append(out, fixedMs("Bank UAT testing", "Onboarding Changes", 10,
			"Fixed external dependency — bank provides test creds, runs UAT, issues prod credentials 8 working days post sign-off", phase))
	}

	// ── MERCHANT DASHBOARD (Phase 1 only) ─────────────────────────────────────
	if phase == 1 {
		out = append(out, milestone("Merchant dashboard — bank URL config + feature flags", "Merchant Dashboard",
			9, 0.35,
			[]string{"frontend-design", "code-review:blade-reviewer", "webapp-testing"},
			"Admin dashboard bulk upload for SFDC/WL/Fiserv; Blade compliance auto-checked",
			phase, sc, aiScale))
	}

	// ── SOUNDBOX & DEVICES (Phase 1 only) ─────────────────────────────────────
	if req.HasDevices && phase == 1 {
		out = append(out, milestone("SoundBox language files + device config (CS60→WD10)", "SoundBox & Devices",
			10, 0.50,
			[]string{"backend-engineer:workflows:work", "rzp-discover:pos-platform"},
			"Language audio files, zero-day integration; A910, ET389, A99, A50D device support",
			phase, sc, aiScale))
		out = append(out, milestone("Device E2E testing", "SoundBox & Devices",
			5, 0.60,
			[]string{"slit-generator-v2:test-executor"},
			"Physical device loops limit automation; AI assists test planning and defect triage",
			phase, sc, aiScale))
		out = append(out, milestone("QR Sync App — install/deinstall/breakfix flows", "QR Sync App",
			14, 0.50,
			[]string{"backend-engineer:workflows:plan", "backend-engineer:workflows:work", "slit-generator-v2:create-slits"},
			"NFC tag mapping, device↔merchant mapping, dummy transaction verification, SQR static QR in DQR",
			phase, sc, aiScale))
	}

	// ── RAZORPAY APP (Phase 1 only) ───────────────────────────────────────────
	if req.HasMobileApp && phase == 1 {
		out = append(out, milestone("Bank addition in Razorpay App flows", "Razorpay App",
			6, 0.42,
			[]string{"frontend-design", "code-review:blade-reviewer"},
			"",
			phase, sc, aiScale))
	}

	// ── REPORTING ─────────────────────────────────────────────────────────────
	out = append(out, milestone("Reporting (settlement/recon/offline collection)", "Reporting",
		8, 0.42,
		[]string{"backend-engineer:workflows:work", "rzp-discover:financial-infra", "rzp-discover:data-engineering"},
		"Daily settlement file, offline collection, LACR reports, Fiserv misc debit adjustment — data engineering needed",
		phase, sc, aiScale))

	// ── NPCI CERTIFICATION ────────────────────────────────────────────────────
	if req.HasNPCICert {
		out = append(out, milestone("NPCI Certification (~30 test cases)", "NPCI Certification",
			15, 0.53,
			[]string{"slit-generator-v2:create-slits", "backend-engineer:workflows:work"},
			"Software testing only; hardware and bank/NPCI coordination is manual",
			phase, sc, aiScale))
	}

	// ── E2E TESTING + GO-LIVE ─────────────────────────────────────────────────
	out = append(out, milestone("E2E testing + observability + prod go-live", "End-to-End Testing",
		13, 0.40,
		[]string{"slit-generator-v2:create-slits", "slit-generator-v2:test-executor", "webapp-testing"},
		"Automated E2E suite, prod pilot with bank, 1st merchant go-live; QA sign-off remains manual",
		phase, sc, aiScale))

	return out
}

// ── Reference data ─────────────────────────────────────────────────────────────

func referenceCalibrations() []ReferenceCalib {
	return []ReferenceCalib{
		{
			Name:            "IDFC Omni migration — Phase I (WL + Montran, actual)",
			TraditionalDays: 97,
			AIDays:          50,
			SavingsPct:      48.5,
			IntegrationType: "migration",
			Notes: "Likith Reddy's actual: 61d Phase I + 36d Phase II = 97d (incl. 19d buffer). " +
				"AI projection: ~50d. Source: IDFC Omni migration_ Items to develop.xlsx",
		},
		{
			Name:            "IDFC Bank greenfield (tech spec baseline)",
			TraditionalDays: 185,
			AIDays:          83,
			SavingsPct:      55.1,
			IntegrationType: "new_gateway",
			Notes: "14-step onboarding, AES-256-CBC + OAuth2 RS256, Cards + UPI, file-based CSV (390 cols). " +
				"Source: Tech Spec – IDFC Bank.pdf (Feb 2026)",
		},
		{
			Name:            "Axis Bank Fastag H2H (NPCI cert)",
			TraditionalDays: 210,
			AIDays:          95,
			SavingsPct:      54.8,
			IntegrationType: "new_gateway",
			Notes: "90-day spec, ISO8583, key exchange ceremony, NPCI certification, hardware. " +
				"Source: Gateway_Integration_CheckList.xlsx Rough Sample sheet",
		},
		{
			Name:            "Kotak Juspay UPI (API-based)",
			TraditionalDays: 120,
			AIDays:          52,
			SavingsPct:      56.7,
			IntegrationType: "new_payment_method",
			Notes: "UPI-only, API-based onboarding, existing gateway infra leveraged. " +
				"Estimated from repo pattern analysis across integrations-upi and pg-router",
		},
	}
}

// ── Assumptions ────────────────────────────────────────────────────────────────

func buildAssumptions(req EstimationRequest) []string {
	base := []string{
		"Engineer is proficient with Claude Code CLI and has all relevant skills installed (slit-generator-v2, backend-engineer:workflows, code-review, rzp-discover)",
		"Bank tech spec / BRD is available before sprint begins",
		"Devstack environment is operational and accessible",
		"Bank UAT credentials provided on time — external dependency, not reducible by AI",
		"Estimates in working days: 1 senior engineer, 8h/day",
		"AI reduction factors calibrated against IDFC Bank tech spec (Feb 2026) and Likith Reddy's actual Phase I/II breakdown",
		"Bank UAT (10d), prod release windows, and hardware sign-offs are fixed — AI cannot accelerate these",
	}

	switch req.CodebaseReuse {
	case "migration":
		base = append(base,
			"Migration mode: traditional baseline deflated to 40% of greenfield (calibrated against IDFC Omni actual 97d vs greenfield 240d)",
			"AI gains 12% more effectiveness in migration mode — existing patterns give AI more context to work with")
	case "partial":
		base = append(base,
			"Partial reuse: traditional baseline deflated to 60% of greenfield — some service patterns already exist")
	}

	if req.NumPhases == 2 {
		base = append(base,
			"Phase 2 applies 55% of Phase 1 scope — patterns, tests, and CI fully transfer; no re-doing devices/dashboard/mobile app",
			"Phase 2 AI is 8% more effective than Phase 1 — engineer knows the codebase deeply by then")
	}

	if req.BufferPct > 0 {
		base = append(base, fmt.Sprintf(
			"Buffer of %.0f%% of net work added; AI reduces buffer needs by 40%% (fewer rework loops, auto-generated SLITs catch bugs earlier)",
			req.BufferPct))
	}

	if req.EncryptionType == "dukpt" {
		base = append(base, "DUKPT: HSM hardware and BDK key available on day 1; AI implements code but not the physical key ceremony")
	}
	if req.HasDevices {
		base = append(base, "Physical device stock (A910, ET389, A99, A50D, CS60, WD10) available from day 1; audio assets pre-recorded")
	}
	if req.OnboardingType == "both" {
		base = append(base, "Both file-based and API-based onboarding in scope — 1.5x single-mode; includes SFDC, WL, Fiserv, NTB/ETB SoundBox flows")
	}
	return base
}
