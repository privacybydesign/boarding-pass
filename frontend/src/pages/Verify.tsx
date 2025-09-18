/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useMemo, useState } from "react";
import { apiEndpoint } from "../config";

type Ticket = { firstName: string; lastName: string; documentNumber: string };

export default function Verify() {
  const [error, setError] = useState<string>("");
  const [sessionDone, setSessionDone] = useState(false);

  const ticket: Ticket | null = useMemo(() => {
    try {
      const raw = localStorage.getItem("ticket");
      return raw ? (JSON.parse(raw) as Ticket) : null;
    } catch {
      return null;
    }
  }, []);

  const ticketId: string | null = useMemo(() => {
    try {
      return localStorage.getItem("ticketId");
    } catch {
      return null;
    }
  }, []);

  useEffect(() => {
    if (!ticket || !ticketId) {
      setError("Ticket details missing. Please create a ticket again.");
      return;
    }

    let cancelled = false;
    let web: any;

    setSessionDone(false);
    setError("");

    import("@privacybydesign/yivi-frontend")
      .then((yivi: any) => {
        if (cancelled) return;

        const payload = {
          ticketId,
          firstName: ticket.firstName.trim(),
          lastName: ticket.lastName.trim(),
          documentNumber: ticket.documentNumber.trim().toUpperCase(),
        };

        web = yivi.newWeb({
          debugging: true,
          element: "#yivi-web-form",
          language: "en",

          session: {
            url: apiEndpoint,
            start: {
              url: (o: any) => `${o.url}/start`,
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify(payload),
            },
            result: {
              url: (o: any, { sessionPtr }: any) => {
                if (!sessionPtr || !sessionPtr.u)
                  return `${o.url}/result?sessionID=`;
                const sessionID = sessionPtr.u.split("/").pop();
                return `${o.url}/result?sessionID=${sessionID}`;
              },
              method: "GET",
            },
          },
        });

        web
          .start()
          .then(() => {
            setSessionDone(true);
          })
          .catch((e: any) => {
            console.error(e);
            setError("Verification failed to start.");
          });
      })
      .catch((e: any) => {
        console.error(e);
        setError("Unable to load verification client.");
      });

    return () => {
      cancelled = true;
      web?.abort?.();
    };
  }, [ticket, ticketId]);

  return (
    <div
      className="card"
      style={{
        display: "grid",
        gap: "1rem",
        maxWidth: 820,
        margin: "0 auto",
        width: "100%",
        height: 680,
        overflow: "hidden",
      }}
    >
      <h2 style={{ marginTop: 0 }}>Verify and Get Boarding Pass</h2>
      {!ticket && (
        <div
          style={{
            color: "#b45309",
            background: "#fff7ed",
            border: "1px solid #fed7aa",
            padding: "0.75rem 1rem",
            borderRadius: 8,
          }}
        >
          No ticket details found. Please create one on the Buy page.
        </div>
      )}

      <div style={{ display: "grid", gap: "1rem" }}>
        <div>
          <BoardingPass ticket={ticket} ready={sessionDone} />
        </div>
        <div style={{ display: "grid", gap: "0.5rem", placeItems: "center" }}>
          <p style={{ marginTop: 0, textAlign: "center" }}>
            Scan the QR code with your Yivi app to authenticate. If you don’t
            have the app, get it from{" "}
            <a href="https://yivi.app/#download">https://yivi.app/#download</a>
          </p>
          <div id="yivi-web-form" />
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
              <strong>Error:</strong> {error}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function BoardingPass({
  ticket,
  ready,
}: {
  ticket: Ticket | null;
  ready: boolean;
}) {
  const name = ticket ? `${ticket.firstName} ${ticket.lastName}` : "—";
  const doc = ticket ? ticket.documentNumber : "—";

  return (
    <div
      style={{
        borderRadius: 16,
        border: "1px solid #e5e7eb",
        background: "linear-gradient(135deg, #eff6ff, #ffffff)",
        color: "#111827",
        padding: "1rem",
      }}
    >
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          marginBottom: 8,
        }}
      >
        <div style={{ fontWeight: 700 }}>OpenSky Air</div>
        <div style={{ fontSize: 12, color: "#6b7280" }}>Demo Boarding Pass</div>
      </div>
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
        <Field label="Name" value={name} />
        <Field label="Document" value={doc} />
        <Field label="From" value="AMS" />
        <Field label="To" value="BCN" />
        <Field label="Flight" value="OS123" />
        <Field label="Seat" value="12A" />
      </div>
      <div
        style={{
          marginTop: 12,
          fontSize: 12,
          color: ready ? "#065f46" : "#92400e",
        }}
      >
        {ready
          ? "Verification complete. Boarding pass ready."
          : "Complete verification to activate your pass."}
      </div>
    </div>
  );
}

function Field({ label, value }: { label: string; value: string }) {
  return (
    <div style={{ display: "grid", gap: 2 }}>
      <div
        style={{
          fontSize: 10,
          letterSpacing: 0.6,
          textTransform: "uppercase",
          color: "#6b7280",
        }}
      >
        {label}
      </div>
      <div style={{ fontSize: 16, fontWeight: 600 }}>{value}</div>
    </div>
  );
}
