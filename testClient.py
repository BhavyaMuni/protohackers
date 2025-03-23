import socket

host = socket.gethostname()
port = 10006                   # The same port as used by the server
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect((host, port))

# s.sendall(b'\x10\x00\x00\x00')
s.sendall(b'\x40\x00\x00\x00\x0A')
s.sendall(b'\x80\x00\x42\x00\x64\x00\x3c') # IAmCamera 1
s.sendall(b'\x80\x00\x42\x00\x64\x00\x3c') # IAmCamera 2
s.sendall(b'\x81\x01\x00\x42') # IAmDispatcher
s.sendall(b'\x20\x04\x55\x4e\x31\x58\x00\x00\x03\xe8')
s.sendall(b'\x20\x04\x55\x4e\x31\x58\x00\x00\x03\xe8')

while True:
    response = s.recv(1024)
    print(''.join('{:02x}'.format(x) for x in response))
