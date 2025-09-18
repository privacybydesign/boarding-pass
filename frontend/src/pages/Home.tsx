import { useTranslation } from "react-i18next";

export default function Home() {
  const { t } = useTranslation();

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
      <h2 style={{ marginTop: 0 }}>{t("home.title")}</h2>
      <p>{t("home.description")}</p>
      <ul style={{ margin: 0, paddingLeft: "1.25rem" }}>
        <li>{t("home.step_buy")}</li>
        <li>{t("home.step_verify")}</li>
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
          <button>{t("home.buy_button")}</button>
        </a>
        <a href="#/verify">
          <button>{t("home.verify_button")}</button>
        </a>
      </div>
    </div>
  );
}
