#!/bin/bash
echo "Testing JWT Auth protection (without token)..."
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/admin/mesajlar
echo " - Should be 401"

echo ""
echo "Testing File Upload Extension restriction (this expects 401 if token is invalid, but let's test if it handles the fake extension when bypass isn't active)"
echo "dummy file" > fake_image.exe
curl -s -X POST -F "image=@fake_image.exe" http://localhost:8080/api/admin/upload -H "Authorization: Bearer mock"
echo ""
rm fake_image.exe
