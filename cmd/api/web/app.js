const $ = (selector, root = document) => root.querySelector(selector);
const $$ = (selector, root = document) => [...root.querySelectorAll(selector)];

const demoBase = [
  { id: "PP-2026-0142", item: "Northstar ceramic desk lamp", amount: 8400, currency: "USD", state: "unknown", customer: "Maya Chen", localID: "ord_71K9•••2F8", providerID: "8BC4•••1X29", attemptID: "att_3d18•••a07", eventID: "evt_9f2a•••83c", last: "2 min ago", orderStatus: "pending", attemptStatus: "unknown", events: [] },
  { id: "PP-2026-0140", item: "Linen duvet cover · Queen", amount: 12900, currency: "USD", state: "pending", customer: "Jordan Lee", localID: "ord_83A1•••6K2", providerID: "5LU7•••9C14", attemptID: "att_61cd•••e5b", eventID: "evt_—", last: "18 min ago", orderStatus: "pending", attemptStatus: "created", events: [] },
  { id: "PP-2026-0138", item: "Walnut monitor stand", amount: 6800, currency: "USD", state: "review", customer: "Ari Patel", localID: "ord_02F4•••1M8", providerID: "2VE3•••A801", attemptID: "att_90ad•••23c", eventID: "evt_6a14•••d91", last: "1 hr ago", orderStatus: "review_required", attemptStatus: "failed", events: [] },
];

const state = {
  mode: "demo",
  scenario: "ready",
  selectedId: demoBase[0].id,
  demoOrders: structuredClone(demoBase),
  liveOrders: [],
  liveSelected: null,
  captureReadyID: "",
  analysis: null,
  running: false,
  filter: "attention",
};

const glyph = {
  check: '<svg viewBox="0 0 14 14"><path d="m3 7 2.5 2.5L11 4"/></svg>',
};

function escapeHTML(value = "") {
  return String(value).replace(/[&<>"']/g, (char) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" })[char]);
}

function money(minor, currency = "USD") {
  const amount = Number(minor || 0) / 100;
  return new Intl.NumberFormat("en-US", { style: "currency", currency, minimumFractionDigits: 2 }).format(amount);
}

function currentOrders() {
  return state.mode === "demo" ? state.demoOrders : state.liveOrders;
}

function selectedOrder() {
  return currentOrders().find((order) => order.id === state.selectedId) || currentOrders()[0] || null;
}

function statusInfo(status) {
  const normalized = String(status || "unknown").toLowerCase();
  if (["paid", "completed", "success", "captured"].includes(normalized)) return { label: "COMPLETED", cls: "completed" };
  if (["review", "review_required", "failed"].includes(normalized)) return { label: normalized === "failed" ? "FAILED" : "REVIEW", cls: normalized === "failed" ? "failed" : "review" };
  if (["created", "pending", "approved"].includes(normalized)) return { label: "PENDING", cls: "pending" };
  return { label: "UNKNOWN", cls: "unknown" };
}

function ordersForRender() {
  let orders = currentOrders();
  if (state.filter === "attention") orders = orders.filter((order) => !["paid", "completed"].includes(String(order.state || order.status).toLowerCase()));
  const query = $("#search-input").value.trim().toLowerCase();
  if (query) orders = orders.filter((order) => [order.id, order.item, order.description, order.merchant_reference].some((value) => String(value || "").toLowerCase().includes(query)));
  return orders;
}

function renderList() {
  const list = $("#payment-list");
  const orders = ordersForRender();
  if (!orders.length) {
    list.innerHTML = '<div class="queue-empty">No payments match this view.</div>';
    return;
  }
  list.innerHTML = orders.map((order) => {
    const stateValue = order.state || order.status;
    const badge = statusInfo(stateValue);
    const dot = badge.cls === "unknown" ? "" : badge.cls;
    const reference = order.id || order.merchant_reference || "Payment";
    const item = order.item || order.description || "PayPal payment";
    const amount = order.amount ?? order.amount_minor;
    const formatted = typeof amount === "string" && amount.includes("$") ? amount : money(amount, order.currency || "USD");
    return `<button class="payment-row ${reference === state.selectedId ? "selected" : ""}" role="option" aria-selected="${reference === state.selectedId}" data-order-id="${escapeHTML(reference)}"><span class="payment-primary"><span class="payment-reference"><i class="payment-dot ${dot}"></i>${escapeHTML(reference)}</span><span class="payment-name">${escapeHTML(item)}</span></span><span class="payment-amount">${escapeHTML(formatted)}</span><span><span class="status-badge ${badge.cls}">${badge.label}</span></span><span class="last-event">${escapeHTML(order.last || relativeTime(order.updated_at || order.created_at))}</span><span class="row-arrow">›</span></button>`;
  }).join("");
  $("#queue-count").textContent = String(orders.length).padStart(2, "0");
  $("#queue-footer-text").textContent = `Showing ${orders.length} ${state.filter === "attention" ? "payments needing attention" : "payments"}`;
  $("#nav-payment-count").textContent = String(currentOrders().filter((order) => !["paid", "completed"].includes(String(order.state || order.status).toLowerCase())).length).padStart(2, "0");
}

