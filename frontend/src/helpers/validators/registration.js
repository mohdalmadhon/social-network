export function validateName(name = "") {
    if (name.length < 3 || name.length > 15) {
        return "Name must be between 2 and 15 characters"
    }

    const reg = /^[A-Z][a-zA-Z]*$/

    if (!reg.test(name)) {
        return "Name format is incorrect"
    }

    return ""
}

export function validateUsername(username = "") {
    if (username.length === 0) {
        return ""
    }

    if (username.length < 3 || username.length > 12) {
        return "Username must be between 3 and 12 characters"
    }

    const reg = /^[a-z0-9_-]+$/

    if (!reg.test(username)) {
        return "Username can only contain a-z, 0-9, '-' and '_'"
    }

    return ""
}

export function validateEmail(email = "") {
    if (email.length < 5 || email.length > 75) {
        return "Email must be between 5 and 75 characters"
    }

    const reg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/

    if (!reg.test(email)) {
        return "Invalid email"
    }

    return ""
}

export function validatePassword(password = "") {
    if (password.length < 8) {
        return "Password must be at least 8 characters long"
    }

    if (password.length > 75) {
        return "Password must be less than 75 characters long"
    }

    const numberReg = /[0-9]/
    const specialReg = /[!@#$%*&^?]/
    const alphaReg = /[a-zA-Z]/
    const forbiddenReg = /[^a-zA-Z0-9!@#$%*&^?]/

    if (!numberReg.test(password)) {
        return "Password must contain a number"
    }

    if (!specialReg.test(password)) {
        return "Password must contain a special character: ! @ # $ % * & ^ ?"
    }

    if (!alphaReg.test(password)) {
        return "Password must contain alphabetic characters"
    }

    if (forbiddenReg.test(password)) {
        return "Password contains forbidden characters"
    }

    return ""
}

export function validateAbout(about = "") {
    if (about.length === 0) {
        return ""
    }

    if (about.length > 1000) {
        return "About is too long"
    }

    return ""
}

export function validateDOB(dob = "") {
    if (!dob) {
        return "Date of birth is invalid"
    }

    const date = new Date(dob)
    const now = new Date()

    const tooOld = new Date()
    tooOld.setFullYear(now.getFullYear() - 120)

    const tooYoung = new Date()
    tooYoung.setFullYear(now.getFullYear() - 12)

    if (date > now || date < tooOld || date > tooYoung) {
        return "Date of birth is invalid"
    }

    return ""
}

export function handleNameInput(value = "") {
    let result = value.replace(/\s+/g, " ").trimStart()

    if (result.length > 0) {
        result = result.charAt(0).toUpperCase() + result.slice(1).toLowerCase()
    }

    return result
}

export function handleUsernameInput(value = "") {
    return value.toLowerCase().trim()
}

export function handleEmailInput(value = "") {
    return value.toLowerCase().trim()
}

export function validateAvatar(file) {
    if (!file || file.size === 0) {
        return ""
    }

    const allowedTypes = ["image/jpeg", "image/png", "image/gif"]

    if (!allowedTypes.includes(file.type)) {
        return "Avatar must be a JPG, PNG, or GIF file"
    }

    return ""
}