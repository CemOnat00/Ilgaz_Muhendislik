document.addEventListener('DOMContentLoaded', () => {
    // --- MOBILE MENU ---
    const mobileBtn = document.querySelector('.mobile-menu-btn');
    const nav = document.querySelector('.nav'); // Using class for consistency if adjusted
    const header = document.querySelector('.header');

    // For Tailwind pages (iletisim, etc.) the selector might be different
    const tailwindMobileBtn = document.querySelector('header button.lg\\:hidden');
    
    function toggleMobileMenu() {
        // Implement a simple mobile menu or just show a message for now
        // In a real app, this would toggle a sidebar or dropdown
        console.log('Mobile menu toggled');
    }

    if (mobileBtn) {
        mobileBtn.addEventListener('click', toggleMobileMenu);
    }
    
    if (tailwindMobileBtn) {
        tailwindMobileBtn.addEventListener('click', () => {
            alert('Mobil menü yakında eklenecektir.');
        });
    }

    // --- HEADER SCROLL EFFECT ---
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            header?.classList.add('scrolled');
            if (header) header.style.boxShadow = "0 4px 6px -1px rgb(0 0 0 / 0.1)";
        } else {
            header?.classList.remove('scrolled');
            if (header) header.style.boxShadow = "none";
        }
    });

    // --- SEARCH BAR MOCK ---
    const searchBtn = document.querySelector('button i.fa-magnifying-glass')?.parentElement;
    if (searchBtn) {
        searchBtn.addEventListener('click', () => {
            const query = prompt('Aramak istediğiniz ürün veya hizmeti yazın:');
            if (query) {
                window.location.href = `/urunler.html?search=${encodeURIComponent(query)}`;
            }
        });
    }
});
