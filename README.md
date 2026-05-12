# Ilgaz Mühendislik - Kurumsal Web Portalı

İklimlendirme ve mühendislik sektöründe öncü çözümler sunan Ilgaz Mühendislik için geliştirilmiş; modern, dinamik ve tam fonksiyonel kurumsal web sitesi ve yönetim paneli projesidir.

## 🚀 Proje Hakkında
Bu proje, bir mühendislik firmasının tüm dijital ihtiyaçlarını karşılamak üzere uçtan uca tasarlanmıştır. Müşteriler için zengin bir ürün kataloğu ve iletişim kanalları sunarken, yönetici kadrosu için tüm içeriği kontrol edebilecekleri kapsamlı bir admin paneli barındırır.

## 🛠 Teknik Altyapı
### Backend
* **Dil:** Go (Golang) 1.26+
* **Framework:** Gin Gonic (Yüksek performanslı HTTP web framework)
* **Veritabanı:** MongoDB Atlas (NoSQL) - **MongoDB Go Driver v2** entegrasyonu.
* **Görsel Yönetimi:** Yerel dosya yükleme sistemi (`multipart/form-data`)

### Frontend
* **UI/UX:** HTML5, modern CSS3 (Custom Variables & Modern Layouts).
* **Dinamizm:** Vanilla JavaScript (ES6+) ile asenkron API entegrasyonu (Fetch API).
* **İkonlar:** FontAwesome 6.4.0
* **Fontlar:** Inter (Google Fonts)

## ✨ Öne Çıkan Özellikler
*   **Dinamik Ürün Kataloğu:** Kategori bazlı filtreleme, anlık arama, fiyat sıralaması ve sayfalama sistemi.
*   **Ürün Detay Entegrasyonu:** Her ürün için otomatik "Teklif İste" yönlendirmesi.
*   **Gelişmiş Admin Paneli:**
    *   **Dashboard:** Gerçek zamanlı istatistikler ve son gelen mesajlar.
    *   **Ürün Yönetimi:** Drag-and-drop sıralama, görsel yükleme, CRUD işlemleri.
    *   **Proje Takibi:** Mühendislik projelerinin ilerleme durumunu ve detaylarını yönetme.
    *   **Mesaj Merkezi:** Gelen müşteri taleplerini okuma, silme ve yönetme.
*   **Mobil Uyumluluk:** Tüm cihazlarda sorunsuz çalışan tam ekran mobil menü ve responsive tasarım.
*   **Güvenlik:** 
    *   **JWT Auth:** Stateless kimlik doğrulama ile güvenli API erişimi.
    *   **Protected Routes:** Sadece yetkili kullanıcıların erişebileceği endpoint'ler.
    *   **Secure Uploads:** Dosya tipi ve boyutu doğrulamalı güvenli görsel yükleme sistemi.

## 🧪 Test ve Güvenlik
*   **Automated Tests:** `internal/handlers/auth_test.go` üzerinden otomatik auth testleri.
*   **Security Scripts:** `scripts/security_test.sh` ile API güvenlik taramaları.

## 📂 Klasör Yapısı
*   `cmd/api/`: Uygulamanın giriş noktası ve sunucu yapılandırması.
*   `internal/`: Backend mantığı, modeller, veritabanı bağlantısı ve API işleyicileri.
*   `public/`: Tüm frontend dosyaları (HTML, CSS, JS).
*   `assets/uploads/`: Kullanıcılar tarafından yüklenen ürün görselleri.

## 🛠 Kurulum ve Çalıştırma
1. Bağımlılıkları yükleyin: `go mod tidy`
2. Uygulamayı başlatın: `go run ./cmd/api/main.go`
3. Tarayıcıda açın: `http://localhost:8080`

---
&copy; 2026 Ilgaz Mühendislik A.Ş. - Tüm hakları saklıdır.
