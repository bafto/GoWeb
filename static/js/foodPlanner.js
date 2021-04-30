const generateFoodButton = document.getElementById("generateFood")
const labelGrid = document.getElementById("labelGrid")

async function generateFood() {
    let labelConstraints = []
    document.querySelectorAll('.labelInput').forEach((el) => {
        if (el.checked) {
            labelConstraints.push(el.value)
        }
    })
    let resp = await fetch("/api/getFoodConstrained", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "Label": labelConstraints,
            "Count": 5
        })
    }).then(async (r) => {return await r.json()})
    let foodList = document.getElementById('foodList').querySelector('ul')
    foodList.innerHTML = ''
    if (resp !== null) {
        resp.forEach(el => {
            let item = document.createElement('li')
            item.innerHTML = el.Name
            foodList.appendChild(item)
        })
    } else {
        let item = document.createElement('li')
        item.innerHTML = "No Food with the specified restrictions found"
        foodList.appendChild(item)
    }
}

generateFoodButton.addEventListener("click", async (ev) => {
    await generateFood()
})

async function setup() {
    let label = await fetch("/api/getLabel", {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    }).then( async (r) => {return await r.json()})
    label.forEach(async (el) => {
        let newLabelInput = document.createElement('input')
        newLabelInput.classList.add('labelInput')
        newLabelInput.type = 'checkbox'
        newLabelInput.value = el
        newLabelInput.checked = false
        let inpLabel = document.createElement('label')
        inpLabel.appendChild(newLabelInput)
        inpLabel.innerHTML += el
        labelGrid.appendChild(inpLabel)
    })
}
setup()