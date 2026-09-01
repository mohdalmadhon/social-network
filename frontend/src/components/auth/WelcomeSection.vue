<template>
    <section class="welcome-section">

        <div class="welcome-content">

            <div class="logo">N</div>

            <p class="eyebrow">
                WELCOME BACK
            </p>

            <h1>
                Your people.<br>
                Your <span>place.</span>
            </h1>

            <p class="welcome-text">
                Connect with people, share your thoughts,
                and keep up with the things that matter to you.
            </p>

            <div class="welcome-line"></div>

            <p class="small-text">
                A simple place to stay connected.
            </p>

        </div>

        <div class="decoration decoration-one"></div>
        <div class="decoration decoration-two"></div>
        <div class="decoration decoration-three"></div>

        <div class="circle circle-one"></div>
        <div class="circle circle-two"></div>

        <div class="grid-pattern"></div>

    </section>
</template>

<style scoped>
.welcome-section {
    position: relative;

    min-height: 100vh;

    display: flex;
    align-items: center;

    padding: 80px;

    background:
        radial-gradient(circle at 20% 35%, rgba(75, 63, 160, 0.12), transparent 35%),
        #100d2b;

    overflow: hidden;

    font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.welcome-content {
    position: relative;
    z-index: 5;

    max-width: 650px;
}

.logo {
    width: 55px;
    height: 55px;

    display: flex;
    align-items: center;
    justify-content: center;

    margin-bottom: 55px;

    background: linear-gradient(135deg, #805cff, #ff6288);

    border-radius: 14px;

    box-shadow: 0 12px 30px rgba(0, 0, 0, 0.3);

    color: white;

    font-family: Inter, sans-serif;
    font-size: 24px;
    font-weight: 700;
}

.eyebrow {
    margin: 0 0 15px;

    color: #ff6288;

    font-family: Arial, sans-serif;

    font-size: 13px;
    font-weight: 600;

    letter-spacing: 3px;
}

h1 {
    margin: 0;

    color: #ffffff;

    font-family: Inter, sans-serif;

    font-size: clamp(52px, 6vw, 88px);

    line-height: 1.05;

    letter-spacing: -1.7px;

    font-weight: 750;
}

h1 span {
    color: #7b55ff;
}

.welcome-text {
    max-width: 480px;

    margin: 30px 0 25px;

    color: #9c98b3;

    font-size: 17px;

    line-height: 1.7;
}

.welcome-line {
    width: 70px;
    height: 4px;

    margin-bottom: 20px;

    border-radius: 3px;

    background: linear-gradient(90deg, #805cff, #ff6288);
}

.small-text {
    margin: 0;

    color: #69749a;

    font-family: Arial, sans-serif;

    font-size: 12px;
}

.decoration {
    position: absolute;

    border-radius: 50%;

    pointer-events: none;

    z-index: 0;
}

.decoration-one {
    width: 300px;
    height: 300px;

    right: -90px;
    top: 12%;

    background: rgba(128, 92, 255, 0.22);
}

.decoration-two {
    width: 220px;
    height: 220px;

    right: 14%;
    bottom: -60px;

    background: rgba(255, 98, 136, 0.14);
}

.decoration-three {
    width: 120px;
    height: 120px;

    left: 8%;
    bottom: 10%;

    background: rgba(75, 63, 160, 0.28);
}

.circle {
    position: absolute;

    border: 1px solid rgba(128, 92, 255, 0.3);

    border-radius: 50%;

    pointer-events: none;

    z-index: 0;
}

.circle-one {
    width: 230px;
    height: 230px;

    right: 10%;
    top: -90px;
}

.circle-two {
    width: 90px;
    height: 90px;

    left: 40%;
    bottom: -45px;

    background: rgba(131, 92, 255, 0.08);
}

.grid-pattern {
    position: absolute;

    right: 0;
    bottom: 0;

    width: 320px;
    height: 320px;

    opacity: 0.12;

    background-image:
        linear-gradient(
            rgba(128, 92, 255, 0.6) 1px,
            transparent 1px
        ),
        linear-gradient(
            90deg,
            rgba(128, 92, 255, 0.6) 1px,
            transparent 1px
        );

    background-size: 25px 25px;

    mask-image: linear-gradient(
        to top left,
        black,
        transparent
    );
}

@media (max-width: 900px) {
    .welcome-section {
        min-height: 430px;

        padding: 50px;
    }

    h1 {
        font-size: clamp(42px, 10vw, 68px);
    }

    .decoration-one {
        right: -110px;
    }

    .circle-one {
        right: -70px;
    }
}

@media (max-width: 500px) {
    .welcome-section {
        min-height: 390px;

        padding: 35px 25px;
    }

    .logo {
        width: 46px;
        height: 46px;

        margin-bottom: 35px;

        border-radius: 12px;
        font-size: 20px;
    }

    h1 {
        font-size: 42px;

        letter-spacing: -1px;
    }

    .welcome-text {
        margin-top: 22px;

        font-size: 14px;
    }

    .decoration-one {
        right: -140px;
    }

    .circle-one {
        right: -100px;
    }
}
</style>
