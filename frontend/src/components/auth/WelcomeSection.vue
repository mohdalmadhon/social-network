<template>
  <section class="welcome-section">
    <div class="welcome-content">
      <div class="logo" aria-label="Orbit">◎</div>
      <p class="eyebrow">WELCOME BACK</p>
      <h1>Your people.<br />Your <span>place.</span></h1>
      <p class="welcome-text">Connect with people, share your thoughts, and keep up with the things that matter to you.</p>
      <div class="welcome-line"></div>
      <p class="small-text">A simple place to stay connected.</p>
    </div>
    <div class="decoration decoration-one" aria-hidden="true"></div>
    <div class="decoration decoration-two" aria-hidden="true"></div>
    <div class="decoration decoration-three" aria-hidden="true"></div>
    <div class="circle circle-one" aria-hidden="true"></div>
    <div class="circle circle-two" aria-hidden="true"></div>
    <div class="grid-pattern" aria-hidden="true"></div>
  </section>
</template>

<style scoped>
.welcome-section {
  position: relative;
  display: flex;
  min-height: 26rem;
  align-items: center;
  overflow: hidden;
  padding: clamp(2rem, 8vw, 5rem) clamp(1.5rem, 8vw, 5rem);
  background: #0b1220;
  isolation: isolate;
}
.welcome-section::after {
  position: absolute;
  inset: 0;
  z-index: -1;
  background: linear-gradient(135deg, transparent 0 55%, rgb(72 217 193 / 6%) 55% 55.2%, transparent 55.2%);
  content: '';
}
.welcome-content { position: relative; z-index: 2; max-width: 38rem; }
.logo { display: grid; width: 3.25rem; height: 3.25rem; margin-bottom: clamp(2rem, 8vw, 4.5rem); place-items: center; border: 1px solid rgb(72 217 193 / 70%); border-radius: var(--radius-small); background: var(--color-mint); box-shadow: 0 .75rem 2rem rgb(0 0 0 / 28%); color: #071713; font-size: 1.5rem; font-weight: 800; }
.eyebrow { margin: 0 0 .75rem; color: var(--color-coral); font-family: var(--font-meta); font-size: .75rem; letter-spacing: .2em; }
h1 { margin: 0; color: var(--color-text); font-family: var(--font-display); font-size: clamp(2.75rem, 10vw, 5.5rem); font-weight: 700; letter-spacing: -.05em; line-height: .98; }
h1 span { color: var(--color-mint); }
.welcome-text { max-width: 38ch; margin: 1.5rem 0 1.25rem; color: var(--color-text-soft); font-size: 1rem; line-height: 1.6; }
.welcome-line { width: 4rem; height: .2rem; margin-bottom: 1rem; background: var(--color-coral); }
.small-text { margin: 0; color: var(--color-text-muted); font-family: var(--font-meta); font-size: .75rem; }
.decoration, .circle { position: absolute; pointer-events: none; }
.decoration { border-radius: 50%; }
.decoration-one { top: 14%; right: -8rem; width: 18rem; aspect-ratio: 1; background: rgb(158 140 255 / 15%); }
.decoration-two { right: 12%; bottom: -8rem; width: 16rem; aspect-ratio: 1; background: rgb(255 111 145 / 12%); }
.decoration-three { bottom: 8%; left: 8%; width: 7rem; aspect-ratio: 1; background: rgb(72 217 193 / 10%); }
.circle { border: 1px solid rgb(105 183 232 / 35%); border-radius: 50%; }
.circle-one { top: -7rem; right: 10%; width: 15rem; aspect-ratio: 1; }
.circle-two { bottom: -3rem; left: 40%; width: 6rem; aspect-ratio: 1; }
.grid-pattern { position: absolute; right: 0; bottom: 0; width: 20rem; height: 20rem; opacity: .35; background-image: linear-gradient(rgb(105 183 232 / 12%) 1px, transparent 1px), linear-gradient(90deg, rgb(105 183 232 / 12%) 1px, transparent 1px); background-size: 2rem 2rem; mask-image: linear-gradient(to top left, black, transparent); }
@media (min-width: 56.25rem) { .welcome-section { min-height: 100vh; } }
</style>
