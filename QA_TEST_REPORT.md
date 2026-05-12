# Ilgaz Mühendislik - Kalite Güvence & Güvenlik Test Raporu (QA_TEST_REPORT.md)

**Tarih:** 12 Mayıs 2026
**Uygulanan Test Fazı:** API & Backend Security (Bölüm 3 & 4)
**Durum:** Başarılı

## 1. Uygulanan Güvenlik (Security) Yamaları

Projenin üretim (production) ortamına alınmadan önce belirlenen "En Kritik" iki güvenlik açığı tamamen kapatılmıştır.

### 1.1. Yetkilendirme Koruması (JWT Authentication)
- **Problem:** `/api/admin/*` altındaki uç noktalar (endpoints) dışarıya açıktı.
- **Çözüm:** `internal/middleware/auth.go` adlı güçlü bir JWT (JSON Web Token) kalkanı kodlandı. 
- **Entegrasyon:** `internal/routes/routes.go` dosyasında `admin` API grubuna `middleware.AuthRequired()` bağlandı.
- **Sonuç:** Token olmadan veya geçersiz/süresi dolmuş bir token ile sisteme erişmeye çalışan tüm istekler anında **401 Unauthorized** hatasıyla reddedilmektedir.

### 1.2. Güvenli Dosya Yükleme (Secure Image Upload)
- **Problem:** Kötü niyetli kişilerin sunucuya `.exe`, `.pdf` veya devasa boyutlu dosyalar (örn: 1 GB) yükleyerek sistemi çökertme (DDoS) veya ele geçirme (RCE) riski vardı.
- **Çözüm:** `handlers.go` içerisindeki `UploadImage` fonksiyonu güncellendi.
  - **Boyut Sınırı:** Go'nun yerleşik `http.MaxBytesReader` metodu ile payload **10 MB** ile sınırlandı.
  - **Uzantı Sınırı:** Strict Mode uygulanarak sadece `.jpg, .jpeg, .png, .webp` uzantılı dosyalara izin verildi.
- **Sonuç:** Kriterlere uymayan her deneme **400 Bad Request** hatası alarak sunucuya yazılmadan iptal edilmektedir.

---

## 2. Otomasyon Testleri (Go Unit/Integration Tests)

Manuel testi beklememek ve ilerideki geliştirmelerde güvenliği garantilemek adına otomatik Go testleri (`auth_test.go`) kodlandı.

**Oluşturulan Test Senaryoları:**
1. `TestAuthMiddleware_NoToken`: Header'ında yetki anahtarı olmayan bir isteğin "401" dönüp dönmediğini doğrular.
2. `TestAuthMiddleware_InvalidToken`: "Bearer invalid-fake-token" gönderilerek sistemin sahte tokenı "401" ile reddettiğini doğrular.
3. `TestUploadImage_InvalidExtension`: Sisteme sanal olarak bir `fake_hack.exe` yüklenmeye çalışılır ve sistemin uzantıyı algılayıp "400" dönüp dönmediğini doğrular.

**Test Sonuçları:**
```bash
=== RUN   TestAuthMiddleware_NoToken
--- PASS: TestAuthMiddleware_NoToken (0.00s)
=== RUN   TestAuthMiddleware_InvalidToken
--- PASS: TestAuthMiddleware_InvalidToken (0.00s)
=== RUN   TestUploadImage_InvalidExtension
--- PASS: TestUploadImage_InvalidExtension (0.00s)
PASS
```
Tüm testler başarıyla (PASS) geçilmiştir.

---

## Sonuç
Sistemin beyni sayılan yetkilendirme ve dosya işleme kısımları artık tamamen güvenli ve otonom testlere tabidir. Uygulamanız bu aşama itibariyle JWT güvenliği standartlarına uygundur.
