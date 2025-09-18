import { useEffect, useState } from "react";
import "./App.css";
import Nav from "./components/Nav";
import logoUrl from "./assets/air-logo.svg";
import Router from "./router";
import i18n from "./i18n";
import { useTranslation } from "react-i18next";

type Language = "nl" | "en";

function detectLanguageFromPath(pathname: string): Language {
  const first = pathname.replace(/^\/+/, "").split("/")[0];
  return first === "nl" || first === "en" ? (first as Language) : "nl";
}

function ensureLanguagePath(language: Language) {
  const segment = language === "nl" ? "nl" : "en";
  const firstSegment = window.location.pathname
    .replace(/^\/+/, "")
    .split("/")[0];
  if (firstSegment === segment) return;

  const rest = window.location.pathname
    .replace(/^\/+/, "")
    .split("/")
    .slice(1)
    .join("/");
  const hash = window.location.hash;
  const search = window.location.search;
  const newUrl = `/${segment}${rest ? "/" + rest : ""}${search || ""}${
    hash || ""
  }`;
  window.history.replaceState({}, "", newUrl);
}

function AppContent() {
  const { t } = useTranslation();
  const year = new Date().getFullYear();
  const copyright = t("app.footer.copyright").replace("{year}", String(year));

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="app-header__brand">
          <img src={logoUrl} alt={t("brand.name")} width={56} height={56} />
          <div className="app-header__copy">
            <h1>{t("brand.name")}</h1>
            <p>{t("brand.tagline")}</p>
          </div>
        </div>
        <div className="app-header__meta">{t("brand.notice")}</div>
      </header>

      <Nav />

      <main className="app-main">
        <Router />
      </main>

      <footer className="app-footer" aria-label={`${t("brand.name")} footer`}>
        <div className="app-footer__inner">
          <span>{copyright}</span>
          <span>{t("app.footer.tagline")}</span>
        </div>
      </footer>
    </div>
  );
}

function App() {
  const [language, setLanguage] = useState<Language>(() =>
    typeof window !== "undefined"
      ? detectLanguageFromPath(window.location.pathname)
      : "nl"
  );

  // Keep state in sync with URL (back/forward)
  useEffect(() => {
    if (typeof window === "undefined") return;

    const updateLanguage = () => {
      const next = detectLanguageFromPath(window.location.pathname);
      setLanguage(next);
    };

    window.addEventListener("popstate", updateLanguage);
    return () => window.removeEventListener("popstate", updateLanguage);
  }, []);

  // Ensure URL has the lang segment and keep i18n synced with `language`
  useEffect(() => {
    if (typeof window === "undefined") return;

    // Force /nl or /en prefix
    ensureLanguagePath(language);

    // Sync i18next
    if (i18n.language !== language) {
      i18n.changeLanguage(language);
    }
  }, [language]);

  return <AppContent />;
}

export default App;
