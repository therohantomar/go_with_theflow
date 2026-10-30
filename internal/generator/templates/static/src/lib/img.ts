// All photos live in src/assets/<section>/ (not public/) so Astro can optimise them with sharp.
import type { ImageMetadata } from "astro";

const files = import.meta.glob<{ default: ImageMetadata }>("/src/assets/**/*.{jpg,jpeg,png,webp,avif,svg}", { eager: true });

/** Returns the image metadata for src/assets/<section>/<file>, or null when blank or missing. */
export function img(section: string, file?: string): ImageMetadata | null {
  if (!file) return null;
  return files[`/src/assets/${section}/${file}`]?.default ?? null;
}
