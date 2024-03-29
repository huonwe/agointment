async function login() {
    username = document.querySelector("#username").value;
    password = document.querySelector("#password").value;
    const formData = new FormData();
    formData.append('username', username);
    formData.append('password', password);
    // console.log(formData.getAll("username"));
    const response = await fetch("/login" , {
        method: 'post',
        body: formData
    })
    const res = await response.json();
    if(res.status == "Success") {
        window.location = "/home/index"
    }else {
        alert(res.msg)
    }
}