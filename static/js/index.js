const input = document.getElementById("formInput")
const submit = document.getElementById("submit")
const foodList = document.getElementById("foodList")

submit.addEventListener("click", () => {
    let foodName = input.value;
    const listItem = document.createElement("div")
    listItem.classList.add("listItem")
    listItem.innerHTML = `<p>${foodName}</p><button>close</button>`
    foodList.appendChild(listItem)
    listItem.querySelector("button").addEventListener("click", async (e) => {
        e.target.parentElement.remove()
        food = await getAllFood()
        console.log(food[0])
    })
})

const buttons = foodList.querySelectorAll("button")

buttons.forEach((btn) => {
    btn.addEventListener('click', async (e) => {
        e.target.parentElement.remove()
        food = await getAllFood()
        console.log(food[0])
    })
})

async function getAllFood() {
    let resp = await fetch("/static/foods.json")
    let json = await resp.json()
    return json
}