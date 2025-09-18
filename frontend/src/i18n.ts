import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    detection: {
      order: ["path", "navigator"],
      lookupFromPathIndex: 0,
    },
    resources: {
      en: {
        translation: {
          "brand.name": "Yivi International Airlines",
          "brand.tagline": "Online check-in & boarding pass demo",
          "brand.notice": "Check-in opens 24h before departure",
          "app.footer.tagline":
            "Privacy-first journeys from check-in to arrival.",
          "app.footer.copyright": "© {year} Yivi International Airlines",
          "nav.label": "Online check-in",
          "nav.home": "Home",
          "nav.buy": "Buy",
          "nav.verify": "Verify",
          "nav.about": "About",
          "home.title": "Welcome to Yivi International Airlines",
          "home.description":
            "This simple demo lets you “buy” a ticket with your name and passport/document number, then verify your identity with Yivi to display a boarding pass. It’s a minimal flow to showcase privacy-preserving verification.",
          "home.step_buy": "Buy: create a fake ticket with basic info.",
          "home.step_verify":
            "Verify: scan a QR in Yivi to complete check-in and get a boarding pass.",
          "home.buy_button": "Buy Ticket",
          "home.verify_button": "Verify & Boarding Pass",
          "about.title": "About",
          "about.description":
            "This demo starts a Yivi verification session and shows the raw session pointer JSON. Use the Home page to start a session.",
          "buy.title": "Buy a Ticket",
          "buy.subtitle": "Enter basic details for a demo ticket.",
          "buy.first_name_label": "First name",
          "buy.first_name_placeholder": "Jane",
          "buy.last_name_label": "Last name",
          "buy.last_name_placeholder": "Doe",
          "buy.document_label": "Passport / Document number",
          "buy.document_placeholder": "X1234567",
          "buy.error_generic": "Failed to create ticket.",
          "buy.submit_processing": "Processing…",
          "buy.submit_cta": "Continue to verification",
          "verify.title": "Verify and Get Boarding Pass",
          "verify.missing_ticket":
            "Ticket details missing. Please create a ticket again.",
          "verify.banner_missing":
            "No ticket details found. Please create one on the Buy page.",
          "verify.instructions":
            "Scan the QR code with your Yivi app to authenticate. If you don’t have the app, get it from",
          "verify.error_label": "Error",
          "verify.error_start": "Verification failed to start.",
          "verify.error_client": "Unable to load verification client.",
          "boardingpass.title": "Yivi International Airlines",
          "boardingpass.subtitle": "Demo Boarding Pass",
          "boardingpass.field_name": "Name",
          "boardingpass.field_document": "Document",
          "boardingpass.field_from": "From",
          "boardingpass.field_to": "To",
          "boardingpass.field_flight": "Flight",
          "boardingpass.field_seat": "Seat",
          "boardingpass.ready": "Verification complete. Boarding pass ready.",
          "boardingpass.pending":
            "Complete verification to activate your pass.",
        },
      },
      nl: {
        translation: {
          "brand.name": "Yivi International Airlines",
          "brand.tagline": "Online inchecken & demo voor boardingpass",
          "brand.notice": "Inchecken opent 24 uur voor vertrek",
          "app.footer.tagline":
            "Privacy-first reizen van inchecken tot aankomst.",
          "app.footer.copyright": "© {year} Yivi International Airlines",
          "nav.label": "Online inchecken",
          "nav.home": "Start",
          "nav.buy": "Ticket kopen",
          "nav.verify": "Verifiëren",
          "nav.about": "Over",
          "home.title": "Welkom bij Yivi International Airlines",
          "home.description":
            'Deze eenvoudige demo laat je een ticket "kopen" met je naam en paspoort- of documentnummer. Daarna verifieer je je identiteit met Yivi om een boardingpass te tonen. Het is een minimalistische flow die privacyvriendelijke verificatie demonstreert.',
          "home.step_buy": "Kopen: maak een nep-ticket aan met basisgegevens.",
          "home.step_verify":
            "Verifiëren: scan een QR-code in Yivi om in te checken en je boardingpass te ontvangen.",
          "home.buy_button": "Ticket kopen",
          "home.verify_button": "Verifiëren & boardingpass",
          "about.title": "Over",
          "about.description":
            "Deze demo start een Yivi-verificatiesessie en toont de ruwe session pointer JSON. Gebruik de startpagina om een sessie te beginnen.",
          "buy.title": "Ticket kopen",
          "buy.subtitle": "Voer basisgegevens in voor een demoticket.",
          "buy.first_name_label": "Voornaam",
          "buy.first_name_placeholder": "Jan",
          "buy.last_name_label": "Achternaam",
          "buy.last_name_placeholder": "Jansen",
          "buy.document_label": "Paspoort- of documentnummer",
          "buy.document_placeholder": "X1234567",
          "buy.error_generic": "Ticket aanmaken mislukt.",
          "buy.submit_processing": "Bezig…",
          "buy.submit_cta": "Verder naar verificatie",
          "verify.title": "Verifiëren en boardingpass ophalen",
          "verify.missing_ticket":
            "Ticketgegevens ontbreken. Maak opnieuw een ticket.",
          "verify.banner_missing":
            "Geen ticketgegevens gevonden. Maak er een op de pagina Ticket kopen.",
          "verify.instructions":
            "Scan de QR-code met je Yivi-app om je te identificeren. Heb je de app nog niet? Download hem via",
          "verify.error_label": "Fout",
          "verify.error_start": "Verificatie kon niet starten.",
          "verify.error_client": "Verificatieclient kon niet worden geladen.",
          "boardingpass.title": "Yivi International Airlines",
          "boardingpass.subtitle": "Demo-boardingpass",
          "boardingpass.field_name": "Naam",
          "boardingpass.field_document": "Document",
          "boardingpass.field_from": "Van",
          "boardingpass.field_to": "Naar",
          "boardingpass.field_flight": "Vlucht",
          "boardingpass.field_seat": "Stoel",
          "boardingpass.ready": "Verificatie voltooid. Boardingpass is klaar.",
          "boardingpass.pending":
            "Maak de verificatie af om je pass te activeren.",
        },
      },
    },
    lng: "en", // default language (will be overridden if path/navigator detection finds one)
    fallbackLng: "en",
    interpolation: {
      escapeValue: false, // react already escapes
    },
  });

export default i18n;
