const labelInput = document.getElementById("addLabelInput")
const labelSubmit = document.getElementById("addLabelSubmit")
const labelList = document.getElementById("labelList")

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
            const listItem = document.createElement('div')
            listItem.classList.add('listItem')
            listItem.innerHTML = `
                <p>${labelName}</p>
                <button class="removeBtn"><img src="static/assets/removeBtn.png" height="30" width="30"></button>
                `
            labelList.appendChild(listItem)
            listItem.querySelector('.removeBtn').addEventListener("click", async (el) => {
            let resp = await fetch("/api/editLabel", {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(labelName)
            })
            if (resp.status == 200) {
                let parent = el.target.parentElement
                while (!parent.classList.contains('listItem')) {
                    parent = parent.parentElement
                }
                parent.remove()
            }
        })
        }
    }
    labelInput.value = ''
}

labelSubmit.addEventListener("click", async () => {
    addLabel()
})
labelInput.addEventListener("keyup", async (e) => {
    if (e.key === 'Enter') {
        addLabel()
    }
})

async function setup() {
    let resp = await fetch("/api/getLabel", {
        method: 'GET'
    }).then((r) => {
       return r.json()
    })
    console.log(resp)
    resp.forEach((e) => {
        const listItem = document.createElement('div')
        listItem.classList.add('listItem')
        listItem.innerHTML = `
            <p>${e}</p>
            <button class="removeBtn"><img src="static/assets/removeBtn.png" height="30" width="30"></button>
        `
        labelList.appendChild(listItem)
        listItem.querySelector('.removeBtn').addEventListener("click", async (el) => {
            let resp = await fetch("/api/editLabel", {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(e)
            })
            if (resp.status == 200) {
                let parent = el.target.parentElement
                while (!parent.classList.contains('listItem')) {
                    parent = parent.parentElement
                }
                parent.remove()
            }
        })
    })
}
setup()