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
      0. init.go : 마켓정보 저장
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
            1. 코드별 price 컬렉션의 max dt값
         2. 종료기준은 오늘일자
      3. 코드조회
      4. 저장
         1. replace로 처리
      5. 장열림 저장
      ---
   + 테이블
     + DB_DATA
        + DB_DATA_COLL_PRICE
      + DB_PUB
        + DB_PUB_COLL_NOTE
          + DB_PUB_COLL_NOTE_PRICE_UPDATED_KR: 가격정보 업데이트 일자


### stockbot/src/task/asmb/line/bound 
   + goal
     + 코드의 종시저고가별 누적 몇퍼센트인지 찾기
   + process
      ---
      1. 코드목록조회
      2. 코드별 종시저고가별 마지막 바운드점의 일자조회
      3. 마지막 바운드점의 일자 이후의 가격목록 조회
      4. 마지막 바운드점의 일자 이후부터 가격목록으로 바운스점 찾기
      5. 바운스점 저장 
     
      ---
      ---
      1. 코드목록조회
      2. 코드별 종시저고가별 마지막 바운드점의 일자조회
      3. 마지막 바운드점의 일자 이후의 가격목록 조회
      4. 마지막 바운드점의 일자 이후부터 가격목록으로 바운스점 찾기
      5. 바운스점 저장 
     
      ---      
   + 테이블
     + DB_PUB
        + DB_PUB_COLL_COMPANY
      + DB_DATA
        + DB_DATA_COLL_BOUND_POINT


### stockbot/src/task/asmb/line/ymxb 
   + goal
     + 내일 가격 찾기
   + desc  
     + 직선의 방정식을 이용하여 p1마지막 반등과 p2 현재가격을 직선으로 이어 p3인 내일 가격 찾기
   + process
      ---
      1. 코드목록조회
      2. 가격분류별 마지막 반등 POINT인 P1 조회
      3. 가격분류별 마지막 가격 POINT인 P2 조회
      4. P1과 P2를 이용해 기울기인 M과 B를 구한후 Y값에 해당하는 호가를 호가테이블에서 가져오기
      5. 저장하기
      ---      
   + 테이블
     + DB_PUB
        + DB_PUB_COLL_COMPANY
      + DB_DATA
        + DB_DATA_COLL_PRICE
        + DB_DATA_COLL_BOUND_POINT
        + DB_DATA_COLL_YMXB
        + DB_DATA_COLL_YMXB_QUOTE_UNIT
   + 기능 추가시
     + ymxb_type1 은 p1,p2가 마지막 반등, 마지막 가격이라면
     + ymxb_type2 은 p1,p2가 저가기존 뒤에서2번째 반등, 1번째 반등
     + ymxb_type3 은 p1,p2가 고가기존 뒤에서2번째 반등, 1번째 반등

### stockbot/src/task/asmb/line/ymxb_hist 
   + goal
     + 종목의 기울기와 y절편의 변동내역 확인하기
   + desc  
     + 종목의 기울기와 y절편의 변동내역 확인하기

### stockbot/src/task/asmb/agg/vol 
   + goal
     + 몰리는 거래량 찾기
   + desc
      종목의 거래량 그래프를 연도별로 나누어 
      주별로 거래량을 합하여 제일 큰 주를 찾고
      월별로 거래량을 합하여 제일 큰 월을 찾고
      분기별로 거래량을 합하여 제일 큰 분기을 찾는다.

      전체연도에서 찾은 값들의 비중을 확인한다.
      

      기간구분 주, 월, 분기
   + process
      ---
       1. 계산1) 새로운 데이터의 거래량를 새로운 데이터의 일자에 해당하는 누적테이블에 저장한다.
          1. 새로운데이터목록에서 데이터의 연도를 뽑아내서 다음 계산을 위한 범위를 정한다.
          2. 새로운 데이터를 누적테이블에 저장한다.
             1. 기존 데이터를 조회한다 
             2. 더한다
             3. 업데이트한다.
                1. 한 종목의 신규 가격목록이 1000개 라면? 
                2. 1000개를 반복문 돌려서 GOLANG으로 계산한다.
                3. 기존데이터를 조회한다 더한다 업데이트한다. 기간종류에 따라 3줄이라면 3번해야 한다.
                   1. 그러므로 기간종류를 칼럼으로 두고 처리한다면?
                      1. 칼럼으로 둔다면 가격테이블을 쿼리로 집계하여 거래량의 합을 구하는건?
                      2. 
       2. 계산2) 계산1)에서 구한 연도들을 반복하여 연도별 기간별 최소,최대,평균의 기간별 값을 새로 구한다.
          1. 2022년도 데이터를 추가하였는데 1999년를 다시 계산할 필요는 없지만 초기실행시 1999~2022 까지 온경우를 대비하여 input을 배열로 처리. 
          2. 계산 범위는 계산1)에서 정한 범위로 정한다.
       3. 계산3) 누적테이블에서 전체연도의 최소,최대,평균,백분율을 다시 구한다.
                  
        
      ---      
   + 테이블
     + DB_PUB
        + DB_PUB_COLL_COMPANY
      + DB_DATA
        + DB_DATA_COLL_AGG_VOL_SUM
        + DB_DATA_COLL_AGG_VOL
