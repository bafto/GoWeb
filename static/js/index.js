const input = document.getElementById("formInput")
const submit = document.getElementById("submit")
const foodList = document.getElementById("foodList")

submit.addEventListener("click", () => {
    let foodName = input.value;
    const listItem = document.createElement("div")
    listItem.classList.add("listItem")
    listItem.innerHTML = `<p>${foodName}</p><button>close</button>`
    foodList.appendChild(listItem)
    listItem.querySelector("button").addEventListener("click", (e) => {
        e.target.parentElement.remove()
    })
})

const buttons = foodList.querySelectorAll("button")

buttons.forEach((btn) => {
    btn.addEventListener('click', (e) => {
        e.target.parentElement.remove()
    })
})
