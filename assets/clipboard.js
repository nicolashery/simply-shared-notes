document.addEventListener("DOMContentLoaded", function () {
  const copyButtons = document.querySelectorAll("button[data-target]");

  copyButtons.forEach((button) => {
    button.addEventListener("click", function () {
      const targetId = this.getAttribute("data-target");
      const input = document.getElementById(targetId);

      if (!input) {
        console.error("Target input not found:", targetId);
        return;
      }

      if (!navigator.clipboard) {
        console.error("Clipboard API not available");
        return;
      }

      navigator.clipboard
        .writeText(input.value)
        .then(() => {
          const originalText = this.textContent;
          const copiedText = this.getAttribute("data-copied-text");
          this.textContent = copiedText;
          this.disabled = true;

          setTimeout(() => {
            this.textContent = originalText;
            this.disabled = false;
          }, 2000);
        })
        .catch((err) => {
          console.error("Failed to copy to clipboard:", err);
          const originalText = this.textContent;
          const errorText = this.getAttribute("data-error-text");
          this.textContent = errorText;
          setTimeout(() => {
            this.textContent = originalText;
          }, 2000);
        });
    });
  });
});
