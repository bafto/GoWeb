const foodInput = document.getElementById("addFoodInput")
const foodSubmit = document.getElementById("addFoodSubmit")
const labelInput = document.getElementById("addLabelInput")
const labelSubmit = document.getElementById("addLabelSubmit")
const foodList = document.getElementById("foodList")

//called when a food is added
async function addFood() {
    let foodName = foodInput.value;
    if (foodName.length != 0)
    {
        //Post the new Food to the backend
        let resp = await fetch("/api/editFood", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: '',
                Name: foodName,
                Label: {}
            })
        })
        let respJson = await resp.json()
        //if the post was a success we add the new Food to the DOM
        if (resp.status == 200) {
            const listItem = document.createElement("div") //the div that will hold the Food
            listItem.classList.add("listItem")
            //add boilerplate html
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
                <button class="removeBtn">remove</button>
            </div>
            `
            //add the labels with checkboxes
            let labelList = listItem.querySelector(".labelList")
            for (let key of Object.keys(respJson.Label)) {
                labelList.innerHTML = labelList.innerHTML + `
                <div class="Label">
                    <label class="checkLabel">
                        <input type="checkbox" class="labelInput" value="${key}">
                        ${key}
                    </label>
                </div>
                `
            }
            foodList.appendChild(listItem)
            listItem.querySelector(".listItemContent").classList.add("collapsed") //by default the new food is collapsed
            listItem.querySelector(".listItemHeader").addEventListener("click", () => { //the event listener to open and close the accordion
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
            listItem.querySelector(".removeBtn").addEventListener("click", async (e) => { //the event listener to remove the food
                //Delete the food from the backend
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
                if (response.status == 200) { //on success we remove the food from the DOM
                    e.target.parentElement.parentElement.remove()
                }
            })
            listItem.querySelectorAll(".labelInput").forEach((el) => { //if a label state is changed we inform the backend
                el.addEventListener("change", async () => {
                    let Label = {}
                    listItem.querySelectorAll(".labelInput").forEach((e) => {
                        Label[e.value] = e.checked
                    })
        
                    let response = await fetch("/api/changeFood", {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            ID: respJson.ID,
                            Name: respJson.Name,
                            Label: Label
                        })
                    })
                    if (response.status != 200) {
                        console.log('Error saving changes. Add a revert feature')
                    }
                })
            })
            foodInput.value = ''
    }
    }
}

//add the event listeners to add a food
foodSubmit.addEventListener("click", async () => {
    addFood()
})
foodInput.addEventListener("keyup", async (e) => {
    if (e.key === 'Enter') {
        addFood()
    }
})

//called when a label is added to the label list
async function addLabel() {
    let labelName = labelInput.value
    if (labelName.length != 0) {
        let resp = await fetch("/api/editLabel", { //post the new Label to the backend
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(labelName)
        })
        if (resp.status == 200) { //on success we add the label to every food in the DOM
            let labelLists = foodList.querySelectorAll(".labelList")
            labelLists.forEach((e, i) => {
                let label = {}
                labelLists[i].querySelectorAll('.labelInput').forEach((e) => {
                    label[e.value] = e.checked
                })
                labelLists[i].innerHTML = labelLists[i].innerHTML + `
                <div class="Label">
                    <label class="checkLabel">
                        <input type="checkbox" class="labelInput" value="${labelName}">
                        ${labelName}
                    </label>
                </div>
                `
                let labelInputs = labelLists[i].querySelectorAll('.labelInput')
                    labelInputs.forEach((e) => {
                        e.checked = label[e.value]
                    })
            })
        }
    }
    labelInput.value = ''
}

//add the event listeners to add a label
labelSubmit.addEventListener("click", async () => {
    addLabel()
})
labelInput.addEventListener("keyup", async (e) => {
    if (e.key === 'Enter') {
        addLabel()
    }
})

//Get all food as json array from the backend
async function getAllFood() {
    let resp = await fetch("/api/getFood")
    return resp.json()
}

//called on page load
async function setup() {
    let food = await getAllFood()
    console.log(food)
    food.forEach((el) => { //we add each food to the DOM
        const listItem = document.createElement("div") //the div to hold the food 
            listItem.classList.add("listItem")
            //add the boilerplate html
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
                <button class="removeBtn">remove</button>
            </div>
            `
            //add the labels to the food
            let labelList = listItem.querySelector(".labelList")
            for (let key of Object.keys(el.Label)) {
                labelList.innerHTML = labelList.innerHTML + `
                <div class="Label">
                    <label class="checkLabel">
                        <input type="checkbox" class="labelInput" value="${key}">
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
        listItem.querySelector(".listItemContent").classList.add("collapsed") //by default everything is collapsed
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
        listItem.querySelector(".removeBtn").addEventListener("click", async (e) => { //event listener to remove the food
            let response = await fetch("/api/editFood", { //we inform the backend
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
            if (response.status == 200) { //on success we remove the food from the DOM
                e.target.parentElement.parentElement.remove()
            }
        })
        listItem.querySelectorAll(".labelInput").forEach((element) => { //event listener when a label is changed
            element.addEventListener("change", async () => {
                let Label = {}
                listItem.querySelectorAll(".labelInput").forEach((e) => {
                    Label[e.value] = e.checked
                })
                let response = await fetch("/api/changeFood", {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        ID: el.ID,
                        Name: el.Name,
                        Label: Label
                    })
                })
                if (response.status != 200) {
                    console.log('Error saving changes. Add a revert feature')
                }
            })
        })
    })
}
setup()