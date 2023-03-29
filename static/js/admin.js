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
                            let detail = res.detail
                            mask.innerHTML = `
                            <div class="detail">
                                <div>详细信息</div>
                                <div>
                                    <span>设备名称</span>
                                    <span>${detail.Equipment.Name}</span>
                                </div>
                                <div>
                                    <span>设备编号</span>
                                    <span>${detail.EquipmentUnit.ID}</span>
                                </div>
                                <div>
                                    <span>序列号</span>
                                    <span>${detail.EquipmentUnit.SerialNumber}</span>
                                </div>
                                <div>
                                    <span>型号</span>
                                    <span>${detail.Equipment.Type}</span>
                                </div>
                                <div>
                                    <span>品牌</span>
                                    <span>${detail.EquipmentUnit.Brand}</span>
                                </div>
                                <div>
                                    <span>单价</span>
                                    <span>${detail.EquipmentUnit.Price}</span>
                                </div>
                                <div>
                                    <span>所在科室</span>
                                    <span>${detail.User.DepartmentName}</span>
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
                                    <span>${detail.EquipmentUnit.Status}</span>
                                </div>
                                <div>
                                    <span>设备分类</span>
                                    <span>${detail.Equipment.Class}</span>
                                </div>
                                <div>
                                    <span>设备标注</span>
                                    <span>${detail.EquipmentUnit.Label}</span>
                                </div>
                                <div>
                                    <span>厂家信息</span>
                                    <span>${detail.EquipmentUnit.Factory}</span>
                                </div>
                                <div>
                                    <span>备注</span>
                                    <span>${detail.EquipmentUnit.Remark}</span>
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