import './styles/bootstrap.min.css'
import './styles/style.css'
import 'dotenv/config'

export const server_url = process.env.SERVER_URL

const loginButton = document.getElementById("login-button")
loginButton.addEventListener("click", (event) => {
    event.preventDefault()
    window.location.href = '/login'
})

const signupButton = document.getElementById("signup-button")
signupButton.addEventListener("click", (event) => {
    event.preventDefault()
    alert('signup')
})