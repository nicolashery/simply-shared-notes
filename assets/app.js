function setTimezoneCookie() {
  try {
    var tz = Intl.DateTimeFormat().resolvedOptions().timeZone;
    if (!tz) return;
    var cookie = "tz=" + tz + "; Path=/; Max-Age=31536000; SameSite=Lax";
    if (location.protocol === "https:") cookie += "; Secure";
    // Only set if changed
    var m = document.cookie.match(/(?:^|;\s*)tz=([^;]+)/);
    if (!m || m[1] !== tz) document.cookie = cookie;
  } catch (_) {}
}

setTimezoneCookie();
