export default function Home() {
  return (
    <div
      className="card"
      style={{
        display: "grid",
        gap: "1rem",
        maxWidth: 760,
        margin: "0 auto",
        textAlign: "left",
      }}
    >
      <h2 style={{ marginTop: 0 }}>Welcome to Yivi International Airlines</h2>
      <p>
        This simple demo lets you “buy” a ticket with your name and
        passport/document number, then verify your identity with Yivi to display
        a boarding pass. It’s a minimal flow to showcase privacy-preserving
        verification.
      </p>
      <ul style={{ margin: 0, paddingLeft: "1.25rem" }}>
        <li>Buy: create a fake ticket with basic info.</li>
        <li>
          Verify: scan a QR in Yivi to complete check-in and get a boarding
          pass.
        </li>
      </ul>
      <div
        style={{
          display: "flex",
          gap: "0.75rem",
          flexWrap: "wrap",
          marginTop: "0.5rem",
        }}
      >
        <a href="#/buy">
          <button>Buy Ticket</button>
        </a>
        <a href="#/verify">
          <button>Verify & Boarding Pass</button>
        </a>
      </div>
    </div>
  );
}
