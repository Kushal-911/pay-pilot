# PayPal AI Hackathon dossier

**Research date:** 19 October 2026  
**Competition:** [PayPal AI Hackathon: Build what’s next with PayPal and AI](https://paypalaihackathon.devpost.com/)  
**Purpose:** turn public evidence into a project direction, submission plan, and concise outreach approach.  
**Scope note:** this is a public-source research dossier, not an account of PayPal's private architecture, internal plans, unpublished work, or confidential team documents.

## Recommendation in one page

Keep the PayPilot codebase, but narrow its proposition to **transaction truth and safe recovery for AI-assisted commerce**. The distinctive demo should show how PayPilot resolves a payment whose outcome is ambiguous across a merchant's order, PayPal's Orders API, capture attempt, and verified webhook history. It should return a causal timeline, separate facts from hypotheses, cite each fact to stored/provider evidence, and make the next permitted action explicit. It must not retry a capture merely because the first request timed out.

This territory is differentiated from the early public submissions we found, which cluster around pre-purchase spending policies, agent liability, group expense management, and payment routing. It is also adjacent to, rather than a copy of, PayPal's public agent toolkit: that toolkit exposes commerce actions; PayPilot is an independently governed reliability and evidence layer around the *outcome* of those actions. This is a positioning hypothesis, not evidence of a gap in PayPal's internal roadmap.

**Blunt caveat:** “AI payment reconciliation” by itself is not novel. PayPal's AI developer pages explicitly mention payment reconciliation as a use case; the official agent toolkit includes reporting/insights and other merchant actions. The entry only becomes compelling if the functional demo proves a hard edge case: delayed/duplicated/out-of-order evidence, no false “paid” state, no unsafe duplicate capture, and grounded explanation. If the team cannot deliver that end-to-end with real Sandbox calls, switch to a smaller but functioning single scenario rather than claiming a broad recovery platform.

## Contents

- [Hackathon rules and submission checklist](hackathon-requirements.md)
- [PayPal public direction and product map](paypal-direction.md)
- [Project and competition landscape](project-landscape.md)
- [Opportunity assessment and hard audit](opportunity-analysis.md)
- [People and appropriate outreach](people-and-outreach.md)
- [Pitch, demo, and Devpost text outline](pitch-and-demo.md)
- [Sources, caveats, and research log](sources.md)

## Facts versus inferences

- **Fact** means stated in an official rule, PayPal publication, API documentation, or project page linked in the dossier.
- **Inference** means a product or business interpretation drawn from public facts. It is marked as such and should not be represented as PayPal's internal intent.
- **Unknown** means no credible public source located, access not available, or the claim needs direct validation.

The current Devpost project gallery says the managers have not published it. Search indexing nevertheless exposes some individual public project pages submitted to this event. Those are a useful early sample, not a complete competitor census. The project list and count can change; re-check the gallery and search for new submissions before finalizing the pitch/video.

## Two separate meanings of “architecture”

The technical guide describes (1) public PayPal APIs and documented integration flows, and (2) PayPilot's code architecture. PayPal does not publish its complete production internal architecture, proprietary models, risk systems, or operational playbooks. The dossier deliberately does not invent those internals.
