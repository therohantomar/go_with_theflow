// Resolves "src/assets/<section>/<file>" to a URL. Returns null when the file is blank or missing.
const files = import.meta.glob("/src/assets/**/*.{jpg,jpeg,png,webp,avif,svg}", {
  eager: true,
  query: "?url",
  import: "default",
}) as Record<string, string>;

export function img(section: string, file?: string): string | null {
  if (!file) return null;
  return files[`/src/assets/${section}/${file}`] ?? null;
}
