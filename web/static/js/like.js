// いいねボタンの AJAX 処理。
// ページリロードなしでいいね/いいね解除し、UI を即座に更新する。
// イベント委譲で .like-form のサブミットをすべて捕捉する。

document.addEventListener("submit", async function (e) {
  var form = e.target;
  if (!form.classList.contains("like-form")) return;
  e.preventDefault();

  var micropostId = form.dataset.micropostId;
  var isLiked     = form.dataset.liked === "true";
  var csrfInput   = form.querySelector('[name="csrf_token"]');
  if (!csrfInput) return;
  var csrfToken   = csrfInput.value;

  var url   = isLiked ? "/likes/" + micropostId : "/likes";
  var parts = ["csrf_token=" + encodeURIComponent(csrfToken)];
  if (isLiked) {
    parts.push("_method=DELETE");
  } else {
    parts.push("micropost_id=" + micropostId);
  }

  try {
    var response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type":      "application/x-www-form-urlencoded",
        "X-Requested-With":  "XMLHttpRequest",
      },
      body: parts.join("&"),
    });

    if (!response.ok) return;
    var data = await response.json();

    // このページ上で同じ micropostId を持つすべての like-info を更新
    document.querySelectorAll('.like-info[data-micropost-id="' + micropostId + '"]').forEach(function (info) {
      // カウント更新
      var countEl = info.querySelector(".like-count");
      if (countEl) countEl.textContent = "\u2665 " + data.count;

      // フォームを新しいいいね状態に合わせて書き換える
      var f = info.querySelector(".like-form");
      if (!f) return;

      if (data.liked) {
        f.action              = "/likes/" + micropostId;
        f.dataset.liked       = "true";
        f.innerHTML =
          '<input type="hidden" name="csrf_token" value="' + escHtml(csrfToken) + '"/>' +
          '<input type="hidden" name="_method" value="DELETE"/>' +
          '<button type="submit" class="btn-like btn-like--active" title="Unlike">\u2665</button>';
      } else {
        f.action              = "/likes";
        f.dataset.liked       = "false";
        f.innerHTML =
          '<input type="hidden" name="csrf_token" value="' + escHtml(csrfToken) + '"/>' +
          '<input type="hidden" name="micropost_id" value="' + micropostId + '"/>' +
          '<button type="submit" class="btn-like" title="Like">\u2661</button>';
      }
    });
  } catch (err) {
    console.error("Like error:", err);
  }
});

function escHtml(s) {
  return s.replace(/&/g, "&amp;").replace(/"/g, "&quot;");
}
