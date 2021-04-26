const generateFoodButton = document.getElementById("generateFood")
const labelGrid = document.getElementById("labelGrid")

async function generateFood() {
    let labelConstraints = []
    document.querySelectorAll('.labelInput').forEach((el) => {
        if (el.checked) {
            labelConstraints.push(el.value)
        }
    })
    console.log(JSON.stringify({
        method: 'Get',
        headers: {
            'Content-Type': 'application/json'
        },
        body: labelConstraints
    }))
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