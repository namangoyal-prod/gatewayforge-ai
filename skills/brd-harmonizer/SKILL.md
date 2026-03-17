# BRD Harmonizer Skill

Validates uploaded Business Requirements Documents (BRDs) against Razorpay's gateway integration standards. Ensures completeness, technical accuracy, and implementation readiness before engineering resources are committed.

## Purpose

Transform inconsistent, incomplete BRDs into validated, implementation-ready specifications by:
1. Checking completeness against mandatory sections
2. Validating technical accuracy (ISO 8583, API specs, encryption)
3. Ensuring conformance to Razorpay's BRD template
4. Flagging ambiguities and gaps
5. Auto-generating fix suggestions based on similar integrations

## Input

- **BRD Document**: PDF, DOCX, or Google Doc link
- **Integration Context**: Partner name, type (gateway/aggregator/direct), payment methods

## Output

```json
{
  "quality_score": 85,
  "status": "GREEN",
  "validation_report": {
    "completeness": {
      "score": 90,
      "missing_sections": [],
      "incomplete_sections": ["Error Handling"]
    },
    "technical_accuracy": {
      "score": 80,
      "issues": [
        {
          "section": "ISO 8583 Mapping",
          "severity": "MEDIUM",
          "issue": "MTI 0200 missing field 39 (response code)",
          "fix": "Add field 39 mapping for transaction status"
        }
      ]
    },
    "conformance": {
      "score": 85,
      "deviations": ["Non-standard timeout values"]
    },
    "clarity": {
      "score": 90,
      "ambiguities": []
    },
    "regulatory_compliance": {
      "score": 95,
      "flags": []
    }
  },
  "gap_analysis": [
    {
      "section": "Error Handling",
      "gap": "Missing retry logic for timeout scenarios",
      "severity": "HIGH",
      "suggested_fix": "Add retry configuration: max_retries=3, backoff=exponential, timeout=30s"
    }
  ],
  "auto_fix_suggestions": [...],
  "comparison_matrix": {
    "reference_brd": "Amazon Pay DQR UPI",
    "similarity_score": 82,
    "key_differences": [...]
  }
}
```

## Validation Dimensions

### 1. Completeness (Weight: 25%)

**Mandatory Sections:**
- ✅ Partner Profile (name, type, contact)
- ✅ API Specifications (endpoints, methods, auth)
- ✅ Message Formats (ISO 8583 / JSON / XML)
- ✅ Authentication (OAuth, API Key, DUKPT, 3DES, AES)
- ✅ Error Codes & Handling
- ✅ Settlement Logic
- ✅ Reconciliation Specs

**Scoring:** `(present_sections / 7) * 25`

### 2. Technical Accuracy (Weight: 25%)

**Checks:**
- ISO 8583 field mappings are valid
- MTI codes are correct (0200, 0210, 0400, 0420, etc.)
- UIPP parameters conform to NPCI specs
- Encryption specs are valid (DUKPT/3DES/AES with proper key lengths)
- Network protocol is specified (TCP/IP, HTTP/REST, SOAP)
- Timeout values are reasonable (not < 5s or > 180s)

### 3. Conformance (Weight: 20%)

**Template Alignment:**
- Follows Razorpay BRD template structure
- Uses standard naming conventions
- Field mappings align with internal data models
- Configuration follows existing patterns

### 4. Clarity & Specificity (Weight: 15%)

**Anti-Patterns to Flag:**
- Vague language: "appropriate", "reasonable", "as needed"
- Missing conditional logic details
- Timeout/retry values marked as "TBD"
- Error handling described as "standard" without specifics

### 5. Regulatory Compliance (Weight: 15%)

**Requirements:**
- NPCI circulars referenced (for UPI)
- RBI guidelines compliance
- PCI-DSS requirements mentioned
- Data residency requirements (India/UAE/SEA)

## Implementation Logic

