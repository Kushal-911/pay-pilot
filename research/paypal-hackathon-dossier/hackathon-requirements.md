# Hackathon requirements and entry audit

**Reviewed:** 19 October 2026. Rules may change; the [official rules page](https://paypalaihackathon.devpost.com/rules) is controlling.

## Event facts

| Item | Officially published information |
|---|---|
| Event | PayPal AI Hackathon — “Build what’s next with PayPal and AI” |
| Submission window | 1 October 2026, 9:15 a.m. PT through 12 November 2026, noon PT |
| Deadline | 12 November 2026, noon PST on Devpost. For India, that is 13 November 2026, 1:30 a.m. IST (confirm the platform countdown near submission). |
| Judging | 1–15 December 2026; winners expected around 21 December 2026 |
| Core requirement | Build/update an application meaningfully integrating PayPal Developer Platform in Sandbox and an AI tool/model/platform. PayPal must be central. |
| Tracks | No prescribed problem statements or tracks |
| Functional threshold | Working prototype/PoC, installable and consistent with video and description. Static mockup alone is disallowed. |
| Submission | Text description; working demo URL or complete setup instructions; public code repository; visible open-source license; tool/model disclosure; public YouTube demo under three minutes |
| Evaluation | Stage 1 pass/fail viability/theme/API fit; then equally weighted implementation, design, potential impact, innovation, and presentation |
| Existing projects | Allowed only if significantly updated after the 1 October submission-period start; explain what changed during the period. Sponsor/administrator decide significance. |
| Eligibility | Age of majority and exclusions in rules; excluded jurisdictions and conflict rules apply to each team member. Check rules for your own eligibility. |
| Sponsor tools | AG Grid, APIMatic, Astropods, Bryntum, Channel3, Elastic, KERNEL, Postman, Render, Zapier are optional. Each sponsor prize has extra requirements. |
| Project ownership | Entry must be original, solely owned by entrant/team, and respect third-party IP/licensing. Open-source dependencies are allowed with compliance and meaningful enhancements. |
| Testing availability | Keep project available free and without access restrictions for the judging period; provide test instructions/credentials if needed. Never put production secrets in a public repo. |

The landing page inconsistently describes the total prize pool ($67,500 in one section and $69,750 in another). Treat the rules and prize-specific section as authoritative; do not put a total in the pitch without reconfirming it.

## Equal-weight scorecard translated into acceptance tests

| Judge criterion | What the project must visibly prove | Current status / gap |
|---|---|---|
| Technological implementation | Successful PayPal Sandbox order/capture and a meaningful AI call grounded in persisted records; deterministic safety state survives model failure | Code exists for PayPal Orders v2, capture, webhook verification, and OpenAI function calls. Prior audit says no real Sandbox/DB run and Docker daemon was unavailable. Reverify. |
| Design | One coherent merchant workflow from exception queue to timeline, explanation, evidence, and action | Missing: no frontend existed at last repo audit. |
| Potential impact | Named merchant operator/user, painful moment, believable business value, and measurable proxy | Hypothesis only. Need short merchant interviews or credible sourced research; don't promise reduced chargebacks or revenue uplift without data. |
| Innovation | Clear difference from generic AI payment agent and similar submissions | Requires distinction on cross-system causal reconciliation and safe post-timeout handling. Existing entries overlap with agent spend controls. |
| Presentation | Under-three-minute end-to-end demo, understandable without narration-heavy explanation | Not started; demo must use the actual functional app and name synthetic data. |

## Entry-specific issue: PayPilot predates the event

The project brief says PayPilot had already been in development for roughly two days. That does not automatically disqualify it, but the official rules require significant work after the period began and tell entrants to explain it. The visible Git history has substantive checkout, payment-reliability, investigation, and webhook commits dated 18 October 2026, after the 1 October start; this is useful supporting evidence, but not a guarantee of eligibility. Preserve a dated, honest change record: feature list, commits, release notes, screenshots, and demo. Do not imply the project was started during the event. Have the entrant confirm the project's eligibility against Section 4 of the rules; only Sponsor/Devpost can decide “significant change.”

## Submission completion checklist

- [ ] Entrant registered and eligibility checked for every teammate.
- [ ] Clearly explain the meaningful post-1 October work in the submission.
- [ ] Public source repository; add a suitable open-source license at repo root after deciding ownership.
- [ ] Pin exact setup instructions and confirm clean install/run from a fresh clone.
- [ ] Functional demo runs consistently; publish a hosted Sandbox demo if feasible or make local steps exceptionally short.
- [ ] PayPal Sandbox path tested with the matching business/buyer test accounts; document environment setup without leaking credentials.
- [ ] List model/AI service and explain its bounded use, retrieval sources, and fallback.
- [ ] Text submission answers what it does, who benefits, how PayPal and AI are central, and how it differs.
- [ ] Public YouTube video is under 3 minutes, shows the real application functioning on its target device, uses no unlicensed music/marks, and link is on Devpost.
- [ ] Keep hosted app / test build available through judging period; rotate/revoke temporary credentials after judging.
- [ ] Re-check gallery, event FAQ/updates, prize requirements, and the deadline countdown immediately before submission.

## Schedule working backwards

Deadline is 12 November PT. On 19 October there are just over three weeks. Suggested risk-first timeline: by 22 Oct make the Sandbox flow reproducible; by 27 Oct complete one UI journey; by 31 Oct validate the seeded ambiguous-outcome scenario with DB persistence; by 4 Nov complete independent clean-run QA and submission text; by 7 Nov record and review the demo; by 10 Nov freeze scope; submit early on 11 Nov local time. This is an execution proposal, not an event-prescribed schedule.
