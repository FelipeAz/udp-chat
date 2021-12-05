# udp-chat

UDP-CHAT is a chat made in Golang which simulates a real-time chat using a UDP server. 

#### Business Rules
- The last 20 messages are stored into a Redis with an expiry time of 20 minutes.
- If the server took more than 5 seconds to reply, the connection is closed (timeout set as 5 sec)

## PROJECT REQUIREMENTS

This project requires only ``DOCKER`` installed in your machine.

## Ports Availability
``6380`` - REDIS

``8000`` - SERVER

## INSTALLATION

To install & run the chat is pretty simple, all the steps have been added to a shell script and we can run it as a makefile. 

### STEP-BY-STEP

On the main project folder ``{project_path}/udp-chat``, run the following commands:

 ``make server`` -> Initialize Docker adding the ``REDIS`` container and running the ``SERVER`` on port ``8000``
 
``make client`` -> Runs a ``CLIENT`` instance which listen on server port 8000. 

If you want to stop the server, go to the project folder again and type ``make server-stop``

### Debugging cache operations

You can see every action made on REDIS using the client ``redis-cli``. 

Just open a new terminal and type: ``redis-cli -p 6380`` then ``monitor``. If you get an ``OK`` you will see every action made on REDIS.