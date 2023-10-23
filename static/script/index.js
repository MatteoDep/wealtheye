if (modal !== undefined) {
    const modal = htmx.find("#modal")

    showModal = () => {
        modal.style.display = "flex"
    }

    hideModal = () => {
        modal.style.display = "none"
    }
}
