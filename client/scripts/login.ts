function loginUser(event) {
    event.preventDefault()
    alert('now logged in')
}

const form = document.getElementById('login-form')
form?.addEventListener("submit", loginUser)

