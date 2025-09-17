import { useState } from "react";

export default function BuyTicket() {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [documentNumber, setDocumentNumber] = useState("");

  function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    const ticket = { firstName, lastName, documentNumber };
    try {
      localStorage.setItem("ticket", JSON.stringify(ticket));
    } catch {}
    window.location.hash = "#/verify";
  }

  return (
    <div className="card" style={{ maxWidth: 520, margin: "0 auto", textAlign: "left" }}>
      <h2 style={{ marginTop: 0 }}>Buy a Ticket</h2>
      <p style={{ marginTop: 0, color: "#888" }}>Enter basic details for a demo ticket.</p>
      <form onSubmit={onSubmit} style={{ display: "grid", gap: "0.75rem" }}>
        <label style={{ display: "grid", gap: "0.35rem" }}>
          <span>First name</span>
          <input
            required
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            placeholder="Jane"
          />
        </label>
        <label style={{ display: "grid", gap: "0.35rem" }}>
          <span>Last name</span>
          <input
            required
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            placeholder="Doe"
          />
        </label>
        <label style={{ display: "grid", gap: "0.35rem" }}>
          <span>Passport / Document number</span>
          <input
            required
            value={documentNumber}
            onChange={(e) => setDocumentNumber(e.target.value)}
            placeholder="X1234567"
          />
        </label>
        <button type="submit">Continue to verification</button>
      </form>
    </div>
  );
}