function relativeTime(value) {
  if (!value) return "just now";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "recently";
  const minutes = Math.max(0, Math.floor((Date.now() - date.getTime()) / 60000));
  if (minutes < 1) return "just now";
  if (minutes < 60) return `${minutes} min ago`;
  const hours = Math.floor(minutes / 60);
  return `${hours} hr${hours > 1 ? "s" : ""} ago`;
}

function timelineEvents(order) {
  if (state.mode === "demo") {
    const events = [
      { title: "Agent purchase intent received", detail: "Mandate MD-2091 · max spend $90.00", time: "10:42:03", type: "success", ref: "intent_2bd1•••a61" },
      { title: "PayPal order approved", detail: `Order ${order.providerID || "8BC4•••1X29"} · ${money(order.amount, order.currency)}`, time: "10:42:11", type: "success", ref: order.providerID },
      { title: "Capture request sent", detail: `Idempotency key cap_${String(order.id).slice(-4).toLowerCase()}`, time: "10:42:14", type: "current", ref: order.attemptID },
    ];
    if (["capturing", "lost", "event", "verified", "complete"].includes(state.scenario) && order.id === demoBase[0].id) {
      events.push({ title: "Merchant request timed out", detail: "No response received within 30 seconds", time: "10:42:44", type: "alert", ref: "trace_71d•••09f" });
    }
    if (["event", "verified", "complete"].includes(state.scenario) && order.id === demoBase[0].id) {
      events.push({ title: "Completion event received · processing delayed", detail: "PAYMENT.CAPTURE.COMPLETED · signature verified", time: "10:43:02", type: "current", ref: order.eventID });
    }
    if (["verified", "complete"].includes(state.scenario) && order.id === demoBase[0].id) {
      events.push({ title: "Provider status refreshed", detail: "PayPal API reports COMPLETED · capture IDs match", time: "10:43:19", type: "success", ref: "refresh_4c3a•••7df" });
    }
    if (state.scenario === "complete" && order.id === demoBase[0].id) {
      events.push({ title: "Merchant order reconciled", detail: "Local state updated once · no second capture", time: "10:43:20", type: "success", ref: order.localID });
    }
    return events;
  }
  const events = [];
  if (order.created_at) events.push({ title: "Merchant order created", detail: order.description || "PayPal checkout", time: formatClock(order.created_at), type: "success", ref: shortId(order.id), timestamp: order.created_at });
  for (const attempt of state.liveSelected?.attempts || []) {
    events.push({ title: `Payment attempt ${attempt.status}`, detail: `PayPal order ${shortId(attempt.paypal_order_id)}`, time: formatClock(attempt.created_at), type: attempt.status === "completed" ? "success" : attempt.status === "unknown" ? "alert" : "current", ref: shortId(attempt.id), timestamp: attempt.created_at });
    if (attempt.updated_at && attempt.updated_at !== attempt.created_at) events.push({ title: `Attempt updated · ${attempt.status}`, detail: `Capture ${shortId(attempt.paypal_capture_id)}`, time: formatClock(attempt.updated_at), type: attempt.status === "completed" ? "success" : "current", ref: shortId(attempt.paypal_capture_id), timestamp: attempt.updated_at });
  }
  for (const event of state.liveSelected?.providerEvents || []) {
    const amount = event.amount_minor ? ` · ${money(event.amount_minor, event.currency || order.currency || "USD")}` : "";
    const providerStatus = event.resource_status ? ` · ${event.resource_status}` : "";
    const processing = event.processing_status === "failed" ? "Not applied; review required" : `Verified · ${event.processing_status}`;
    events.push({
      title: `PayPal event · ${event.event_type}`,
      detail: `${processing}${providerStatus}${amount}`,
      time: formatClock(event.received_at),
      type: event.processing_status === "failed" ? "alert" : event.event_type.endsWith("COMPLETED") ? "success" : "current",
      ref: shortId(event.provider_event_id),
      timestamp: event.received_at,
    });
  }
  return events.sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
}

function shortId(value) {
  if (!value) return "—";
  const text = String(value);
  return text.length > 16 ? `${text.slice(0, 6)}•••${text.slice(-4)}` : text;
}

function formatClock(value) {
  if (!value) return "—";
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? "—" : d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false });
}

