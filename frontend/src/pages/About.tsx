import { useTranslation } from "react-i18next";

export default function About() {
  const { t } = useTranslation();

  return (
    <div
      className="card"
      style={{ display: "flex", gap: "1rem", flexDirection: "column" }}
    >
      <h2>{t("about.title")}</h2>
      <p>{t("about.description")}</p>
    </div>
  );
}

