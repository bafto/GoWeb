const generateFoodButton = document.getElementById("generateFood")
const labelGrid = document.getElementById("labelGrid")

async function generateFood() {
    let labelConstraints = []
    document.querySelectorAll('.labelInput').forEach((el) => {
        if (el.checked) {
            labelConstraints.push(el.value)
        }
    })
    let foodCount = document.getElementById("foodCount").value
    foodCount = foodCount > 7 ? 7 : parseInt(foodCount)
    foodCount = foodCount < 1 ? 1 : parseInt(foodCount)
    let resp = await fetch("/api/getFoodConstrained", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "Label": labelConstraints,
            "Names": [],
            "Count": foodCount
        })
    }).then(async (r) => {return await r.json()})
    let foodList = document.getElementById('foodList').querySelector('ul')
    foodList.innerHTML = ''
    if (resp !== null) {
        resp.forEach(el => {
            let item = document.createElement('li')
            let listItem = document.createElement('div')
            listItem.classList.add('listItem')
            listItem.innerHTML = `
                <p>${el.Name}</p>
                <button class="rerollBtn"><img src="static/assets/rerollBtn.png" height="30" width="30"></button>
                `
            listItem.querySelector('.rerollBtn').addEventListener("click", async (el) => {
                let nameConstrains = []
                foodList.querySelectorAll('p').forEach(el => {
                    nameConstrains.push(el.innerText)
                })
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
                        "Names": nameConstrains,
                        "Count": 1
                    })
                }).then(async (r) => {return await r.json()})
                let p = listItem.querySelector('p').innerText = resp[0].Name
            })
            item.appendChild(listItem)
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
        let span = document.createElement('span')
        span.innerHTML = el
        let inpLabel = document.createElement('label')
        inpLabel.classList.add('grid-item')
        inpLabel.appendChild(newLabelInput)
        inpLabel.appendChild(document.createElement('br'))
        inpLabel.appendChild(span)
        labelGrid.appendChild(inpLabel)
    })
}
setup()