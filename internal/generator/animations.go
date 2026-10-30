package generator

// animationScript returns src/scripts/animations.ts for the chosen level.
// No template literals are used so the strings stay valid inside Go raw strings.
func animationScript(level string) string {
	switch level {
	case "subtle":
		return subtleTS
	case "expressive":
		return expressiveTS
	default:
		return "// Animations were turned off in the wizard.\nexport {};\n"
	}
}

const countersTS = `
  // stat counters: "50k+" counts up to 50 and keeps the "k+" suffix
  document.querySelectorAll<HTMLElement>("[data-count]").forEach((el) => {
    const m = (el.dataset.count ?? "").match(/^([\d,.]+)(.*)$/);
    if (!m) return;
    const target = parseFloat(m[1].replace(/,/g, ""));
    const suffix = m[2];
    const dec = m[1].includes(".") ? m[1].split(".")[1].length : 0;
    inView(el, () => {
      animate(0, target, {
        duration: 1.6,
        ease: "easeOut",
        onUpdate: (v) =>
          (el.textContent = v.toLocaleString(undefined, { minimumFractionDigits: dec, maximumFractionDigits: dec }) + suffix),
      });
    });
  });
`

const subtleTS = `// Subtle animations: fade-up reveals. Library: motion (Framer Motion's vanilla API).
import { animate, inView, stagger } from "motion";

const reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

if (!reduce) {
  animate("[data-hero-item]", { opacity: [0, 1], y: [16, 0] }, { duration: 0.6, delay: stagger(0.1) });

  inView("[data-reveal]", (el) => {
    animate(el, { opacity: [0, 1], y: [16, 0] }, { duration: 0.6 });
  });
` + countersTS + `}
`

const expressiveTS = `// Expressive animations: anime.js hero entrance + motion scroll reveals and counters.
import { animate, inView } from "motion";
import { animate as anime, stagger } from "animejs";

const reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

if (!reduce) {
  anime("[data-hero-item]", {
    opacity: [0, 1],
    translateY: [32, 0],
    delay: stagger(120),
    duration: 900,
    ease: "outExpo",
  });

  inView("[data-reveal]", (el) => {
    animate(el, { opacity: [0, 1], y: [40, 0], scale: [0.98, 1] }, { duration: 0.8, ease: "easeOut" });
  });
` + countersTS + `}
`
