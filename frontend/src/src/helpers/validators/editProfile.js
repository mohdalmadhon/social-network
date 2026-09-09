export function validateAboutField(value = '') {
    if (value.length > 200) {
        return 'Must not exceed 200 characters';
    }

    return '';
}