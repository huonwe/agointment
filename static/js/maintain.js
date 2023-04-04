function maintain(unitID){
    let formData = new FormData()
    formData.append("id",unitID)
    fetch(`/admin/maintain`,{
        method: "POST",
        body: formData
    }).then((res)=>{res.json().then((res)=>{
        if(res.status != "Success"){
            alert(res.msg)
        }else {
            location.reload()
        }
    })})
}

function getMaintain(page, pageSize, op) {
    page = page || 1
    pageSize = pageSize || 20
    if(op == "prev") {
        page = page - 1
    }else if(op == "next"){
        page = page + 1
    }
    window.location = `/admin/maintain?page=${page}&pageSize=${pageSize}`
}