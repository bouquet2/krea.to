// Function to dynamically load theme CSS
const loadThemeCSS = (theme) => {
    if (!theme) return;
    
    const themeId = `theme-css-${theme}`;
    if (!document.getElementById(themeId)) {
        const link = document.createElement('link');
        link.id = themeId;
        link.rel = 'stylesheet';
        link.href = `/css/themes/${theme}.css`;
        document.head.appendChild(link);
    }
};

// Apply theme immediately to prevent flash
(function() {
    const body = document.body;
    const schemes = ['mocha', 'frappe', 'latte', 'macchiato', 'gruvbox', 'nord', 'tokyonight', 'monokai', 'onedark', 'solarized', 'kanagawa', 'pinkie'];
    
    // Store the original server-set theme
    const originalServerTheme = body.getAttribute('data-theme') || 'pinkie';
    
    // Check for saved theme preference first
    let savedScheme = localStorage.getItem('colorScheme');
    if (!savedScheme) {
        savedScheme = originalServerTheme;
    }
    if (!savedScheme || !schemes.includes(savedScheme)) {
        savedScheme = 'pinkie';
    }
    
    // Apply the theme immediately
    if (schemes.includes(savedScheme)) {
        body.setAttribute('data-theme', savedScheme);
        loadThemeCSS(savedScheme);
    }
    
    // Apply light-theme class if using latte
    if (savedScheme === 'latte') {
        body.classList.add('light-theme');
    }
})();

