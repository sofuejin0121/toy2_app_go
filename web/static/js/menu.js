// メニュー操作

// トグルリスナーを追加する
function addToggleListener(selected_id, menu_id, toggle_class) {
  let selected_element = document.querySelector(`#${selected_id}`);
  if (selected_element) {
    selected_element.addEventListener("click", function(event) {
      event.preventDefault();
      let menu = document.querySelector(`#${menu_id}`);
      menu.classList.toggle(toggle_class);
    });
  }
}

// クリックをリッスンするトグルリスナーを追加する
document.addEventListener("DOMContentLoaded", function() {
  addToggleListener("hamburger", "navbar-menu",   "collapse");
  addToggleListener("account",   "dropdown-menu", "active");
});
