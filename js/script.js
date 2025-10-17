document.addEventListener('DOMContentLoaded', function() {

    // Accessibility mode functionality
    const accessibilityButton = document.getElementById('accessibility-button');
    const backgroundButton = document.getElementById('background-button');
    const transparencyButton = document.getElementById('transparency-button');
    const luckyButton = document.getElementById('lucky-button');
    const themeSelect = document.getElementById('theme-select');
    const fontSizeRange = document.getElementById('font-size-range');
    const fontSizeValue = document.getElementById('font-size-value');
    const body = document.body;
    const terminal = document.querySelector('.terminal');
    
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
    const savedBackground = localStorage.getItem('background');
    if (savedBackground === 'disabled') {
        body.classList.add('no-background');
        backgroundButton.textContent = 'Show Background';
    } else {
        backgroundButton.textContent = 'Hide Background';
    }
    
    // Check for saved transparency preference
    const savedTransparency = localStorage.getItem('transparency');
    if (savedTransparency === 'disabled') {
        if (terminal) terminal.classList.add('no-transparency');
        transparencyButton.textContent = 'Enable Transparency';
    } else {
        transparencyButton.textContent = 'Disable Transparency';
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
    
    // Toggle transparency on button click
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

    // Font size control
    const savedFontSize = localStorage.getItem('fontSize') || '1.1';
    if (fontSizeRange) {
        fontSizeRange.value = savedFontSize;
        if (fontSizeValue) fontSizeValue.textContent = `${savedFontSize}em`;
        body.style.fontSize = `${savedFontSize}em`;
        
        fontSizeRange.addEventListener('input', function() {
            const size = this.value;
            body.style.fontSize = `${size}em`;
            if (fontSizeValue) fontSizeValue.textContent = `${size}em`;
            localStorage.setItem('fontSize', size);
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
            body.style.setProperty('--background-image', `url('../assets/${backgroundImage}.jpg')`);
        }
    };
    
    // Check for saved color scheme or default to mocha
    const savedScheme = localStorage.getItem('colorScheme') || 'mocha';
    if (themeSelect) {
        themeSelect.value = savedScheme;
    }
    
    // Set initial background image
    let initialBackground = savedScheme;
    if (['mocha', 'frappe', 'macchiato', 'latte'].includes(savedScheme)) {
        initialBackground = 'catppuccin-dark';
    }
    body.style.setProperty('--background-image', `url('../assets/${initialBackground}.jpg')`);
    
    // Apply saved scheme
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
        // Track current directory
        let currentDir = '~';
        
        // Track if terminal is being actively used
        let isTerminalActive = false;
        
        // List of available commands
        // 'popo' command is intentionally omitted to hide from tab completion/help
        const commands = ['help', 'about', 'echo', 'clear', 'date', 'ls', 'cd', 'yes', 'cat', 'ifconfig', 'upower', 'scheme', 'background', 'transparency'];
        
        // List of available files/directories
        const files = ['blog.md', 'github.txt', 'infra.tf', 'status.sh'];
        const directories = ['.secret', '.kube', 'blog', 'github', 'projects'];
        const schemes = ['mocha', 'frappe', 'latte', 'macchiato', 'gruvbox', 'nord', 'tokyonight', 'monokai', 'onedark', 'solarized', 'kanagawa'];
        
        // Color configurations for each scheme
        // Lucky button functionality
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
        
        // Function to update all prompts with current directory
        function updatePrompts() {
            // Only update the last prompt (current input) and any new prompts
            const lastPrompt = document.querySelector('.terminal-section:last-child .prompt');
            if (lastPrompt) {
                const parts = lastPrompt.innerHTML.split('$');
                lastPrompt.innerHTML = `kreato@akiri:${currentDir}$ ${parts[1]}`;
            }
        }
        
        // Function to create a new input field
        function createNewInput() {
            const terminalContent = document.querySelector('.terminal-content');
            const newLastSection = document.createElement('div');
            newLastSection.className = 'terminal-section';
            newLastSection.innerHTML = `
                <div class="prompt">kreato@akiri:${currentDir}$ <input type="text" id="terminal-input" style="background: transparent; border: none; outline: none; color: inherit; font-family: inherit; font-size: inherit; width: 70%;"></div>
            `;
            terminalContent.appendChild(newLastSection);
            
            const input = document.getElementById('terminal-input');
            
            // Only focus if terminal is active
            if (isTerminalActive) {
                input.focus();
            }
            
            input.addEventListener('keydown', function(e) {
                if (e.key === 'Enter') {
                    const command = input.value.trim();
                    if (command) {
                        handleCommand(command);
                    }
                    input.value = '';
                } else if (e.key === 'Tab') {
                    e.preventDefault(); // Prevent default tab behavior
                    
                    const currentInput = input.value.trim();
                    const words = currentInput.split(' ');
                    const currentWord = words[words.length - 1];
                    
                    let completions = [];
                    let showAllOptions = false;
                    
                    // Check if this is a double tab (within 500ms of the last tab)
                    const now = Date.now();
                    if (input.lastTabTime && now - input.lastTabTime < 500) {
                        showAllOptions = true;
                    }
                    input.lastTabTime = now;
                    
                    // If we're completing the first word (command)
                    if (words.length === 1) {
                        completions = commands.filter(cmd => cmd.startsWith(currentWord));
                    } 
                    // If we're completing after 'cd', 'cat', or 'scheme'
                    else if (words.length === 2) {
                        if (words[0] === 'cd') {
                            completions = directories.filter(dir => dir.startsWith(currentWord));
                        } else if (words[0] === 'cat') {
                            completions = files.filter(file => file.startsWith(currentWord));
                        } else if (words[0] === 'scheme') {
                            completions = schemes.filter(scheme => scheme.startsWith(currentWord));
                        }
                    }
                    
                    // If we have exactly one completion, use it
                    if (completions.length === 1 && !showAllOptions) {
                        words[words.length - 1] = completions[0];
                        input.value = words.join(' ') + (words[0] === 'cd' ? '/' : ' ');
                    } 
                    // If we have multiple completions or it's a double tab, show them
                    else if (completions.length > 1 || showAllOptions) {
                        // For double tab, show all available options for the command
                        let optionsToShow = completions;
                        if (showAllOptions) {
                            if (words[0] === 'scheme') {
                                optionsToShow = schemes;
                            } else if (words[0] === 'cd') {
                                optionsToShow = directories;
                            } else if (words[0] === 'cat') {
                                optionsToShow = files;
                            } else if (words.length === 1) {
                                optionsToShow = commands;
                            }
                        }
                        
                        const newSection = document.createElement('div');
                        newSection.className = 'terminal-section';
                        newSection.innerHTML = `
                            <div class="prompt">kreato@akiri:${currentDir}$ <span>${currentInput}</span></div>
                            <div class="output"><p>${optionsToShow.join('  ')}</p></div>
                        `;
                        
                        // Insert the completions before the last section
                        const lastSection = document.querySelector('.terminal-section:last-child');
                        terminalContent.insertBefore(newSection, lastSection);
                    }
                }
            });

            // Add focus and blur events to track active state
            input.addEventListener('focus', function() {
                isTerminalActive = true;
            });

            input.addEventListener('blur', function() {
                isTerminalActive = false;
            });
        }
        
        // Create a terminal command input for the last prompt
        const terminalContent = document.querySelector('.terminal-content');
        
        // Add click event listener to terminal content
        terminalContent.addEventListener('click', function() {
            isTerminalActive = true;
            const input = document.getElementById('terminal-input');
            if (input) {
                input.focus();
            }
        });
        
        // Add blur event listener to document to track when terminal loses focus
        document.addEventListener('click', function(e) {
            if (!terminalContent.contains(e.target)) {
                isTerminalActive = false;
            }
        });
        
        // Create initial input
        createNewInput();
        
        // Simple command handler
        function handleCommand(command) {
            let response = '';
            
            // Check for command prefixes first
            if (command.startsWith('echo ')) {
                response = command.substring(5);
            } else if (command.startsWith('cd ')) {
                const dir = command.substring(3).trim();
                switch (dir) {
                    case '~':
                        currentDir = '~';
                        response = 'Changed to home directory';
                        break;
                    case 'blog':
                    case 'github':
                    case 'projects':
                        currentDir = dir;
                        response = `Changed directory to ${dir}`;
                        break;
                    case '.secret':
                        window.location.href = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';
                        return;
                    case '.kube':
                        currentDir = '.kube';
                        response = `<div class="links">
                            <a href="#" class="file">config</a>
                        </div>`;
                        break;
                    case '/':
                        response = `cd: i ate it`;
                        break;
                    case '..':
                        response = `Where is bro going?`;
                        break;
                    default:
                        response = `cd: no such directory: ${dir}`;
                }
                updatePrompts();
            } else if (command.startsWith('cat ')) {
                const file = command.substring(4).trim();
                switch (file) {
                    case 'config':
                        if (currentDir === '.kube') {
                            window.location.href = 'https://www.youtube.com/watch?v=9wvEwPLcLcA';
                            return;
                        }
                        response = `cat: ${file}: No such file or directory`;
                        break;
                    case 'blog.md':
                        response = 'Redirecting to blog...';
                        setTimeout(() => {
                            window.location.href = 'blog/index.html';
                        }, 1000);
                        break;
                    case 'github.txt':
                        response = 'Redirecting to GitHub...';
                        setTimeout(() => {
                            window.location.href = 'https://github.com/kreatoo';
                        }, 1000);
                        break;
                    case 'infra.tf':
                        response = 'Redirecting to infrastructure repository...';
                        setTimeout(() => {
                            window.location.href = 'https://github.com/kreatoo/bouquet2';
                        }, 1000);
                        break;
                    case 'status.sh':
                        response = 'Redirecting to status page...';
                        setTimeout(() => {
                            window.location.href = 'https://status.krea.to';
                        }, 1000);
                        break;
                    default:
                        response = `cat: ${file}: No such file or directory`;
                }
            } else {
                switch (command) {
                    case 'popo':
                        response = '<img src="assets/popo.webp" alt="popo" style="max-width:300px;max-height:300px;">';
                        break;
                    case 'help':
                        response = 'Available commands: help, about, echo [text], clear, date, ls, cd [dir], yes, cat [file], ifconfig, upower, scheme [theme], background [on|off], transparency [on|off]';
                        break;
                    case 'about':
                        response = 'Kreato - Tinkerer and Developer';
                        break;
                    case 'clear':
                        // Clear the terminal except for the last prompt
                        const sections = document.querySelectorAll('.terminal-section');
                        for (let i = 0; i < sections.length - 1; i++) {
                            sections[i].style.display = 'none';
                        }
                        return;
                    case 'date':
                        response = new Date().toString();
                        break;
                    case 'yes':
                        response = 'y';
                        break;
                    case 'yes no':
                        response = 'Do you have nothing else to do other than look for things in this site? Your life must be boring.';
                        break;
                    case 'yes please':
                        response = '🥺'
                        break;
                    case 'ls':
                        if (currentDir === '.kube') {
                            response = `<div class="links">
                                <a href="#" class="file">config</a>
                            </div>`;
                        } else {
                            response = `<div class="links">
                                <a href="#" class="folder">.secret/</a>
                                <a href="#" class="folder">.kube/</a>
                                <a href="blog/index.html" class="file">blog.md</a>
                                <a href="https://github.com/kreatoo" target="_blank" class="file">github.txt</a>
                                <a href="https://github.com/kreatoo/bouquet2" target="_blank" class="file">infra.tf</a>
                                <a href="https://status.krea.to" target="_blank" class="file">status.sh</a>
                            </div>`;
                        }
                        break;
                    case 'ifconfig':
                        response = `<div class="gif-output"><img src="assets/out.gif" alt="Network configuration animation" style="max-width: 100%; height: auto;"></div>`;
                        break;
                    case 'upower':
                        response = `<div class="gif-output"><img src="assets/discord-this.gif" alt="Discord this animation" style="max-width: 100%; height: auto;"></div>`;
                        break;
                    case 'scheme':
                        response = `Available schemes: ${schemes.map((s, i) => i === 0 ? `${s} (default)` : s).join(', ')}`;
                        break;
                    case 'background':
                        response = 'Usage: background [on|off] - Toggle background image display';
                        break;
                    case 'transparency':
                        response = 'Usage: transparency [on|off] - Toggle terminal transparency';
                        break;
                    default:
                        if (command.startsWith('scheme ')) {
                            const scheme = command.substring(7).trim();
                            
                            // Check if the scheme exists in our schemes array
                            if (schemes.includes(scheme)) {
                                // Apply the selected scheme
                                applyScheme(scheme);
                                
                                // Remove light-theme class and update dropdown when switching schemes
                                body.classList.remove('light-theme');
                                if (themeSelect) themeSelect.value = scheme;
                                
                                response = scheme === schemes[0] ? 
                                    `Switched to ${scheme} scheme (default)` : 
                                    `Switched to ${scheme} scheme`;
                            } else {
                                response = `Unknown scheme: ${scheme}. Use 'scheme' to see available schemes.`;
                            }
                        } else if (command.startsWith('background ')) {
                            const action = command.substring(11).trim();
                            
                            if (action === 'on') {
                                setTimeout(() => {
                                    body.classList.remove('no-background');
                                }, 50);
                                localStorage.setItem('background', 'enabled');
                                backgroundButton.textContent = 'Hide Background';
                                response = 'Background enabled';
                            } else if (action === 'off') {
                                setTimeout(() => {
                                    body.classList.add('no-background');
                                }, 50);
                                localStorage.setItem('background', 'disabled');
                                backgroundButton.textContent = 'Show Background';
                                response = 'Background disabled';
                            } else {
                                response = `Invalid argument: ${action}. Use 'background on' or 'background off'`;
                            }
                        } else if (command.startsWith('transparency ')) {
                            const action = command.substring(13).trim();
                            
                            if (action === 'on') {
                                if (terminal) terminal.classList.remove('no-transparency');
                                localStorage.setItem('transparency', 'enabled');
                                transparencyButton.textContent = 'Disable Transparency';
                                response = 'Transparency enabled';
                            } else if (action === 'off') {
                                if (terminal) terminal.classList.add('no-transparency');
                                localStorage.setItem('transparency', 'disabled');
                                transparencyButton.textContent = 'Enable Transparency';
                                response = 'Transparency disabled';
                            } else {
                                response = `Invalid argument: ${action}. Use 'transparency on' or 'transparency off'`;
                            }
                        } else {
                            response = `Command not found: ${command}`;
                        }
                }
            }
            
            // Create a new section with the response
            const newSection = document.createElement('div');
            newSection.className = 'terminal-section';
            newSection.innerHTML = `
                <div class="prompt">kreato@akiri:${currentDir}$ <span>${command}</span></div>
                <div class="output"><p>${response}</p></div>
            `;
            
            // Insert the new section before the last section
            const lastSection = document.querySelector('.terminal-section:last-child');
            terminalContent.insertBefore(newSection, lastSection);
            
            // Remove the last section and create a new input
            lastSection.remove();
            createNewInput();
        }
    }
}); 
