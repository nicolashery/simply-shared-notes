import { OverType } from "overtype";

window.OverType = OverType;

const themeLight = "solar";
const themeDark = "cave";

// Read once; SSR guarantees this won't change at runtime
const rootTheme = document.documentElement.getAttribute("data-theme");
const isForcedTheme = rootTheme !== null;
const isForcedDark = rootTheme === "dark";

const mql = window.matchMedia("(prefers-color-scheme: dark)");

function isEffectiveDark() {
  // Forced mode wins; otherwise fall back to OS preference
  return isForcedTheme ? isForcedDark : mql.matches;
}

function applyEditorTheme() {
  OverType.setTheme(isEffectiveDark() ? themeDark : themeLight);
}

function createEditor() {
  // Get initial value
  const editorElement = document.querySelector("#editor");
  let value = editorElement?.getAttribute("data-value") || "";
  // Remove the data-value attribute to clean up the DOM
  if (editorElement && editorElement.hasAttribute("data-value")) {
    editorElement.removeAttribute("data-value");
  }

  // Set initial theme
  applyEditorTheme();

  // Create editor
  const [editor] = new OverType("#editor", {
    value: value,
    padding: "12px", // Match daisyUI intput/textarea padding
    textareaProps: {
      id: "content",
      name: "content",
    },
  });

  window.overTypeEditor = editor;
}

window.addEventListener("DOMContentLoaded", () => {
  createEditor();

  // In System mode, track OS changes
  if (!isForcedTheme) {
    const onChange = () => applyEditorTheme();
    if (mql.addEventListener) {
      mql.addEventListener("change", onChange);
    } else if (mql.addListener) {
      // Safari fallback
      mql.addListener(onChange);
    }
  }
});
