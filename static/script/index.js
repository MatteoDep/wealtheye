const canHideMap = new(Map)
function showOverlay(id) {
    canHideMap.set(id) = true
    document.getElementById(id).style.display = "block";
}

function hideOverlay(id) {
    if (canHideMap.get(id)) {
        document.getElementById(id).style.display = "none";
    }
}

function blockHiding(id) {
    canHideMap.set(id) = false
}

function allowHiding(id) {
    canHideMap.set(id) = true
}
