#!/bin/bash

API_URL="http://localhost/api/v1"

echo "🔑 Login..."
TOKEN=$(curl -s -X POST ${API_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

echo "✅ Token: ${TOKEN:0:30}..."

# Get first member ID
MEMBER_ID=$(curl -s -X GET ${API_URL}/anggota \
  -H "Authorization: Bearer $TOKEN" | jq -r '.data[0].id')

echo "👤 Testing with member: $MEMBER_ID"
echo ""
echo "🔄 Creating 10 concurrent deposits on the SAME date..."
echo "This should test the race condition fix"
echo ""

# Run 10 concurrent requests with the SAME timestamp
for i in {1..10}; do
  (
    RESPONSE=$(curl -s -X POST ${API_URL}/simpanan \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"idAnggota\": \"$MEMBER_ID\",
        \"tipeSimpanan\": \"sukarela\",
        \"tanggalTransaksi\": \"2025-02-01T10:00:00Z\",
        \"jumlahSetoran\": 10000,
        \"keterangan\": \"Concurrent test $i\"
      }")
    
    STATUS=$(echo $RESPONSE | jq -r '.status // "error"')
    if [ "$STATUS" = "success" ]; then
      REF=$(echo $RESPONSE | jq -r '.data.nomorReferensi')
      echo "  ✅ $i: $REF"
    else
      ERROR=$(echo $RESPONSE | jq -r '.error // .message // "unknown"')
      echo "  ❌ $i: $ERROR"
    fi
  ) &
done

wait
echo ""
echo "📊 Checking journal numbers created..."