function renderTimeline(order) {
  const events = timelineEvents(order);
  $("#timeline").innerHTML = events.map((event) => `<div class="timeline-event ${event.type}"><span class="event-dot"></span><div class="event-copy"><div class="event-title">${escapeHTML(event.title)}</div><div class="event-detail">${escapeHTML(event.detail)}${event.ref ? ` · ${escapeHTML(event.ref)}` : ""}</div></div><span class="event-time">${escapeHTML(event.time)}</span></div>`).join("") || '<div class="queue-empty">No events recorded yet.</div>';
  $("#timeline-subtitle").textContent = `${events.length} events · latest ${state.mode === "demo" ? (state.scenario === "ready" ? "2 min ago" : "just now") : "from API"}`;
}

function renderEvidence(order) {
  const callout = $("#evidence-callout");
  const summary = $("#evidence-summary");
  const links = $("#evidence-links");
  if (state.mode === "demo" && order.id === demoBase[0].id) {
    if (["verified", "complete"].includes(state.scenario)) {
      summary.textContent = "PayPal reports COMPLETED and the capture ID matches the stored event. The local order is reconciled; no second capture is needed.";
      links.innerHTML = `<button data-evidence="${escapeHTML(order.eventID)}">EVENT <span>${escapeHTML(order.eventID)}</span></button><button data-evidence="${escapeHTML(order.attemptID)}">ATTEMPT <span>${escapeHTML(order.attemptID)}</span></button><button data-evidence="refresh_4c3a•••7df">REFRESH <span>refresh_4c3a•••7df</span></button>`;
      $(".confidence-pill").textContent = "Confirmed in simulation";
      $("#detail-state").textContent = "COMPLETED";
      $("#detail-state").className = "detail-state state-complete";
      callout.classList.add("resolved-evidence");
    } else if (["event", "lost", "capturing"].includes(state.scenario)) {
      summary.textContent = "The capture response was lost. A signature-verified completion event is waiting to be reconciled. Refresh provider state before retrying.";
      links.innerHTML = `<button data-evidence="${escapeHTML(order.eventID)}">EVENT <span>${escapeHTML(order.eventID)}</span></button><button data-evidence="${escapeHTML(order.attemptID)}">ATTEMPT <span>${escapeHTML(order.attemptID)}</span></button>`;
      $(".confidence-pill").textContent = "2 records · consistent";
      $("#detail-state").textContent = "UNKNOWN";
      $("#detail-state").className = "detail-state state-unknown";
      callout.classList.remove("resolved-evidence");
    } else {
      summary.textContent = "A capture timeout is ambiguous: it does not prove the payment failed. Check provider state before attempting another capture.";
      $(".confidence-pill").textContent = "High · outcome unknown";
      links.innerHTML = `<button data-evidence="${escapeHTML(order.attemptID)}">ATTEMPT <span>${escapeHTML(order.attemptID)}</span></button>`;
      $("#detail-state").textContent = "UNKNOWN";
      $("#detail-state").className = "detail-state state-unknown";
      callout.classList.remove("resolved-evidence");
    }
    return;
  }
  if (state.mode === "demo") {
    summary.textContent = order.state === "review" ? "A prior attempt failed. Review the provider response and payment history before deciding whether a new checkout is appropriate." : "This order is awaiting payer approval. No capture request has been made.";
    links.innerHTML = `<button data-evidence="${escapeHTML(order.attemptID)}">ATTEMPT <span>${escapeHTML(order.attemptID)}</span></button>`;
    $(".confidence-pill").textContent = "Record status";
    $("#detail-state").textContent = statusInfo(order.state).label;
    $("#detail-state").className = `detail-state state-${statusInfo(order.state).cls}`;
    callout.classList.remove("resolved-evidence");
    return;
  }
  const report = state.liveSelected?.investigation?.report;
  if (report) {
    summary.textContent = report.summary || report.finding || "Investigation loaded from PayPilot API.";
    links.innerHTML = (report.evidence || []).slice(0, 3).map((item) => `<button data-evidence="${escapeHTML(item.id)}">${escapeHTML(item.record_type || "RECORD").toUpperCase()} <span>${escapeHTML(shortId(item.record_id || item.id))}</span></button>`).join("") || "";
    $(".confidence-pill").textContent = `${Math.round((report.confidence || 0) * 100)}% confidence`;
    $("#detail-state").textContent = statusInfo(order.status).label;
    $("#detail-state").className = `detail-state state-${statusInfo(order.status).cls}`;
  } else {
    summary.textContent = "Run an investigation to load backend-validated evidence and a safe next-step recommendation.";
    links.innerHTML = "";
    $(".confidence-pill").textContent = "Not investigated";
    $("#detail-state").textContent = statusInfo(order.status).label;
    $("#detail-state").className = `detail-state state-${statusInfo(order.status).cls}`;
  }
  callout.classList.remove("resolved-evidence");
}

