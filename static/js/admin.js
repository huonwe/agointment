function adminRequestingOp(requestID, op, equipmentID, requestorID) {
    switch (op) {
        case "assign":
            const mask = showMask()
            fetch(`/admin/assignUnits?requestID=${requestID}&equipmentID=${equipmentID}`).then((res) => {
                res.text().then((res) => {
                    mask.innerHTML = res;

                    const btn = document.createElement("button");
                    btn.onclick = hideMask;
                    btn.innerText = "取消";
                    mask.appendChild(btn);
                    const form = mask.querySelector("#assign");
                    form.addEventListener("submit", (ev) => {
                        ev.preventDefault();
                        const data = new FormData(form);
                        data.append("requestID", requestID);
                        data.append("equipmentID", requestID);
                        data.append("requestorID", requestorID);
                        fetch("/admin/assignUnits", {
                            method: "post",
                            body: data
                        }).then((res) => {
                            res.json().then((res) => {

                                if (res.status != "Success") {
                                    alert(res.msg)
                                } else {    // Success
                                    location.reload()
                                }

                            }
                            )
                        })
                        console.log(data.getAll("unitID"))
                    })
                })

            })
            break
        case "reject":
            fetch(`/admin/requestingsOp?op=${op}&requestID=${requestID}&equipmentID=${equipmentID}`)
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

function adminAssignUnit() {

}


function showMask(text) {
    if (document.querySelector("#c")) {
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