document.addEventListener('DOMContentLoaded', function() {

    // Accessibility mode functionality
    const accessibilityButton = document.getElementById('accessibility-button');
    const backgroundButton = document.getElementById('background-button');
    const transparencyButton = document.getElementById('transparency-button');
    const luckyButton = document.getElementById('lucky-button');
    const themeSelect = document.getElementById('theme-select');
    const fontSelect = document.getElementById('font-select');
    const fontSizeRange = document.getElementById('font-size-range');
    const fontSizeValue = document.getElementById('font-size-value');
    const body = document.body;
    const terminal = document.querySelector('.terminal');
    
    // Color schemes available
    const schemes = ['mocha', 'frappe', 'latte', 'macchiato', 'gruvbox', 'nord', 'tokyonight', 'monokai', 'onedark', 'solarized', 'kanagawa', 'pinkie'];

    const updateAccessibilityButtonState = () => {
        if (!accessibilityButton) return;
        const isEnabled = body.classList.contains('accessibility-mode');
        const label = isEnabled ? 'Disable accessibility mode' : 'Enable accessibility mode';
        accessibilityButton.setAttribute('aria-pressed', isEnabled ? 'true' : 'false');
        accessibilityButton.setAttribute('aria-label', label);
        accessibilityButton.setAttribute('title', label);
    };
    
    // Check for saved accessibility preference
    const savedAccessibility = localStorage.getItem('accessibility');
    if (savedAccessibility === 'enabled') {
        body.classList.add('accessibility-mode');
        if (terminal) terminal.classList.add('accessibility-mode');
    }

    updateAccessibilityButtonState();
    
    // Check for saved background preference
    if (backgroundButton) {
        const savedBackground = localStorage.getItem('background');
        if (savedBackground === 'disabled') {
            body.classList.add('no-background');
            backgroundButton.textContent = 'Show Background';
        } else {
            backgroundButton.textContent = 'Hide Background';
        }
    }
    
    // Check for saved transparency preference
    if (transparencyButton) {
        const savedTransparency = localStorage.getItem('transparency');
        if (savedTransparency === 'disabled') {
            if (terminal) terminal.classList.add('no-transparency');
            transparencyButton.textContent = 'Enable Transparency';
        } else {
            transparencyButton.textContent = 'Disable Transparency';
        }
    }
    
    // Toggle accessibility mode on button click
    if (accessibilityButton) {
        accessibilityButton.addEventListener('click', function() {
            // Add a small delay to ensure smooth transition
            setTimeout(() => {
                body.classList.toggle('accessibility-mode');
                if (terminal) terminal.classList.toggle('accessibility-mode');

                if (body.classList.contains('accessibility-mode')) {
                    localStorage.setItem('accessibility', 'enabled');
                } else {
                    localStorage.setItem('accessibility', 'disabled');
                }

                updateAccessibilityButtonState();
            }, 50);
        });
    }
    
    // Toggle background on button click
    if (backgroundButton) {
        backgroundButton.addEventListener('click', function() {
        // Add a small delay for better visual feedback
        setTimeout(() => {
            body.classList.toggle('no-background');
            
            if (body.classList.contains('no-background')) {
                localStorage.setItem('background', 'disabled');
                backgroundButton.textContent = 'Show Background';
            } else {
                localStorage.setItem('background', 'enabled');
                backgroundButton.textContent = 'Hide Background';
            }
        }, 50);
        });
    }
    
    // Toggle transparency on button click
    if (transparencyButton) {
        transparencyButton.addEventListener('click', function() {
        if (terminal) {
            terminal.classList.toggle('no-transparency');
            
            if (terminal.classList.contains('no-transparency')) {
                localStorage.setItem('transparency', 'disabled');
                transparencyButton.textContent = 'Enable Transparency';
            } else {
                localStorage.setItem('transparency', 'enabled');
                transparencyButton.textContent = 'Disable Transparency';
            }
        }
        });
    }

    // Font size control
    const savedFontSize = localStorage.getItem('fontSize') || '1.1';
    if (fontSizeRange) {
        fontSizeRange.value = savedFontSize;
        if (fontSizeValue) fontSizeValue.textContent = `${savedFontSize}em`;
        document.documentElement.style.setProperty('--font-size', `${savedFontSize}em`);
        
        fontSizeRange.addEventListener('input', function() {
            const size = this.value;
            document.documentElement.style.setProperty('--font-size', `${size}em`);
            if (fontSizeValue) fontSizeValue.textContent = `${size}em`;
            localStorage.setItem('fontSize', size);
        });
    }

    // Terminal blur intensity control
    const blurIntensityRange = document.getElementById('blur-intensity-range');
    const blurIntensityValue = document.getElementById('blur-intensity-value');
    const savedBlurIntensity = localStorage.getItem('terminalBlurIntensity') || '20';
    
    if (blurIntensityRange) {
        blurIntensityRange.value = savedBlurIntensity;
        if (blurIntensityValue) blurIntensityValue.textContent = `${savedBlurIntensity}px`;
        document.documentElement.style.setProperty('--terminal-blur', `${savedBlurIntensity}px`);
        
        blurIntensityRange.addEventListener('input', function() {
            const intensity = this.value;
            document.documentElement.style.setProperty('--terminal-blur', `${intensity}px`);
            if (blurIntensityValue) blurIntensityValue.textContent = `${intensity}px`;
            localStorage.setItem('terminalBlurIntensity', intensity);
        });
    }

    // Font selection control
    const savedFont = localStorage.getItem('fontFamily') || 'scientifica';
    if (fontSelect) {
        fontSelect.value = savedFont;
        
        // Apply saved font
        if (savedFont === 'pokemon') {
            document.documentElement.style.setProperty('--font-family', "'Pokemon DP Pro', sans-serif");
        } else {
            document.documentElement.style.setProperty('--font-family', "'Scientifica', sans-serif");
        }
        
        fontSelect.addEventListener('change', function() {
            const selectedFont = this.value;
            
            if (selectedFont === 'pokemon') {
                document.documentElement.style.setProperty('--font-family', "'Pokemon DP Pro', sans-serif");
            } else {
                document.documentElement.style.setProperty('--font-family', "'Scientifica', sans-serif");
            }
            
            localStorage.setItem('fontFamily', selectedFont);
        });
    }

    // Theme toggle functionality
    const themeButton = document.getElementById('theme-button');

    // Get the original server-set theme (already stored in IIFE)
    const originalServerTheme = body.getAttribute('data-theme') || 'pinkie';

    const updateThemeButtonState = () => {
        if (!themeButton) return;
        const isLight = body.classList.contains('light-theme');
        const label = isLight ? 'Switch to dark mode' : 'Switch to light mode';
        themeButton.setAttribute('aria-pressed', isLight ? 'true' : 'false');
        themeButton.setAttribute('aria-label', label);
        themeButton.setAttribute('title', label);
    };

    // Check for saved theme preference first (light/dark)
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'light') {
        body.classList.add('light-theme');
    }
    
    // Helper function to apply a color scheme
    const applyScheme = (schemeName) => {
        if (schemes.includes(schemeName)) {
            body.setAttribute('data-theme', schemeName);
            localStorage.setItem('colorScheme', schemeName);
            loadThemeCSS(schemeName);
        }
    };
    
    // Get the current scheme (already applied by inline script)
    const savedScheme = localStorage.getItem('colorScheme') || body.getAttribute('data-theme') || 'pinkie';
    if (themeSelect) {
        themeSelect.value = savedScheme;
    }

    updateThemeButtonState();
    
    // Theme select dropdown handler
    if (themeSelect) {
        themeSelect.addEventListener('change', function() {
            const selectedScheme = this.value;
            applyScheme(selectedScheme);
            
            // Update light-theme class based on scheme
            if (selectedScheme === 'latte') {
                body.classList.add('light-theme');
                localStorage.setItem('theme', 'light');
            } else {
                body.classList.remove('light-theme');
                localStorage.setItem('theme', 'dark');
            }
            
            updateThemeButtonState();
        });
    }
    
    // Toggle theme on button click
    if (themeButton) {
        themeButton.addEventListener('click', function() {
            if (body.classList.contains('light-theme')) {
                // Switch to dark mode (use original server theme
                body.classList.remove('light-theme');
                applyScheme(originalServerTheme);
                localStorage.setItem('theme', 'dark');
                if (themeSelect) themeSelect.value = originalServerTheme;
            } else {
                // Switch to light mode (latte)
                body.classList.add('light-theme');
                applyScheme('latte');
                localStorage.setItem('theme', 'light');
                if (themeSelect) themeSelect.value = 'latte';
            }
            updateThemeButtonState();
        });
    }
    
    // Blog search functionality
    const searchInput = document.getElementById('search-input');
    const noResults = document.getElementById('no-results');
    const searchShortcut = document.getElementById('search-shortcut');
    let searchIndex = null;
    
    if (searchInput) {
        // Detect OS and update shortcut hint
        const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0 || 
                      navigator.userAgent.toUpperCase().indexOf('MAC') >= 0;
        if (searchShortcut && !isMac) {
            searchShortcut.textContent = 'Ctrl K';
        }
        
        // Keyboard shortcut: Cmd+K (Mac) or Ctrl+K (Windows/Linux)
        document.addEventListener('keydown', function(e) {
            if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
                e.preventDefault();
                searchInput.focus();
                searchInput.select();
            }
            // Escape to blur search
            if (e.key === 'Escape' && document.activeElement === searchInput) {
                searchInput.blur();
            }
        });
        
        // Load search index
        fetch('search-index.json')
            .then(response => response.json())
            .then(data => {
                searchIndex = data;
                // Create a map for quick lookup by link
                searchIndex.forEach(entry => {
                    entry.linkKey = entry.link;
                });
            })
            .catch(err => {
                console.log('Search index not available, using basic search');
            });
        
        // Helper function to highlight matched text
        const highlightMatch = (text, query) => {
            if (!query) return text;
            const regex = new RegExp(`(${query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
            return text.replace(regex, '<mark>$1</mark>');
        };
        
        // Helper function to get context around the match
        const getMatchContext = (text, query, contextLength = 50) => {
            if (!text) return null;
            const lowerText = text.toLowerCase();
            const lowerQuery = query.toLowerCase();
            const matchIndex = lowerText.indexOf(lowerQuery);
            
            if (matchIndex === -1) return null;
            
            const start = Math.max(0, matchIndex - contextLength);
            const end = Math.min(text.length, matchIndex + query.length + contextLength);
            
            let context = text.substring(start, end);
            if (start > 0) context = '...' + context;
            if (end < text.length) context = context + '...';
            
            return highlightMatch(context, query);
        };
        
        // Find search index entry by link
        const findIndexEntry = (link) => {
            if (!searchIndex) return null;
            return searchIndex.find(entry => entry.link === link || entry.linkKey === link);
        };
        
        searchInput.addEventListener('input', function() {
            const query = this.value.toLowerCase().trim();
            const blogPosts = document.querySelectorAll('.blog-list .blog-post');
            const directories = document.querySelectorAll('.directory-list .directory');
            let visiblePosts = 0;
            let visibleDirs = 0;
            
            // Filter blog posts
            blogPosts.forEach(post => {
                const title = post.dataset.title || '';
                const description = post.dataset.description || '';
                const postLink = post.querySelector('a')?.getAttribute('href') || '';
                
                // Get content from search index if available
                const indexEntry = findIndexEntry(postLink);
                const content = indexEntry?.content || '';
                
                const titleMatch = title.toLowerCase().includes(query);
                const descMatch = description.toLowerCase().includes(query);
                const contentMatch = content.toLowerCase().includes(query);
                const matches = titleMatch || descMatch || contentMatch;
                
                const searchMatchEl = post.querySelector('.search-match');
                const descriptionEl = post.querySelector('.post-description');
                
                if (matches && query) {
                    post.style.display = '';
                    visiblePosts++;
                    
                    // Show matched context
                    if (searchMatchEl) {
                        if (titleMatch) {
                            searchMatchEl.innerHTML = '<strong>Title:</strong> ' + highlightMatch(title, query);
                        } else if (descMatch) {
                            const context = getMatchContext(description, query);
                            searchMatchEl.innerHTML = '<strong>Description:</strong> ' + context;
                        } else if (contentMatch) {
                            const context = getMatchContext(content, query);
                            searchMatchEl.innerHTML = '<strong>Content:</strong> ' + context;
                        }
                        searchMatchEl.style.display = 'block';
                    }
                    
                    // Hide original description when showing match
                    if (descriptionEl) descriptionEl.style.display = 'none';
                } else if (!query) {
                    // No search query - show all posts normally
                    post.style.display = '';
                    visiblePosts++;
                    if (searchMatchEl) searchMatchEl.style.display = 'none';
                    if (descriptionEl) descriptionEl.style.display = '';
                } else {
                    // No match
                    post.style.display = 'none';
                    if (searchMatchEl) searchMatchEl.style.display = 'none';
                    if (descriptionEl) descriptionEl.style.display = '';
                }
            });
            
            // Filter directories
            directories.forEach(dir => {
                const name = dir.querySelector('a')?.textContent.toLowerCase() || '';
                const matches = name.includes(query);
                
                dir.style.display = matches ? '' : 'none';
                if (matches) visibleDirs++;
            });
            
            // Show/hide section headers based on visible items
            const blogList = document.querySelector('.blog-list');
            const directoryList = document.querySelector('.directory-list');
            
            if (blogList) {
                blogList.style.display = visiblePosts > 0 ? '' : 'none';
            }
            if (directoryList) {
                directoryList.style.display = visibleDirs > 0 ? '' : 'none';
            }
            
            // Show "no results" message if nothing matches
            if (noResults) {
                noResults.style.display = (visiblePosts === 0 && visibleDirs === 0 && query !== '') ? 'block' : 'none';
            }
        });
    }

    // Only run terminal-specific code if terminal exists (main site)
    if (terminal) {
        // Lucky button functionality
        if (luckyButton) {
            luckyButton.addEventListener('click', function() {
            // Get current scheme from data attribute
            let currentScheme = body.getAttribute('data-theme');
            if (!currentScheme || !schemes.includes(currentScheme)) {
                currentScheme = 'pinkie';
            }

            // Get available schemes excluding current one
            const availableSchemes = schemes.filter(scheme => scheme !== currentScheme);
            
            // Randomly select a new scheme
            const randomScheme = availableSchemes[Math.floor(Math.random() * availableSchemes.length)];
            
            // Apply the new scheme
            applyScheme(randomScheme);
            
            // Update dropdown and remove light-theme class
            if (themeSelect) themeSelect.value = randomScheme;
            body.classList.remove('light-theme');
            });
        }
        
        // Terminal typing effect
        const typingElements = document.querySelectorAll('.typing-effect');
        
        // Function to create a typing effect with a delay
        function createTypingEffect(element, delay) {
            const originalText = element.textContent;
            element.textContent = '';
            
            setTimeout(() => {
                let i = 0;
                const typingInterval = setInterval(() => {
                    if (i < originalText.length) {
                        element.textContent += originalText.charAt(i);
                        i++;
                    } else {
                        clearInterval(typingInterval);
                    }
                }, 50);
            }, delay);
        }
        
        // Apply typing effect to each element with a sequential delay
        typingElements.forEach((element, index) => {
            createTypingEffect(element, 500 * index);
        });
    }
});

// Monaco Editor Integration
(function() {
    function initMonacoEditor() {
        // Check if there are any code blocks to upgrade
        const codeBlocks = document.querySelectorAll('.chroma');
        if (codeBlocks.length === 0) return;

        // Load Monaco Editor loader
        const script = document.createElement('script');
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs/loader.min.js';
        script.onload = () => {
            require.config({ paths: { 'vs': 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min/vs' }});
            
            require(['vs/editor/editor.main'], function() {
                
                // Helper to update theme based on CSS variables
                const updateMonacoTheme = () => {
                    const style = getComputedStyle(document.body);
                    const bgColor = style.getPropertyValue('--bg-color').trim();
                    const textColor = style.getPropertyValue('--text-color').trim();
                    const headerColor = style.getPropertyValue('--terminal-header').trim();
                    const accentColor = style.getPropertyValue('--accent-color').trim();
                    const secondaryColor = style.getPropertyValue('--secondary-color').trim();
                    
                    monaco.editor.defineTheme('site-theme', {
                        base: 'vs-dark', // Always base on dark for this site structure
                        inherit: true,
                        rules: [
                            { token: '', foreground: textColor.replace('#', '') },
                            { token: 'keyword', foreground: accentColor.replace('#', '') },
                            { token: 'string', foreground: 'a3be8c' }, // Greenish
                            { token: 'number', foreground: 'b48ead' }, // Purple-ish
                            { token: 'comment', foreground: '616e87', fontStyle: 'italic' },
                            { token: 'delimiter', foreground: textColor.replace('#', '') },
                            { token: 'type', foreground: '8fbcbb' }, // Teal
                            { token: 'class', foreground: '8fbcbb' },
                            { token: 'function', foreground: '88c0d0' }, // Blue
                            { token: 'variable', foreground: textColor.replace('#', '') },
                            { token: 'operator', foreground: '81a1c1' }
                        ],
                        colors: {
                            'editor.background': headerColor, // Use terminal header color for code blocks
                            'editor.foreground': textColor,
                            'editor.lineHighlightBackground': '#00000000', // Transparent highlight
                            'editorLineNumber.foreground': secondaryColor,
                            'editorLineNumber.activeForeground': accentColor,
                            'scrollbarSlider.background': secondaryColor + '40',
                            'editor.selectionBackground': accentColor + '40'
                        }
                    });
                    monaco.editor.setTheme('site-theme');
                };

                // Apply initial theme
                updateMonacoTheme();

                // Observer for Theme Changes
                const themeObserver = new MutationObserver(() => {
                    updateMonacoTheme();
                });
                themeObserver.observe(document.body, { attributes: true, attributeFilter: ['class', 'data-theme'] });

                codeBlocks.forEach(block => {
                    // Robust Text Extraction: clone and clean
                    const clone = block.cloneNode(true);
                    
                    // Remove line numbers (.lnt, .ln, or first cell of table)
                    const lineNumbers = clone.querySelectorAll('.lnt, .ln, .lntable tr td:first-child');
                    lineNumbers.forEach(el => el.remove());

                    let code = clone.textContent;
                    // Clean up leading/trailing whitespace
                    code = code.replace(/^\s*\n/, '').replace(/\s*$/, '');

                    if (!code.trim()) return;
                    
                    // Detect language
                    let language = 'plaintext';
                    
                    // Check the wrapper first (most reliable)
                    const wrapper = block.closest('.code-wrapper');
                    if (wrapper && wrapper.dataset.lang) {
                        language = wrapper.dataset.lang;
                    }

                    // Fallback to internal classes if wrapper fails
                    if (language === 'plaintext') {
                        const checkClasses = (el) => {
                            if (!el || !el.classList) return;
                            el.classList.forEach(cls => {
                                if (cls.startsWith('language-')) {
                                    language = cls.replace('language-', '');
                                }
                            });
                            if (language === 'plaintext' && el.dataset.lang) {
                                language = el.dataset.lang;
                            }
                        };
                        checkClasses(block.querySelector('code'));
                        checkClasses(block.querySelector('.lntd:last-child code'));
                        checkClasses(block);
                    }

                    // Map aliases
                    const langMap = {
                        'bash': 'shell', 'sh': 'shell', 'zsh': 'shell', 'console': 'shell',
                        'golang': 'go', 'go': 'go',
                        'js': 'javascript', 'javascript': 'javascript',
                        'ts': 'typescript', 'typescript': 'typescript',
                        'py': 'python', 'python': 'python',
                        'c': 'c', 'cpp': 'cpp', 'c++': 'cpp',
                        'java': 'java',
                        'html': 'html',
                        'css': 'css',
                        'json': 'json',
                        'yaml': 'yaml', 'yml': 'yaml',
                        'rust': 'rust', 'rs': 'rust'
                    };
                    if (langMap[language]) language = langMap[language];

                    console.log(`[Monaco] Initializing block. Detected lang: '${language}'`);

                    // Create container
                    const container = document.createElement('div');
                    container.className = 'monaco-editor-container';
                    
                    // Height calculation
                    // Use a slightly larger line height for pixel fonts to breathe
                    const lineCount = code.split('\n').length;
                    const lineHeight = 22; 
                    const height = Math.max(lineCount * lineHeight + 24, 100);
                    
                    container.style.height = height + 'px';
                    container.style.marginBottom = '20px';
                    container.style.borderRadius = '8px';
                    container.style.overflow = 'hidden';
                    container.style.border = '1px solid var(--terminal-header)';

                    block.parentNode.insertBefore(container, block);
                    
                    // Initialize Editor
                    try {
                        const editor = monaco.editor.create(container, {
                            value: code,
                            language: language,
                            theme: 'site-theme',
                            readOnly: true,
                            automaticLayout: true,
                            minimap: { enabled: false },
                            scrollBeyondLastLine: false,
                            renderLineHighlight: 'none',
                            contextmenu: false,
                            // Pixel fonts often look best at specific sizes
                            fontFamily: "'Scientifica', 'Pokemon DP Pro', monospace",
                            fontSize: 15,
                            lineHeight: 22,
                            padding: { top: 12, bottom: 12 },
                            lineNumbers: 'on',
                            renderLineNumbers: (lineNumber) => lineNumber.toString(),
                            scrollbar: {
                                vertical: 'hidden',
                                horizontal: 'auto',
                                useShadows: false,
                                alwaysConsumeMouseWheel: false
                            },
                            overviewRulerLanes: 0,
                            hideCursorInOverviewRuler: true
                        });
                        
                        block.style.display = 'none';
                    } catch (e) {
                        console.error('Monaco Editor init failed:', e);
                        container.remove();
                        block.style.display = '';
                    }
                });
            });
        };
        document.body.appendChild(script);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initMonacoEditor);
    } else {
        initMonacoEditor();
    }
})();
