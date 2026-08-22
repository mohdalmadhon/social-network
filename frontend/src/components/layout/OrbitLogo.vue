<template>
  <span class="orbit-logo" aria-hidden="true">
    <span class="orbit-logo__top"></span>
    <span class="orbit-logo__bottom"></span>
  </span>
</template>

<style scoped>
.orbit-logo {
  position: relative;
  width: 1.75rem;
  height: 2.5rem;
  flex: 0 0 1.75rem;
}

.orbit-logo__top,
.orbit-logo__bottom {
  position: absolute;
  left: 50%;
  border-radius: 50%;
  transform: translateX(-50%);
}

.orbit-logo__top {
  top: 0;
  width: 1.25rem;
  height: 1.25rem;
  background: var(--color-violet);
}

.orbit-logo__bottom {
  bottom: 0;
  width: 1.75rem;
  height: 1.75rem;
  background: var(--color-coral);
}

.orbit-logo__bottom::after {
  position: absolute;
  inset: 0.5rem;
  border-radius: 50%;
  background: var(--color-text);
  content: "";
}
</style>