function renderDetail() {
  const order = selectedOrder();
  if (!order) {
    $("#detail-reference").textContent = "No payment selected";
    return;
  }
  const demo = state.mode === "demo";
  const reference = order.id || order.merchant_reference || "Payment";
  const amount = order.amount ?? order.amount_minor;
  $("#detail-reference").textContent = reference;
  $("#detail-amount").textContent = typeof amount === "string" && amount.includes("$") ? amount : money(amount, order.currency || "USD");
  $("#detail-local-id").textContent = demo ? order.localID : shortId(order.id);
  const attempt = state.liveSelected?.attempts?.[0];
  $("#detail-provider-id").textContent = demo ? order.providerID : shortId(attempt?.paypal_order_id);
  $("#detail-state").textContent = statusInfo(demo ? order.state : order.status).label;
  $("#detail-state").className = `detail-state state-${statusInfo(demo ? order.state : order.status).cls}`;
  renderTimeline(order);
  renderEvidence(order);
  if (demo && reference === demoBase[0].id && state.scenario !== "ready" && state.scenario !== "verified" && state.scenario !== "complete") {
    $("#refresh-label").textContent = "Verify with PayPal";
  } else if (demo && ["verified", "complete"].includes(state.scenario)) {
    $("#refresh-label").textContent = "Provider verified";
  } else if (!demo && attempt?.status === "completed") {
    $("#refresh-label").textContent = "Payment completed";
  } else if (!demo && attempt?.status === "created" && state.captureReadyID === attempt.id) {
    $("#refresh-label").textContent = "Capture approved payment";
  } else {
    $("#refresh-label").textContent = demo ? "Verify with PayPal" : attempt?.status === "created" ? "Check approval with PayPal" : "Refresh provider status";
  }
  $("#refresh-button").disabled = (demo && ["ready", "complete", "verified"].includes(state.scenario) && reference === demoBase[0].id) || (!demo && attempt?.status === "completed");
  $("#investigate-button").disabled = false;
}

function renderMetrics() {
  const orders = currentOrders();
  if (state.mode === "demo") {
    const needs = orders.filter((order) => !["paid", "completed"].includes(order.state)).length;
    $("#metric-attention").textContent = String(needs).padStart(2, "0");
    $("#queue-count").textContent = String(ordersForRender().length).padStart(2, "0");
  } else {
    const needs = orders.filter((order) => !["paid", "completed"].includes(String(order.status).toLowerCase())).length;
    $("#metric-attention").textContent = String(needs).padStart(2, "0");
  }
}

function renderAll() {
  renderMetrics();
  renderList();
  renderDetail();
}

function setProgress(step) {
  const order = ["intent", "checkout", "capture", "reconcile"];
  const current = order.indexOf(step);
  $$(".progress-step").forEach((element) => {
    const index = order.indexOf(element.dataset.progress);
    element.classList.toggle("done", index < current || step === "complete");
    element.classList.toggle("active", index === current && step !== "complete");
  });
}

function delay(ms) { return new Promise((resolve) => setTimeout(resolve, ms)); }

async function runScenario() {
  if (state.running) return;
  if (state.scenario !== "ready") resetScenario(false);
  state.running = true;
  const button = $("#run-scenario");
  button.disabled = true;
  button.innerHTML = '<span class="play-icon">◌</span> Running transaction trace…';
  $("#journey-card").classList.add("scenario-running");
  const steps = [
    { phase: "intent", visual: "INTENT SIGNED", caption: "Agent intent scoped to a $90 spend limit.", toast: "Mandate MD-2091 validated. Amount ceiling: $90.00." },
    { phase: "checkout", visual: "ORDER APPROVED", caption: "PayPal Sandbox simulation approved an $84.00 order.", toast: "Checkout approved. The simulated PayPal order is 8BC4•••1X29." },
    { phase: "capture", visual: "RESPONSE LOST", caption: "The capture may have succeeded. Local state stays UNKNOWN.", toast: "Simulated 30-second timeout. PayPilot pauses retry; it does not mark this failed." },
    { phase: "event", visual: "EVENT ARRIVED", caption: "A completion event is verified but not yet reconciled.", toast: "Late simulated completion event received. Local order is still pending." },
  ];
  for (const step of steps) {
    await delay(1450);
    if (step.phase === "capture") state.scenario = "lost";
    if (step.phase === "event") state.scenario = "event";
    if (step.phase === "checkout") state.scenario = "capturing";
    setProgress(step.phase === "event" ? "capture" : step.phase);
    $("#visual-state").textContent = step.visual;
    $("#scenario-caption").textContent = step.caption;
    renderAll();
    toast(step.toast);
  }
  state.running = false;
  $("#journey-card").classList.remove("scenario-running");
  button.disabled = false;
  button.innerHTML = '<span class="play-icon">↻</span> Replay scenario';
  toast("Now investigate or verify the provider state to complete the trace.");
}

