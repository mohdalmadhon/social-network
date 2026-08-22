<script setup>
const followRequests = [
  { id: 1, name: 'Sana Iqbal', color: '#5d43b4' },
  { id: 2, name: 'Theo Almeida', color: '#a23857' },
]

const onlinePeople = [
  { id: 1, name: 'Kiko Tanaka', color: '#19585b' },
  { id: 2, name: 'Mara Voss', color: '#685022' },
  { id: 3, name: 'Leo Marchetti', color: '#1d5472' },
  { id: 4, name: 'Aya Benali', color: '#42368b', status: '32m ago' },
]
</script>

<template>
  <aside class="feed-sidebar" aria-label="Feed shortcuts">
    <section class="feed-panel orbit-surface">
      <header class="feed-panel__header">
        <h2>Follow requests</h2>
        <span class="feed-panel__count">2</span>
      </header>

      <div v-for="request in followRequests" :key="request.id" class="request-row">
        <span class="mini-avatar" :style="{ background: request.color }" aria-hidden="true">
          {{ request.name.charAt(0) }}
        </span>
        <span class="request-row__name">
          <strong>{{ request.name }}</strong>
          <small>wants to follow you</small>
        </span>
        <button class="request-button request-button--accept" type="button" :aria-label="`Accept ${request.name}`">✓</button>
        <button class="request-button" type="button" :aria-label="`Decline ${request.name}`">×</button>
      </div>

      <p class="feed-panel__note">Your profile is <strong>private</strong> — requests need your approval.</p>
    </section>

    <section class="feed-panel orbit-surface">
      <header class="feed-panel__header">
        <h2>In orbit now</h2>
        <span class="online-dot" aria-label="Online"></span>
      </header>

      <div v-for="person in onlinePeople" :key="person.id" class="person-row">
        <span class="mini-avatar mini-avatar--online" :style="{ background: person.color }" aria-hidden="true">
          {{ person.name.charAt(0) }}
        </span>
        <strong>{{ person.name }}</strong>
        <small v-if="person.status">{{ person.status }}</small>
      </div>

      <div class="group-row">
        <span class="mini-avatar mini-avatar--group" aria-hidden="true">≋</span>
        <strong>Weekend Hikers</strong>
        <span class="feed-panel__message-count">3</span>
      </div>
    </section>
  </aside>
</template>

<style scoped>
.feed-sidebar {
  display: none;
}

@media (min-width: 90rem) {
  .feed-sidebar {
    display: grid;
    align-content: start;
    gap: var(--space-4);
    width: 22rem;
  }

  .feed-panel {
    padding: var(--space-5);
  }

  .feed-panel__header,
  .request-row,
  .person-row,
  .group-row {
    display: flex;
    align-items: center;
  }

  .feed-panel__header {
    justify-content: space-between;
    margin-bottom: var(--space-3);
  }

  .feed-panel h2 {
    margin: 0;
    color: var(--color-text);
    font-size: 1rem;
    font-weight: 500;
  }

  .feed-panel__count,
  .feed-panel__message-count {
    display: grid;
    min-width: 1.5rem;
    height: 1.5rem;
    place-items: center;
    border-radius: 999px;
    font-size: 0.75rem;
    font-weight: 700;
  }

  .feed-panel__count {
    background: rgb(255 184 77 / 20%);
    color: var(--color-amber);
  }

  .request-row,
  .person-row,
  .group-row {
    min-height: 3.25rem;
    gap: var(--space-3);
  }

  .mini-avatar {
    position: relative;
    display: grid;
    width: 2.5rem;
    height: 2.5rem;
    flex: 0 0 auto;
    place-items: center;
    border-radius: 50%;
    color: var(--color-text);
    font-size: 0.8125rem;
  }

  .mini-avatar--online::after {
    position: absolute;
    right: -0.05rem;
    bottom: 0.1rem;
    width: 0.55rem;
    height: 0.55rem;
    border: 2px solid var(--color-surface);
    border-radius: 50%;
    background: var(--color-mint);
    content: '';
  }

  .request-row__name,
  .person-row strong,
  .group-row strong {
    min-width: 0;
    flex: 1;
  }

  .request-row__name {
    display: grid;
  }

  .request-row strong,
  .person-row strong,
  .group-row strong {
    color: var(--color-text-soft);
    font-size: 0.875rem;
    font-weight: 500;
  }

  .request-row small,
  .person-row small {
    color: var(--color-text-faint);
    font-size: 0.75rem;
  }

  .request-button {
    display: grid;
    width: var(--touch-target);
    height: var(--touch-target);
    flex: 0 0 auto;
    place-items: center;
    border: 1px solid var(--color-border);
    border-radius: var(--radius-small);
    background: var(--color-input);
    color: var(--color-text-muted);
    cursor: pointer;
  }

  .request-button--accept {
    border-color: var(--color-mint);
    background: rgb(62 230 176 / 8%);
    color: var(--color-mint);
  }

  .feed-panel__note {
    margin: var(--space-3) 0 0;
    color: var(--color-text-faint);
    font-size: 0.75rem;
    line-height: 1.45;
  }

  .feed-panel__note strong {
    color: var(--color-coral);
    font-weight: 500;
  }

  .online-dot {
    width: 0.55rem;
    height: 0.55rem;
    border-radius: 50%;
    background: var(--color-mint);
  }

  .person-row small {
    margin-left: auto;
  }

  .group-row {
    margin-top: var(--space-3);
    padding-top: var(--space-3);
    border-top: 1px solid var(--color-border);
  }

  .mini-avatar--group,
  .feed-panel__message-count {
    background: var(--color-mint);
    color: #071713;
  }
}
</style>
