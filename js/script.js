document.addEventListener('DOMContentLoaded', () => {
    const mobileBtn = document.querySelector('.mobile-menu-btn');
    const nav = document.querySelector('.nav');

    // Simplistic mobile menu toggle for now (styles would need to support a class-based toggle if we had a full mobile menu overlay)
    // Since the CSS provided hides nav on mobile, let's create a basic overlay if needed or just log for now,
    // but to make it functional based on typically expected behavior:

    if (mobileBtn) {
        mobileBtn.addEventListener('click', () => {
            // For a production app, we would toggle a class like 'active' on a mobile menu container.
            // Given the current CSS structure, we'd need to add mobile styles for .nav to show it.
            // Let's alerting for now or adding a simple toggle if I adjust CSS.
            // I'll add a class to body to indicate menu open
            alert('Mobil menü entegrasyonu yapılacak.');
        });
    }

    // Header scroll background effect could be added here
    const header = document.querySelector('.header');
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            header.style.boxShadow = "var(--shadow-sm)";
        } else {
            header.style.boxShadow = "none";
        }
    });
});
