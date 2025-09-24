const envEndpoint =
  typeof import.meta !== "undefined" && import.meta.env
    ? (import.meta.env.VITE_API_ENDPOINT as string | undefined)
    : undefined;

function normalizeEndpoint(value: string) {
  const trimmed = value.trim();
  return trimmed.replace(/\/+$/, "");
}

const placeholderAwareValue = envEndpoint ?? "%VITE_API_ENDPOINT%";
const normalized = placeholderAwareValue
  ? normalizeEndpoint(placeholderAwareValue)
  : "";
const isPlaceholder = normalized.includes("VITE_API_ENDPOINT");

export const apiEndpoint = !normalized || isPlaceholder ? "/api" : normalized;
