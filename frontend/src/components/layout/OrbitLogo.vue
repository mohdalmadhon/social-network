<template>
  <svg class="orbit-logo" viewBox="0 0 42 42" fill="none" aria-hidden="true" focusable="false">
    <defs>
      <linearGradient id="orbit-logo-gradient" x1="4" y1="4" x2="38" y2="38" gradientUnits="userSpaceOnUse">
        <stop stop-color="var(--color-violet)" />
        <stop offset="1" stop-color="var(--color-coral)" />
      </linearGradient>
    </defs>
    <ellipse class="orbit-logo__path" cx="21" cy="21" rx="17" ry="8" transform="rotate(-28 21 21)" stroke="url(#orbit-logo-gradient)" stroke-width="3" />
    <circle class="orbit-logo__ring" cx="21" cy="21" r="10.5" stroke="url(#orbit-logo-gradient)" stroke-width="5" />
    <circle class="orbit-logo__core" cx="21" cy="21" r="3.5" fill="var(--color-text)" />
  </svg>
</template>

<style scoped>
.orbit-logo {
  width: 1.75rem;
  height: 1.75rem;
  flex: 0 0 1.75rem;
  overflow: visible;
}
</style>
