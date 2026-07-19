# Opioid Safety Reference

## CDC 2022 Clinical Practice Guideline for Prescribing Opioids

Source: CDC Clinical Practice Guideline for Prescribing Opioids for Pain -- United States, 2022. MMWR Recomm Rep 2022;71(No. RR-3):1-95.

### Key Recommendations (12 Recommendations)

**Determining Whether to Initiate Opioids:**

1. Non-opioid therapies are preferred for chronic pain. Opioids should not be first-line therapy. Maximize non-pharmacologic and non-opioid pharmacologic treatments first.

2. Before starting opioids, establish realistic treatment goals with the patient focused on function (not just pain scores). Plan to discontinue if benefits do not outweigh risks.

3. Before and periodically during opioid therapy, discuss risks and benefits with the patient. Document the discussion.

**Opioid Selection, Dosage, Duration, Follow-Up:**

4. When starting opioids, prescribe immediate-release (IR) formulations, not extended-release/long-acting (ER/LA). ER/LA should not be used for acute pain or as first opioid.

5. Start with the lowest effective dose. Use caution at any dosage. Reassess at 50 MME/day. Avoid increasing to >= 90 MME/day without careful justification.

6. For acute pain, prescribe the shortest effective duration. 3 days or fewer is often sufficient. Rarely more than 7 days needed.

7. Evaluate benefits and harms within 1-4 weeks of starting or dose escalation. Evaluate ongoing therapy every 3 months at minimum. Taper and discontinue if benefits do not outweigh risks.

**Assessing Risk and Addressing Harms:**

8. Before starting and periodically during opioid therapy, evaluate risk factors for opioid-related harms. Consider naloxone co-prescription.

9. Check PDMP data before starting and periodically during therapy (every prescription or at minimum every 3 months).

10. Use urine drug testing before starting and at least annually during opioid therapy. Consider more frequent testing if risk factors present.

11. Avoid concurrent prescribing of opioids and benzodiazepines whenever possible.

12. Offer or arrange for evidence-based treatment (buprenorphine or methadone) for patients with opioid use disorder.

### MME Thresholds

| Daily MME | Risk Level | Clinical Action |
|-----------|-----------|----------------|
| < 20 | Standard | Routine monitoring |
| 20-49 | Low-moderate | Standard monitoring, ensure multimodal approach |
| 50-89 | Increased risk | Reassess benefits vs harms. Consider dose reduction. Co-prescribe naloxone. |
| >= 90 | High risk | Avoid unless carefully justified. Document rationale. Co-prescribe naloxone. Consider pain specialist consultation. |

## MME Conversion Reference Table

### Standard Conversion Factors (CDC)

| Opioid | Route | MME Conversion Factor | Example: 30mg = ? MME |
|--------|-------|----------------------|----------------------|
| Codeine | Oral | 0.15 | 4.5 MME |
| Hydrocodone | Oral | 1.0 | 30 MME |
| Hydromorphone | Oral | 4.0 | 120 MME |
| Morphine | Oral | 1.0 (reference) | 30 MME |
| Oxycodone | Oral | 1.5 | 45 MME |
| Oxymorphone | Oral | 3.0 | 90 MME |
| Tapentadol | Oral | 0.4 | 12 MME |
| Tramadol | Oral | 0.1 | 3 MME |
### 🚨 CLINICAL SAFETY GUARDRAIL: INCOMPLETE CROSS-TOLERANCE
**Mandatory System Check for Opioid Rotation:**
- **Trigger:** Any AI recommendation involving a switch from one opioid to another.
- **Action:** Calculate standard Equianalgesic Dose based on MME tables.
- **Safety Override:** AUTOMATICALLY REDUCE the calculated dose by **25% to 50%**.
- **Rationale:** Human receptors exhibit incomplete cross-tolerance. A 1:1 mathematical conversion risks fatal respiratory depression. The AI must explicitly state this reduction and prompt the physician to titrate up based on patient response.

### Fentanyl Transdermal Conversion

Fentanyl patches are dosed in mcg/hr. Conversion to daily MME:

| Fentanyl Patch (mcg/hr) | Approximate Daily MME |
|--------------------------|---------------------|
| 12 | 30 MME |
| 25 | 60 MME |
| 37 | 90 MME |
| 50 | 120 MME |
| 75 | 180 MME |
| 100 | 240 MME |

Formula: fentanyl mcg/hr * 2.4 = daily MME

### Methadone Conversion (Variable -- Use Extreme Caution)

Methadone has a variable conversion ratio that increases with dose due to its unique pharmacology (long half-life, NMDA receptor activity, potential for accumulation).

