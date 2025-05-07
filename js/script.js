document.addEventListener('DOMContentLoaded', function() {
    // Accessibility mode functionality
    const accessibilityButton = document.getElementById('accessibility-button');
    const terminal = document.querySelector('.terminal');
    
    // Check for saved accessibility preference
    const savedAccessibility = localStorage.getItem('accessibility');
    if (savedAccessibility === 'enabled') {
        terminal.classList.add('accessibility-mode');
        accessibilityButton.textContent = 'Standard Font';
    }
    
    // Toggle accessibility mode on button click
    accessibilityButton.addEventListener('click', function() {
        terminal.classList.toggle('accessibility-mode');
        
        if (terminal.classList.contains('accessibility-mode')) {
            localStorage.setItem('accessibility', 'enabled');
            accessibilityButton.textContent = 'Standard Font';
        } else {
            localStorage.setItem('accessibility', 'disabled');
            accessibilityButton.textContent = 'Accessibility Mode';
        }
    });

    // Theme toggle functionality
    const themeButton = document.getElementById('theme-button');
    const body = document.body;
    
    // Track current directory
    let currentDir = '~';
    
    // Track if terminal is being actively used
    let isTerminalActive = false;
    
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
    
    // Check for saved theme preference or default to dark
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'light') {
        body.classList.add('light-theme');
        themeButton.textContent = 'Dark Mode';
    } else {
        themeButton.textContent = 'Light Mode';
    }
    
    // Toggle theme on button click
    themeButton.addEventListener('click', function() {
        body.classList.toggle('light-theme');
        
        if (body.classList.contains('light-theme')) {
            localStorage.setItem('theme', 'light');
            themeButton.textContent = 'Dark Mode';
        } else {
            localStorage.setItem('theme', 'dark');
            themeButton.textContent = 'Light Mode';
        }
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
        
        // Simple command parsing
        if (command.startsWith('echo ')) {
            response = command.substring(5);
        } else if (command === 'help') {
            response = 'Available commands: help, about, echo [text], clear, date, ls, cd [dir], yes, cat [file]';
        } else if (command === 'about') {
            response = 'Kreato - Tinkerer and Developer';
        } else if (command === 'clear') {
            // Clear the terminal except for the last prompt
            const sections = document.querySelectorAll('.terminal-section');
            for (let i = 0; i < sections.length - 1; i++) {
                sections[i].style.display = 'none';
            }
            return;
        } else if (command === 'date') {
            response = new Date().toString();
        } else if (command === 'yes') {
            response = 'y';
        } else if (command === 'yes no') {
            response = 'Do you have nothing else to do other than look for things in this site? Your life must be boring.';
        } else if (command === 'ls') {
            if (currentDir === '.kube') {
                response = `<div class="links">
                    <a href="#" class="file">config</a>
                </div>`;
            } else {
                response = `<div class="links">
                    <a href="#" class="folder">.secret/</a>
                    <a href="#" class="folder">.kube/</a>
                    <a href="https://kreato.dev/Blogs/" target="_blank" class="file">blog.md</a>
                    <a href="https://github.com/kreatoo" target="_blank" class="file">github.txt</a>
                    <a href="https://github.com/kreatoo/bouquet2" target="_blank" class="file">infra.tf</a>
                    <a href="https://status.krea.to" target="_blank" class="file">status.sh</a>
                </div>`;
            }
        } else if (command.startsWith('cd ')) {
            const dir = command.substring(3).trim();
            if (dir === '~') {
                currentDir = '~';
                response = 'Changed to home directory';
            } else if (dir === 'blog' || dir === 'github' || dir === 'projects') {
                currentDir = dir;
                response = `Changed directory to ${dir}`;
            } else if (dir === '.secret') {
                window.location.href = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';
                return;
            } else if (dir === '.kube') {
                currentDir = '.kube';
                response = `<div class="links">
                    <a href="#" class="file">config</a>
                </div>`;
            } else if (dir === '/') {
                response = `cd: i ate it`;
            } else if (dir === '..') {
                response = `Where is bro going?`;
            } else {
                response = `cd: no such directory: ${dir}`;
            }
            updatePrompts();
        } else if (command.startsWith('cat ')) {
            const file = command.substring(4).trim();
            if (file === 'config' && currentDir === '.kube') {
                window.location.href = 'https://www.youtube.com/watch?v=9wvEwPLcLcA';
                return;
            } else if (file === 'blog.md') {
                response = 'Redirecting to blog...';
                setTimeout(() => {
                    window.location.href = 'https://kreato.dev/Blogs/';
                }, 1000);
            } else if (file === 'github.txt') {
                response = 'Redirecting to GitHub...';
                setTimeout(() => {
                    window.location.href = 'https://github.com/kreatoo';
                }, 1000);
            } else {
                response = `cat: ${file}: No such file or directory`;
            }
        } else {
            response = `Command not found: ${command}`;
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
        
        // Only scroll to bottom if terminal is being actively used
        if (isTerminalActive) {
            terminalContent.scrollTop = terminalContent.scrollHeight;
        }
    }
}); 