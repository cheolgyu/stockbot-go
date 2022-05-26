# stockbot

stockbot/src/task/kr/company
---
1. DB 조회: 기존 회사
2. 엑셀다운: 기존 회사 갱신정보, new 회사
3. 기존회사와 새회사 합치기
   1. 회사code일치시 새회사정보로 변경( object id 유지 ) 
   2. 회사code 불일치시 new회사로 objectid 부여 
4. 저장: replace 와 upsert로 처리
---

stockbot/src/task/kr/price
---
1. 네이버차트 가격데이터 조회시 필요한것 조회기간+code
2. 조회시간
   1. 시작기준
      1. 코드별에서 나라별 시작일자로 수정 db_public.note 필드명: kr_price_updated_date
   2. 종료기준은 오늘일자
3. 코드조회
4. 저장
   1. replace로 처리
---

