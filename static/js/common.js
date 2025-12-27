// Common JavaScript functions used across the app

function toggleVipFields() {
    const vipCheckbox = document.getElementById('vip-checkbox');
    const vipFields = document.querySelectorAll('.vip-field');

    vipFields.forEach(field => {
        field.disabled = !vipCheckbox.checked;
        if (vipCheckbox.checked) {
            field.classList.remove('bg-gray-100');
        } else {
            field.classList.add('bg-gray-100');
        }
    });
}
