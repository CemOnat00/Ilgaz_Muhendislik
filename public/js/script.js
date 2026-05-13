document.addEventListener('DOMContentLoaded', () => {
    // --- MOBILE MENU LOGIC ---
    const mobileBtn = document.querySelector('.lg\\:hidden i.fa-bars')?.parentElement;
    
    // Create Mobile Menu Overlay if it doesn't exist
    if (!document.getElementById('mobile-menu')) {
        const menuOverlay = document.createElement('div');
        menuOverlay.id = 'mobile-menu';
        menuOverlay.className = 'fixed inset-0 bg-dark-navy/95 z-[100] flex flex-col items-center justify-center gap-8 transition-all duration-300 opacity-0 pointer-events-none';
        menuOverlay.innerHTML = `
            <button class="absolute top-6 right-6 text-white text-3xl" id="close-mobile-menu">
                <i class="fa-solid fa-xmark"></i>
            </button>
            <a href="index.html" class="text-white text-xl font-bold uppercase tracking-widest hover:text-primary transition">Ana Sayfa</a>
            <a href="hakkimizda.html" class="text-white text-xl font-bold uppercase tracking-widest hover:text-primary transition">Hakkımızda</a>
            <a href="hizmetler.html" class="text-white text-xl font-bold uppercase tracking-widest hover:text-primary transition">Hizmetlerimiz</a>
            <a href="urunler.html" class="text-white text-xl font-bold uppercase tracking-widest hover:text-primary transition">Ürünlerimiz</a>
            <a href="iletisim.html" class="text-white text-xl font-bold uppercase tracking-widest hover:text-primary transition">İletişim</a>
            <div class="mt-8 flex gap-6 text-white text-2xl">
                <a href="#"><i class="fa-brands fa-linkedin"></i></a>
                <a href="#"><i class="fa-brands fa-instagram"></i></a>
                <a href="#"><i class="fa-brands fa-twitter"></i></a>
            </div>
        `;
        document.body.appendChild(menuOverlay);
    }

    const mobileMenu = document.getElementById('mobile-menu');
    const closeMenuBtn = document.getElementById('close-mobile-menu');

    if (mobileBtn) {
        mobileBtn.addEventListener('click', () => {
            mobileMenu.classList.remove('opacity-0', 'pointer-events-none');
            document.body.style.overflow = 'hidden';
        });
    }

    if (closeMenuBtn) {
        closeMenuBtn.addEventListener('click', () => {
            mobileMenu.classList.add('opacity-0', 'pointer-events-none');
            document.body.style.overflow = 'auto';
        });
    }

    // --- SEARCH OVERLAY LOGIC ---
    const searchBtn = document.querySelector('button i.fa-magnifying-glass')?.parentElement;
    
    if (!document.getElementById('search-overlay')) {
        const searchOverlay = document.createElement('div');
        searchOverlay.id = 'search-overlay';
        searchOverlay.className = 'fixed inset-0 bg-white/98 z-[100] flex flex-col items-center justify-center p-6 transition-all duration-300 opacity-0 pointer-events-none';
        searchOverlay.innerHTML = `
            <button class="absolute top-6 right-6 text-dark-navy text-3xl" id="close-search">
                <i class="fa-solid fa-xmark"></i>
            </button>
            <div class="w-full max-w-2xl">
                <h2 class="text-3xl font-black text-dark-navy mb-8 text-center uppercase tracking-tighter">Ne aramıştınız?</h2>
                <div class="relative">
                    <input type="text" id="search-input" placeholder="Ürün veya hizmet adı yazın..." class="w-full border-b-4 border-dark-navy py-4 text-2xl font-bold focus:outline-none focus:border-primary transition-colors bg-transparent">
                    <button id="execute-search" class="absolute right-0 bottom-4 text-2xl text-dark-navy hover:text-primary transition-colors">
                        <i class="fa-solid fa-arrow-right"></i>
                    </button>
                </div>
                <p class="mt-4 text-gray-400 text-sm font-medium">Örn: Pompa sistemleri, Klima bakımı, Endüstriyel çözümler...</p>
            </div>
        `;
        document.body.appendChild(searchOverlay);
    }

    const searchOverlay = document.getElementById('search-overlay');
    const closeSearchBtn = document.getElementById('close-search');
    const searchInput = document.getElementById('search-input');
    const executeSearchBtn = document.getElementById('execute-search');

    if (searchBtn) {
        searchBtn.addEventListener('click', () => {
            searchOverlay.classList.remove('opacity-0', 'pointer-events-none');
            searchInput.focus();
            document.body.style.overflow = 'hidden';
        });
    }

    if (closeSearchBtn) {
        closeSearchBtn.addEventListener('click', () => {
            searchOverlay.classList.add('opacity-0', 'pointer-events-none');
            document.body.style.overflow = 'auto';
        });
    }

    const performSearch = () => {
        const query = searchInput.value.trim();
        if (query) {
            window.location.href = `/urunler.html?search=\${encodeURIComponent(query)}`;
        }
    };

    if (executeSearchBtn) executeSearchBtn.addEventListener('click', performSearch);
    if (searchInput) {
        searchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') performSearch();
        });
    }

    // --- HEADER SCROLL EFFECT ---
    const header = document.querySelector('header');
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            header?.classList.add('shadow-lg');
            header?.classList.add('h-16');
            header?.classList.remove('h-20');
        } else {
            header?.classList.remove('shadow-lg');
            header?.classList.add('h-20');
            header?.classList.remove('h-16');
        }
    // --- CONTACT FORM LOGIC ---
    const contactForm = document.getElementById('contactForm') || document.querySelector('form[action="#"]');
    if (contactForm) {
        if (!contactForm.id) contactForm.id = 'contactForm';
        contactForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            
            const formData = {
                ad_soyad: contactForm.querySelector('input[type="text"], [name="name"]')?.value || '',
                email: contactForm.querySelector('input[type="email"], [name="email"]')?.value || '',
                konu: contactForm.querySelector('input[placeholder*="Konu"], [name="subject"]')?.value || 'Genel İletişim',
                mesaj: contactForm.querySelector('textarea, [name="message"]')?.value || ''
            };

            try {
                const response = await fetch('/api/mesajlar', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(formData)
                });

                if (response.ok) {
                    alert('Mesajınız başarıyla gönderildi!');
                    contactForm.reset();
                } else {
                    const errorData = await response.json();
                    alert('Hata: ' + (errorData.error || 'Mesaj gönderilemedi.'));
                }
            } catch (error) {
                console.error('Contact Form Error:', error);
                alert('Sistem hatası. Lütfen daha sonra tekrar deneyiniz.');
            }
        });
    }

    // --- ADMIN LOGIN LOGIC ---
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            const loginData = {
                kullanici_adi: loginForm.querySelector('#username')?.value || '',
                sifre: loginForm.querySelector('#password')?.value || ''
            };

            try {
                const response = await fetch('/api/admin/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(loginData)
                });

                const result = await response.json();

                if (response.ok) {
                    localStorage.setItem('adminToken', result.token);
                    window.location.href = 'admin-dashboard.html';
                } else {
                    alert('Kullanıcı adı veya şifre hatalı!');
                }
            } catch (error) {
                console.error('Login Error:', error);
                alert('Giriş yapılırken bir hata oluştu.');
            }
        });
    }
});

