function getQueryVariable(variable)
{
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i=0;i<vars.length;i++) {
        var pair = vars[i].split("=");
        if(pair[0] == variable){return pair[1];}
    }
    return(false);
}

function getHTML(key)
{
    document.querySelector("#content").innerHTML = "请求中...";
    fetch('/?page='+key)
    .then((res)=>{
        res.text().then((res)=>{
            document.querySelector("#content").innerHTML = res;
            
            switch (key) {
                case "appoint":
                    getAvailiable('')
                    break;
                case "status":
                    getMyRequest('')
                    break;
            };

        })
    }
    ).catch((e)=>{
        console.log(e)
    })
    
}

function getAvailiable(name) {
    document.querySelector("#manifest").innerHTML = "请求中...";
    fetch('/equipment/availiable?name='+name)
    .then((res)=>{
        res.text().then((res)=>{
            document.querySelector("#manifest").innerHTML = res
        })
    }
    )
}

async function requestEquipment(equipmentID){
    const response = await fetch("/equipment/makeRequest?equipmentID="+equipmentID);
    const res = await response.json()
    getAvailiable('')
    alert(res.msg)
}

function getMyRequest(name, page, pageSize, op) {
    page = parseInt(page) || 1;
    pageSize = parseInt(pageSize) || 15;
    if(op == "prev") page = page -1;
    if(op == "next") page = page +1;
    document.querySelector("#manifest").innerHTML = "请求中...";
    fetch( `/user/myRequest?name=${name}&page=${page}&pageSize=${pageSize}`)
    .then((res)=>{
        res.text().then((res)=>{
            document.querySelector("#manifest").innerHTML = res
        })
    }
    )
}

function myRequestOp(requestID, op) {
    fetch('/user/myRequestOp?'+'op='+op+'&requestID='+requestID)
    .then((res)=>{
        res.json().then((res)=>{
            if(res.status != "Success"){
                alert(res.msg)
            }
            getHTML('status')
        })
    }
    )
}

// window.addEventListener("load", ()=>{
//     key = getQueryVariable("key");
//     if(){}
// })