// Smooth scrolling for navigation links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
  anchor.addEventListener('click', function (e) {
    e.preventDefault();
    const target = document.querySelector(this.getAttribute('href'));
    if (target) {
      target.scrollIntoView({
        behavior: 'smooth',
        block: 'start'
      });
    }
  });
});

// Navbar background on scroll
window.addEventListener('scroll', () => {
  const navbar = document.querySelector('.navbar');
  if (window.scrollY > 50) {
    navbar.style.background = 'rgba(30, 41, 59, 0.95)';
    navbar.style.backdropFilter = 'blur(10px)';
  } else {
    navbar.style.background = 'var(--bg-secondary)';
    navbar.style.backdropFilter = 'none';
  }
});

// Animate elements on scroll
const observerOptions = {
  threshold: 0.1,
  rootMargin: '0px 0px -50px 0px'
};

const observer = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      entry.target.style.opacity = '1';
      entry.target.style.transform = 'translateY(0)';
    }
  });
}, observerOptions);

document.querySelectorAll('.feature-card, .download-card, .doc-card').forEach(el => {
  el.style.opacity = '0';
  el.style.transform = 'translateY(20px)';
  el.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
  observer.observe(el);
});

// Copy code to clipboard
document.querySelectorAll('.code-block').forEach(codeBlock => {
  codeBlock.style.cursor = 'pointer';
  codeBlock.title = 'Click to copy';
  
  codeBlock.addEventListener('click', async () => {
    const code = codeBlock.textContent.trim();
    try {
      await navigator.clipboard.writeText(code);
      const originalText = codeBlock.textContent;
      codeBlock.textContent = '✓ Copied!';
      setTimeout(() => {
        codeBlock.textContent = originalText;
      }, 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  });
});

// Add download button animation
document.querySelectorAll('.btn').forEach(btn => {
  btn.addEventListener('click', function(e) {
    const href = this.getAttribute('href');
    if (href && href.startsWith('#')) {
      return;
    }
    
    // Add click effect
    this.style.transform = 'scale(0.95)';
    setTimeout(() => {
      this.style.transform = '';
    }, 150);
  });
});

// Mobile menu toggle (for future enhancement)
const createMobileMenu = () => {
  const nav = document.querySelector('.navbar');
  const container = nav.querySelector('.container');
  
  const menuButton = document.createElement('button');
  menuButton.className = 'mobile-menu-btn';
  menuButton.innerHTML = '☰';
  menuButton.style.display = 'none';
  menuButton.style.background = 'none';
  menuButton.style.border = 'none';
  menuButton.style.color = 'var(--text-primary)';
  menuButton.style.fontSize = '1.5rem';
  menuButton.style.cursor = 'pointer';
  
  const navLinks = document.querySelector('.nav-links');
  
  menuButton.addEventListener('click', () => {
    navLinks.style.display = navLinks.style.display === 'flex' ? 'none' : 'flex';
  });
  
  // Show menu button on mobile
  const checkMobile = () => {
    if (window.innerWidth <= 768) {
      menuButton.style.display = 'block';
      navLinks.style.display = 'none';
      navLinks.style.flexDirection = 'column';
      navLinks.style.position = 'absolute';
      navLinks.style.top = '100%';
      navLinks.style.left = '0';
      navLinks.style.right = '0';
      navLinks.style.background = 'var(--bg-secondary)';
      navLinks.style.padding = '1rem';
      navLinks.style.borderBottom = '1px solid var(--border)';
    } else {
      menuButton.style.display = 'none';
      navLinks.style.display = 'flex';
      navLinks.style.flexDirection = 'row';
      navLinks.style.position = 'static';
    }
  };
  
  container.insertBefore(menuButton, navLinks);
  checkMobile();
  window.addEventListener('resize', checkMobile);
};

// Initialize mobile menu
createMobileMenu();

// Console message
console.log('%c🌊 SmirnovTorrent', 'font-size: 24px; font-weight: bold; color: #0ea5e9;');
console.log('%cLightweight BitTorrent Client', 'font-size: 14px; color: #94a3b8;');
console.log('%cVisit https://github.com/ainsonet/smirnovtorrent for more info', 'font-size: 12px; color: #cbd5e1;');
