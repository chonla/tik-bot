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

**Syntax**: [สวัสดี|hi|hello]

```
<me> สวัสดี
<tik> ชื่ออะไรเหรอ
<me> อู
<tik> สวัสดีจ้ะอู
<me> สวัสดี
<tik> สวัสดีจ้ะอู
```

## การลงชื่อเข้าทำงาน

**Syntax**: [checkin|check-in|เข้าทำงาน|ลงชื่อ] [<ชื่อสถานที่>]

**หมายเหตุ**: ถ้าไม่ได้ระบุที่ทำงานและยังไม่เคยลงชื่อเข้าทำงาน Tik-bot จะถามว่าทำงานที่ไหน
แต่ถ้าเคยลงชื่อเข้าทำงานแล้วและมีที่ทำงานที่เดียว Tik-bot จะ checkin ที่ทำงานนั้นให้ทันที
ถ้าต้องการระบุที่ทำงานที่อื่น ให้ระบุชื่อสถานที่ด้วย
ถ้าไม่ระบุและถ้ามีที่ทำงานหลายที่ Tik-bot จะให้เลือกว่าทำงานที่ไหน

```
<me> checkin
<tik> วันนี้เข้าทำงานที่ไหนเหรอ
<me> บ้าน
<tik> ลงชื่อเข้าทำงานที่ บ้าน เรียบร้อยจ้ะ
<me> checkin
<tik> ลงชื่อเข้าทำงานที่ บ้าน เรียบร้อยจ้ะ
<me> checkin โรงเรียน
<tik> ลงชื่อเข้าทำงานที่ โรงเรียน เรียบร้อยจ้ะ
```

## TODO: ดูสรุปจำนวนวันทำงาน

**Syntax**: [สรุป|sum]

```
<me> สรุป
<tik> สรุปรอบเงินเดือน 26 เม.ย. 2561 - 25 พ.ค. 2561
โรงเรียน : 3
บ้าน : 2
โรงแรม : 1.5
```