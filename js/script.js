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
    
    // Hamburger menu toggle functionality
    const hamburgerMenu = document.querySelector('.hamburger-menu');
    const hamburgerIcon = document.querySelector('.hamburger-icon');
    const menuItems = document.querySelector('.menu-items');
    
    if (hamburgerIcon) {
        hamburgerIcon.addEventListener('click', function(e) {
            e.stopPropagation();
            hamburgerMenu.classList.toggle('menu-open');
        });
    }
    
    // Close menu when clicking outside
    document.addEventListener('click', function(e) {
        if (hamburgerMenu && !hamburgerMenu.contains(e.target)) {
            hamburgerMenu.classList.remove('menu-open');
        }
    });
    
    // Prevent menu from closing when clicking inside menu items
    if (menuItems) {
        menuItems.addEventListener('click', function(e) {
            e.stopPropagation();
        });
    }
    
    // Color scheme configurations
    const schemeConfigs = {
        'frappe': {
            '--bg-color': '#303446',
            '--text-color': '#c6d0f5',
            '--accent-color': '#f4b8e4',
            '--secondary-color': '#ca9ee6',
            '--terminal-header': '#292c3c',
            '--link-color': '#8caaee'
        },
        'mocha': {
            '--bg-color': '#1e1e2e',
            '--text-color': '#cdd6f4',
            '--accent-color': '#f5c2e7',
            '--secondary-color': '#cba6f7',
            '--terminal-header': '#181825',
            '--link-color': '#89b4fa'
        },
        'latte': {
            '--bg-color': '#eff1f5',
            '--text-color': '#4c4f69',
            '--accent-color': '#ea76cb',
            '--secondary-color': '#8839ef',
            '--terminal-header': '#e6e9ef',
            '--link-color': '#1e66f5'
        },
        'macchiato': {
            '--bg-color': '#24273a',
            '--text-color': '#cad3f5',
            '--accent-color': '#f5bde6',
            '--secondary-color': '#c6a0f6',
            '--terminal-header': '#1e2030',
            '--link-color': '#8aadf4'
        },
        'gruvbox': {
            '--bg-color': '#282828',
            '--text-color': '#ebdbb2',
            '--accent-color': '#b8bb26',
            '--secondary-color': '#fabd2f',
            '--terminal-header': '#1d2021',
            '--link-color': '#83a598'
        },
        'nord': {
            '--bg-color': '#2e3440',
            '--text-color': '#eceff4',
            '--accent-color': '#88c0d0',
            '--secondary-color': '#81a1c1',
            '--terminal-header': '#3b4252',
            '--link-color': '#5e81ac'
        },
        'tokyonight': {
            '--bg-color': '#1a1b26',
            '--text-color': '#a9b1d6',
            '--accent-color': '#bb9af7',
            '--secondary-color': '#7aa2f7',
            '--terminal-header': '#24283b',
            '--link-color': '#7dcfff'
        },
        'monokai': {
            '--bg-color': '#1B1E1C',
            '--text-color': '#F5F5F5',
            '--accent-color': '#FF1493',
            '--secondary-color': '#AF87FF',
            '--terminal-header': '#333333',
            '--link-color': '#5FD7FF'
        },
        'onedark': {
            '--bg-color': '#282c34',
            '--text-color': '#abb2bf',
            '--accent-color': '#c678dd',
            '--secondary-color': '#61afef',
            '--terminal-header': '#21252b',
            '--link-color': '#e06c75'
        },
        'solarized': {
            '--bg-color': '#002b36',
            '--text-color': '#93a1a1',
            '--accent-color': '#d33682',
            '--secondary-color': '#268bd2',
            '--terminal-header': '#073642',
            '--link-color': '#2aa198'
        },
        'kanagawa': {
            '--bg-color': '#1F1F28',
            '--text-color': '#DCD7BA',
            '--accent-color': '#957FB8',
            '--secondary-color': '#7FB4CA',
            '--terminal-header': '#16161D',
            '--link-color': '#98BB6C'
        }
    };

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
        body.style.setProperty('font-size', `${savedFontSize}em`, 'important');
        
        fontSizeRange.addEventListener('input', function() {
            const size = this.value;
            body.style.setProperty('font-size', `${size}em`, 'important');
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
            body.style.setProperty('font-family', "'Pokemon DP Pro', sans-serif", 'important');
        } else {
            body.style.setProperty('font-family', "'Scientifica', sans-serif", 'important');
        }
        
        fontSelect.addEventListener('change', function() {
            const selectedFont = this.value;
            
            if (selectedFont === 'pokemon') {
                body.style.setProperty('font-family', "'Pokemon DP Pro', sans-serif", 'important');
            } else {
                body.style.setProperty('font-family', "'Scientifica', sans-serif", 'important');
            }
            
            localStorage.setItem('fontFamily', selectedFont);
        });
    }

    // Theme toggle functionality
    const themeButton = document.getElementById('theme-button');

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
        const config = schemeConfigs[schemeName];
        if (config) {
            for (const [property, value] of Object.entries(config)) {
                document.documentElement.style.setProperty(property, value);
            }
            body.setAttribute('data-theme', schemeName);
            localStorage.setItem('colorScheme', schemeName);
            
            // Update background image based on scheme
            // Catppuccin variants use the same background
            let backgroundImage = schemeName;
            if (['mocha', 'frappe', 'macchiato', 'latte'].includes(schemeName)) {
                backgroundImage = 'catppuccin-dark';
            }
            
            // Determine file extension based on theme
            let extension = 'jpg'; // default for most themes
            if (['mocha', 'frappe', 'macchiato', 'latte', 'catppuccin-dark', 'monokai', 'onedark', 'kanagawa'].includes(backgroundImage)) {
                extension = 'avif';
            } else if (['solarized', 'gruvbox', 'tokyonight', 'nord'].includes(backgroundImage)) {
                extension = 'png';
            }
            
            // Only set background if terminal exists (main site only)
            if (terminal) {
                document.documentElement.style.setProperty('--background-image', `url('../assets/${backgroundImage}.${extension}')`);
            }
        }
    };
    
    // Check for saved color scheme or default to mocha
    const savedScheme = localStorage.getItem('colorScheme') || 'mocha';
    if (themeSelect) {
        themeSelect.value = savedScheme;
    }
    
    // Apply saved scheme (which also sets the background image)
    applyScheme(savedScheme);
    
    // Update light-theme class if using latte
    if (savedScheme === 'latte') {
        body.classList.add('light-theme');
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
                // Switch to dark mode (mocha)
                body.classList.remove('light-theme');
                applyScheme('mocha');
                localStorage.setItem('theme', 'dark');
                if (themeSelect) themeSelect.value = 'mocha';
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
    
    // Only run terminal-specific code if terminal exists (main site)
    if (terminal) {
        const schemes = ['mocha', 'frappe', 'latte', 'macchiato', 'gruvbox', 'nord', 'tokyonight', 'monokai', 'onedark', 'solarized', 'kanagawa'];
        
        // Lucky button functionality
        if (luckyButton) {
            luckyButton.addEventListener('click', function() {
            // Get current scheme from data attribute
            const currentScheme = body.getAttribute('data-theme') || 'mocha';

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
