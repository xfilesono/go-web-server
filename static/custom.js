var list = document.getElementById("list")
var submit = document.getElementById("submit")

var all = fetch("/users")
.then(response => response.json())
.then(data => {
    if (data.id != "") {
        list.innerHTML += `Registration Succesful..<br>Your informations added to our system.<br>`
        list.innerHTML += `<td>${data.fullname}</td><br>`
        list.innerHTML += `<td>${data.phone}</td><br>`
        list.innerHTML += `<td>${data.mail}</td><br><br><br>`
        list.innerHTML += `<a href="#">Sign-in</a>
        <a href="./movies.html">Movies</a>
        <a href="./addmovie.html">Add a Movie</a>`
    }
    else {
        list.innerHTML += `Welcome to CapsCode's Movie Database<br>
        Use below form to sign-up<br><br>
        <a href="#">Sign-in</a>
        <a href="./movies.html">Movies</a>
        <a href="./addmovie.html">Add a Movie</a>`
    }
})

