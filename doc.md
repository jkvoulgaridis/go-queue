## Gqueue -- Project for tasks queues in golang 

## Componenets

* TaskRegistry: Adds/ removes tasks available for broker and worker to execute. Initialized once on application startup 
* TaskQueue: Keeps all tasks pending to be executed. Supports:
    * Submit: Adds a tasks to the end of the Queue. Task name should be present in hte registered task
    * Asssign: Worker sends a me


## Networking: 

Messages are user facing (from application level) and system  (sent from broker to worker and vice verca). Messages: 

* Submit: Public Facing, sent from anyome using the library to the broker to add a new tasks to the Queue. 
* Pull: System msg, worker sends to master to ask for messages. One param N, up to how many tasks worker will pick up
* Assign: System: Master sends one or more tasks to the worker. Each task is picjed up rom the queue, added to a temp "in flight" queue, 
    is assigmed a UIID7 id and a timestamp. 
* ACK: System, worker exeutes task, if successfull sends ACK with the task id and task is remoevd from the in flight. 
* NACK: System, on task execution failure, worker sends this and master in response moves the task to a. failed tasks list named DLQ 


## Setup
* The whole logic for the library is in the /internal. Here all the message handling should be handled. 
* /example contains a sample use case, handlers defined. How to run it? 

## Questions

* How can the appp configure which taskss to expose as messages and how to pass them to broker to register them as handlers and the worker to know as well? 

