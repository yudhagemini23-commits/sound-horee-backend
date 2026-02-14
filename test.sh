#!/bin/bash

# ==========================================
# KONFIGURASI DEFAULT & WARNA
# ==========================================
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api/v1"

# Default Values (Jika tidak ada parameter)
TARGET_UID="user_$(date +%s)"
TARGET_EMAIL="default@soundhoree.com"
TARGET_STORE="Toko Default AI"
SCENARIO="normal"

# ==========================================
# FUNGSI BANTUAN (HELPER)
# ==========================================
usage() {
    echo -e "${YELLOW}Usage: $0 [-u USER_ID] [-e EMAIL] [-n STORE_NAME] [-s SCENARIO]${NC}"
    echo ""
    echo "Options:"
    echo "  -u   Set Custom User ID (ex: yudha_pro_01)"
    echo "  -e   Set Email (ex: owner@toko.com)"
    echo "  -n   Set Nama Toko (ex: 'Toko Kelontong')"
    echo "  -s   Pilih Skenario Transaksi: [normal | rich | locked | empty]"
    echo ""
    exit 1
}

# ==========================================
# PARSING ARGUMEN (Ganti-ganti Akun)
# ==========================================
while getopts "u:e:n:s:h" opt; do
  case $opt in
    u) TARGET_UID="$OPTARG" ;;
    e) TARGET_EMAIL="$OPTARG" ;;
    n) TARGET_STORE="$OPTARG" ;;
    s) SCENARIO="$OPTARG" ;;
    h) usage ;;
    *) usage ;;
  esac
done

echo -e "${BLUE}==================================================${NC}"
echo -e "${BLUE}   🤖 SOUND HOREE BOT SIMULATOR v2.0${NC}"
echo -e "   User ID  : ${YELLOW}$TARGET_UID${NC}"
echo -e "   Email    : ${YELLOW}$TARGET_EMAIL${NC}"
echo -e "   Toko     : ${YELLOW}$TARGET_STORE${NC}"
echo -e "   Skenario : ${YELLOW}$SCENARIO${NC}"
echo -e "${BLUE}==================================================${NC}"
echo ""

# ==========================================
# STEP 1: LOGIN / REGISTER
# ==========================================
echo -e "${GREEN}[STEP 1] Login ke System...${NC}"

LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
-H "Content-Type: application/json" \
-d "{
    \"uid\": \"$TARGET_UID\",
    \"email\": \"$TARGET_EMAIL\",
    \"store_name\": \"$TARGET_STORE\",
    \"phone_number\": \"081299990000\"
}")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Gagal Login! Cek server Golang.${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Login Sukses! Token aman.${NC}"
fi
echo "--------------------------------------------------"

# ==========================================
# STEP 2: UPDATE PROFIL
# ==========================================
echo -e "${GREEN}[STEP 2] Sync Profil Toko...${NC}"

curl -s -X POST "$BASE_URL/profile/sync" \
-H "Authorization: Bearer $TOKEN" \
-H "Content-Type: application/json" \
-d "{
    \"uid\": \"$TARGET_UID\",
    \"email\": \"$TARGET_EMAIL\",
    \"store_name\": \"$TARGET_STORE\",
    \"phone_number\": \"081299998888\",
    \"category\": \"Retail Configurable\"
}" > /dev/null

echo -e "${GREEN}✅ Profil '$TARGET_STORE' berhasil disinkronkan.${NC}"
echo "--------------------------------------------------"

# ==========================================
# STEP 3: DATA TRANSAKSI (Sesuai Skenario)
# ==========================================
echo -e "${GREEN}[STEP 3] Mengirim Transaksi (Skenario: $SCENARIO)...${NC}"

NOW=$(date +%s)000
YESTERDAY=$(date -v-1d +%s)000

PAYLOAD=""

if [ "$SCENARIO" == "normal" ]; then
    PAYLOAD="[
        {\"user_id\": \"$TARGET_UID\", \"source_app\": \"id.dana\", \"amount\": 50000, \"raw_message\": \"DANA Masuk 50rb\", \"timestamp\": $NOW, \"is_trial_limited\": false},
        {\"user_id\": \"$TARGET_UID\", \"source_app\": \"com.bca\", \"amount\": 150000, \"raw_message\": \"BCA Masuk 150rb\", \"timestamp\": $NOW, \"is_trial_limited\": false},
        {\"user_id\": \"$TARGET_UID\", \"source_app\": \"com.gojek\", \"amount\": 10000, \"raw_message\": \"Gopay 10rb (Kemarin)\", \"timestamp\": $YESTERDAY, \"is_trial_limited\": false}
    ]"
elif [ "$SCENARIO" == "rich" ]; then
    PAYLOAD="[
        {\"user_id\": \"$TARGET_UID\", \"source_app\": \"com.bca\", \"amount\": 50000000, \"raw_message\": \"Transfer Besar 50jt\", \"timestamp\": $NOW, \"is_trial_limited\": false},
        {\"user_id\": \"$TARGET_UID\", \"source_app\": \"mandiri.online\", \"amount\": 25000000, \"raw_message\": \"Livin 25jt\", \"timestamp\": $NOW, \"is_trial_limited\": false}
    ]"
elif [ "$SCENARIO" == "locked" ]; then
    PAYLOAD="[
        {\"user_id\": \"$TARGET_UID\", \"source_app\": \"id.dana\", \"amount\": 50000, \"raw_message\": \"DANA 50rb (LOCKED)\", \"timestamp\": $NOW, \"is_trial_limited\": true},
        {\"user_id\": \"$TARGET_UID\", \"source_app\": \"com.bca\", \"amount\": 100000, \"raw_message\": \"BCA 100rb (LOCKED)\", \"timestamp\": $NOW, \"is_trial_limited\": true}
    ]"
elif [ "$SCENARIO" == "empty" ]; then
    PAYLOAD="[]"
fi

# Kirim Request
if [ "$SCENARIO" != "empty" ]; then
    curl -s -X POST "$BASE_URL/transactions/sync" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD" > /dev/null
    echo -e "${GREEN}✅ Data transaksi terkirim!${NC}"
else
    echo -e "${YELLOW}⚠️ Skenario Empty: Tidak ada transaksi dikirim.${NC}"
fi

echo "--------------------------------------------------"

# ==========================================
# STEP 4: CEK LAPORAN (REPORTING)
# ==========================================
echo -e "${GREEN}[STEP 4] Cek Laporan Hari Ini...${NC}"

START_TODAY=$(date -v0H -v0M -v0S +%s)000
END_TODAY=$(date -v23H -v59M -v59S +%s)999

REPORT=$(curl -s -X GET "$BASE_URL/transactions?user_id=$TARGET_UID&start=$START_TODAY&end=$END_TODAY" \
-H "Authorization: Bearer $TOKEN")

# Tampilkan Hasil Pretty Print (Jika ada jq, kalau tidak raw text)
if command -v jq &> /dev/null; then
    echo $REPORT | jq '. | {status: .status, total_amount: .total_amount, count: .count}'
else
    echo $REPORT
fi

echo ""
echo -e "${BLUE}==================================================${NC}"
echo -e "${GREEN}✅ FINISH.${NC}"