function resetScenario(showToast = true) {
  state.scenario = "ready";
  state.demoOrders = structuredClone(demoBase);
  state.analysis = null;
  state.running = false;
  state.selectedId = demoBase[0].id;
  $("#journey-card").classList.remove("scenario-running", "scenario-complete");
  $("#visual-state").textContent = "READY";
  $("#scenario-caption").textContent = "~15 sec · see each system handoff";
  const button = $("#run-scenario");
  button.disabled = false;
  button.innerHTML = '<span class="play-icon">▶</span> Run payment scenario';
  setProgress("intent");
  renderAll();
  if (showToast) toast("Scenario reset. No provider calls were made.");
}

function completeDemoReconciliation() {
  const order = state.demoOrders.find((item) => item.id === demoBase[0].id);
  if (!order || !["event", "lost", "capturing"].includes(state.scenario)) {
    toast("Run the scenario first so there is an ambiguous capture to reconcile.");
    return;
  }
  state.scenario = "verified";
  state.scenario = "complete";
  order.state = "completed";
  order.orderStatus = "paid";
  order.attemptStatus = "completed";
  order.last = "just now";
  state.analysis = {
    source: "simulation",
    summary: "PayPal's simulated provider response confirms a completed $84.00 capture. The amount and capture references match. The local order is reconciled; a second capture is blocked.",
    facts: ["Capture request timed out in the merchant application.", "Provider refresh reports COMPLETED for the same PayPal order.", "Amount and currency match the local order: USD 84.00."],
    evidence: [order.attemptID, order.eventID, "refresh_4c3a•••7df"],
  };
  setProgress("complete");
  $("#journey-card").classList.add("scenario-complete");
  $("#visual-state").textContent = "RECONCILED";
  $("#scenario-caption").textContent = "One confirmed capture · no duplicate attempt · merchant record in sync.";
  $("#run-scenario").innerHTML = '<span class="play-icon">✓</span> Scenario complete';
  $("#run-scenario").disabled = false;
  renderAll();
  toast("Provider truth confirmed in simulation. Order reconciled once; no retry issued.");
}

async function investigateDemo() {
  const order = selectedOrder();
  const isMain = order?.id === demoBase[0].id;
  const result = $("#investigation-result");
  result.classList.remove("hidden", "error");
  result.innerHTML = "<strong>Checking the simulated record set…</strong>Applying deterministic amount, attempt, event, and state checks.";
  $("#run-investigation").disabled = true;
  await delay(850);
  if (isMain && ["event", "verified", "complete"].includes(state.scenario)) {
    result.innerHTML = `<strong>Outcome is recoverable · retry is not permitted</strong><div>Confirmed in this simulation: the capture request timed out, a completion event arrived, and provider refresh reports <b>COMPLETED</b>. The $84.00 USD amount matches. Reconcile the local order; do not create another capture.<br><br>Evidence: <span class="evidence-ref">${escapeHTML(order.attemptID)}</span> · <span class="evidence-ref">${escapeHTML(order.eventID)}</span> · <span class="evidence-ref">refresh_4c3a•••7df</span><br><br><span style="color:#8a958b">Simulation rules executed locally; no LLM or PayPal endpoint was called.</span></div>`;
  } else if (isMain && ["lost", "capturing"].includes(state.scenario)) {
    result.innerHTML = `<strong>Payment outcome remains unknown</strong><div>The merchant request timed out. There is not yet enough evidence to say that capture succeeded or failed. Safest next step: refresh provider state; do not retry yet.<br><br>Evidence: <span class="evidence-ref">${escapeHTML(order.attemptID)}</span><br><br><span style="color:#8a958b">Simulation rules executed locally; no LLM or PayPal endpoint was called.</span></div>`;
  } else if (isMain) {
    result.innerHTML = `<strong>Insufficient provider evidence to decide</strong><div>A capture attempt is in an unknown state. Refresh the provider before retrying.<br><br>Evidence: <span class="evidence-ref">${escapeHTML(order.attemptID)}</span><br><br><span style="color:#8a958b">Simulation rules executed locally; no LLM or PayPal endpoint was called.</span></div>`;
  } else {
    result.innerHTML = `<strong>${order.state === "review" ? "Manual review recommended" : "Payment is awaiting buyer approval"}</strong><div>${order.state === "review" ? "The previous attempt is marked failed. Inspect the provider response before starting another checkout." : "No capture attempt has been made. Wait for buyer approval before requesting capture."}<br><br>Evidence: <span class="evidence-ref">${escapeHTML(order.attemptID)}</span><br><br><span style="color:#8a958b">Simulation rules executed locally; no LLM or PayPal endpoint was called.</span></div>`;
  }
  $("#run-investigation").disabled = false;
  $("#run-investigation").innerHTML = "Run again <span>→</span>";
}

async function apiRequest(path, options = {}) {
  const response = await fetch(path, { ...options, headers: { ...(options.headers || {}), ...(options.body ? { "Content-Type": "application/json" } : {}) } });
  const text = await response.text();
  let body;
  try { body = text ? JSON.parse(text) : {}; } catch { body = { message: text }; }
  if (!response.ok) {
    const message = body?.error?.message || body?.message || `Request failed (${response.status})`;
    throw new Error(message);
  }
  return body;
}

