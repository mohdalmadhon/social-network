<template>
    <header class="top-navigation">
        <div class="nav-left">
            <a href="/" class="logo">N</a>

            <div class="search">
                <span>⌕</span>
                <input type="text" placeholder="Search" />
            </div>
        </div>

        <nav class="nav-links">
            <a href="/">Home</a>
            <a href="/friends">Friends</a>
            <a href="/groups">Groups</a>
        </nav>

        <div class="nav-right">
            <button class="nav-button">+</button>
            <button class="nav-button">♧</button>
            <button class="nav-button">☼</button>

            <a href="/profile" class="nav-avatar">
                A
            </a>
        </div>
    </header>
</template>

<style scoped>
.top-navigation {
    position: fixed;
    z-index: 100;
    top: 0;
    left: 0;
    width: 100%;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 clamp(12px, 3vw, 28px);
    background: var(--bg-color);
    border-bottom: 2px solid var(--main-color);
    box-sizing: border-box;
    gap: 12px;
}

.nav-left,
.nav-right,
.nav-links {
    display: flex;
    align-items: center;
    min-width: 0;
}

.nav-left {
    gap: clamp(10px, 2vw, 18px);
    flex: 1;
    min-width: 0;
}

.logo {
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--input-focus);
    border: 2px solid var(--main-color);
    border-radius: 6px;
    box-shadow: 4px 4px var(--main-color);
    color: white;
    font-family: "JetBrains Mono", monospace;
    font-size: 19px;
    font-weight: 600;
}

.search {
    width: clamp(120px, 20vw, 220px);
    height: 38px;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 12px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--page-background);
    box-sizing: border-box;
    min-width: 0;
}

.search span {
    flex-shrink: 0;
    font-size: 20px;
}

.search input {
    width: 100%;
    min-width: 0;
    border: 0;
    outline: 0;
    background: transparent;
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
}

.nav-links {
    gap: clamp(14px, 2.5vw, 35px);
}

.nav-links a {
    color: var(--font-color);
    font-family: "JetBrains Mono", monospace;
    font-size: 11px;
    font-weight: 600;
    white-space: nowrap;
}

.nav-links a:hover {
    color: var(--input-focus);
}

.nav-right {
    flex-shrink: 0;
    gap: clamp(6px, 1.5vw, 10px);
}

.nav-button {
    flex-shrink: 0;
    width: 35px;
    height: 35px;
    border: 2px solid var(--main-color);
    border-radius: 5px;
    background: var(--page-background);
    color: var(--main-color);
    font-size: 16px;
    box-shadow: 3px 3px var(--main-color);
}

.nav-button:active {
    transform: translate(3px, 3px);
    box-shadow: none;
}

.nav-avatar {
    flex-shrink: 0;
    width: 38px;
    height: 38px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-left: 5px;
    border: 2px solid var(--main-color);
    border-radius: 50%;
    background: var(--input-focus);
    color: white;
    font-family: "JetBrains Mono", monospace;
    font-weight: 600;
}

@media (max-width: 900px) {
    .nav-links {
        display: none;
    }
}

@media (max-width: 640px) {
    .search {
        width: 100%;
    }
}

@media (max-width: 480px) {
    .search span {
        display: none;
    }

    .nav-button:nth-of-type(2) {
        display: none;
    }
}

</style>