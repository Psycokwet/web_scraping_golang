# web_scraping_golang

At the time, this project is a work in progress. Here I will detail the current situation of the different required parts of the project.

Launching instruction :
consumer.go and producer.go must be compiled (go build ....go) and launched. consumer must be start in first, since, producer will
send data that consumer will... consume.

The producer :
It pass through all pages of listings from the first, and exract for each, each offer URL.
It lack a blacklist, since one URL that does not lead to an offer, look like one that does, but should not be treated further.
It does send each url through a kafka broker, with the topic web-adresses.
A kafka broker must be started to work.

The consumer :
It does only read from a topic and print the URL in the standard output.
If I may not be able to complete this part, the next potential questions would have been which criteria would I have considered to 
save data from an offer?
Since some data may be irrelevant to the offer itself. The project subject says to save a maximum of data from an offert, so, 
one possible approach to this question may have been to save everything, and to worry about the treatment later (outside of the project). 
I think I would have chose this last possibility.


MongoDB Atlas :
Currently does not allow connection for an unknown reason. Neither from go drivers, or from MongoDB graphic official tools.
Once this problem is solved, adding data to the database wouldn't be a problem.

Global note :
I would have made a more comfortable code if I had the time I was supposed to get for the project.
I'm not sure if the current state of the project is of any relevance to you to be the judge of my qualities as a coder,
but regarding my current situation, it's better than nothing.
