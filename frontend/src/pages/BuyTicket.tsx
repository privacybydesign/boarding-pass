import { useMemo, useState } from "react";
import { apiEndpoint } from "../config";

export default function BuyTicket() {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [documentNumber, setDocumentNumber] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string>("");

  const canSubmit = useMemo(
    () =>
      firstName.trim() !== "" &&
      lastName.trim() !== "" &&
      documentNumber.trim() !== "",
    [firstName, lastName, documentNumber]
  );

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!canSubmit || isSubmitting) return;

    setIsSubmitting(true);
    setError("");
    const payload = {
      firstName: firstName.trim(),
      lastName: lastName.trim(),
      documentNumber: documentNumber.trim(),
    };

    try {
      const response = await fetch(`${apiEndpoint}/tickets`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const message = await response.text();
        throw new Error(
          message || `Failed to create ticket (status ${response.status})`
        );
      }

      const ticket = await response.json();
      try {
        localStorage.setItem("ticketId", ticket.id);
        localStorage.setItem("ticket", JSON.stringify(ticket));
      } catch {}
      window.location.hash = "#/verify";
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to create ticket.";
      setError(message);
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div
      className="card"
      style={{ maxWidth: 520, margin: "0 auto", textAlign: "left" }}
    >
      <h2 style={{ marginTop: 0 }}>Buy a Ticket</h2>
      <p style={{ marginTop: 0, color: "#888" }}>
        Enter basic details for a demo ticket.
      </p>
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
        {error && (
          <div
            style={{
              color: "#b91c1c",
              background: "#fef2f2",
              border: "1px solid #fecaca",
              padding: "0.75rem 1rem",
              borderRadius: 8,
            }}
          >
            {error}
          </div>
        )}
        <button type="submit" disabled={!canSubmit || isSubmitting}>
          {isSubmitting ? "Processing…" : "Continue to verification"}
        </button>
      </form>
    </div>
  );
}
