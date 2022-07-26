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


   ### stockbot/src/task/fetch/
      + process
         ---
            1. 거래소 저장
            2. company 실행(나라)
            3. price 실행(나라) 
         ---

   ### stockbot/src/task/fetch/company   
      + process
         ---
         input(나라)
            1. DB 조회: 기존회사(나라)
            2. 크롤링시작 
               1. kr: krx.data   엑셀파일
               2. us: nasdaq.com json
            3. 파일로 저장
            4. 파일을 []stuct로 변환
            5. 기존회사와 새회사 합치기
            6. 저장: replace 와 upsert로 처리
         ---
      + 테이블
      + DB_PUB
         + DB_PUB_COLL_COMPANY
         + DB_PUB_COLL_NOTE

   ### stockbot/src/task/fetch/price   
      + process
         ---
         1. init()
            price_index 체크
         2. 종목별 조회시작 종료 기간 만들기 
            1. 시작: 코드별 price 컬렉션의 max dt값
            2. 종료: 오늘
         3. input(나라)
            1. 나라의 종목코드,마켓코드 조회
            2. loop 코드
               1. 나라의 종목코드의 기간의 가격목록 파일저장
               2. 파일의 []stuct로 변환
               3. 저장
           ---
        + 테이블
          + open,close,high,low의 data type는 decimal
        + DB_DATA
           + DB_DATA_COLL_PRICE
        + DB_PUB
           + DB_PUB_COLL_NOTE
            


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
      + 나중에
      + 연도의 월별 거래량의 퍼센트를 구한후 전체연도에서 월별 거래량의 퍼센트의 변화를 확인하자
      + goal
      + 몰리는 거래량 찾기
      + 용어
      + 기간구분
         + 주,월,분기
      + 정책
      + agg_vol_sum의 계산범위
            + 코드별 마지막 가격데이터의 년도
      + desc
         종목의 거래량 그래프를 연도별로 나누어 
         주별로 거래량을 합하여 제일 큰 주를 찾고
         월별로 거래량을 합하여 제일 큰 월을 찾고
         분기별로 거래량을 합하여 제일 큰 분기을 찾는다.

         전체연도에서 찾은 값들의 비중을 확인한다.
         

         기간구분 주, 월, 분기
      + process
         ---
         agg_vol_sum 계산범위:  코드별 마지막 가격데이터의 년도
         1. agg_vol_sum에서 테이블의 코드별 마지막 년도를 구한다 없으면0을 반환한다.
            1.1 agg_vol_sum에서 테이블의 코드별 마지막 년도를 구한다 없으면0을 반환한다.
            1.2 agg_vol_sum의 year가 0이면 code의 price 데이터중 가장 작은 year은 구한다.
            1.3 특정연도의 가격테이터를 조회한다.
         2. 코드별 연도별 가격데이터에서 주별 거래량의합, 월별 거래량의 합, 분기별 거래량의합을 구한다.
         3. 코드별 해당연도의 agg_vol_sum을 upsert 한다.
         4. 코드별 전체연도의 agg_vol_sum을 조회한다.
         5. 코드별 코드의 전체연도의 agg_vol_sum데이터로 기간별 표준편차를 구한다.
            5.1 표준편차를 구한다.
               5.1.1 관찰값들의 평균을 구한다. (편차값을 구하기 위해서)
               5.1.2 편차를 구한다. (편차: 관측값에서 평균을 뺀것)
               5.1.3 표준편차를 구한다.
            5.2 빈도를 분석한다.
         6. 코드별 코드의 전체연도의 기간별 표준편차를 저장한다.

         
         ---      
      + 테이블
      + DB_PUB
         + DB_PUB_COLL_COMPANY
         + DB_DATA
         + DB_DATA_COLL_PRICE
         + DB_DATA_COLL_AGG_VOL_SUM
         + DB_DATA_COLL_AGG_VOL