async function loadLiveOrders() {
  setConnection("Connecting to PayPilot API…", "amber");
  try {
    const result = await apiRequest("/api/v1/orders");
    state.liveOrders = (result.data || []).map((order) => ({ ...order, id: order.id, item: order.description, amount: order.amount_minor, state: order.status, last: relativeTime(order.updated_at || order.created_at) }));
    if (!state.liveOrders.some((order) => order.id === state.selectedId)) state.selectedId = state.liveOrders[0]?.id || "";
    state.liveSelected = null;
    setConnection("PayPilot API connected", "green");
    $("#mode-banner").className = "mode-banner live";
    $("#mode-banner span:last-of-type").innerHTML = "<strong>Live API mode</strong> — this view reads the connected PayPilot API. Provider actions use your server configuration.";
    $("#mode-banner .banner-link").textContent = "Safety notes ↗";
    renderAll();
    if (state.selectedId) await loadLiveDetail(state.selectedId);
  } catch (error) {
    setConnection("API unavailable", "red");
    $("#mode-banner").className = "mode-banner error";
    $("#mode-banner span:last-of-type").innerHTML = `<strong>Could not connect to the PayPilot API</strong> — ${escapeHTML(error.message)}. Start the Go API and PostgreSQL, or switch back to Demo.`;
    $("#mode-banner .banner-link").textContent = "Connection help ↗";
    state.liveOrders = [];
    state.liveSelected = null;
    renderAll();
    toast("Live API is unavailable. Demo mode remains ready.", "error");
  }
}

async function loadLiveDetail(id) {
  if (!id) return;
  try {
    const [order, attempts, investigations, providerEvents] = await Promise.all([
      apiRequest(`/api/v1/orders/${encodeURIComponent(id)}`),
      apiRequest(`/api/v1/orders/${encodeURIComponent(id)}/payments`),
      apiRequest(`/api/v1/orders/${encodeURIComponent(id)}/investigations`),
      apiRequest(`/api/v1/orders/${encodeURIComponent(id)}/events`),
    ]);
    let latest = null;
    if (investigations.data?.length) {
      const first = investigations.data[0];
      try { latest = await apiRequest(`/api/v1/investigations/${encodeURIComponent(first.id)}`); } catch { latest = null; }
    }
    state.liveSelected = { ...order, attempts: attempts.data || [], providerEvents: providerEvents.data || [], investigations: investigations.data || [], investigation: latest };
    renderDetail();
  } catch (error) {
    toast(`Could not load payment detail: ${error.message}`, "error");
  }
}

function setConnection(label, color) {
  $("#connection-label").textContent = label;
  const dot = $(".connection-dot");
  dot.style.background = color === "green" ? "#4b9a67" : color === "red" ? "#c76d53" : "#df9f3f";
  dot.style.boxShadow = color === "green" ? "0 0 0 3px #4b9a6720" : color === "red" ? "0 0 0 3px #c76d5320" : "0 0 0 3px #df9f3f20";
}

async function refreshSelected() {
  if (state.mode === "demo") {
    completeDemoReconciliation();
    return;
  }
  const order = selectedOrder();
  const attempt = state.liveSelected?.attempts?.[0];
  if (!order) return;
  if (attempt?.status === "created" && state.captureReadyID === attempt.id) {
    const button = $("#refresh-button");
    button.disabled = true;
    try {
      const result = await apiRequest(`/api/v1/payments/${encodeURIComponent(attempt.id)}/capture`, { method: "POST", headers: { "Idempotency-Key": crypto.randomUUID() }, body: "{}" });
      state.captureReadyID = "";
      toast(`PayPal capture: ${result.status || "received"}.`);
      await loadLiveOrders();
    } catch (error) {
      state.captureReadyID = "";
      toast(error.message, "error");
      await loadLiveOrders();
      await loadLiveDetail(order.id);
    } finally {
      button.disabled = false;
    }
    return;
  }
  if (!attempt) {
    toast("Create a PayPal checkout before refreshing provider state.");
    return;
  }
  const button = $("#refresh-button");
  button.disabled = true;
  try {
    const result = await apiRequest(`/api/v1/payments/${encodeURIComponent(attempt.id)}/refresh`);
    state.captureReadyID = result.provider_status === "APPROVED" && result.status === "created" ? attempt.id : "";
    if (result.provider_status === "APPROVED") {
      toast("PayPal confirms buyer approval. Capture is now available as a separate action.");
    } else {
      toast(`Provider refresh: ${result.provider_status || result.status || "complete"}.`);
    }
    await loadLiveOrders();
    await loadLiveDetail(order.id);
  } catch (error) {
    toast(error.message, "error");
  } finally {
    button.disabled = false;
  }
}

