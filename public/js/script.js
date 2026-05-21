document.addEventListener('DOMContentLoaded', () => {
    // --- CENTRAL SCROLL MANAGEMENT ---
    window.updateBodyScroll = () => {
        const catalogModal = document.getElementById('catalogModal');
        const kvkkModal = document.getElementById('kvkk-modal');
        const searchOverlay = document.getElementById('search-overlay');
        const mobileMenu = document.getElementById('mobile-menu');

        const catalogOpen = catalogModal && !catalogModal.classList.contains('hidden');
        const kvkkOpen = kvkkModal && !kvkkModal.classList.contains('hidden');
        const searchOpen = searchOverlay && !searchOverlay.classList.contains('opacity-0');
        const mobileOpen = mobileMenu && !mobileMenu.classList.contains('opacity-0');

        if (catalogOpen || kvkkOpen || searchOpen || mobileOpen) {
            document.body.style.overflow = 'hidden';
        } else {
            document.body.style.overflow = '';
        }
    };

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
            window.updateBodyScroll();
        });
    }

    if (closeMenuBtn) {
        closeMenuBtn.addEventListener('click', () => {
            mobileMenu.classList.add('opacity-0', 'pointer-events-none');
            window.updateBodyScroll();
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
            window.updateBodyScroll();
        });
    }

    if (closeSearchBtn) {
        closeSearchBtn.addEventListener('click', () => {
            searchOverlay.classList.add('opacity-0', 'pointer-events-none');
            window.updateBodyScroll();
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
    });
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

    // --- SMOOTH SCROLL FOR ANCHOR LINKS ---
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            const href = this.getAttribute('href');
            if (href.startsWith('#') && href.length > 1) {
                const target = document.querySelector(href);
                if (target) {
                    e.preventDefault();
                    const header = document.querySelector('header');
                    // Offset exactly equal to the sticky header height (80px)
                    const headerHeight = header ? 80 : 0;
                    const targetPosition = target.getBoundingClientRect().top + window.pageYOffset - headerHeight;
                    window.scrollTo({
                        top: targetPosition,
                        behavior: 'smooth'
                    });
                }
            }
        });
    });

    // --- KVKK POPUP LOGIC ---
    if (!document.getElementById('kvkk-modal')) {
        const kvkkModal = document.createElement('div');
        kvkkModal.id = 'kvkk-modal';
        kvkkModal.className = 'fixed inset-0 bg-black/60 backdrop-blur-sm z-[150] hidden items-center justify-center p-4 transition-all duration-300 opacity-0';
        kvkkModal.innerHTML = `
            <div class="bg-white rounded-2xl w-full max-w-2xl shadow-2xl overflow-hidden transform scale-95 transition-all duration-300">
                <div class="bg-primary px-6 py-5 flex justify-between items-center text-white">
                    <div class="flex items-center gap-3">
                        <i class="fa-solid fa-shield-halved text-2xl"></i>
                        <h3 class="font-bold text-lg leading-tight">KVKK Aydınlatma Metni</h3>
                    </div>
                    <button type="button" id="close-kvkk-btn" class="text-white/80 hover:text-white text-xl cursor-pointer bg-transparent border-0 outline-none"><i class="fa-solid fa-xmark"></i></button>
                </div>
                <div class="p-6 max-h-[60vh] overflow-y-auto space-y-4 text-gray-700 text-sm leading-relaxed" id="kvkk-content">
                    <p class="font-medium text-gray-900"><strong>Ilgaz Mühendislik A.Ş.</strong> ("Şirket") olarak, 6698 sayılı Kişisel Verilerin Korunması Kanunu ("Kanun") uyarınca, veri sorumlusu sıfatıyla, kişisel verilerinizin güvenliğine önem veriyor ve bu verilerin işlenmesinde hukuka ve dürüstlük kurallarına uygun hareket ediyoruz.</p>
                    
                    <h4 class="font-bold text-gray-900 border-b pb-1 mt-4">1. Kişisel Verilerin Hangi Amaçla İşleneceği</h4>
                    <p>Web sitemiz üzerinden paylaştığınız kişisel verileriniz (ad, soyad, e-posta, telefon numarası gibi iletişim verileri), ürün kataloglarının iletilmesi, taleplerinizin yanıtlanması, müşteri değerlendirmelerinizin yayınlanması, hizmet kalitemizin artırılması ve sizinle gerekli durumlarda iletişime geçilmesi amaçlarıyla sınırlı olarak işlenmektedir.</p>

                    <h4 class="font-bold text-gray-900 border-b pb-1 mt-4">2. İşlenen Kişisel Verilerin Aktarılması</h4>
                    <p>Kişisel verileriniz, yasal yükümlülüklerimizin yerine getirilmesi amacıyla yetkili kamu kurum ve kuruluşları ile yasal sınırlar dahilinde paylaşılabilecektir. Verileriniz ticari veya pazarlama amaçlarıyla üçüncü şahıslara veya şirket dışı yapılara aktarılmamakta ve satılmamaktadır.</p>

                    <h4 class="font-bold text-gray-900 border-b pb-1 mt-4">3. Kişisel Veri Toplamanın Yöntemi ve Hukuki Sebebi</h4>
                    <p>Kişisel verileriniz, tamamen veya kısmen otomatik yollarla (web sitemizdeki formlar aracılığıyla) toplanmaktadır. Bu veriler, Kanun’un 5. maddesinde belirtilen "sözleşmenin kurulması veya ifasıyla doğrudan doğruya ilgili olması kaydıyla, sözleşmenin taraflarına ait kişisel verilerin işlenmesinin gerekli olması" ve "veri sorumlusunun hukuki yükümlülüğünü yerine getirebilmesi için zorunlu olması" hukuki sebeplerine dayalı olarak işlenmektedir.</p>

                    <h4 class="font-bold text-gray-900 border-b pb-1 mt-4">4. Veri Sahibinin Hakları</h4>
                    <p>Kanun’un 11. maddesi kapsamında; kişisel verilerinizin işlenip işlenmediğini öğrenme, işlenmişse buna ilişkin bilgi talep etme, işlenme amacını ve amacına uygun kullanılıp kullanılmadığını öğrenme, verilerinizin eksik veya yanlış işlenmiş olması hâlinde düzeltilmesini isteme ve silinmesini talep etme haklarına sahipsiniz. Bu haklarınızı kullanmak için firmamızla doğrudan iletişim kurabilirsiniz.</p>
                </div>
                <div class="p-4 bg-gray-50 border-t flex justify-end">
                    <button type="button" id="close-kvkk-bottom" class="bg-gray-800 hover:bg-gray-700 text-white font-bold px-6 py-2.5 rounded-lg text-sm transition-colors cursor-pointer border-0">Kapat</button>
                </div>
            </div>
        `;
        document.body.appendChild(kvkkModal);

        const modalContent = kvkkModal.querySelector('.transform');

        const showKvkkModal = (e) => {
            if (e) e.preventDefault();
            kvkkModal.classList.remove('hidden');
            kvkkModal.classList.add('flex');
            setTimeout(() => {
                kvkkModal.classList.remove('opacity-0');
                modalContent.classList.remove('scale-95');
                modalContent.classList.add('scale-100');
            }, 10);
            window.updateBodyScroll();
        };

        const hideKvkkModal = () => {
            kvkkModal.classList.add('opacity-0');
            modalContent.classList.remove('scale-100');
            modalContent.classList.add('scale-95');
            setTimeout(() => {
                kvkkModal.classList.add('hidden');
                kvkkModal.classList.remove('flex');
                window.updateBodyScroll();
            }, 300);
        };

        document.getElementById('close-kvkk-btn').addEventListener('click', hideKvkkModal);
        document.getElementById('close-kvkk-bottom').addEventListener('click', hideKvkkModal);
        kvkkModal.addEventListener('click', (e) => {
            if (e.target === kvkkModal) hideKvkkModal();
        });

        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && !kvkkModal.classList.contains('hidden')) {
                hideKvkkModal();
            }
        });

        // Global link interceptor
        document.addEventListener('click', (e) => {
            const target = e.target.closest('a');
            if (target && target.textContent.includes('KVKK Aydınlatma Metni')) {
                showKvkkModal(e);
            }
        });
    }
});

