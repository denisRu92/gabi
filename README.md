A small web service for printing similar words in the English language.
Two words w_1 and w_2 are considered similar if w_1 is a letter permutation of w_2 (e.g., "stressed" and "desserts").

We provide a DB of the English dictionary you should work with. Please find it in the attached archive. The service can expect the DB (the txt file) to be in the local directory with the same name.

The web service should listen on port 8000 and support the following two HTTP endpoints:

- GET /api/v1/similar?word=stressed
  Returns all words in the dictionary that have the same permutation as the word "stressed". The word in the query should not be returned.

The result format is a JSON object as follows:
{
similar:[list,of,words,that,are,similar,to,provided,word]
}

For example:
http://localhost:8000/api/v1/similar?word=apple
{"similar":["appel","pepla"]}


- GET /api/v1/stats
  Return general statistics about the program:
  Total number of words in the dictionary
  Total number of requests (not including "stats" requests)
  Average time for request handling in nano seconds (not including "stats" requests)

The output is a JSON object structured as follows:
{
totalWords:int
totalRequests:int
avgProcessingTimeNs:int
}

For example:
http://localhost:8000/api/v1/stats
{"totalWords":351075,"totalRequests":9,"avgProcessingTimeNs":45239}


Misc / hints:
Please shortly describe the algorithm you use to solve the problem as part of the submission.
Use any popular programming language you like, but please do think about CPU and memory optimization or the possible trade-offs.
Please write the task under Linux dev environment.
Please do expect a high rate of requests and requests arriving in parallel.
Please write all errors/logs to stdout/stderr.
Program should be submitted as a complete runnable unit with instructions on how to build and run it (including the source code). You can use Github but please make the repo private in this case. 