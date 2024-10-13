/* State variables */

const FILESTATE = {
  saved: false
};


/* Global functions */

// Update the "synched/unsynched" icon next to the button.
function updateIcon() {
  const buttonText = document.getElementById("button-text");
  const icon = document.getElementById("icon");

  if (FILESTATE.saved) {
      // change to checkmark if file is saved
      icon.className = "fas fa-check-circle has-text-success";
      buttonText.innerText = "Synched";
  } else {
      // change to exclamation if file is unsaved
      icon.className = "fas fa-exclamation-triangle has-text-warning";
      buttonText.innerText = "Unsynched";
  }
}

// Fetch the saved clipboard from the server
// and set it has the textarea value.
async function setClipboardText() {
  const response = await fetch("/read");
  const text = await response.text();

  document.getElementById("textarea").value = text;

  FILESTATE.saved = true;
  updateIcon();
}

// Send the text within the textarea to the server to be saved.
function saveClipboardText() {
  const text = document.getElementById("textarea").value;

  fetch("/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ data: text }),
  });

  FILESTATE.saved = true;
  updateIcon();
}


/* Event listeners */

// On page load, fetch the updated clipboard.
window.addEventListener("load", () => setClipboardText());

// On any keypress, check for (ctrl | meta) + s.
window.addEventListener("keydown", (event) => {
  if ((event.ctrlKey || event.metaKey) && event.key === 's') {
    event.preventDefault();
    saveClipboardText();
  }
});

// Force refresh when app becomes visible (mobile apps).
if (/Mobi|Android|iPhone|iPad/i.test(navigator.userAgent)) {
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "visible") {
      location.reload(true);
    }
  });
}

// On any keypress within the 'textarea':
// 1) From a list of keys, update the "sync" icon button.
// 2) If it's a tab, introduce that tab char into the textarea.
document.getElementById("textarea").addEventListener("keydown", (event) => {
  const nonEditingKeys = [
    "ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown",
    "Shift", "Control", "Alt", "Meta", "CapsLock", "Escape"
  ];
  if (!nonEditingKeys.includes(event.key)) {
    FILESTATE.saved = false;
    updateIcon();
  }

  if (event.key === "Tab") {
    event.preventDefault();

    // insert 2 spaces instead of tab
    const spaces = "  ";
    this.value = (
      this.value.substring(0, this.selectionStart) + spaces + this.value.substring(this.selectionEnd)
    );

    // move caret to the right position after spaces are inserted
    this.selectionStart = this.selectionEnd = start + spaces.length;
  }
});
