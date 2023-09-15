const historyButton = document.getElementById("history-button")
historyButton.addEventListener("click", (event) => {
    event.preventDefault()
    window.location.href = '/history'
})

const createButton = document.getElementById("create-button")
createButton.addEventListener("click", (event) => {
    event.preventDefault()
    window.location.href = '/create'
})