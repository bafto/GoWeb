const input = document.getElementById("formInput")
const submit = document.getElementById("submit")
const foodList = document.getElementById("foodList")

async function addFood() {
    let foodName = input.value;
    if (foodName.length != 0)
    {
        let resp = await fetch("/api/editFood", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: '',
                Name: foodName,
                Label: {
                    'Label1': false,
                    'Label2': false,
                    'Label3': false
                }
            })
        })
        let respJson = await resp.json()
        if (resp.status == 200) {
            const listItem = document.createElement("div")
            listItem.classList.add("listItem")
            listItem.innerHTML = `
            <button class="listItemHeader">
                <p>${respJson.Name}</p>
                <svg
                    class="arrow"
                    viewBox="0 0 266 438"
                    xmlns="http://www.w3.org/2000/svg">
                    <path
                        class="pathColored"
                        d="m258.476 235.971-194.344 194.343c-9.373 9.373-24.569 9.373-33.941 0l-22.667-22.667c-9.357-9.357-9.375-24.522-.04-33.901l154.021-154.746-154.021-154.745c-9.335-9.379-9.317-24.544.04-33.901l22.667-22.667c9.373-9.373 24.569-9.373 33.941 0l194.343 194.343c9.373 9.372 9.373 24.568.001 33.941z"
                    />
                </svg>
            </button>
            <div class="listItemContent">
                <div class="labelList"></div>
                <button class="saveChanges">save changes</button>
                <br>
                <button class="removeBtn">remove</button>
            </div>
            `
            let labelList = listItem.querySelector(".labelList")
            for (let key of Object.keys(respJson.Label)) {
                labelList.innerHTML = labelList.innerHTML + `
                <div class="Label">
                    <label class="checkLabel">
                        <input type="checkbox">
                        ${key}
                    </label>
                </div>
                `
            }
            foodList.appendChild(listItem)
            listItem.querySelector(".listItemContent").classList.add("collapsed")
            listItem.querySelector(".listItemHeader").addEventListener("click", () => {
                let content = listItem.querySelector(".listItemContent")
                if (content.classList.contains("collapsed")) {
                    content.classList.remove("collapsed")
                } else {
                    content.classList.add("collapsed")
                }

                let svg = listItem.querySelector("svg")
                if (svg.classList.contains("rotated90")) {
                    svg.classList.remove("rotated90")
                } else {
                    svg.classList.add("rotated90")
                }
            })
            listItem.querySelector(".removeBtn").addEventListener("click", async (e) => {
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
    console.log(food)
    food.forEach((el) => {
        const listItem = document.createElement("div")
            listItem.classList.add("listItem")
            listItem.innerHTML = `
            <button class="listItemHeader">
                <p>${el.Name}</p>
                <svg
                    class="arrow"
                    viewBox="0 0 266 438"
                    xmlns="http://www.w3.org/2000/svg">
                    <path
                        class="pathColored"
                        d="m258.476 235.971-194.344 194.343c-9.373 9.373-24.569 9.373-33.941 0l-22.667-22.667c-9.357-9.357-9.375-24.522-.04-33.901l154.021-154.746-154.021-154.745c-9.335-9.379-9.317-24.544.04-33.901l22.667-22.667c9.373-9.373 24.569-9.373 33.941 0l194.343 194.343c9.373 9.372 9.373 24.568.001 33.941z"
                    />
                </svg>
            </button>
            <div class="listItemContent">
                <div class="labelList"></div>
                <button class="saveChanges">save changes</button>
                <br>
                <button class="removeBtn">remove</button>
            </div>
            `
            let labelList = listItem.querySelector(".labelList")
            for (let key of Object.keys(el.Label)) {
                labelList.innerHTML = labelList.innerHTML + `
                <div class="Label">
                    <label class="checkLabel">
                        <input type="checkbox">
                        ${key}
                    </label>
                </div>
                `
            }
            let divs = labelList.querySelectorAll('.Label')
            let i = 0
            for (let key of Object.keys(el.Label)) {
                divs[i].querySelector('input').checked = el.Label[key]
                i++
            }
        foodList.appendChild(listItem)
        listItem.querySelector(".listItemContent").classList.add("collapsed")
        listItem.querySelector(".listItemHeader").addEventListener("click", () => {
            let content = listItem.querySelector(".listItemContent")
            if (content.classList.contains("collapsed")) {
                content.classList.remove("collapsed")
            } else {
                content.classList.add("collapsed")
            }

            let svg = listItem.querySelector("svg")
            if (svg.classList.contains("rotated90")) {
                svg.classList.remove("rotated90")
            } else {
                svg.classList.add("rotated90")
            }
        })
        listItem.querySelector(".removeBtn").addEventListener("click", async (e) => {
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