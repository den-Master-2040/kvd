## 2026-10-16 Ч пайплайнинг: Flush только при пустом входном буфере

ѕроблема: при -P 16 рост всего в 7,8 раза. strace: write 200 002 на 200 000 команд Ч
Flush после каждого ответа.
–ешение: Flush, только если r.Buffered() == 0.
–езультат (strace, 200k SET, -P 16): write 200 002 ? 12 501, все syscalls 225 898 ? 32 724.
Ѕенчмарк: <цифры до/после>.

ƒќ 

redis-benchmark -p 6380 -t set,get -n 200000 -c 50 -q
redis-benchmark -p 6380 -t set,get -n 200000 -c 50 -P 16 -q

WARNING: Could not fetch server CONFIG
SET: 93984.96 requests per second, p50=0.263 msec                   
GET: 95374.35 requests per second, p50=0.263 msec                   

WARNING: Could not fetch server CONFIG
SET: 740740.69 requests per second, p50=0.551 msec                    
GET: 740740.69 requests per second, p50=0.543 msec      

ѕќ—Ћ≈

redis-benchmark -p 6380 -t set,get -n 200000 -c 50 -q
redis-benchmark -p 6380 -t set,get -n 200000 -c 50 -P 16 -q
redis-benchmark -p 6380 -t set,get -n 1000000 -c 50 -P 16 --threads 4 -q

SET: 93545.37 requests per second, p50=0.271 msec                   
GET: 94607.38 requests per second, p50=0.263 msec                   

WARNING: Could not fetch server CONFIG
SET: 1265822.75 requests per second, p50=0.335 msec        
GET: 1282051.25 requests per second, p50=0.319 msec                    

WARNING: Could not fetch server CONFIG
SET: 1992032.00 requests per second, p50=0.175 msec                     
GET: 1996008.00 requests per second, p50=0.135 msec      
 лиент и сервер работают на одной машине и дел€т €дра. ÷ифра 2 млн Ч это нижн€€ оценка возможностей сервера, а не его предел.
∆елезо: Acer Extensa, Ubuntu, клиент на той же машине.