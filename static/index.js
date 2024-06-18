const createLink = (url, text) => {
    let link = document.createElement("a");

    link.setAttribute("href", url);
    link.setAttribute("target", "_blank");
    link.setAttribute("style", "text-overflow: ellipsis;white-space: nowrap;overflow: hidden;");

    let linkText = document.createTextNode(text);
    link.appendChild(linkText);

    return link
};

let tableUrls = new Set();

const addShortToTable = (url, code) => {
    tableUrls.add(url);
    let table = document.getElementById("table-body");
    var row = table.insertRow(0);
    let cell1 = row.insertCell(0);
    let cell2 = row.insertCell(1);

    let codeLink = "http://shorturl.space/" + code;
    let codeLinkText = "shorturl.space/" + code;
    cell1.appendChild(createLink(codeLink, codeLinkText));
    cell2.appendChild(createLink(url, url));
};


const showError = (t) => {
    alert(t);
};

const fetchNewCode = async function (url) {
    const requestBody = {
        urlOriginal: url
    };

    const response = await fetch("http://localhost:8000/api/new", {
                           method: "POST", 
                           headers: {
                               "Content-Type": "application/json",
                           },
                           body: JSON.stringify(requestBody)
                          });
    if (response.ok) {
        let j = await response.json();
        addShortToTable(url, j["urlCode"]);
    } else {
        let t = await response.text();
        showError("API Error: " + t);
    }
};

const alreadyShortened = (url) => {
    if (tableUrls.has(url)) {
        return true;
    }

    return false;
};

let urlInput = document.getElementById("url-input");
let form = document.getElementById("form");

form.addEventListener('submit', function(e) {
        e.preventDefault();
        let url = urlInput.value;
        if (!alreadyShortened(url)) {
            fetchNewCode(url);
        }
});

