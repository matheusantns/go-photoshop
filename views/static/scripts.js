function validateForm(event, elementName) {
    var checkboxes = document.getElementsByName(elementName);
    var isChecked = false;
    for (var i = 0; i < checkboxes.length; i++) {
        if (checkboxes[i].checked) {
            isChecked = true;
            break;
        }
    }

    if (!isChecked) {
        document.getElementById('form-error').style.display = 'flex';
        event.preventDefault();
        return false;
    }

    return true;
}
