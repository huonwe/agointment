function adminRequestingOp(requestID, op, equipmentID, requestorID) {
    equipmentID = equipmentID || ""
    requestorID = requestorID || ""

    switch (op) {
        case "assign":
            var mask = showMask()
            fetch(`/admin/assignUnits?requestID=${requestID}&equipmentID=${equipmentID}`).then((res) => {
                res.text().then((res) => {
                    mask.innerHTML = res;

                    let btn = document.createElement("button");
                    btn.className = "hideBtn"
                    btn.onclick = hideMask;
                    btn.innerText = "取消";
                    mask.appendChild(btn);
                    let form = mask.querySelector("#assign");
                    form.addEventListener("submit", (ev) => {
                        ev.preventDefault();
                        let data = new FormData(form);
                        data.append("requestID", requestID);
                        data.append("equipmentID", equipmentID);
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
        case "finish":
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
        case "detail":
            var mask = showMask("加载中");
            fetch(`/admin/requestingsOp?op=${op}&requestID=${requestID}&equipmentID=${equipmentID}`)
                .then((res) => {
                    res.json().then((res) => {
                        if (res.status != "Success") {
                            alert(res.msg)
                        } else {    // Success
                            // console.log(res.detail)
                            let detail = res.detail
                            mask.innerHTML = `
                            <div class="detail">
                                <div>详细信息</div>
                                <div>
                                    <span>设备名称</span>
                                    <span>${detail.EquipmentName}</span>
                                </div>
                                <div>
                                    <span>设备编号</span>
                                    <span>${detail.UnitUID}</span>
                                </div>
                                <div>
                                    <span>序列号</span>
                                    <span>${detail.UnitSerialNumber}</span>
                                </div>
                                <div>
                                    <span>型号</span>
                                    <span>${detail.EquipmentType}</span>
                                </div>
                                <div>
                                    <span>品牌</span>
                                    <span>${detail.EquipmentBrand}</span>
                                </div>
                                <div>
                                    <span>单价</span>
                                    <span>${detail.UnitPrice}</span>
                                </div>
                                <div>
                                    <span>所在科室</span>
                                    <span>${detail.User.Department.Name}</span>
                                </div>
                                <div>
                                    <span>联系人</span>
                                    <span>${detail.User.Name}</span>
                                </div>
                                <div>
                                    <span>借用日期</span>
                                    <span>${detail.BeginAtStr}</span>
                                </div>
                                <div>
                                    <span>完成日期</span>
                                    <span>${detail.EndAtStr}</span>
                                </div>
                                <div>
                                    <span>设备状态</span>
                                    <span>${detail.UnitStatus}</span>
                                </div>
                                <div>
                                    <span>设备分类</span>
                                    <span>${detail.EquipmentClass}</span>
                                </div>
                                <div>
                                    <span>设备标注</span>
                                    <span>${detail.UnitLabel}</span>
                                </div>
                                <div>
                                    <span>厂家信息</span>
                                    <span>${detail.UnitFactory}</span>
                                </div>
                                <div>
                                    <span>备注</span>
                                    <span>${detail.UnitRemark}</span>
                                </div>
                            </div>
                            
                            `
                            let btn = document.createElement("button");
                            btn.className = "hideBtn"
                            btn.onclick = hideMask;
                            btn.innerText = "返回";
                            mask.appendChild(btn);
                        }
                    })
                }
                )
            break
    }
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

function adminAll(name, page, pageSize, op) {
    page = parseInt(page) || 1;
    pageSize = parseInt(pageSize) || 15;
    if(op == "prev") page = page -1;
    if(op == "next") page = page +1;
    window.location = `/admin/all?name=${name}&page=${page}&pageSize=${pageSize}`
    // document.querySelector("#manifest").innerHTML = "请求中...";
    // fetch( `/admin/all?name=${name}&page=${page}&pageSize=${pageSize}`)
    // .then((res)=>{
    //     res.text().then((res)=>{
    //         document.querySelector("#manifest").innerHTML = res
    //     })
    // }
    // )
}

function deptOp(deptName, op) {
    // console.log(deptName)
    let formData = new FormData()
    if(op == "deptNew"){
        let name = document.querySelector("#dept_add_name").value;
        let description = document.querySelector("#dept_add_dscrpt").value;
        console.log(name,description)
        formData.append("deptName",name);
        formData.append("deptDescpt",description);
    }else {
        formData.append("deptName", deptName)
        // console.log(formData.getAll("deptName"))
    }
    // console.log(formData.get("deptName"))
    fetch(`/admin/users/${op}`,{
        method: "post",
        body: formData
    }).then((res)=>{
        res.json().then((res)=>{
            // alert(res.msg)
            if(res.status != "Success"){
                alert(res.msg);
            }else {
                location.reload();
            }
        })
    })
}


function userOp(op, username, userid) {
    switch(op) {
        case "userNew":
            let user_dpt = document.querySelector("#user_add_dept").value;
            let user_name = document.querySelector("#user_add_name").value;
            let user_password = document.querySelector("#user_add_password").value;
            let formData = new FormData();
            formData.append("user_dept", user_dpt);
            formData.append("user_name", user_name);
            formData.append("user_password", user_password);
            // console.log(formData.get("user_dept"))
            fetch(`/admin/users/${op}`,{
                method: "post",
                body: formData
            }).then((res)=>{
                res.json().then((res)=>{
                    alert(res.msg)
                })
            })
            break
        case "userSearch":
            let user_search_dept = document.querySelector("#user_search_dept").value || "";
            let user_search_name = document.querySelector("#user_search_name").value || "";

            fetch(`/admin/users/userSearch?dept=${user_search_dept}&name=${user_search_name}`,{method: "post"}).then((res)=>{
                // console.log(res)
                res.text().then((res)=>{
                document.querySelector("#users_found").innerHTML = res
            })})
            break;
        case "userDel":
            let formData1 = new FormData()
            formData1.append("name", username)
            formData1.append("id",userid)
            fetch(`/admin/users/userDel`,{
                method: "post",
                body: formData1
            }).then((res)=>{
                console.log(res)
                res.json().then((res)=>{
                alert(res.msg)
            })})
            break;
        case "userPasswd":
            let new_passwd = prompt("请输入新的密码");
            if(new_passwd == "" || new_passwd == null){
                return
            }
            let formData2 = new FormData()
            formData2.append("name", username)
            formData2.append("id",userid)
            formData2.append("newPasswd", new_passwd)
            fetch(`/admin/users/userPasswd`,{
                method: "post",
                body: formData2
            }).then((res)=>{res.json().then((res)=>{
                alert(res.msg)
            })})
            break;
        case "userSetAdmin":
            let formData3 = new FormData()
            formData3.append("name", username)
            formData3.append("id",userid)
            fetch(`/admin/users/userSetAdmin`,{
                method: "post",
                body: formData3
            }).then((res)=>{
                console.log(res)
                res.json().then((res)=>{
                alert(res.msg)
            })})
            break;
        case "userUnsetAdmin":
            let formData4 = new FormData()
            formData4.append("name", username)
            formData4.append("id",userid)
            fetch(`/admin/users/userUnsetAdmin`,{
                method: "post",
                body: formData4
            }).then((res)=>{
                console.log(res)
                res.json().then((res)=>{
                alert(res.msg)
            })})
            break;
    }
    // location.reload()
}

function adminEmptyEnds() {
    let yes = confirm("确定要一键清空已结束(取消、拒绝、完成)的请求吗？\n强烈建议您在使用本功能前先导出excel, 并且定期执行本功能.")
    if(!yes){
        return
    }
    fetch(`/admin/emptyEnded`).then((res)=>{res.json().then((res)=>{
        alert(res.msg)
        location.reload()
    })})
}

function attentionUpdate(content) {
    let formData = new FormData()
    formData.append("content", content)
    fetch(`/admin/attentionUpdate`,{
        method: "POST",
        body: formData
    }).then((res)=>{res.json().then((res)=>{
        alert(res.msg)
    })})
}