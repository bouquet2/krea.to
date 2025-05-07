document.addEventListener('DOMContentLoaded', function() {
    // Theme toggle functionality
    const themeButton = document.getElementById('theme-button');
    const body = document.body;
    
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
    
    // Blinking cursor effect
    const cursor = document.querySelector('.cursor');
    setInterval(() => {
        cursor.style.opacity = cursor.style.opacity === '0' ? '1' : '0';
    }, 500);
    
    // Create a terminal command input for the last prompt
    const terminalContent = document.querySelector('.terminal-content');
    const lastPrompt = document.querySelector('.terminal-section:last-child .prompt');
    
    // Allow typing in the last terminal prompt
    lastPrompt.addEventListener('click', function() {
        const existingInput = document.getElementById('terminal-input');
        if (!existingInput) {
            const input = document.createElement('input');
            input.type = 'text';
            input.id = 'terminal-input';
            input.style.background = 'transparent';
            input.style.border = 'none';
            input.style.outline = 'none';
            input.style.color = 'inherit';
            input.style.fontFamily = 'inherit';
            input.style.fontSize = 'inherit';
            input.style.width = '70%';
            
            const cursorSpan = document.querySelector('.cursor');
            cursorSpan.parentNode.replaceChild(input, cursorSpan);
            input.focus();
            
            // Handle command execution on Enter
            input.addEventListener('keydown', function(e) {
                if (e.key === 'Enter') {
                    const command = input.value.trim();
                    if (command) {
                        handleCommand(command);
                    }
                    // Reset the prompt
                    input.value = '';
                }
            });
        }
    });
    
    // Simple command handler
    function handleCommand(command) {
        let response = '';
        
        // Simple command parsing
        if (command.startsWith('echo ')) {
            response = command.substring(5);
        } else if (command === 'help') {
            response = 'Available commands: help, about, echo [text], clear, date, ls, cd [dir], yes';
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
            response = `<div class="links">
                <a href="https://kreato.dev/Blogs/" target="_blank" class="file">blog.md</a>
                <a href="https://github.com/kreatoo" target="_blank" class="file">github.txt</a>
            </div>`;
        } else if (command.startsWith('cd ')) {
            const dir = command.substring(3).trim();
            if (dir === '..' || dir === '/' || dir === '~') {
                response = 'Already in root directory';
            } else if (dir === 'blog' || dir === 'github' || dir === 'projects') {
                response = `Changed directory to ${dir}`;
            } else if (dir === 'secret') {
                window.location.href = 'https://www.youtube.com/watch?v=dQw4w9WgXcQ';
                return;
            } else {
                response = `cd: no such directory: ${dir}`;
            }
        } else {
            response = `Command not found: ${command}`;
        }
        
        // Create a new section with the response
        const newSection = document.createElement('div');
        newSection.className = 'terminal-section';
        newSection.innerHTML = `
            <div class="prompt">kreato@akiri:~$ <span>${command}</span></div>
            <div class="output"><p>${response}</p></div>
        `;
        
        // Insert the new section before the last section
        const lastSection = document.querySelector('.terminal-section:last-child');
        terminalContent.insertBefore(newSection, lastSection);
        
        // Scroll to the bottom
        terminalContent.scrollTop = terminalContent.scrollHeight;
    }
}); 