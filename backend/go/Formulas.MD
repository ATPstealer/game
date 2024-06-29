## Production formula:

One production cycle = blueprint time * ( workers building demand / ( workers hired + equipment impact))  
The equipment adds speed to production. One worker is taken as a unit of change.  
The efficiency of using equipment depends on the number of workers. To get maximum efficiency you need to have a ratio:
Number of workers >= number of equipment.
If there is a shortage of workers, efficiency decreases proportionately.

## Hiring formula:    
SUM for n buildings    
∑ hn (1 - a * (sn / sm)) = p    
hn - demand in one building (HiringNeeds).    
a - alpha - coefficient which determined how much low salary buildings will have employers.    
sn - salary building    
sm - salary max in cell    
p - population in cell     

alpha = (∑ hn - p) / ∑ (hn * (sn / sm - 1))     

## Cell population demand and max spend    
Counting one time a day.    

Average price in cell determined as revenue for this product / sell count last day.    

It's determined by steps. For example, we have 3 product in cell Water, Rise, Apples.    
Water demand is 10, average price is 1    
Rise demand is 3, average price is 10    
Apples demand is 1, average price is 40    

Population is 1000, Population savings 500000, spend rate is 0.1.    
That means that next day population in cell can spend 50000.    
Total demand:    
Water demand is 10000, average price is 1, total price 10000    
Rise demand is 3000, average price is 10, total price 30000    
Apples demand is 1000, average price is 40, total price 40000    

Make steps of demand:    
First step 10000:    
Water demand satisfied, we can sell 10000 units of Water by total price 10000    
Rise demand particular satisfied, we can sell 300 units of Rise by total price 10000    
Apples demand particular satisfied, we can sell 250 units of Rise by total price 10000    
On first step we count total spend 30000    

Second stop 30000:    
We have 50000 - 30000 (last step) = 20000    
It last step because we can't satisfy demand on this step. Now we split 20000 on 2 products:    
Rise demand particular satisfied, we can sell 300 units of Rise by total price 10000    
Apples demand particular satisfied, we can sell 250 units of Rise by total price 10000    

Total spend in this cell:    
Water 10000    
Rise 20000    
Apples 20000    

## Store day sell    

Capacity is how many products can sell this store with 1 level 1 square if they have maximum employers.    
Day sell limit = capacity * 0.75 * priceAverage / price    
priceAverage is average price in cell    
price - price of this product    



