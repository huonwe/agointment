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
            };

        })
    }
    ).catch((res)=>{
        document.querySelector("#content").innerHTML = res;
    })
    
}

function getAvailiable(name) {
    document.querySelector("#manifest").innerHTML = "请求中...";
    fetch('/equipment/getAvailiable?name='+name)
    .then((res)=>{
        res.text().then((res)=>{
            document.querySelector("#manifest").innerHTML = res
        })
    }
    )
}

function requestEquipment(){
    
}


// window.addEventListener("load", ()=>{
//     key = getQueryVariable("key");
//     if(){}
// })