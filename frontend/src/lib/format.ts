export function fmtDate(s?: string | null) {
  if (!s) return "—";
  const d = new Date(s);
  if (isNaN(d.getTime())) return s;
  return d.toLocaleDateString(undefined, {
    year: "numeric", month: "short", day: "numeric",
  });
}

export function fmtDateRange(start?: string, end?: string) {
  if (start && end) {
    const s = toInputDate(start);
    const e = toInputDate(end);
    if (s && e && s === e) return fmtDate(start);
  }
  if (start && !end) return fmtDate(start);
  if (!start && end) return fmtDate(end);
  return `${fmtDate(start)} → ${fmtDate(end)}`;
}

export function toInputDate(s?: string | null) {
  if (!s) return "";
  const d = new Date(s);
  if (isNaN(d.getTime())) return "";
  // YYYY-MM-DD
  return d.toISOString().slice(0, 10);
}

export function fromInputDate(s: string): string {
  // ISO at midnight UTC
  if (!s) return "";
  return new Date(s + "T00:00:00.000Z").toISOString();
}

export function fmtRub(value?: number | null) {
  const amount = Number(value ?? 0);
  return `${amount.toLocaleString("ru-RU")} ₽`;
}
