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
    fetch('/?page='+key)
    .then((res)=>{
        res.text().then((res)=>{
            document.querySelector("#content").innerHTML = res
        })
    }
    )
}

function test() {
    alert("test")
}