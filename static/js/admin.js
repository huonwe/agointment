function adminRequestingOp(requestID, op) {
    switch (op) {
        case "assign":
            const mask = showMask()
            
            break
        case "reject":
            fetch('/admin/requestingsOp?' + 'op=' + op + '&requestID=' + requestID)
                .then((res) => {
                    res.json().then((res) => {
                        if (res.status != "Success") {
                            alert(res.msg)
                        } else {    // Success
                            location.reload()
                        }
                    })
                }
                )
            break
    }
}


// const loading = `<div id="loading" class="full_window" style=""></div>`

function showMask(text) {
    if(document.querySelector("#c")) {
        return
    }
    const div = document.createElement('div');
    div.id = "c"
    div.className = "mask"
    div.innerText = text || "加载中"
    document.body.appendChild(div)
    return div
}

function hideMask() {
    document.body.removeChild(document.querySelector("#c"))
}