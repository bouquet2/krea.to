// krea.to — theme, light/dark, accessibility, search

const SCHEMES = ['pinkie', 'tokyonight', 'mocha', 'gruvbox', 'latte'];
const LIGHT_THEME = 'latte';
const STORAGE_THEME = 'colorScheme';
const STORAGE_LAST_DARK = 'lastDarkTheme';
const STORAGE_A11Y = 'accessibility';

function loadThemeCSS(theme) {
    if (!SCHEMES.includes(theme)) theme = 'pinkie';
    const id = `theme-css-${theme}`;
    if (!document.getElementById(id)) {
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = `/css/themes/${theme}.css`;
        link.id = id;
        document.head.appendChild(link);
    }
}

function currentTheme() {
    return document.body.getAttribute('data-theme') || 'pinkie';
}

function isLight(theme) {
    return theme === LIGHT_THEME;
}

function applyTheme(theme) {
    if (!SCHEMES.includes(theme)) theme = 'pinkie';
    document.body.setAttribute('data-theme', theme);
    document.body.classList.toggle('is-light', isLight(theme));
    loadThemeCSS(theme);
    localStorage.setItem(STORAGE_THEME, theme);
    const select = document.getElementById('theme-select');
    if (select) select.value = theme;
}

(function applyImmediately() {
    const saved = localStorage.getItem(STORAGE_THEME);
    const serverTheme = document.body.getAttribute('data-theme') || 'pinkie';
    applyTheme(saved && SCHEMES.includes(saved) ? saved : serverTheme);
})();

document.addEventListener('DOMContentLoaded', function () {
    const body = document.body;

    const themeSelect = document.getElementById('theme-select');
    if (themeSelect) {
        themeSelect.value = currentTheme();
        themeSelect.addEventListener('change', function () {
            const theme = this.value;
            if (!isLight(theme)) localStorage.setItem(STORAGE_LAST_DARK, theme);
            applyTheme(theme);
        });
    }

    const lightdarkButton = document.getElementById('lightdark-button');
    if (lightdarkButton) {
        lightdarkButton.addEventListener('click', function () {
            const cur = currentTheme();
            if (isLight(cur)) {
                const lastDark = localStorage.getItem(STORAGE_LAST_DARK) || 'pinkie';
                applyTheme(lastDark);
            } else {
                localStorage.setItem(STORAGE_LAST_DARK, cur);
                applyTheme(LIGHT_THEME);
            }
        });
    }

    const a11yButton = document.getElementById('accessibility-button');
    const savedA11y = localStorage.getItem(STORAGE_A11Y);
    if (savedA11y === 'enabled') {
        body.classList.add('accessibility-mode');
        if (a11yButton) a11yButton.setAttribute('aria-pressed', 'true');
    }
    if (a11yButton) {
        a11yButton.addEventListener('click', function () {
            const enabled = body.classList.toggle('accessibility-mode');
            localStorage.setItem(STORAGE_A11Y, enabled ? 'enabled' : 'disabled');
            a11yButton.setAttribute('aria-pressed', enabled ? 'true' : 'false');
            body.classList.remove('a11y-swap');
            void body.offsetWidth; // restart animation
            body.classList.add('a11y-swap');
        });
        body.addEventListener('animationend', function () {
            body.classList.remove('a11y-swap');
        });
    }

    const searchInput = document.getElementById('search-input');
    const blogList = document.getElementById('blog-list');
    if (searchInput && blogList) {
        const items = Array.from(blogList.querySelectorAll('.post'));
        const noResults = document.getElementById('no-results');
        searchInput.addEventListener('input', function () {
            const q = this.value.trim().toLowerCase();
            let visible = 0;
            items.forEach(item => {
                const title = (item.getAttribute('data-title') || item.textContent || '').toLowerCase();
                const match = !q || title.includes(q);
                item.style.display = match ? '' : 'none';
                if (match) visible++;
            });
            if (noResults) noResults.style.display = visible === 0 && q ? 'block' : 'none';
        });
    }
});
