# Data collection
## Launching the data collection
To execute the data collection you can execute `runserver.sh` in the docker container for the server and `collectdata.sh` in the client container. Then the data collection will begin.

## Transforming the data from pcap to CSV
For this end, from any terminal in the correct folder run `dataToCSV.sh`.

**Word of caution**: Sometimes there can be some issues arising with folder ownership because of the containers, this has to be solved manually using `chown`.

## Seeing the creation the model 
The computation of the model can be witnessed by running the jupyter notebook `classifier.ipynb` and see the results and processing in it. There can be some errors arising there because of the way the data was handled by us (some issued had to be mitigated after the data collection).