| Current Daily Morphine Equivalent | Methadone Conversion Ratio (oral morphine:methadone) |
|----------------------------------|---------------------------------------------------|
| < 30 MME/day | 2:1 (30mg morphine = 15mg methadone) |
| 30-99 MME/day | 4:1 (60mg morphine = 15mg methadone) |
| 100-299 MME/day | 8:1 (160mg morphine = 20mg methadone) |
| 300-499 MME/day | 12:1 (360mg morphine = 30mg methadone) |
| 500-999 MME/day | 15:1 |
| >= 1000 MME/day | 20:1 |

Methadone conversions should be managed by experienced clinicians. QTc monitoring required.

### Buprenorphine Conversion

- Buprenorphine transdermal (Butrans): MME factor varies. 5 mcg/hr patch ~ 12 MME/day. 10 mcg/hr ~ 24 MME/day. 20 mcg/hr ~ 48 MME/day.
- Buprenorphine buccal (Belbuca): 150mcg ~ 15 MME, 300mcg ~ 30 MME, 450mcg ~ 45 MME, 600mcg ~ 60 MME, 750mcg ~ 75 MME, 900mcg ~ 90 MME.
- Suboxone (buprenorphine/naloxone for OUD): not typically converted to MME; used for addiction treatment.

### Parenteral to Oral Conversions

| Opioid | IV/IM to Oral Ratio |
|--------|-------------------|
| Morphine | 1:3 (10mg IV = 30mg oral) |
| Hydromorphone | 1:5 (2mg IV = 10mg oral) |
| Oxycodone | No standard IV form |

## Opioid Taper Protocols

### When to Taper

- Benefits no longer outweigh risks
- Patient requests tapering
- No functional improvement despite adequate trial
- Adverse effects (sedation, respiratory depression, falls)
- Aberrant behaviors (lost prescriptions, early refills, escalating doses)
- Concurrent high-risk medications (benzodiazepines) that cannot be discontinued
- Evidence of opioid use disorder (consider transition to buprenorphine or methadone)

### Taper Rate Guidelines

**Gradual taper (preferred):**
- Reduce by 5-10% of total daily dose every 2-4 weeks
- Slower tapers (5% per month) for patients on long-term therapy (>1 year)
- Pause or slow the taper if withdrawal symptoms are intolerable
- Typical taper duration: months to years depending on starting dose and duration of use

**Moderate taper:**
- Reduce by 10-20% every 2-4 weeks
- Used when more rapid reduction is clinically indicated but not urgent

**Rapid taper (only when clinically necessary):**
- Reduce by 20-25% every 1-2 weeks
- Indicated for: severe adverse effects, active substance abuse, patient safety concerns
- Higher risk of withdrawal; provide supportive medications

### Withdrawal Management During Taper

| Symptom | Management |
|---------|-----------|
| Anxiety/insomnia | Hydroxyzine 25-50mg, trazodone 50-100mg, clonidine 0.1mg BID |
| GI symptoms (nausea, diarrhea, cramping) | Loperamide 2-4mg PRN, ondansetron 4-8mg PRN, dicyclomine 20mg QID |
| Muscle aches | NSAIDs (ibuprofen 400-600mg TID), heating pad |
| Restless legs | Clonidine 0.1mg BID-TID, gabapentin 300-600mg TID |
| Diaphoresis | Clonidine 0.1mg BID |
| Lacrimation/rhinorrhea | Supportive care |

Clonidine is the most broadly useful withdrawal adjunct (addresses autonomic hyperactivity).

### Taper Monitoring

- Check in by phone or visit every 1-2 weeks during active taper
- Monitor for withdrawal symptoms, pain flares, and mood changes
- PHQ-9 for depression (common during taper)
- PEG scale for function
- Reassure patient that temporary increase in pain during taper is expected and usually resolves

## Risk Mitigation Strategies

### Urine Drug Testing (UDT)

**Frequency:**
- Before initiating opioid therapy (baseline)
- At least annually during ongoing therapy
- More frequently (every 3-6 months) if risk factors: history of SUD, aberrant behaviors, concurrent benzodiazepines, high MME

**Types of testing:**
- Immunoassay (point-of-care): screening test, detects drug classes, subject to false positives/negatives
- Confirmatory (GC-MS, LC-MS/MS): quantitative, identifies specific drugs, used for unexpected results

**Interpreting results:**

| Scenario | Possible Explanations | Action |
|----------|---------------------|--------|
| Prescribed opioid ABSENT | Non-adherence, diversion, metabolism variation, timing of last dose | Confirm timing, retest, consider confirmatory. Discuss adherence. |
| Non-prescribed opioid PRESENT | Use of non-prescribed opioid, obtaining from multiple sources | Discuss. Check PDMP. Consider OUD evaluation. |
| Illicit substance PRESENT (cocaine, methamphetamine, heroin) | Active substance use | Screen for SUD. Consider addiction referral. Reassess opioid prescribing. |
| Benzodiazepine PRESENT (not prescribed) | Non-prescribed use | Discuss. Check PDMP. Counsel on risks of combination. |
| THC PRESENT | Cannabis use | Discuss. Not necessarily contraindication but document. Assess for SUD if concern. |

