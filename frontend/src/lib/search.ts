export function normalizeSearchText(value: string): string {
  return value.trim().toLocaleLowerCase("ru-RU");
}

export function includesCaseInsensitive(haystack: string | null | undefined, needle: string): boolean {
  const q = normalizeSearchText(needle);
  if (!q) return true;
  return normalizeSearchText(haystack ?? "").includes(q);
}
