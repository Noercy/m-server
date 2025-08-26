history.replaceState({ page: "index" }, "", "/");


document.getElementById("scanBtn").addEventListener("click", () => {
    fetch("/scan", {method: "POST" })
    .then(res => res.json())
    .then(data => {
        console.log(data.message);
    })
})

function loadSeries(id) {
    fetch(`/series/${id}`)
        .then(res => res.text())
        .then(html => {
            document.getElementById("main-content").innerHTML = html;
            history.pushState({}, "", `/series/${id}`); // updates URL
        });
}

function loadIndex() {
    fetch(`/`)
        .then(res => res.text())
        .then(html => {
            document.getElementById("main-content".innerHTML = html);
            history.pushState({ page: "index" }, "", `/`);
        })
}

window.addEventListener("popstate", function(event) {
    if (event.state) {
        if (event.state.page === "series") {
            loadSeries(event.state.id); // re-fetch series page
        } else if (event.state.page === "index") {
            loadIndex(); // re-fetch homepage
        }
    } else {
        // No state → assume index
        loadIndex();
    }
});