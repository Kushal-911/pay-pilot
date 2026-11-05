# Hackathon requirements and entry audit

**Reviewed:** 5 November 2026. Rules may change; the [official rules page](https://paypalaihackathon.devpost.com/rules) is controlling. Rechecked the event page, official rules, FAQ, project gallery, and event updates on this date.

## Event facts

| Item | Officially published information |
|---|---|
| Event | PayPal AI Hackathon — “Build what’s next with PayPal and AI” |
| Submission window | 1 October 2026, 9:15 a.m. PT through 12 November 2026, noon PT |
| Deadline | 12 November 2026, noon Pacific Time (PST on the stated date) on Devpost. For India, that is 13 November 2026, 1:30 a.m. IST. Recheck Devpost countdown immediately before submission. |
| Judging | 1 December 2026, 2:00 p.m. PT through 15 December 2026, 8:00 a.m. PT; winners expected around 21 December 2026, 2:00 p.m. PT |
| Core requirement | Build/update an application meaningfully integrating PayPal Developer Platform in Sandbox and an AI tool/model/platform. PayPal must be central. |
| Tracks | No prescribed problem statements or tracks |
| Functional threshold | Working prototype/PoC, installable and consistent with video and description. Static mockup alone is disallowed. |
| Submission | Text description; functional demo URL or complete setup instructions; public code repository with visible open-source license; AI tool/model disclosure; public YouTube demo under three minutes showing the working project on its target device |
| Evaluation | Stage 1 pass/fail viability/theme/API fit; then equally weighted implementation, design, potential impact, innovation, and presentation |
| Existing projects | Allowed if significantly updated after the 1 October submission-period start; explain what changed during the period. Sponsor/administrator decide significance. |
| Eligibility | Age of majority and exclusions in rules; excluded jurisdictions and conflict rules apply to each team member. Check rules for your own eligibility. |
| Sponsor tools | AG Grid, APIMatic, Astropods, Bryntum, Channel3, Elastic, KERNEL, Postman, Render, Zapier are optional. Each sponsor prize has extra requirements. |
| Project ownership | Entry must be original, solely owned by entrant/team, and respect third-party IP/licensing. Open-source dependencies are allowed with compliance and meaningful enhancements. |
| Testing availability | Keep project available free and without access restrictions for the judging period; provide test instructions/credentials if needed. Never put production secrets in a public repo. |
| Functional demo | A hosted URL is optional if the repository gives complete instructions for judges to build and run the working application; static mockups do not qualify. |

The landing page still shows different total prize figures ($67,500 in one section and $69,750 in another); the rules enumerate prize items separately. Avoid claiming a total prize pool. No sponsor-tool prize is currently selected or pursued for PayPilot; sponsor tools are optional. Recheck only if the owner later chooses a sponsor-specific prize.

## Equal-weight scorecard translated into acceptance tests

| Judge criterion | What the project must visibly prove | Current status / gap |
|---|---|---|
| Technological implementation | Successful PayPal Sandbox order/capture and a meaningful AI call grounded in persisted records; deterministic safety state survives model failure | PayPal Orders v2, capture, webhook verification, and OpenAI Responses code exist. PostgreSQL-backed tests and a clean-clone Compose order smoke test passed. No real Sandbox checkout/webhook or live OpenAI run has been verified. |
| Design | One coherent merchant workflow from exception queue to timeline, explanation, evidence, and action | Responsive workbench and synthetic recovery flow are implemented; keyboard, focus, contrast, and phone/tablet checks are recorded in PR #71 and issues #22/#23. A complete Live API operator walkthrough remains unverified. |
| Potential impact | Named merchant operator/user, painful moment, believable business value, and measurable proxy | Hypothesis only. Need short merchant interviews or credible sourced research; don't promise reduced chargebacks or revenue uplift without data. |
| Innovation | Clear difference from generic AI payment agent and similar submissions | Position on cross-system payment evidence and safe post-timeout recovery. The official gallery remains unpublished; the sampled competitor landscape is incomplete. |
| Presentation | Under-three-minute end-to-end demo, understandable without narration-heavy explanation | Public demo recording/publication remains open. Use the real app, show PayPal and AI roles clearly, identify synthetic data, and avoid claims of unverified live runs. |

## Entry-specific issue: PayPilot predates the event

The project brief says PayPilot had already been in development for roughly two days. That does not automatically disqualify it, but the official rules require significant work after the period began and tell entrants to explain it. The visible Git history has substantive checkout, payment-reliability, investigation, and webhook commits dated 18 October 2026, after the 1 October start; this is useful supporting evidence, but not a guarantee of eligibility. Preserve a dated, honest change record: feature list, commits, release notes, screenshots, and demo. Do not imply the project was started during the event. Have the entrant confirm the project's eligibility against Section 4 of the rules; only Sponsor/Devpost can decide “significant change.”

## Submission completion checklist

- [ ] Entrant registered and eligibility checked for every teammate; each entrant must be of majority age and not excluded by Section 3. Eligibility is an owner confirmation; Devpost registration alone does not establish eligibility.
- [ ] Clearly explain the meaningful post-1 October work in the submission.
- [ ] Public source repository; add a suitable open-source license at repo root after deciding ownership.
- [ ] Pin exact setup instructions and confirm clean install/run from a fresh clone.
- [ ] Functional demo runs consistently; hosting is optional if complete setup/run instructions let judges build and run it. Static mockups are not accepted as a functional demo.
- [ ] PayPal Sandbox path tested with the matching business/buyer test accounts; document environment setup without leaking credentials.
- [ ] List model/AI service and explain its bounded use, retrieval sources, and fallback.
- [ ] Text submission answers what it does, who benefits, how PayPal and AI are central, and how it differs.
- [ ] Public YouTube video is under 3 minutes, shows the real application functioning on its target device, uses no unlicensed music/marks, and link is on Devpost.
- [ ] Keep the project available free and without restrictions through the judging period (ends 15 Dec 2026, 8:00 a.m. PT); provide any necessary test instructions/access without exposing production secrets.
- [ ] Re-check gallery, event FAQ/updates, prize requirements, and the deadline countdown immediately before submission.

## Delivery status as of 5 November 2026

- **Verified complete locally:** fresh-clone Compose startup and database-backed order smoke test; PostgreSQL-backed integration suite; workbench keyboard/accessibility and narrow viewport checks. See `../../PROJECT_STATUS.md`, merged PR #71, and closed issues #22/#23.
- **Still required for V1:** actual PayPal Sandbox buyer-approved checkout and reachable webhook proof; complete Live API UI journey; live OpenAI run or truthful deterministic-only disclosure; owner-selected license; accurate final submission, public video, and go/no-go review.
- **Near-term plan:** request an early eligibility review if still available; complete Sandbox and Live API proof as soon as credentials/runtime permit; feature freeze by 10 Nov; record/review the video and prepare the entry before 11 Nov IST; submit early. These are project targets, not additional event rules.

The Devpost event update reviewed on 5 November recommends a public, under-three-minute YouTube video, opening with the user problem, showing PayPal and AI working end to end, and uploading early. The update also invited builders to submit by 5 November for early eligibility review while noting entries can continue to be edited until the deadline; confirm in the entrant's account whether that review window remains available.
