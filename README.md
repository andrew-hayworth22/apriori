# Apriori Algorithm
An implementation of the Apriori Algorithm for CS63015 Data Mining Techniques

This algorithm is used to efficiently scan transactional databases for frequent item sets. This example takes in a file of structured data and returns all itemsets that are above the provided minimum support level.

## How to Run
To run this program, first build the binary with this command:

```
go build
```

Then, run the binary with:
```
./apriori [filename] [minimum support level]
```