async function runLiveInvestigation() {
  const order = selectedOrder();
  if (!order) return;
  const output = $("#investigation-result");
  output.classList.remove("hidden", "error");
  output.innerHTML = "<strong>Loading scoped payment records…</strong>Running read-only investigation through the PayPilot API.";
  $("#run-investigation").disabled = true;
  try {
    const result = await apiRequest("/api/v1/investigations", {
      method: "POST",
      headers: { "Idempotency-Key": crypto.randomUUID() },
      body: JSON.stringify({ order_id: order.id, question: $("#investigation-question").value.trim() || "Did this payment complete, and is it safe to retry?" }),
    });
    const report = result.report || {};
    const references = (report.evidence || []).slice(0, 4).map((item) => `<span class="evidence-ref">${escapeHTML(item.record_type)}:${escapeHTML(shortId(item.record_id))}</span>`).join(" · ");
    output.innerHTML = `<strong>${escapeHTML(report.finding || "Investigation complete")}</strong><div>${escapeHTML(report.summary || report.explanation?.text || "The report has no summary.")}<br><br><b>Recommendation:</b> ${escapeHTML(report.recommendation || "manual_review")}<br><b>Confidence:</b> ${Math.round((report.confidence || 0) * 100)}% · ${escapeHTML(report.confidence_explanation || "See evidence references.")}<br><br>${references || "No evidence references returned."}<br><br><span style="color:#8a958b">Source: ${escapeHTML(report.analysis_source || "deterministic")} · model status: ${escapeHTML(report.model_status || "not reported")}</span></div>`;
    await loadLiveDetail(order.id);
    toast("Investigation saved by the connected PayPilot API.");
  } catch (error) {
    output.classList.add("error");
    output.innerHTML = `<strong>Investigation unavailable</strong><div>${escapeHTML(error.message)}</div>`;
    toast(error.message, "error");
  } finally {
    $("#run-investigation").disabled = false;
  }
}

async function createLiveOrder() {
  const description = $("#order-description").value.trim();
  const amount = Math.round(Number($("#order-amount").value) * 100);
  const reference = $("#order-reference").value.trim();
  const resultBox = $("#order-result");
  resultBox.classList.remove("hidden", "error");
  if (!description || !Number.isFinite(amount) || amount <= 0 || !reference) {
    resultBox.classList.add("error");
    resultBox.innerHTML = "Enter a description, valid positive amount, and merchant reference.";
    return;
  }
  $("#create-order").disabled = true;
  try {
    const created = await apiRequest("/api/v1/orders", { method: "POST", headers: { "Idempotency-Key": crypto.randomUUID() }, body: JSON.stringify({ merchant_reference: reference, amount_minor: amount, currency: "USD", description }) });
    const checkout = await apiRequest(`/api/v1/orders/${encodeURIComponent(created.id)}/paypal-order`, { method: "POST", headers: { "Idempotency-Key": crypto.randomUUID() }, body: "{}" });
    resultBox.innerHTML = `<strong>PayPal checkout created</strong><div>Order ${escapeHTML(shortId(created.id))} · PayPal ${escapeHTML(shortId(checkout.paypal_order_id))}. ${checkout.approval_url ? `<a href="${escapeHTML(checkout.approval_url)}" target="_blank" rel="noreferrer">Open Sandbox approval ↗</a>` : "No approval URL was returned."}<br><br>After buyer approval, select the order and capture through this workbench.</div>`;
    await loadLiveOrders();
    state.selectedId = created.id;
    await loadLiveDetail(created.id);
  } catch (error) {
    resultBox.classList.add("error");
    resultBox.innerHTML = `<strong>Could not create checkout</strong><div>${escapeHTML(error.message)}</div>`;
  } finally {
    $("#create-order").disabled = false;
  }
}

function toast(message, type = "success") {
  const element = document.createElement("div");
  element.className = `toast ${type === "error" ? "error" : ""}`;
  element.textContent = message;
  $("#toast-region").append(element);
  setTimeout(() => element.remove(), 4200);
}

function openDialog(id) { const dialog = document.getElementById(id); if (dialog && !dialog.open) dialog.showModal(); }

function switchMode(mode) {
  state.mode = mode;
  $("#demo-mode").classList.toggle("selected", mode === "demo");
  $("#live-mode").classList.toggle("selected", mode === "live");
  $("#mode-banner").className = `mode-banner ${mode === "live" ? "live" : ""}`;
  if (mode === "demo") {
    $("#connection-label").textContent = "Demo environment";
    setConnection("Demo environment", "amber");
    $("#mode-banner span:last-of-type").innerHTML = "<strong>Guided simulation</strong> — every order and provider event below is synthetic. No money is moved.";
    $("#mode-banner .banner-link").innerHTML = 'How this works <span>↗</span>';
    state.liveSelected = null;
    state.selectedId = demoBase[0].id;
    renderAll();
  } else {
    $("#connection-label").textContent = "Connecting to PayPilot API…";
    $("#mode-banner span:last-of-type").innerHTML = "<strong>Live API mode</strong> — connecting to this PayPilot server. No synthetic fallback is shown as live data.";
    $("#mode-banner .banner-link").innerHTML = 'Safety notes <span>↗</span>';
    loadLiveOrders();
  }
}

