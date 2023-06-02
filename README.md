project requirements golang version 1.20

in order to run the project go to the main dir where you saved the project and run in the terminal "go run main.go"

project port 8000, in order to change the port just change the value of PORT in ".env" file

project end points

GET /api/v1/stats
GET /api/v1/similar?word={wanted word}
short description of the algorithm used to solve the problem: my data source is a file of words in the project i create a map in which the key is a word with its characters sorted and the value is a slice of all the permutations of the word found in the file so basiclly on loadtime i read all the words in the file and each word i create a sorted version of it for the key in my map by which i add the word to the wanted slice in my map by doing so when i will get rest requests for "/api/v1/similar" end point all i need to do is to sort the word to create the key to get all the permutations of the word