import random
import socket


# create an INET, STREAMing socket
# bind the socket to a public host, and a well-known port
# become a server socket
# accept connections from outside
# now do something with the clientsocket
# in this case, we'll pretend this is a threaded server

def getWord(length):
    word = ""
    for i in range(length):
        word += chr(ord('a') + random.randint(0, 25))
    return word


def getLine(length):
    line = ""
    for i in range(length):
        wordLength = random.randint(3, 5)
        line += getWord(wordLength)
        line += ' '
    line += '\n'

    return line


random.seed(9999)

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sockfd:
    sevrAddr = ("localhost", 9999)
    sockfd.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sockfd.bind(sevrAddr)
    sockfd.listen(5)
    while True:
        clientsocket, address = sockfd.accept()
        with clientsocket:
            for i in range(100):
                try:
                    lineLength = random.randint(5, 15)
                    line = getLine(lineLength)
                    clientsocket.send(line.encode())
                except BrokenPipeError:
                    exit(1)

        break
