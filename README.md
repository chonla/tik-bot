# Tik-Bot

Tik-Bot คือ Slack Bot ทำหน้าที่ลงเวลาทำงานสำหรับ ODDS team Tik-Bot ใช้ Firebase สำหรับเก็บข้อมูล ใช้การ authen ด้วย service account

# วิธีการทดสอบ

```
GOOGLE_APPLICATION_CREDENTIALS=<service-account-file> go run main.go
```

# ไฟล์ oddsy.json

```
{
    "slack-token": "<string>",
    "ignore-bot-message": <boolean>,
    "debug": <boolean>,
    "gcp-token": "<string>",
    "firebase-project-id": "<string>"
}
```

* slack-token: token ที่ได้จาก Slack
* ignore-bot-message: (default: true) บอกให้ bot ไม่สนใจ message ที่มาจาก bot ด้วยกัน รวมถึงตัวเอง
* debug: (default: false) แสดง debug message ใน log
* gcp-token: Google Cloud Platform token (ตอนนี้ยังใช้ไม่ได้ ให้ใช้ service account แทน)
* firebase-project-id: ID ของ firebase project

# การพูดคุย

Tik-bot ถูกตั้งให้พูดคุยผ่าน direct message อย่างเดียว

## การลงทะเบียนชื่อให้ Tik-bot รู้จัก

ให้ทัก Tik-bot ไปก่อน ถ้า tik-bot ยังไม่รู้จัก Tik-bot จะถามชื่อกลับมา

**Alias**: สวัสดี, hi, hello

```
<me> สวัสดี
<tik> ชื่ออะไรเหรอ
<me> อู
<tik> สวัสดีจ้ะอู
<me> สวัสดี
<tik> สวัสดีจ้ะอู
```

## TODO: การลงชื่อเข้าทำงาน

ให้บอก Tik-bot ว่า checkin

**Alias**: checkin, check-in, ทำงาน, มาแล้ว

```
<me> checkin
<tik> วันนี้เข้าทำงานที่ไหนเหรอ
<me> บ้าน
<tik> ลงชื่อเข้าทำงานที่ บ้าน เรียบร้อยจ้ะ
```