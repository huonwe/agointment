function adminRequestingOp(requestID, op) {
    fetch('/admin/myRequestOp?'+'op='+'cancel&requestID='+requestID)
    .then((res)=>{
        res.json().then((res)=>{
            alert(res.msg)
            getHTML('status')
        })
    }
    )
}