const productsPage = () => {
    let dataTable
    const productsListTable = document.querySelector("#productsList")
    if (productsListTable) {
        dataTable = new simpleDatatables.DataTable(productsListTable, {
            fixedHeight: false,
        })
    }
}

document.addEventListener("DOMContentLoaded", productsPage)