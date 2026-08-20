<script setup>
</script>

<template>
  <main class="landing">
    <div class="orb orb-1"></div>
    <div class="orb orb-2"></div>
    <div class="orb orb-3"></div>
    <div class="orb orb-4"></div>
    <div class="orb orb-5"></div>

    <header class="header">
      <div class="brand">
        
        <span class="brand-name">orbit</span>
      </div>
    </header>

    <section class="hero">
      <h1>
        Find your<br />
        people. Keep<br />
        your <span>orbit.</span>
      </h1>

      <p class="description">
        Post to everyone, your followers, or a hand-picked few.<br />
        Privacy is a choice you see — not a setting you forget.
      </p>

      <div class="features">
        <div class="feature">
          <div class="feature-icon people">
            <span></span>
            <span></span>
          </div>
          <span>Followers you actually control</span>
        </div>

        <div class="feature">
          <div class="feature-icon groups">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <span>Groups with events and RSVP</span>
        </div>

        <div class="feature">
          <div class="feature-icon chat">↗</div>
          <span>Real-time chat over websockets</span>
        </div>
      </div>
    </section>

    <div class="accent-dot"></div>
  </main>
</template>

<style scoped>
* {
  box-sizing: border-box;
}

.landing {
  position: relative;
  min-height: 100vh;
  width: 35%;
  overflow: hidden;
  background:
    radial-gradient(circle at 20% 35%, rgba(75, 63, 160, 0.12), transparent 35%),
    #100d2b;
  color: #fff;
  font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.orb {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  z-index: 0;
}

.orb-1 {
  width: 360px;
  height: 360px;
  left: -80px;
  top: 55px;
  background: rgba(71, 61, 143, 0.28);
}

.orb-2 {
  width: 420px;
  height: 420px;
  left: -110px;
  top: 300px;
  background: rgba(48, 45, 119, 0.28);
}

.orb-3 {
  width: 270px;
  height: 270px;
  left: -42px;
  bottom: 2px;
  background: rgba(133, 72, 106, 0.32);
}

.orb-4 {
  width: 220px;
  height: 220px;
  left: -120px;
  top: 260px;
  background: rgba(116, 65, 112, 0.26);
}

.orb-5 {
  width: 500px;
  height: 500px;
  right: -370px;
  bottom: -250px;
  background: rgba(83, 52, 157, 0.08);
}

.header {
  position: relative;
  z-index: 2;
  padding: 12px 0 0 47px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-name {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.logo-placeholder {
  position: relative;
  width: 30px;
  height: 54px;
}

.logo-placeholder span {
  position: absolute;
  left: 3px;
  width: 25px;
  height: 25px;
  border-radius: 50%;
}

.logo-placeholder span:first-child {
  top: 0;
  background: linear-gradient(135deg, #8059ff, #8b54ee);
}

.logo-placeholder span:last-child {
  bottom: 0;
  background: linear-gradient(135deg, #ff5f99, #aa58e9);
}

.logo-placeholder span:last-child::after {
  content: "";
  position: absolute;
  width: 10px;
  height: 10px;
  left: 7.5px;
  top: 7.5px;
  border-radius: 50%;
  background: #f4eefb;
}

.hero {
  position: relative;
  z-index: 2;
  line-height: 1.5;
  margin-left: 46px;
  margin-top: 108px;
}

.hero h1 {
  margin: 0;
  font-size: clamp(42px, 10vw, 52px);
  line-height: 1.5;
  letter-spacing: -1.7px;
  font-weight: 750;
}

.hero h1 span {
  color: #7b55ff;
}

.description {
  margin: 5px 0 0 5px;
  margin-top: 5px;
  color: #9c98b3;
  font-size: 18px;
  line-height: 20px;
  font-weight: 400;
}

.features {
  display: flex;
  flex-direction: column;
  gap: 13px;
  margin-top: 21px;
}

.feature {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #e9e6f1;
  font-size: 11px;
  line-height: 1;
}

.feature-icon {
  position: relative;
  width: 17px;
  height: 17px;
  flex: 0 0 17px;
}

.feature-icon.people {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 2px;
}

.feature-icon.people::before,
.feature-icon.people::after {
  content: "";
  position: absolute;
  top: 0;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #8055ff;
}

.feature-icon.people::before {
  left: 1px;
}

.feature-icon.people::after {
  right: 1px;
}

.feature-icon.people span {
  width: 8px;
  height: 7px;
  border-radius: 7px 7px 2px 2px;
  background: #8055ff;
}

.feature-icon.groups {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 2px;
}

.feature-icon.groups span {
  display: block;
  width: 17px;
  height: 4px;
  background: #ff5c87;
  clip-path: polygon(0 20%, 50% 0, 100% 20%, 50% 45%);
}

.feature-icon.groups span:nth-child(2) {
  transform: translateY(-1px);
}

.feature-icon.groups span:nth-child(3) {
  transform: translateY(-2px);
}

.feature-icon.chat {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #36e0ad;
  font-size: 23px;
  line-height: 1;
  font-weight: 900;
  transform: rotate(-8deg);
}

.accent-dot {
  position: absolute;
  z-index: 2;
  width: 17px;
  height: 17px;
  right: 140px;
  bottom: 215px;
  border-radius: 50%;
  background: #ff6385;
}

@media (max-width: 500px) {
  .header {
    padding-left: 47px;
  }

  .hero {
    margin-left: 46px;
    margin-top: 108px;
  }

  .hero h1 {
    font-size: 42px;
  }

  .description {
    max-width: 330px;
  }
}

@media (max-width: 380px) {
  .hero {
    margin-left: 32px;
  }

  .header {
    padding-left: 32px;
  }

  .hero h1 {
    font-size: 38px;
  }

  .description {
    font-size: 11px;
  }
}
</style>