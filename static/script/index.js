if (modal !== undefined) {
    const modal = htmx.find("#modal")
    const dialog = htmx.find("#dialog")

    showModal = () => {
        modal.style.display = "flex"
    }

    hideModal = () => {
        modal.style.display = "none"
    }
}
