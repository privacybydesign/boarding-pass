import { useTranslation } from "react-i18next";

export default function Nav() {
  const { t } = useTranslation();

  return (
    <nav className="main-nav">
      <div className="main-nav__inner">
        <span className="main-nav__label">{t("nav.label")}</span>
        <div className="main-nav__links">
          <a href="#/">{t("nav.home")}</a>
          <a href="#/verify">{t("nav.verify")}</a>
          <a href="#/about">{t("nav.about")}</a>
        </div>
      </div>
    </nav>
  );
}
