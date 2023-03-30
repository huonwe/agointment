window.onload = function(){
    let file_upload = document.querySelector("#file")
    file_upload.addEventListener("change",(ev)=>{
        ev.preventDefault()
        let formData = new FormData()
        formData.append("file",file_upload.files[0])
        fetch(`/equipment/import`,{
            method: "POST",
            body: formData,
        }).then((response) => response.json()).then(result => {
            alert(result.msg)
            if(result.status == "Success"){
                location.reload()
            }
        })
    })
}