let form = document.querySelector("#signupForm")
form.addEventListener("submit",(ev)=>{
    ev.preventDefault()
    let data = new FormData(form);
    let username = data.get("username");
    let password = data.get("pass");
    let password_again = data.get("repass");
    
    if(data.get("dept") && data.get("username") && data.get("pass") && data.get("repass") == ""){
        alert("还有未填项哦");
        return
    }

    if(username.length < 2) {
        alert("用户名长度至少为2");
        return
    }
    if(password != password_again){
        alert("两次密码不一致");
        return
    }
    
    if(password.length < 6) {
        alert("密码至少为6位");
        return
    }

    
    fetch("/signup",{
        method: "post",
        body: data
    }).then((res)=>{
        res.json().then((res)=>{
            alert(res.msg)
        })
    })
})