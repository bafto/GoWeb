const input = document.getElementById("formInput")
const submit = document.getElementById("submit")
const foodList = document.getElementById("foodList")

async function addFood() {
    let foodName = input.value;
    if(foodName.length != 0)
    {
        let resp = await fetch("/api/editFood", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: '',
                Name: foodName,
                Label: []
            })
        })
        let respJson = await resp.json()
        console.log(respJson.ID)
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
                    body: JSON.stringify({
                        ID: respJson.ID,
                        Name: respJson.Name,
                        Label: respJson.Label
                    })
                })
                if (response.status == 200) {
                    e.target.parentElement.remove()
                }
            })
            input.value = ''
    }
    }
}

submit.addEventListener("click", async () => {
    addFood()
    console.log("added Food")
})
input.addEventListener("keyup", async (e) => {
    if (e.key === 'Enter') {
        addFood()
        console.log("added Food")
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
                body: JSON.stringify({
                    ID: el.ID,
                    Name: el.Name,
                    Label: el.Label
                })
            })
            if (response.status == 200) {
                e.target.parentElement.remove()
            }
        })
    })
}
setup()