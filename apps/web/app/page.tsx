export default async function Home() {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";
  let health: any = null;

  try {
    const res = await fetch(`${apiBase}/healthz`, { cache: "no-store" });
    health = await res.json();
  } catch (e) {
    health = { status: "unknown", error: "API not reachable" };
  }

  return (
    <main style={{ padding: 24, fontFamily: "system-ui, sans-serif" }}>
      <h1>The Hive x Anti-Fraud</h1>
      <p>Reference UI scaffold (replace with your real design).</p>

      <h2>API health</h2>
      <pre>{JSON.stringify(health, null, 2)}</pre>

      <hr />
      <h2>Next steps</h2>
      <ul>
        <li>Implement report form + triage dashboard</li>
        <li>Add training modules + quiz</li>
        <li>Wire to OpenAPI at <code>packages/openapi/openapi.yaml</code></li>
      </ul>
    </main>
  );
}
