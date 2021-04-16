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
                Label: [
                    'Label1',
                    'Label2',
                    'Label3'
                ]
            })
        })
        let respJson = await resp.json()
        if (resp.status == 200) {
            const listItem = document.createElement("div")
            listItem.classList.add("listItem")
            listItem.innerHTML = `
            <div class="listItemHeader">
            <p>${respJson.Name}</p><button>close</button>
            </div>
            <div class="labelList"></div>
            `
            let labelList = listItem.querySelector(".labelList")
            respJson.Label.forEach((e) => {
                labelList.innerHTML = labelList.innerHTML + `<p>${e}</p>`
            })
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
                    e.target.parentElement.parentElement.remove()
                }
            })
            input.value = ''
    }
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
    food.forEach((el) => {
        const listItem = document.createElement("div")
            listItem.classList.add("listItem")
            listItem.innerHTML = `
            <div class="listItemHeader">
            <p>${el.Name}</p><button>close</button>
            </div>
            <div class="labelList"></div>
            `
            let labelList = listItem.querySelector(".labelList")
            el.Label.forEach((e) => {
                labelList.innerHTML = labelList.innerHTML + `<p>${e}</p>`
            })
        foodList.appendChild(listItem)
        listItem.querySelector("button").addEventListener("click", async (e) => {
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
                e.target.parentElement.parentElement.remove()
            }
        })
    })
}
setup()