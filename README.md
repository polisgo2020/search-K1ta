# search-K1ta

# Benchmarks

## Index.Build

No goroutines:
```
BenchmarkBuild_SimpleTest/few_texts,_lot_identical_words-4         	      10	 100572572 ns/op
BenchmarkBuild_SimpleTest/few_texts,_lot_different_words-4         	       4	 306221418 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_identical_words-4         	       5	 205309372 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_different_words-4         	       2	 595846354 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_identical_words-4         	       5	 217976917 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_different_words-4         	       4	 293552779 ns/op
```

Goroutines for each text with mutexes:
```
BenchmarkBuild_SimpleTest/few_texts,_lot_identical_words-4         	       7	 146568713 ns/op
BenchmarkBuild_SimpleTest/few_texts,_lot_different_words-4         	       4	 317066828 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_identical_words-4         	       3	 447675132 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_different_words-4         	       1	1550963764 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_identical_words-4         	       2	 535026034 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_different_words-4         	       2	 778628443 ns/op
```

Goroutines for each text with merging:
```
BenchmarkBuild_SimpleTest/few_texts,_lot_identical_words-4         	      24	  47877207 ns/op
BenchmarkBuild_SimpleTest/few_texts,_lot_different_words-4         	       2	 633592157 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_identical_words-4         	       8	 128835411 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_different_words-4         	       1	1343868627 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_identical_words-4         	       2	 631027274 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_different_words-4         	       2	 572615981 ns/op
```

Goroutines for each text with straight-forward merging:
```
BenchmarkBuild_SimpleTest/few_texts,_lot_identical_words-4         	      21	  49184865 ns/op
BenchmarkBuild_SimpleTest/few_texts,_lot_different_words-4         	       2	 603161220 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_identical_words-4         	       8	 130800840 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_lot_different_words-4         	       1	1415940423 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_identical_words-4         	       2	 586184468 ns/op
BenchmarkBuild_SimpleTest/lot_texts,_few_different_words-4         	       2	 588764880 ns/op
```

## Index.Find

No goroutines
```
BenchmarkIndex_Find/lot_texts,_few_words,_100-words_phrase-4         	     123	   9972986 ns/op
BenchmarkIndex_Find/lot_texts,_few_words,_500-words_phrase-4         	      24	  49493559 ns/op
BenchmarkIndex_Find/few_texts,_lot_words,_100-words_phrase-4         	       9	 122293310 ns/op
BenchmarkIndex_Find/few_texts,_lot_words,_500-words_phrase-4         	       2	 606485413 ns/op
```

With goroutines
```
BenchmarkIndex_Find/lot_texts,_few_words,_100-words_phrase-4         	      43	  27244357 ns/op
BenchmarkIndex_Find/lot_texts,_few_words,_500-words_phrase-4         	      10	 110494054 ns/op
BenchmarkIndex_Find/few_texts,_lot_words,_100-words_phrase-4         	       4	 317114024 ns/op
BenchmarkIndex_Find/few_texts,_lot_words,_500-words_phrase-4         	       1	1335210932 ns/op
```

## Main.GetTextsAndTitlesFromDir

No goroutines
```
BenchmarkGetTextsAndTitlesFromFiles/read_a_lot_of_small_files-4             14      77076962 ns/op
BenchmarkGetTextsAndTitlesFromFiles/read_a_few_of_large_files-4             655     1772980 ns/op
```

With goroutines
```
BenchmarkGetTextsAndTitlesFromDir/read_a_lot_of_small_files-4               20      55243215 ns/op
BenchmarkGetTextsAndTitlesFromDir/read_a_few_of_large_files-4               694     1723631 ns/op
```
