<script>
</script>

<template>
    <div class="login-container">
        <form action="" method="post">
            <div class="login-row">
                <label for="email">EMAIL *</label>
                <input type="email" name="email" id="email">
            </div>
            <div class="login-row">
                <label for="passowrd">PASSWORD *</label>
                <input type="password" name="password" id="passowrd">
            </div>
            <div class="login-row">
                <label for="fName">FIRST NAME *</label>
                <input type="text" name="fName" id="fName">
                <label for="lName">LAST NAME *</label>
                <input type="text" name="lName" id="lName">
            </div>
            <div class="login-row">
                <label for="dob">DATE OF BIRTH *</label>
                <input type="date" name="dob" id="dob">
                <label for="username">USERNAME (optional)</label>
                <input type="text" name="username" id="username">
            </div>
            <div class="login-row">
                <label for="dob">AVATAR (optional) JPG, GIF, PNG</label>
                <input type="file" name="avatar" id="avatar">
            </div>
            <div class="login-row">
                <label for="dob">ABOUT ME (optional)</label>
                <textarea name="about" id="about" maxlength="100"></textarea>
            </div>
        </form>

    </div>
</template>

<style>
    * {
        box-sizing: border-box;
    }

    .login-container {
        width: min(580px, calc(100% - 40px));
        margin: 24px auto;
        padding: 34px 39px 30px;
        background: #171b2d;
        border: 1px solid #303754;
        border-radius: 23px;
        color: #f5f5ff;
        font-family: Arial, sans-serif;
        box-shadow: 0 20px 50px rgba(0, 0, 0, 0.25);
    }

    form {
        width: 100%;
    }

    .login-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        column-gap: 9px;
        margin-bottom: 14px;
    }

    .login-row label {
        display: block;
        margin-bottom: 6px;
        color: #69749a;
        font-size: 10px;
        font-weight: 500;
        letter-spacing: 0.8px;
    }

    .login-row input,
    .login-row textarea {
        width: 100%;
        border: 1px solid transparent;
        border-radius: 7px;
        outline: none;
        background: #222741;
        color: #f4f3ff;
        font-family: inherit;
        font-size: 12px;
        transition: 0.2s ease;
    }

    .login-row input {
        height: 38px;
        padding: 0 13px;
    }

    .login-row textarea {
        height: 61px;
        padding: 12px 13px;
        resize: none;
    }

    .login-row input:focus,
    .login-row textarea:focus {
        border-color: #835cff;
        box-shadow: 0 0 0 1px rgba(131, 92, 255, 0.15);
    }

    .login-row input::placeholder,
    .login-row textarea::placeholder {
        color: #59627f;
    }

    .login-row:nth-child(1),
    .login-row:nth-child(2),
    .login-row:nth-child(5),
    .login-row:nth-child(6) {
        grid-template-columns: 1fr;
    }

    .login-row:nth-child(1) label,
    .login-row:nth-child(2) label,
    .login-row:nth-child(5) label,
    .login-row:nth-child(6) label {
        grid-column: 1;
    }

    .login-row:nth-child(1) input,
    .login-row:nth-child(2) input,
    .login-row:nth-child(5) input,
    .login-row:nth-child(6) textarea {
        grid-column: 1;
    }

    .login-row:nth-child(3) label:nth-of-type(2),
    .login-row:nth-child(4) label:nth-of-type(2) {
        grid-column: 2;
    }

    .login-row:nth-child(3) input:nth-of-type(1),
    .login-row:nth-child(4) input:nth-of-type(1) {
        grid-column: 1;
    }

    .login-row:nth-child(3) input:nth-of-type(2),
    .login-row:nth-child(4) input:nth-of-type(2) {
        grid-column: 2;
        grid-row: 2;
    }

    .login-row:nth-child(3) label:nth-of-type(1),
    .login-row:nth-child(4) label:nth-of-type(1) {
        grid-column: 1;
    }

    .login-row:nth-child(3),
    .login-row:nth-child(4) {
        grid-template-rows: auto 38px;
    }

    .login-row:nth-child(5) {
        margin-top: 1px;
    }

    #avatar {
        height: 48px;
        padding: 0;
        border: 2px dashed #805cff;
        background: #20243c;
        cursor: pointer;
        color: #d8d7e9;
    }

    #avatar::file-selector-button {
        height: 100%;
        margin-right: 12px;
        padding: 0 15px;
        border: 0;
        background: transparent;
        color: #916cff;
        font-weight: 600;
        cursor: pointer;
    }

    #about {
        min-height: 61px;
    }

    @media (max-width: 560px) {
        .login-container {
            width: calc(100% - 24px);
            padding: 28px 22px;
            border-radius: 18px;
        }

        .login-row:nth-child(3),
        .login-row:nth-child(4) {
            display: block;
        }

        .login-row:nth-child(3) label,
        .login-row:nth-child(4) label {
            display: block;
            margin-top: 10px;
        }

        .login-row:nth-child(3) input,
        .login-row:nth-child(4) input {
            display: block;
            margin-bottom: 8px;
        }
    }
</style>