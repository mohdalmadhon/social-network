<script setup>
import OrbitLogo from '@/components/layout/OrbitLogo.vue'
import IconGlyph from '@/components/layout/IconGlyph.vue'

const features = [
  { icon: 'profile', label: 'Followers you actually control', tone: 'violet' },
  { icon: 'groups', label: 'Groups with events and RSVP', tone: 'coral' },
  { icon: 'chat', label: 'Real-time chat over websockets', tone: 'mint' },
]
</script>

<template>
  <section class="welcome-section" aria-label="About Orbit">
    <div class="welcome-inner">
      <a class="welcome-brand" href="/login" aria-label="Orbit home">
        <OrbitLogo />
        <span>orbit</span>
      </a>

      <div class="welcome-content">
        <h1>Find your<br />people. Keep<br /><span>your orbit.</span></h1>
        <p class="welcome-text">
          Post to everyone, your followers, or a hand-picked few. Privacy is a choice you see — not a setting you forget.
        </p>

        <ul class="feature-list" aria-label="Orbit features">
          <li v-for="feature in features" :key="feature.label">
            <span class="feature-icon" :class="`feature-icon--${feature.tone}`"><IconGlyph :name="feature.icon" :size="17" /></span>
            <span>{{ feature.label }}</span>
          </li>
        </ul>
      </div>
    </div>

    <span class="welcome-orbit welcome-orbit--primary" aria-hidden="true"></span>
    <span class="welcome-orbit welcome-orbit--secondary" aria-hidden="true"></span>
    <span class="welcome-orbit welcome-orbit--tertiary" aria-hidden="true"></span>
    <span class="welcome-grid" aria-hidden="true"></span>
  </section>
</template>

<style scoped>
.welcome-section {
  position: relative;
  min-height: 18rem;
  overflow: hidden;
  background: var(--color-background);
  isolation: isolate;
}

.welcome-inner {
  position: relative;
  z-index: 1;
  display: grid;
  align-content: start;
  gap: clamp(2.5rem, 9vh, 6rem);
  height: 100%;
  padding: clamp(1.25rem, 4vw, 3rem) clamp(1.5rem, 5vw, 4rem);
}

.welcome-brand {
  display: inline-flex;
  width: fit-content;
  align-items: center;
  gap: var(--space-2);
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 700;
  text-decoration: none;
}

.welcome-brand :deep(.orbit-logo) {
  transform: scale(.68);
  transform-origin: left center;
  margin-right: -.55rem;
}

.welcome-content {
  max-width: 31rem;
}

h1 {
  margin: 0;
  color: var(--color-text);
  font-family: var(--font-display);
  font-size: clamp(2.75rem, 6vw, 4.25rem);
  font-weight: 700;
  letter-spacing: -.055em;
  line-height: .98;
}

h1 span {
  color: var(--color-violet);
}

.welcome-text {
  max-width: 34rem;
  margin: 1.5rem 0 1.25rem;
  color: var(--color-text-muted);
  font-size: 1rem;
  line-height: 1.45;
}

.feature-list {
  display: grid;
  gap: .65rem;
  margin: 0;
  padding: 0;
  list-style: none;
  color: var(--color-text-soft);
  font-size: .875rem;
}

.feature-list li {
  display: flex;
  align-items: center;
  gap: .7rem;
}

.feature-icon {
  display: inline-grid;
  width: 1.25rem;
  height: 1.25rem;
  flex: 0 0 1.25rem;
  place-items: center;
  font-size: 1.1rem;
  line-height: 1;
}

.feature-icon--violet { color: var(--color-violet); }
.feature-icon--coral { color: var(--color-coral); }
.feature-icon--mint { color: var(--color-mint); }

.welcome-orbit,
.welcome-grid {
  position: absolute;
  pointer-events: none;
}

.welcome-orbit {
  border-radius: 50%;
}

.welcome-orbit--primary {
  top: 12%;
  left: -9rem;
  width: 22rem;
  aspect-ratio: 1;
  background: rgb(124 92 255 / 15%);
}

.welcome-orbit--secondary {
  bottom: -9rem;
  left: -2rem;
  width: 16rem;
  aspect-ratio: 1;
  background: rgb(255 107 138 / 14%);
}

.welcome-orbit--tertiary {
  right: 11%;
  bottom: -12rem;
  width: 25rem;
  aspect-ratio: 1;
  border: 1px solid rgb(76 195 255 / 18%);
}

.welcome-grid {
  right: 0;
  bottom: 0;
  width: 19rem;
  height: 19rem;
  opacity: .28;
  background-image: linear-gradient(rgb(76 195 255 / 14%) 1px, transparent 1px), linear-gradient(90deg, rgb(76 195 255 / 14%) 1px, transparent 1px);
  background-size: 2rem 2rem;
  mask-image: linear-gradient(to top left, black, transparent 72%);
}

@media (min-width: 56.25rem) {
  .welcome-section {
    min-height: 100vh;
  }

  .welcome-inner {
    gap: clamp(4rem, 12vh, 8rem);
  }
}
</style>
