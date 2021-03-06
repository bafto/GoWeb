const foodInput = document.getElementById("addFoodInput")
const foodSubmit = document.getElementById("addFoodSubmit")
const foodList = document.getElementById("foodList")
let allLabel

async function addFoodToDocument(Food) {
    const listItem = document.createElement("div") //the div to hold the food 
    listItem.classList.add("listItem")
    //add the boilerplate html
    listItem.innerHTML = ` 
    <button class="listItemHeader">
        <p>${Food.Name}</p>
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
        <div class="addLabelDiv">
            <div class="autocomplete">
                <input class="addLabelInput" type="text" name="myLabel" placeholder="Add a Label">
            </div>
            <button class="addLabelSubmit styledButton">Submit</button>
        </div>
        <div class="labelList"></div>
        <button class="removeBtn foodRemove"><img src="static/assets/removeBtn.png" height="30" width="30"></button>
    </div>
    `
    autocomplete(listItem.querySelector(".addLabelInput"), allLabel)
    let currentLabel = Food.Label //holds the label state for the food
    let labelList = listItem.querySelector(".labelList")
    Food.Label.forEach((e) => { //add each label to the food
        labelList.innerHTML = labelList.innerHTML + `
        <div class="Label">
            <p>${e}</p>
            <button class="removeBtn labelRemove" value="${e}"><img src="static/assets/removeBtn.png" height="20" width="20"></button>
        </div>
        `
    })
    listItem.querySelector(".addLabelSubmit").addEventListener("click", async (e) => {
        let newLabel = listItem.querySelector(".addLabelInput").value
        let arr = currentLabel
        if (arr.includes(newLabel)) {
            listItem.querySelector(".addLabelInput").value = ''
            return
        }
        arr.push(newLabel)
        let resp = await fetch("/api/changeFood", {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: Food.ID,
                Name: Food.Name,
                Label: arr
            })
        })
        if (resp.status == 200) {
            currentLabel = arr
            listItem.querySelector(".addLabelInput").value = ''
            let div = document.createElement('div')
            div.classList.add('Label')
            div.innerHTML = `
            <p>${newLabel}</p>
            <button class="removeBtn labelRemove" value="${newLabel}"><img src="static/assets/removeBtn.png" height="20" width="20"></button>
            `
            let labelRemoveBtn = div.querySelector(".labelRemove")
            labelRemoveBtn.addEventListener("click", async (t) => {
                //remove the label from the label state
                currentLabel = currentLabel.filter((el) => {return el !== labelRemoveBtn.value})
                //send the POST request to remove the label in the backend
                let response = await fetch("/api/changeFood", {
                    method: 'PATCH',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        ID: Food.ID,
                        Name: Food.Name,
                        Label: currentLabel
                    })
                })
                //on success remove the label from the food in the DOM
                if (response.status == 200) {
                    labelRemoveBtn.parentElement.remove()
                }
            })
            labelList.appendChild(div)
        } else if (resp.status == 304) {
            listItem.querySelector(".addLabelInput").value = ''
        }
    })
    listItem.querySelectorAll(".labelRemove").forEach( async (ev) => {
        ev.addEventListener("click", async (t) => {
            //remove the label from the label state
            currentLabel = currentLabel.filter((el) => {return el !== ev.value})
            //send the POST request to remove the label in the backend
            let response = await fetch("/api/changeFood", {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    ID: Food.ID,
                    Name: Food.Name,
                    Label: currentLabel
                })
            })
            //on success remove the label from the food in the DOM
            if (response.status == 200) {
                ev.parentElement.remove()
            }
        })
    })
    //make the listItem expandable
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
    listItem.querySelector(".foodRemove").addEventListener("click", async (e) => { //event listener to remove the food
        let response = await fetch("/api/editFood", { //we inform the backend
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ID: Food.ID,
                Name: Food.Name,
                Label: Food.Label
            })
        })
        if (response.status == 200) { //on success we remove the food from the DOM
            let parent = e.target.parentElement
            while (!parent.classList.contains('listItem')) {
                parent = parent.parentElement
            }
            parent.remove()
        }
    })
}

//called when a food is added
async function addFood() {
    let foodName = foodInput.value;
    if (foodName.length != 0) {
        //Post the new Food to the backend
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
        //if the post was a success we add the new Food to the DOM
        if (resp.status == 200) {
            await addFoodToDocument(respJson)
            foodInput.value = ''
        }
    }
}

//add the event listeners to add a food
foodSubmit.addEventListener("click", async () => {
    await addFood()
})
foodInput.addEventListener("keyup", async (e) => {
    if (e.key === 'Enter') {
        await addFood()
    }
})

//Get all food as json array from the backend
async function getAllFood() {
    let resp = await fetch("/api/getFood")
    return resp.json()
}

//called on page load
async function setup() {
    allLabel = await fetch("/api/getLabel", {
        method: 'GET'
    }).then(async (r) => {
        return await r.json()
    })
    let food = await getAllFood()
    food.forEach(async (el) => { //we add each food to the DOM
        await addFoodToDocument(el)
    })
}
setup()