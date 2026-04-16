// リプライフォームのトグル表示
document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll(".btn-reply").forEach(function (btn) {
    btn.addEventListener("click", function () {
      var id   = btn.dataset.micropostId;
      var form = document.getElementById("reply-form-" + id);
      if (!form) return;

      var visible = form.style.display !== "none" && form.style.display !== "";
      form.style.display = visible ? "none" : "block";

      if (!visible) {
        // フォームを表示したときにテキストエリアにフォーカスし、カウンターを再初期化
        var textarea = form.querySelector("textarea");
        if (textarea) {
          textarea.focus();
          // カウンターが未初期化の場合は初期化（reply.js ロード後に追加されたフォーム対策）
          initMicropostCounter(form);
        }
      }
    });
  });
});
