const app = () => {
    document.querySelector(".js-trigger-products-list").addEventListener("click", function (e) {
        this.classList.toggle("active")
        document.querySelector(".nav-products-list").classList.toggle("hide")
    })
}

document.addEventListener("DOMContentLoaded", app)