**Common false positives (immunoassay):**
- Amphetamine: pseudoephedrine, bupropion, trazodone, labetalol, ranitidine
- Benzodiazepine: sertraline, efavirenz
- Opiates: poppy seeds, dextromethorphan, quinolone antibiotics
- THC: dronabinol, NSAIDs (rare), hemp products
- PCP: dextromethorphan, ketamine, diphenhydramine

### PDMP (Prescription Drug Monitoring Program) Checks

- Check before every new opioid prescription or at minimum every 3 months
- Look for: multiple prescribers, multiple pharmacies, overlapping opioid prescriptions, concurrent benzodiazepines, escalating quantities
- Document that PDMP was checked in the medical record
- Red flags: obtaining controlled substances from >= 3 prescribers or >= 3 pharmacies in a 90-day period

### Naloxone Co-Prescription

**Prescribe naloxone when:**
- MME >= 50/day
- Concurrent benzodiazepines
- History of overdose
- History of substance use disorder
- Patient request
- Respiratory conditions (COPD, sleep apnea)
- Renal or hepatic impairment
- Elderly (>= 65)

**Formulations:**
- Narcan nasal spray: 4mg intranasal (OTC in most states)
- Generic naloxone nasal spray: 4mg or 8mg
- Naloxone auto-injector: 2mg IM (if available)

**Patient/family education:**
- Recognizing overdose: unresponsive, slow/stopped breathing, blue lips/fingertips, snoring/gurgling
- Call 911 first
- Administer naloxone (1 spray in one nostril)
- If no response in 2-3 minutes, administer second dose in other nostril
- Rescue breathing if trained
- Stay with patient until EMS arrives
- Naloxone wears off in 30-90 minutes; overdose may recur

## Common Adjuvant Medications by Pain Type

### Neuropathic Pain

| Medication | Starting Dose | Therapeutic Range | Key Notes |
|-----------|--------------|------------------|-----------|
| Gabapentin | 100-300mg QHS | 300-1200mg TID (900-3600mg/day) | Titrate slowly. Renal dose adjustment. Dizziness, sedation. |
| Pregabalin | 25-75mg BID | 75-300mg BID (150-600mg/day) | Schedule V. Faster onset than gabapentin. Weight gain. |
| Duloxetine | 30mg daily | 60-120mg daily | FDA-approved for neuropathic pain, fibromyalgia, chronic musculoskeletal. GI side effects. |
| Amitriptyline | 10-25mg QHS | 25-150mg QHS | Effective but anticholinergic. Avoid in elderly (Beers). |
| Nortriptyline | 10-25mg QHS | 25-150mg QHS | Less anticholinergic than amitriptyline. Preferred TCA. |
| Venlafaxine | 37.5mg daily | 150-225mg daily | SNRI. Monitor BP at higher doses. |
| Lidocaine 5% patch | 1-3 patches to area | Up to 12 hours on / 12 hours off | Topical. Localized neuropathic pain. Minimal systemic absorption. |
| Capsaicin 8% patch | Applied by provider | Single 60-minute application q3 months | Requires provider application. Localized neuropathic pain. |
| Carbamazepine | 100mg BID | 200-400mg BID | First-line for trigeminal neuralgia only. Monitor CBC, LFTs, sodium. |

### Musculoskeletal Pain

| Medication | Starting Dose | Key Notes |
|-----------|--------------|-----------|
| Cyclobenzaprine | 5-10mg TID or 5mg QHS | Muscle relaxant. Short-term use (<2-3 weeks). Sedating. Avoid in elderly. |
| Baclofen | 5mg TID | Centrally acting. For spasticity. Taper slowly (withdrawal risk). |
| Tizanidine | 2mg QHS-TID | Alpha-2 agonist. Sedating. Monitor LFTs. |
| Methocarbamol | 500-750mg TID-QID | Less sedating than cyclobenzaprine. No abuse potential. |
| Diclofenac gel 1% | Apply QID to affected area | Topical NSAID. First-line for knee/hand OA. Less systemic exposure. |
| Acetaminophen | 325-650mg q4-6h | Max 3g/day (2g/day if liver disease or alcohol use). Often inadequate alone for chronic pain. |
| Celecoxib | 100-200mg daily-BID | COX-2 selective. Lower GI risk. CV risk similar to other NSAIDs. |

### Headache/Migraine Prevention (if comorbid)

| Medication | Dose | Notes |
|-----------|------|-------|
| Amitriptyline | 10-50mg QHS | Dual benefit for pain and headache |
| Topiramate | 25-100mg BID | Weight loss (may benefit obese patients) |
| Propranolol | 40-160mg daily | Also helps anxiety/tremor |
| Venlafaxine | 150mg daily | Also treats depression/anxiety |
