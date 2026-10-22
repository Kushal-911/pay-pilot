# People and appropriate outreach

This list uses names and roles published on the Devpost event page or public PayPal Developer materials. A listed judge is not a promise of personal access or a recruitment contact. Use the official hackathon Discord or public professional channels; do not seek private contact details, mass-message staff, or ask judges for preferential treatment.

## Relevant listed PayPal judges

| Person | Publicly listed role | Why their perspective could help | Appropriate path |
|---|---|---|---|
| Jo Franchetti | Developer Advocate, PayPal | Public build-session leader; developer integration and Sandbox feedback | Hackathon Discord / public event Q&A; her profile is searchable by name, but search results conflict on current employer, so use the event channel rather than a guessed direct profile URL. |
| Eddie Jaoude | Developer Advocate, PayPal | Developer tooling, open-source and agentic project communication | [LinkedIn](https://www.linkedin.com/in/eddiejaoude/) or public posts; ask one concise technical question, not a job referral. |
| Marco Podien | Developer Advocate, PayPal | Developer education, TypeScript/AI and DX; ran public hackathon build session | [LinkedIn](https://www.linkedin.com/in/marcopodien); best initial contact for feedback on integration/demo clarity. |
| Karthik Ravi | Software Engineer, PayPal | Technical reliability and agentic workflow perspective; public Devpost judge listing | [Public LinkedIn](https://www.linkedin.com/in/kaddynator/). Ask for feedback on evidence/reliability framing, recognizing he may be judging. |
| Himraj Singh | Engineering Manager, PayPal | Engineering quality and productization perspective; has publicly shared this event and spoken about enterprise AI trust | [Public LinkedIn profile](https://www.linkedin.com/in/himraj-singh-5048a266/) or event Discord. Prefer a short, relevant feedback request. |
| Nathaniel Olson | Senior Product Manager, PayPal | Product/customer-value framing; role per Devpost | Use official event channel. Search did not verify a unique public profile, so no direct link is supplied. |

Other judge names on Devpost belong to sponsor companies (AG Grid, APIMatic, Bryntum, Channel3, Elastic, Postman, Render). Contact them only for a specific sponsor-tool question or award criteria.

## Best public channels, ranked

1. **PayPal AI Hackathon Discord** linked in the Devpost resources/FAQ. Most suitable for technical questions about Sandbox, PayPal tools, or judging clarification; keep questions public where possible.
2. **Devpost FAQ/discussion threads and webinar recordings.** The official updates page links the Oct 6 PayPal build session and tools. Search and read existing Q&A before posting.
3. **PayPal Developer public GitHub issues/discussions.** For a concrete SDK/toolkit bug or question, use the relevant repo and include a minimal reproduction without credentials.
4. **A concise LinkedIn note to DevRel** after having a working demo. The PayPal Developer account, Jo, Eddie, and Marco are more appropriate than executive cold outreach. Ask for feedback, not special treatment.
5. **Official PayPal developer support channels** for account/Sandbox access issues. Do not post secrets or customer details in community chat.

## Suggested message (edit to fit your voice)

> Hi [Name] — I’m building PayPilot for the PayPal AI Hackathon. It focuses on a narrow reliability case: after a PayPal capture request times out or webhooks arrive late, a merchant needs evidence of the provider outcome before any retry. The prototype uses PayPal Sandbox Orders/capture and verified event history; AI summarizes only backend-validated evidence. Would you be open to one piece of feedback on whether the demo makes PayPal's role and the merchant value clear? Here’s the short demo/repo: [link]. No worries if you’re busy.

Only send after the linked materials are functional. Do not send this to multiple people at once, imply they endorsed the project, or ask a judge to evaluate it privately. If a conversation could create a conflict-of-interest appearance, keep it to public technical feedback.

## What not to say

- “I found a gap in PayPal's internal system.” We do not know that.
- “PayPal should acquire/adopt this.” Ask for critique of the prototype instead.
- “AI prevents fraud/chargebacks.” The project neither has data nor authority to make that claim.
- “Our model decides whether a payment succeeded.” It must not; verified PayPal state and deterministic backend rules do.
