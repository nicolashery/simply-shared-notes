import { OverType } from "overtype";

window.OverType = OverType;

const themeLight = "solar";
const themeDark = "cave";
const darkModeQuery = window.matchMedia("(prefers-color-scheme: dark)");

function createEditor() {
  // Get initial value
  const editorElement = document.querySelector("#editor");
  let value = editorElement?.getAttribute("data-value") || "";
  // Remove the data-value attribute to clean up the DOM
  if (editorElement && editorElement.hasAttribute("data-value")) {
    editorElement.removeAttribute("data-value");
  }

  // Set initial theme
  OverType.setTheme(darkModeQuery.matches ? themeDark : themeLight);

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

  darkModeQuery.addEventListener("change", (updatedDarkModeQuery) => {
    OverType.setTheme(updatedDarkModeQuery.matches ? themeDark : themeLight);
  });
});
