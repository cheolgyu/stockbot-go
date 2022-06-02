# stockbot

### stockbot/src/common/doc/ 
   + 공통코드정리
      ---
      code.go
         1. 모든회사 조회
         2. 모든코드 조회

      note.go
         1. 공지정보 업데이트
            DB_PUB.DB_PUB_COLL_NOTE
      
      name.go
         디비, 컬렉션명 상수

      ---
   + 테이블
     + DB_PUB
        + DB_PUB_COLL_NOTE
        + DB_PUB_COLL_COMPANY


### stockbot/src/task/fetch/kr/company   
   + process
      ---
      1. DB 조회: 기존회사
      2. 엑셀다운: 기존회사의 갱신정보, 새회사
      3. 기존회사와 새회사 합치기
         1. 회사code일치시 새회사정보로 변경( object id 유지 ) 
         2. 회사code 불일치시 new회사로 objectid 부여 
      4. 저장: replace 와 upsert로 처리
      ---
   + 테이블
     + DB_PUB
        + DB_PUB_COLL_COMPANY

### stockbot/src/task/fetch/kr/price   
   + process
      ---
      1. 네이버차트 가격데이터 조회시 필요한것 조회기간+code
      2. 조회기간
         1. 시작기준
            1. 코드별에서 나라별 시작일자로 수정 db_public.note 필드명: kr_price_updated_date
         2. 종료기준은 오늘일자
      3. 코드조회
      4. 저장
         1. replace로 처리
      ---
   + 테이블
     + DB_DATA
        + DB_DATA_COLL_PRICE
      + DB_PUB
        + DB_PUB_COLL_NOTE
          + DB_PUB_COLL_NOTE_PRICE_UPDATED_KR: 가격정보 업데이트 일자
