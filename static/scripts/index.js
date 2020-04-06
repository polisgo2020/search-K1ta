function onFind() {
    let input = document.getElementsByClassName("search-input")[0];
    let list = document.getElementsByClassName("result")[0];

    let xhr = new XMLHttpRequest();
    xhr.open('GET', "find?phrase=" + input.value, true);
    xhr.send();
    xhr.onreadystatechange = function () {
        if (xhr.readyState !== 4) return;
        if (xhr.status === 200) {
            let response = JSON.parse(xhr.responseText);
            list.innerHTML = "";
            for (let k in response) {
                let div = document.createElement("div");
                div.className = "result-line";
                div.innerHTML = "<div class='result-title'>" + k +
                    ":</div><div class='result-entries'>" + response[k] + " entries</div>";
                list.appendChild(div);
            }
        }
    };
}
