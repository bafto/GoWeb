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
    })
})


async function getAllFood() {
    let resp = await fetch("/static/foods.json")
    let json = await resp.json()
    return json
}

async function setup() {
    let food = await getAllFood()
    food.forEach((e) => {
        const listItem = document.createElement("div")
        listItem.classList.add("listItem")
        listItem.innerHTML = `<p>${e.Name}</p><button>close</button>`
        foodList.appendChild(listItem)
        listItem.querySelector("button").addEventListener("click", (e) => {
            e.target.parentElement.remove()
        })
    })
}
setup()