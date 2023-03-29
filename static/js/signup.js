window.onload = function () {
    // (*^_^*)
    const len_username_atLeast = 2;
    const len_dept_atLeast = 2;
    const len_pass_atLeast = 6;
    let input_username = document.querySelector("#username")
    let input_dept = document.querySelector("#dept")
    let input_pass = document.querySelector("#pass")
    let input_repass = document.querySelector("#repass")
    input_username.addEventListener("input", () => {
        let value = input_username.value;
        if(value.length < len_username_atLeast) {
            input_username.classList.remove("outline_ok");
            input_username.classList.add("outline_attention");
        } else {
            input_username.classList.remove("outline_attention");
            input_username.classList.add("outline_ok");
        }
    })

    input_dept.addEventListener("input", () => {
        let value = input_dept.value;
        if(value.length < len_dept_atLeast) {
            input_dept.classList.remove("outline_ok");
            input_dept.classList.add("outline_attention");
        } else {
            input_dept.classList.remove("outline_attention");
            input_dept.classList.add("outline_ok");
        }
    })

    input_pass.addEventListener("input", () => {
        let value = input_pass.value;
        if(value.length < len_pass_atLeast) {
            input_pass.classList.remove("outline_ok");
            input_pass.classList.add("outline_attention");
        } else {
            input_pass.classList.remove("outline_attention");
            input_pass.classList.add("outline_ok");
        }
    })

    input_repass.addEventListener("input", () => {
        let revalue = input_repass.value;
        let ovalue = input_pass.value;
        if(revalue != ovalue) {
            input_repass.classList.remove("outline_ok");
            input_repass.classList.add("outline_attention");
        } else {
            input_repass.classList.remove("outline_attention");
            input_repass.classList.add("outline_ok");
        }
    })

    let form = document.querySelector("#signupForm")
    form.addEventListener("submit", (ev) => {
        ev.preventDefault()
        let data = new FormData(form);
        let username = data.get("username");
        let password = data.get("pass");
        let password_again = data.get("repass");

        if (username.length < len_username_atLeast) {
            alert("用户名长度至少为2");
            return
        }
        if (dept.length < len_dept_atLeast) {
            alert("部门名称长度至少为2");
            return
        }
        if (password != password_again) {
            alert("两次密码不一致");
            return
        }

        if (password.length < len_pass_atLeast) {
            alert("密码至少为6位");
            return
        }


        fetch("/signup", {
            method: "post",
            body: data
        }).then((res) => {
            res.json().then((res) => {
                alert(res.msg)
                if(res.status == "Success"){
                    window.location = "/login";
                }
            })
        })
    })

}
