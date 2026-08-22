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
        <span class="logo-placeholder" aria-hidden="true">
          <span></span>
          <span></span>
        </span>
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
  width: 100%;
  min-height: 19rem;
  overflow: hidden;
  background:
    radial-gradient(circle at 20% 35%, rgba(75, 63, 160, 0.12), transparent 35%),
    #100d2b;
  color: #fff;
  font-family: var(--font-body);
}

.orb {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
  z-index: 0;
}

.orb-1 {
  width: 16rem;
  height: 16rem;
  left: -5rem;
  top: 2.5rem;
  background: rgba(71, 61, 143, 0.28);
}

.orb-2 {
  width: 20rem;
  height: 20rem;
  left: 35%;
  top: 5rem;
  background: rgba(48, 45, 119, 0.28);
}

.orb-3 {
  width: 12rem;
  height: 12rem;
  right: -4rem;
  bottom: -5rem;
  left: auto;
  background: rgba(133, 72, 106, 0.32);
}

.orb-4 {
  width: 10rem;
  height: 10rem;
  left: -5rem;
  top: 11rem;
  background: rgba(116, 65, 112, 0.26);
}

.orb-5 {
  display: none;
  background: rgba(83, 52, 157, 0.08);
}

.header {
  position: relative;
  z-index: 2;
  padding: var(--space-4) var(--space-5) 0;
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-name {
  font-family: var(--font-display);
  font-size: 1.5rem;
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
  max-width: 34rem;
  margin: var(--space-6) var(--space-5) 0;
}

.hero h1 {
  margin: 0;
  font-family: var(--font-display);
  font-size: clamp(2.25rem, 10vw, 3.25rem);
  line-height: 1.12;
  letter-spacing: -0.04em;
  font-weight: 700;
}

.hero h1 span {
  color: #7b55ff;
}

.description {
  max-width: 34rem;
  margin: var(--space-3) 0 0;
  color: #9c98b3;
  font-size: 1rem;
  line-height: 1.45;
  font-weight: 400;
}

.features {
  display: none;
  flex-direction: column;
  gap: 13px;
  margin-top: 21px;
}

.feature {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #e9e6f1;
  font-size: 0.875rem;
  line-height: 1.3;
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
  width: 1rem;
  height: 1rem;
  right: 12%;
  bottom: 12%;
  border-radius: 50%;
  background: #ff6385;
}

@media (min-width: 48rem) {
  .landing {
    min-height: 25rem;
  }

  .header {
    padding: var(--space-5) var(--space-7) 0;
  }

  .hero {
    margin: var(--space-7) var(--space-7) 0;
  }

  .features {
    display: flex;
  }
}

@media (min-width: 64rem) {
  .landing {
    min-height: 100vh;
  }

  .orb-1 {
    width: 22.5rem;
    height: 22.5rem;
    left: -5rem;
    top: 3.5rem;
  }

  .orb-2 {
    width: 26rem;
    height: 26rem;
    left: -7rem;
    top: 18.75rem;
  }

  .orb-3 {
    width: 17rem;
    height: 17rem;
    right: auto;
    bottom: 0;
    left: -2.5rem;
  }

  .orb-4 {
    width: 14rem;
    height: 14rem;
    left: -7rem;
    top: 16rem;
  }

  .orb-5 {
    display: block;
    width: 31rem;
    height: 31rem;
    right: -23rem;
    bottom: -16rem;
  }

  .header {
    padding: var(--space-5) clamp(2rem, 4vw, 4rem) 0;
  }

  .hero {
    margin: clamp(7rem, 14vh, 12rem) clamp(2rem, 4vw, 4rem) 0;
  }

  .hero h1 {
    font-size: clamp(3rem, 3.2vw, 4.25rem);
  }
}
</style>