function bindEvents() {
  $("#demo-mode").addEventListener("click", () => switchMode("demo"));
  $("#live-mode").addEventListener("click", () => switchMode("live"));
  $("#run-scenario").addEventListener("click", () => state.scenario === "complete" ? resetScenario() : runScenario());
  $("#reset-scenario").addEventListener("click", () => resetScenario());
  $("#payment-list").addEventListener("click", async (event) => {
    const row = event.target.closest("[data-order-id]");
    if (!row) return;
    state.selectedId = row.dataset.orderId;
    state.liveSelected = null;
    renderAll();
    if (state.mode === "live") await loadLiveDetail(state.selectedId);
  });
  $("#search-input").addEventListener("input", renderList);
  $("#filter-button").addEventListener("click", () => {
    state.filter = state.filter === "attention" ? "all" : "attention";
    $("#filter-button").innerHTML = `<svg viewBox="0 0 20 20"><path d="M3 5h14M5.5 10h9M8 15h4"/></svg>${state.filter === "all" ? "All payments" : "Filters"} <span class="filter-count">${state.filter === "all" ? "×" : "2"}</span>`;
    renderList();
  });
  $("#view-all-button").addEventListener("click", () => {
    state.filter = "all";
    $("#filter-button").innerHTML = '<svg viewBox="0 0 20 20"><path d="M3 5h14M5.5 10h9M8 15h4"/></svg>All payments <span class="filter-count">×</span>';
    renderList();
  });
  $("#refresh-button").addEventListener("click", refreshSelected);
  $("#investigate-button").addEventListener("click", () => {
    $("#investigation-result").classList.add("hidden");
    $("#run-investigation").innerHTML = state.mode === "demo" ? "Run simulation analysis <span>→</span>" : "Run investigation <span>→</span>";
    openDialog("investigation-modal");
  });
  $("#run-investigation").addEventListener("click", () => state.mode === "demo" ? investigateDemo() : runLiveInvestigation());
  $("#new-order-button").addEventListener("click", () => {
    if (state.mode === "demo") { toast("Switch to Live API to create a real PayPal Sandbox checkout."); return; }
    $("#order-result").classList.add("hidden"); openDialog("new-order-modal");
  });
  $("#create-order").addEventListener("click", createLiveOrder);
  $("#help-button").addEventListener("click", () => openDialog("help-modal"));
  $("#about-simulation").addEventListener("click", () => openDialog("help-modal"));
  $("#policy-info").addEventListener("click", () => toast("A network timeout cannot prove capture failed. Refresh provider state before considering another capture."));
  $("#evidence-links").addEventListener("click", (event) => {
    const button = event.target.closest("[data-evidence]");
    if (button) toast(`Evidence selected: ${button.dataset.evidence}. Each simulation reference is synthetic.`);
  });
  $$('[data-close]').forEach((button) => button.addEventListener("click", () => document.getElementById(button.dataset.close)?.close()));
  $$("dialog").forEach((dialog) => dialog.addEventListener("click", (event) => { if (event.target === dialog) dialog.close(); }));
  $$(".nav-link").forEach((link) => link.addEventListener("click", (event) => {
    event.preventDefault();
    const view = link.dataset.view;
    if (view === "payments") { state.filter = "all"; renderList(); $("#filter-button").innerHTML = '<svg viewBox="0 0 20 20"><path d="M3 5h14M5.5 10h9M8 15h4"/></svg>All payments <span class="filter-count">×</span>'; }
    const names = { overview: "Overview", payments: "Payments", investigations: "Investigations", timeline: "Event timeline", integrations: "Integrations", playbook: "Runbook" };
    $("#breadcrumb-current").textContent = names[view] || "Overview";
    $$(".nav-link").forEach((item) => item.classList.toggle("active", item === link));
    if (view === "investigations") openDialog("investigation-modal");
    if (view === "timeline") document.querySelector("#detail-card").scrollIntoView({ behavior: "smooth", block: "center" });
    if (view === "integrations") openDialog("help-modal");
    if (view === "playbook") toast("Runbook: verify provider state before retrying; model output never changes payment state.");
    if (view === "overview") window.scrollTo({ top: 0, behavior: "smooth" });
    if (window.innerWidth <= 640) { $("#sidebar").classList.remove("open"); $(".main-area").classList.remove("menu-open"); }
  }));
  $("#mobile-menu").addEventListener("click", () => { $("#sidebar").classList.toggle("open"); $(".main-area").classList.toggle("menu-open"); });
}

bindEvents();
renderAll();
