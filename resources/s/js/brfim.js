function createOwner() {
    owner = localStorage.getItem('owner')
    console.log("current owner", owner)
    if (!owner) {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "/owner", true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
        xhr.onreadystatechange = function () {
            if (this.status == 201) {
                var data = JSON.parse(this.responseText);
                console.log("owner associated with", data.data.owner)
                localStorage.setItem('owner', data.data.owner);
            }
        };
        xhr.send(JSON.stringify({ }));
    }
}

function createShortURL() {
    owner = localStorage.getItem('owner')
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/owner/"+owner+"/url", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
    xhr.onreadystatechange = function () {
        if (this.status == 201) {
            var data = JSON.parse(this.responseText);
            console.log("shortURL ", data.data.shortURL)
            document.getElementById('shortURL').value = data.data.shortURL;
            document.getElementById('qrCodeImage').setAttribute('src', 'data:image/png;base64,'+data.data.qrCode);
            document.getElementById('qrCodeImage').setAttribute('alt', data.data.shortID);
            document.getElementById("qrCodeBlock").style.display = "block";
            document.getElementById('qrCodeLink').setAttribute('href', 'data:image/png;base64,'+data.data.qrCode);
            document.getElementById('qrCodeLink').setAttribute('download', data.data.shortID);
            document.getElementById('qrCodeURL').value = data.data.qrCodeURL;
        }
    };
    xhr.send(JSON.stringify({ "url":  document.getElementById('inputFullURL').value, "prefix": document.getElementById('inputPrefix').value }));
}

document.addEventListener('DOMContentLoaded', function() {
    createOwner()
}, false);
