const input = document.getElementById("formInput")
const submit = document.getElementById("submit")
const foodList = document.getElementById("foodList")

submit.addEventListener("click", async () => {
    let foodName = input.value;
    let resp = await fetch("/api/editFood", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({Name: foodName})
    })
    if (resp.status == 200) {
        const listItem = document.createElement("div")
        listItem.classList.add("listItem")
        listItem.innerHTML = `<p>${foodName}</p><button>close</button>`
        foodList.appendChild(listItem)
        listItem.querySelector("button").addEventListener("click", async (e) => {
            e.target.parentElement.remove()
        })
    }
})

async function postFood(food) {
    let resp = await fetch("/api/editFood", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({Name: food})
    })
    return resp.json()
}

async function getAllFood() {
    let resp = await fetch("/static/foods.json")
    return resp.json()
}

async function setup() {
    let food = await getAllFood()
    console.log(food)
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