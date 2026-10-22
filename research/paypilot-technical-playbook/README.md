# PayPilot technical and business playbook

This folder teaches the current PayPilot architecture, the public PayPal integration model relevant to the prototype, and the proposed opportunity in engineering and business terms.

## Read in this order

1. [Architecture today and proposed target](architecture.md)
2. [Payment reliability concepts](payment-reliability.md)
3. [AI boundary and evidence design](ai-evidence.md)
4. [Business case and validation plan](business-case.md)
5. [Implementation and learning sequence](implementation-guide.md)

## Important boundary

“PayPal architecture” here means publicly documented APIs and behavior, not PayPal's private service topology, payment ledger internals, risk models, or operational systems. The project uses PayPal Sandbox and a small merchant-side database; it does not become part of PayPal's infrastructure.
