window.onload = function () {
    let file_upload = document.querySelector("#file")
    file_upload.addEventListener("change", (ev) => {
        ev.preventDefault()
        let formData = new FormData()
        formData.append("file", file_upload.files[0])
        fetch(`/equipment/import`, {
            method: "POST",
            body: formData,
        }).then((response) => response.json()).then(result => {
            alert(result.msg)
            if (result.status == "Success") {
                location.reload()
            }
        })
    })
}

function getEquipment(name, page, pageSize, op) {
    page = parseInt(page) || 1;
    pageSize = parseInt(pageSize) || 15;
    if (op == "prev") page = page - 1;
    if (op == "next") page = page + 1;

    window.location = `/admin/equipment?name=${name}&page=${page}&pageSize=${pageSize}`
}


function equipmentOp(id, op) {
    fetch(`/equipment/equipmentOp?op=${op}&id=${id}`).then((res) => {
        res.json().then((res) => {
            alert(res.msg)
            location.reload()
        })
    })
}