```python
def harmonize_brd(brd_document, integration_context):
    # Step 1: Extract structured data from document
    structured_data = extract_brd_data(brd_document)

    # Step 2: Load Razorpay BRD template
    template = load_razorpay_brd_template()

    # Step 3: Find most similar approved BRD
    reference_brd = find_similar_brd(integration_context)

    # Step 4: Score each dimension
    completeness_score = check_completeness(structured_data, template)
    technical_score = validate_technical_accuracy(structured_data)
    conformance_score = check_conformance(structured_data, template)
    clarity_score = assess_clarity(structured_data)
    compliance_score = check_regulatory_compliance(structured_data, integration_context.geography)

    # Step 5: Calculate weighted score
    total_score = (
        completeness_score * 0.25 +
        technical_score * 0.25 +
        conformance_score * 0.20 +
        clarity_score * 0.15 +
        compliance_score * 0.15
    )

    # Step 6: Generate gap analysis
    gaps = identify_gaps(structured_data, template, reference_brd)

    # Step 7: Generate auto-fix suggestions
    fixes = generate_fixes(gaps, reference_brd)

    # Step 8: Create comparison matrix
    comparison = compare_with_reference(structured_data, reference_brd)

    # Step 9: Determine status
    status = "GREEN" if total_score >= 70 else "AMBER" if total_score >= 50 else "RED"

    return {
        "quality_score": total_score,
        "status": status,
        "validation_report": {...},
        "gap_analysis": gaps,
        "auto_fix_suggestions": fixes,
        "comparison_matrix": comparison
    }
```

## Knowledge Base Requirements

The skill requires access to:

1. **Razorpay BRD Template** (canonical structure)
2. **Approved Historical BRDs** (for comparison and pattern learning)
3. **ISO 8583 Field Registry** (for message validation)
4. **NPCI UIPP Specifications** (for UPI integrations)
5. **Regulatory Guidelines Database** (RBI, NPCI, CBUAE, PCI-DSS)
6. **Common Error Patterns** (from incident post-mortems)

## Tool Access

- `document_parser`: Extract text from PDF/DOCX
- `iso8583_validator`: Validate message specifications
- `vector_search`: Find similar approved BRDs
- `llm_analyzer`: Assess clarity, identify ambiguities

## Quality Thresholds

| Score | Status | Action |
|-------|--------|--------|
| ≥ 70 | GREEN | Approved - Proceed to PRD generation |
| 50-69 | AMBER | Needs minor fixes - Review and resubmit |
| < 50 | RED | Incomplete - Major rework required |

## Example Gap Analysis

```json
{
  "section": "ISO 8583 Authorization Message",
  "gap": "Field 39 (Response Code) mapping is missing",
  "severity": "HIGH",
  "impact": "Cannot determine transaction success/failure",
  "suggested_fix": {
    "text": "Map Field 39 to internal status codes",
    "mapping": {
      "00": "SUCCESS",
      "05": "DO_NOT_HONOR",
      "51": "INSUFFICIENT_FUNDS",
      "91": "ISSUER_UNAVAILABLE"
    },
    "reference": "See Codec integration Field 39 mapping in docs/codec-iso8583-spec.md"
  }
}
```

## Auto-Fix Strategy

When suggesting fixes, the skill:

1. **Learns from Similar BRDs**: If integrating a UPI gateway, pulls timeout values, retry configs from existing UPI integrations
2. **Applies Platform Defaults**: Uses Razorpay's standard values where partner hasn't specified
3. **Flags for Confirmation**: Auto-suggestions marked with `[AUTO-FILLED]` tag for human review
4. **Maintains Traceability**: Every suggestion includes source reference

## Edge Cases

- **Multi-format BRDs**: Some partners send Excel + PDF combo → Extract from both, merge data
- **Incomplete API Specs**: Missing endpoint URLs → Flag as blocker, cannot auto-fix
- **Ambiguous Field Mappings**: "Field X is optional" → Request clarification on exact conditions
- **Legacy Protocols**: Some partners use SOAP/XML instead of REST → Validate against SOAP schema

## Integration with Pipeline

1. User uploads BRD via frontend
2. Backend stores in object storage, creates DB entry
3. Temporal workflow triggers BRD Harmonizer skill
4. Skill processes document, generates report
5. Report stored in `brd_documents.gap_analysis` JSON field
6. Frontend displays interactive validation report
7. If GREEN: Workflow proceeds to PRD Generation
8. If AMBER/RED: Workflow pauses for human fixes

## Metrics to Track

- **Validation Accuracy**: % of flagged issues that were real problems (vs false positives)
- **Auto-Fix Adoption**: % of suggestions accepted by users
- **Time to Green**: Avg time from first submission to GREEN status
- **Rejection Rate**: % of BRDs rejected on first submission (target: <5%)
