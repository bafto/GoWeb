const foodInput = document.getElementById("addFoodInput")
const foodSubmit = document.getElementById("addFoodSubmit")
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
                Label: []
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
                <div class="addLabelDiv">
                    <h3>Add a Label</h3>
                    <input class="addLabelInput" type="text">
                </div>
                <div class="labelList"></div>
                <button class="removeBtn foodRemove"><img src="static/assets/removeBtn.png" height="30" width="30"></button>
            </div>
            `
            currentLabel = respJson.Label //holds the label state for the food
            let labelList = listItem.querySelector(".labelList")
            respJson.Label.forEach((e) => { //add each label to the food
                labelList.innerHTML = labelList.innerHTML + `
                <div class="Label">
                    <p>${e}</p>
                    <button class="removeBtn labelRemove"><img src="static/assets/removeBtn.png" height="20" width="20"></button>
                </div>
                `
                listItem.querySelectorAll(".labelRemove").forEach( async (e) => {
                    e.addEventListener("click", async (t) => {
                        //remove the label from the label state
                        let i = currentLabel.indexOf(e)
                        if (i != -1) {
                            currentLabel.splice(i, 1)
                        }
                        //send the POST request to remove the label in the backend
                        let response = await fetch("api/changeFood", {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json'
                            },
                            body: JSON.stringify({
                                ID: el.ID,
                                Name: el.Name,
                                Label: currentLabel
                            })
                        })
                        //on success remove the label from the food in the DOM
                        if (response.status == 200) {
                            e.parentElement.remove()
                        }
                    })
                })
            })
            //make the listItem expandable
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
            listItem.querySelector(".foodRemove").addEventListener("click", async (e) => { //the event listener to remove the food
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
                    let parent = e.target.parentElement
                    while (!parent.classList.contains('listItem')) {
                        parent = parent.parentElement
                    }
                    parent.remove()
                }
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
                <div class="addLabelDiv">
                    <h3>Add a Label</h3>
                    <input class="addLabelInput" type="text">
                </div>
                <div class="labelList"></div>
                <button class="removeBtn foodRemove"><img src="static/assets/removeBtn.png" height="30" width="30"></button>
            </div>
            `
            let currentLabel = el.Label //holds the label state for the food
            let labelList = listItem.querySelector(".labelList")
            el.Label.forEach((e) => { //add each label to the food
                labelList.innerHTML = labelList.innerHTML + `
                <div class="Label">
                    <p>${e}</p>
                    <button class="removeBtn labelRemove" value="${e}"><img src="static/assets/removeBtn.png" height="20" width="20"></button>
                </div>
                `
            })
            listItem.querySelectorAll(".labelRemove").forEach( async (ev) => {
                ev.addEventListener("click", async (t) => {
                    //remove the label from the label state
                    currentLabel = currentLabel.filter((el) => {return el !== ev.value})
                    //send the POST request to remove the label in the backend
                    let response = await fetch("api/changeFood", {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            ID: el.ID,
                            Name: el.Name,
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
                    ID: el.ID,
                    Name: el.Name,
                    Label: el.Label
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
    })
}
setup()