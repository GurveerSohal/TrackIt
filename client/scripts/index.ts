const loginButton = document.getElementById("login-button")
loginButton.addEventListener("click", (event) => {
    event.preventDefault()
    window.location.href = '/login'
})

const signupButton = document.getElementById("signup-button")
signupButton.addEventListener("click", (event) => {
    event.preventDefault()
    window.location.href = '/signup'
})