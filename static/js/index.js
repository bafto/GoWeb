const input = document.getElementById("formInput")
const submit = document.getElementById("submit")
const foodList = document.getElementById("foodList")

async function addFood() {
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
            let response = await fetch("/api/editFood", {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({Name: foodName})
            })
            if (response.status == 200) {
                e.target.parentElement.remove()
            }
        })
        input.value = ''
    }
}

submit.addEventListener("click", async () => {
    addFood()
})
input.addEventListener("keyup", async (e) => {
    if (e.key === 'Enter') {
        addFood()
    }
})

async function getAllFood() {
    let resp = await fetch("/api/getFood")
    return resp.json()
}

async function setup() {
    let food = await getAllFood()
    console.log(food)
    food.forEach((el) => {
        const listItem = document.createElement("div")
        listItem.classList.add("listItem")
        listItem.innerHTML = `<p>${el.Name}</p><button>close</button>`
        foodList.appendChild(listItem)
        listItem.querySelector("button").addEventListener("click", async (e) => {
            e.target.parentElement.remove()
            let response = await fetch("/api/editFood", {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({Name: el.Name})
            })
            if (response.status == 200) {
                e.target.parentElement.remove()
            }
        })
    })
}
setup()