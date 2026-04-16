// マイクロポスト投稿フォームのリアルタイム文字数カウンター
// ホームの compose フォームとリプライフォームの両方に対応。
function initMicropostCounter(form) {
  // 二重初期化を防ぐ
  if (form.dataset.counterInitialized) return;
  form.dataset.counterInitialized = "1";

  var textarea = form.querySelector('textarea[name="content"]');
  var counter  = form.querySelector(".micropost-counter");
  var submitBtn = form.querySelector('input[type="submit"], button[type="submit"]');

  if (!textarea || !counter) return;

  var MAX = 140;

  function update() {
    // スプレッド構文でUnicode文字を1文字として数える（Go の []rune と同じ挙動）
    var len = [...textarea.value].length;
    var remaining = MAX - len;

    counter.textContent = len + "/" + MAX;

    if (len > MAX) {
      counter.className = "micropost-counter micropost-counter--over";
      if (submitBtn) submitBtn.disabled = true;
    } else if (remaining <= 20) {
      counter.className = "micropost-counter micropost-counter--warn";
      if (submitBtn) submitBtn.disabled = false;
    } else {
      counter.className = "micropost-counter";
      if (submitBtn) submitBtn.disabled = false;
    }
  }

  textarea.addEventListener("input", update);
  update();
}

document.addEventListener("DOMContentLoaded", function () {
  // .micropost-counter を持つ全フォーム（メインフォーム＋リプライフォーム）を初期化
  document.querySelectorAll(".micropost-counter").forEach(function (counter) {
    var form = counter.closest("form");
    if (form) initMicropostCounter(form);
  });
});
