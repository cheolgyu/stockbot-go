stockbot/src/task/kr
---
엑셀 다운로드 -> db insert
1. 한국거래소에서 종목코드 엑셀 다운로드
   1. 종목코드상세.엑셀 다운로드
   2. 엑셀파일을 구조체로 변환
   3. 종코코드상태.엑셀 다운로드
   4. 엑셀파일을 구조체로 변환
   5. 저장, 저장
   6. 기존 회사의 상태가 갱신되면
      1. 과거 목록을 조회한다 
      2. 과거 목록에 엑셀내용을 집어넣는다.
         1. 과거목록을 key-value로 변환한다
         2. 엑셀에서 내용을  찾으면 key값으로 액세스하여 내용을 수정한다.
         3. 상세 후 상태
      3. update한다.
2. 네이버증권에서 종목코드별 가격목록 다운로드
3. 저장
---
# stockbot
