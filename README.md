# web_scraping_golang

At the time, this project is a work in progress. Here I will detail the current situation of the different required parts of the project.

Launching instruction :
consumer.go and producer.go must be compiled (go build ....go) and launched. consumer must be start in first, since, producer will
send data that consumer will... consume.

The producer :
It pass through all pages of listings from the first, and extract for each, each offer URL.
It lack a blacklist, since one URL that does not lead to an offer, look like one that does, but should not be treated further.
It does send each url through a kafka broker, with the topic web-adresses.
A kafka broker must be started to work.

The consumer :
It does read from a topic, download a page from the url read, and save its data to a mongodb database (local, with docker in my case)
The data currently saved is only the url itself. In the end, I would like to save all the field bellow "<div class="small-12 columns">&nbsp;</div>"
since, everything outside of this box is irrelevant to the offer itself. But, we could go by every field and save them all with the url if wanted.

MongoDB local, thanks to docker :
The database is launched with this line, after installing docker
sudo docker run -d -p 27017-27019:27017-27019 --name mongodb mongo:4.0.4
The mongodb shell cann be accessed with:  sudo docker exec -it mongodb bash
In it, you can see the saved datas.

Global note :
Thank you for the enanced delay.
At least I have been able to make and end to end project.
I hope you can see what you wanted to see in this almost finished work.
This version is still not as complete as I would have want to give you, but I'm still not completely healthy either. So commiting/pushin in case of no more 
improvement.

