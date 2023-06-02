var submit = document.getElementById("submit")

var all = fetch("/users")
.then(response => response.json())
.then(data => {
    if (data.id != "") {
        submit.innerHTML += `<input type="submit" value="submit" disabled>`
    }
    else {
        submit.innerHTML += `<input type="submit" value="submit">`
    }
})

