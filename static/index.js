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

    let codeLink = "http://shorturl.space/" + code;
    let table = document.getElementById("table-body");
    var row = table.insertRow(0);
    let cell1 = row.insertCell(0);
    let cell2 = row.insertCell(1);

    cell1.appendChild(createLink(codeLink, codeLink));
    cell2.appendChild(createLink(url, url));
};


const fetchNewCode = (url) => {
    const requestBody = {
        urlOriginal: url
    };

    fetch("http://localhost:8000/api/new", {
          method: "POST", 
          headers: {
              "Content-Type": "application/json",
          },
          body: JSON.stringify(requestBody)
         })
    .then((response) => response.json())
    .then((responseJSON) => addShortToTable(url, responseJSON["urlCode"]))
    .catch((error) => console.log(error));
};

const urlOk = (url) => {
    if (tableUrls.has(url)) {
        return false;
    }

    if (url.includes("shorturl.space")) {
        return false;
    }

    return true;
};

let urlInput = document.getElementById("url-input");
let form = document.getElementById("form");

form.addEventListener('submit', function(e) {
        e.preventDefault();
        let url = urlInput.value;
        if (urlOk(url)) {
            fetchNewCode(url);
        }
});

