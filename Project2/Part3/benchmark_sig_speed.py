



import timeit
mysetup = """

from your_code import Server,Client
server = Server()
client = Client()

#Generating the keys
pk_serialized, sk_serialized = server.generate_ca("a,b,c")
    
#Registering the user on the server
issuance_request, private_state = client.prepare_registration(pk_serialized, "weewoo", "a,b,c")

response = server.register(pk_serialized, issuance_request, "weewoo", "a,b,c")

credential = client.proceed_registration_response(sk_serialized,response, private_state)
m = b"some message for test"
"""

toExec = """
client.sign_request(pk_serialized, credential, m,"a,b")
"""
from statistics import stdev, mean
data = timeit.repeat(setup=mysetup,stmt=toExec,repeat=500, number=3)

print(mean(data))
print(stdev(data))




mysetup = """

from your_code import Server
server = Server()
"""

toExec = """
pk_serialized, sk_serialized = server.generate_ca("a,b,c")
"""
from statistics import stdev, mean
data = timeit.repeat(setup=mysetup,stmt=toExec,repeat=500, number=3)

print(mean(data))
print(stdev(data))