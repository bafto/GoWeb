const labelInput = document.getElementById("addLabelInput")
const labelSubmit = document.getElementById("addLabelSubmit")
const labelList = document.getElementById("labelList")

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
            <button class="removeBtn">remove</button>
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
                el.target.parentElement.remove()
            }
        })
    })
}
setup()